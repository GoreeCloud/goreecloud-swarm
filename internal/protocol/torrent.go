package protocol

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/GoreeCloud/goreecloud-swarm/internal/protocol/bencode"
)

var ErrInvalidTorrent = errors.New("invalid torrent metainfo")

type TorrentVersion string

const (
	TorrentV1     TorrentVersion = "v1"
	TorrentV2     TorrentVersion = "v2"
	TorrentHybrid TorrentVersion = "hybrid"
)

type TorrentMeta struct {
	Name         string         `json:"name"`
	Version      TorrentVersion `json:"version"`
	Private      bool           `json:"private"`
	PieceLength  int64          `json:"piece_length"`
	FileCount    int            `json:"file_count"`
	TotalSize    int64          `json:"total_size"`
	Trackers     []string       `json:"trackers"`
	InfoHashV1   string         `json:"info_hash_v1,omitempty"`
	InfoHashV2   string         `json:"info_hash_v2,omitempty"`
	MetainfoSize int            `json:"metainfo_size"`
}

func ParseTorrent(data []byte) (TorrentMeta, error) {
	root, err := bencode.Decode(data, bencode.DefaultLimits())
	if err != nil {
		return TorrentMeta{}, fmt.Errorf("%w: %v", ErrInvalidTorrent, err)
	}
	if root.Kind != bencode.KindDictionary {
		return TorrentMeta{}, fmt.Errorf("%w: root must be a dictionary", ErrInvalidTorrent)
	}
	info := root.Dict["info"]
	if info == nil || info.Kind != bencode.KindDictionary {
		return TorrentMeta{}, fmt.Errorf("%w: missing info dictionary", ErrInvalidTorrent)
	}
	if info.Start < 0 || info.End > len(data) || info.Start >= info.End {
		return TorrentMeta{}, fmt.Errorf("%w: invalid info dictionary span", ErrInvalidTorrent)
	}

	nameNode := info.Dict["name"]
	if nameNode == nil || nameNode.Kind != bencode.KindBytes || len(nameNode.Bytes) == 0 {
		return TorrentMeta{}, fmt.Errorf("%w: missing torrent name", ErrInvalidTorrent)
	}

	pieceLengthNode := info.Dict["piece length"]
	if pieceLengthNode == nil || pieceLengthNode.Kind != bencode.KindInteger || pieceLengthNode.Int <= 0 {
		return TorrentMeta{}, fmt.Errorf("%w: invalid piece length", ErrInvalidTorrent)
	}

	private := false
	if privateNode := info.Dict["private"]; privateNode != nil {
		if privateNode.Kind != bencode.KindInteger || (privateNode.Int != 0 && privateNode.Int != 1) {
			return TorrentMeta{}, fmt.Errorf("%w: invalid private flag", ErrInvalidTorrent)
		}
		private = privateNode.Int == 1
	}

	piecesNode := info.Dict["pieces"]
	hasV1 := piecesNode != nil
	if hasV1 {
		if piecesNode.Kind != bencode.KindBytes || len(piecesNode.Bytes) == 0 || len(piecesNode.Bytes)%sha1.Size != 0 {
			return TorrentMeta{}, fmt.Errorf("%w: invalid v1 pieces field", ErrInvalidTorrent)
		}
	}

	metaVersionNode := info.Dict["meta version"]
	hasV2 := metaVersionNode != nil
	if hasV2 {
		if metaVersionNode.Kind != bencode.KindInteger || metaVersionNode.Int != 2 {
			return TorrentMeta{}, fmt.Errorf("%w: unsupported meta version", ErrInvalidTorrent)
		}
		if info.Dict["file tree"] == nil || info.Dict["file tree"].Kind != bencode.KindDictionary {
			return TorrentMeta{}, fmt.Errorf("%w: v2 torrent is missing file tree", ErrInvalidTorrent)
		}
		if pieceLengthNode.Int < 16*1024 || !isPowerOfTwo(pieceLengthNode.Int) {
			return TorrentMeta{}, fmt.Errorf("%w: v2 piece length must be a power of two of at least 16 KiB", ErrInvalidTorrent)
		}
	}

	if !hasV1 && !hasV2 {
		return TorrentMeta{}, fmt.Errorf("%w: unsupported or unrecognized torrent format", ErrInvalidTorrent)
	}

	version := TorrentV1
	switch {
	case hasV1 && hasV2:
		version = TorrentHybrid
	case hasV2:
		version = TorrentV2
	}

	fileCount, totalSize, err := torrentSize(info, hasV2)
	if err != nil {
		return TorrentMeta{}, err
	}

	meta := TorrentMeta{
		Name:         string(nameNode.Bytes),
		Version:      version,
		Private:      private,
		PieceLength:  pieceLengthNode.Int,
		FileCount:    fileCount,
		TotalSize:    totalSize,
		Trackers:     extractTrackers(root),
		MetainfoSize: len(data),
	}

	rawInfo := data[info.Start:info.End]
	if hasV1 {
		hash := sha1.Sum(rawInfo)
		meta.InfoHashV1 = hex.EncodeToString(hash[:])
	}
	if hasV2 {
		hash := sha256.Sum256(rawInfo)
		meta.InfoHashV2 = hex.EncodeToString(hash[:])
	}
	return meta, nil
}

func torrentSize(info *bencode.Node, hasV2 bool) (int, int64, error) {
	if hasV2 {
		return v2TreeSize(info.Dict["file tree"])
	}

	lengthNode, hasLength := info.Dict["length"]
	filesNode, hasFiles := info.Dict["files"]
	if hasLength == hasFiles {
		return 0, 0, fmt.Errorf("%w: v1 info must contain exactly one of length or files", ErrInvalidTorrent)
	}
	if hasLength {
		if lengthNode.Kind != bencode.KindInteger || lengthNode.Int < 0 {
			return 0, 0, fmt.Errorf("%w: invalid file length", ErrInvalidTorrent)
		}
		return 1, lengthNode.Int, nil
	}
	if filesNode.Kind != bencode.KindList || len(filesNode.List) == 0 {
		return 0, 0, fmt.Errorf("%w: invalid files list", ErrInvalidTorrent)
	}
	var total int64
	for _, file := range filesNode.List {
		if file.Kind != bencode.KindDictionary {
			return 0, 0, fmt.Errorf("%w: invalid file entry", ErrInvalidTorrent)
		}
		length := file.Dict["length"]
		path := file.Dict["path"]
		if length == nil || length.Kind != bencode.KindInteger || length.Int < 0 {
			return 0, 0, fmt.Errorf("%w: invalid file length", ErrInvalidTorrent)
		}
		if path == nil || path.Kind != bencode.KindList || len(path.List) == 0 {
			return 0, 0, fmt.Errorf("%w: invalid file path", ErrInvalidTorrent)
		}
		for _, segment := range path.List {
			if segment.Kind != bencode.KindBytes || len(segment.Bytes) == 0 {
				return 0, 0, fmt.Errorf("%w: invalid file path segment", ErrInvalidTorrent)
			}
		}
		if length.Int > 0 && total > (1<<63-1)-length.Int {
			return 0, 0, fmt.Errorf("%w: total size overflow", ErrInvalidTorrent)
		}
		total += length.Int
	}
	return len(filesNode.List), total, nil
}

func v2TreeSize(tree *bencode.Node) (int, int64, error) {
	if tree == nil || tree.Kind != bencode.KindDictionary || len(tree.Dict) == 0 {
		return 0, 0, fmt.Errorf("%w: invalid v2 file tree", ErrInvalidTorrent)
	}
	var count int
	var total int64
	var walk func(*bencode.Node) error
	walk = func(node *bencode.Node) error {
		if node.Kind != bencode.KindDictionary {
			return fmt.Errorf("%w: invalid v2 tree node", ErrInvalidTorrent)
		}
		if leaf, ok := node.Dict[""]; ok {
			if len(node.Dict) != 1 || leaf.Kind != bencode.KindDictionary {
				return fmt.Errorf("%w: invalid v2 file leaf", ErrInvalidTorrent)
			}
			length := leaf.Dict["length"]
			if length == nil || length.Kind != bencode.KindInteger || length.Int < 0 {
				return fmt.Errorf("%w: invalid v2 file length", ErrInvalidTorrent)
			}
			if length.Int > 0 && total > (1<<63-1)-length.Int {
				return fmt.Errorf("%w: total size overflow", ErrInvalidTorrent)
			}
			total += length.Int
			count++
			return nil
		}
		for name, child := range node.Dict {
			if name == "" || strings.ContainsRune(name, '\x00') {
				return fmt.Errorf("%w: invalid v2 path segment", ErrInvalidTorrent)
			}
			if err := walk(child); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(tree); err != nil {
		return 0, 0, err
	}
	if count == 0 {
		return 0, 0, fmt.Errorf("%w: v2 file tree has no files", ErrInvalidTorrent)
	}
	return count, total, nil
}

func extractTrackers(root *bencode.Node) []string {
	seen := make(map[string]struct{})
	trackers := make([]string, 0)
	add := func(node *bencode.Node) {
		if node == nil || node.Kind != bencode.KindBytes {
			return
		}
		candidate := strings.TrimSpace(string(node.Bytes))
		if !validTracker(candidate) {
			return
		}
		if _, exists := seen[candidate]; exists {
			return
		}
		seen[candidate] = struct{}{}
		trackers = append(trackers, candidate)
	}
	add(root.Dict["announce"])
	if tiers := root.Dict["announce-list"]; tiers != nil && tiers.Kind == bencode.KindList {
		for _, tier := range tiers.List {
			if tier.Kind != bencode.KindList {
				continue
			}
			for _, tracker := range tier.List {
				add(tracker)
			}
		}
	}
	return trackers
}

func validTracker(raw string) bool {
	if raw == "" {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return false
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https", "udp":
		return true
	default:
		return false
	}
}

func isPowerOfTwo(v int64) bool {
	return v > 0 && (v&(v-1)) == 0
}
