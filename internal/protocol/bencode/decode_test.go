package bencode

import (
	"errors"
	"testing"
)

func TestDecodeDictionary(t *testing.T) {
	node, err := Decode([]byte("d3:foo3:bar4:spamli1ei2eee"), DefaultLimits())
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if node.Kind != KindDictionary {
		t.Fatalf("kind = %v", node.Kind)
	}
	if got := string(node.Dict["foo"].Bytes); got != "bar" {
		t.Fatalf("foo = %q", got)
	}
	if got := len(node.Dict["spam"].List); got != 2 {
		t.Fatalf("spam items = %d", got)
	}
}

func TestDecodeRejectsNonCanonicalInteger(t *testing.T) {
	for _, input := range []string{"i03e", "i-0e", "i-03e", "i+1e"} {
		if _, err := Decode([]byte(input), DefaultLimits()); err == nil {
			t.Fatalf("expected %q to fail", input)
		}
	}
}

func TestDecodeRejectsDuplicateDictionaryKey(t *testing.T) {
	_, err := Decode([]byte("d1:ai1e1:ai2ee"), DefaultLimits())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("error = %v, want ErrInvalid", err)
	}
}

func TestDecodeHonorsInputLimit(t *testing.T) {
	limits := DefaultLimits()
	limits.MaxInputBytes = 3
	_, err := Decode([]byte("4:test"), limits)
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("error = %v, want ErrLimitExceeded", err)
	}
}

func TestDecodeRecordsRawSpan(t *testing.T) {
	input := []byte("d4:infod4:name4:testee")
	node, err := Decode(input, DefaultLimits())
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	info := node.Dict["info"]
	if got := string(input[info.Start:info.End]); got != "d4:name4:teste" {
		t.Fatalf("raw info = %q", got)
	}
}
