package protocol

import "testing"

func TestParseMagnetBTIH(t *testing.T) {
	raw := "magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567&dn=Example&tr=https%3A%2F%2Ftracker.example%2Fannounce"
	m, err := ParseMagnet(raw)
	if err != nil {
		t.Fatalf("ParseMagnet() error = %v", err)
	}
	if m.DisplayName != "Example" {
		t.Fatalf("display name = %q", m.DisplayName)
	}
	if len(m.Trackers) != 1 {
		t.Fatalf("trackers = %d", len(m.Trackers))
	}
}

func TestParseMagnetRejectsUnsupportedScheme(t *testing.T) {
	_, err := ParseMagnet("https://example.com/file.torrent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseMagnetRejectsMissingTopic(t *testing.T) {
	_, err := ParseMagnet("magnet:?dn=Example")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseMagnetBTMH(t *testing.T) {
	raw := "magnet:?xt=urn:btmh:1220aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if _, err := ParseMagnet(raw); err != nil {
		t.Fatalf("ParseMagnet() error = %v", err)
	}
}
