package protocol

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestParseTorrentV1SingleFile(t *testing.T) {
	info := "d6:lengthi4e4:name4:test12:piece lengthi16384e6:pieces20:aaaaaaaaaaaaaaaaaaaae"
	input := []byte("d4:info" + info + "e")
	meta, err := ParseTorrent(input)
	if err != nil {
		t.Fatalf("ParseTorrent() error = %v", err)
	}
	if meta.Version != TorrentV1 || meta.Name != "test" || meta.FileCount != 1 || meta.TotalSize != 4 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
	hash := sha1.Sum([]byte(info))
	if meta.InfoHashV1 != hex.EncodeToString(hash[:]) {
		t.Fatalf("v1 hash = %q", meta.InfoHashV1)
	}
	if meta.InfoHashV2 != "" {
		t.Fatalf("unexpected v2 hash = %q", meta.InfoHashV2)
	}
}

func TestParseTorrentV2SingleFile(t *testing.T) {
	fileTree := "d4:testd0:d6:lengthi4eeee"
	info := "d9:file tree" + fileTree + "12:meta versioni2e4:name4:test12:piece lengthi16384ee"
	input := []byte("d4:info" + info + "e")
	meta, err := ParseTorrent(input)
	if err != nil {
		t.Fatalf("ParseTorrent() error = %v", err)
	}
	if meta.Version != TorrentV2 || meta.FileCount != 1 || meta.TotalSize != 4 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
	hash := sha256.Sum256([]byte(info))
	if meta.InfoHashV2 != hex.EncodeToString(hash[:]) {
		t.Fatalf("v2 hash = %q", meta.InfoHashV2)
	}
}

func TestParseTorrentHybrid(t *testing.T) {
	fileTree := "d4:testd0:d6:lengthi4eeee"
	info := "d9:file tree" + fileTree + "6:lengthi4e12:meta versioni2e4:name4:test12:piece lengthi16384e6:pieces20:aaaaaaaaaaaaaaaaaaaae"
	input := []byte("d4:info" + info + "e")
	meta, err := ParseTorrent(input)
	if err != nil {
		t.Fatalf("ParseTorrent() error = %v", err)
	}
	if meta.Version != TorrentHybrid || meta.InfoHashV1 == "" || meta.InfoHashV2 == "" {
		t.Fatalf("unexpected hybrid meta: %+v", meta)
	}
}

func TestParseTorrentRejectsInvalidPieces(t *testing.T) {
	input := []byte("d4:infod6:lengthi4e4:name4:test12:piece lengthi16384e6:pieces3:abcee")
	if _, err := ParseTorrent(input); err == nil {
		t.Fatal("expected invalid pieces to fail")
	}
}

func TestParseTorrentRejectsOversizedInput(t *testing.T) {
	input := []byte(strings.Repeat("x", (16<<20)+1))
	if _, err := ParseTorrent(input); err == nil {
		t.Fatal("expected oversized metainfo to fail")
	}
}
