package protocol

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrInvalidMagnet       = errors.New("invalid magnet URI")
	ErrMissingExactTopic   = errors.New("magnet URI is missing an exact topic")
	ErrUnsupportedInfoHash = errors.New("unsupported BitTorrent info-hash encoding")
)

type Magnet struct {
	Raw         string
	ExactTopics []string
	DisplayName string
	Trackers    []string
}

func ParseMagnet(raw string) (Magnet, error) {
	raw = strings.TrimSpace(raw)
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "magnet" {
		return Magnet{}, ErrInvalidMagnet
	}

	q := u.Query()
	ext := q["xt"]
	if len(ext) == 0 {
		return Magnet{}, ErrMissingExactTopic
	}

	validBTIH := false
	validBTMH := false
	for _, topic := range ext {
		lower := strings.ToLower(topic)
		switch {
		case strings.HasPrefix(lower, "urn:btih:"):
			value := topic[len("urn:btih:"):]
			if len(value) != 40 && len(value) != 32 {
				return Magnet{}, fmt.Errorf("%w: btih length %d", ErrUnsupportedInfoHash, len(value))
			}
			validBTIH = true
		case strings.HasPrefix(lower, "urn:btmh:1220"):
			value := topic[len("urn:btmh:1220"):]
			if len(value) != 64 {
				return Magnet{}, fmt.Errorf("%w: btmh sha256 length %d", ErrUnsupportedInfoHash, len(value))
			}
			validBTMH = true
		}
	}

	if !validBTIH && !validBTMH {
		return Magnet{}, ErrUnsupportedInfoHash
	}

	trackers := make([]string, 0, len(q["tr"]))
	for _, tracker := range q["tr"] {
		tracker = strings.TrimSpace(tracker)
		if tracker == "" {
			continue
		}
		tu, err := url.Parse(tracker)
		if err != nil || tu.Scheme == "" || tu.Host == "" {
			return Magnet{}, fmt.Errorf("%w: invalid tracker URL", ErrInvalidMagnet)
		}
		switch strings.ToLower(tu.Scheme) {
		case "http", "https", "udp":
		default:
			return Magnet{}, fmt.Errorf("%w: unsupported tracker scheme %q", ErrInvalidMagnet, tu.Scheme)
		}
		trackers = append(trackers, tracker)
	}

	return Magnet{
		Raw:         raw,
		ExactTopics: append([]string(nil), ext...),
		DisplayName: strings.TrimSpace(q.Get("dn")),
		Trackers:    trackers,
	}, nil
}
