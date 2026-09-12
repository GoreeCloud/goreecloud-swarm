package bencode

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalid       = errors.New("invalid bencode")
	ErrLimitExceeded = errors.New("bencode limit exceeded")
)

type Kind uint8

const (
	KindInteger Kind = iota + 1
	KindBytes
	KindList
	KindDictionary
)

type Node struct {
	Kind  Kind
	Int   int64
	Bytes []byte
	List  []*Node
	Dict  map[string]*Node
	Start int
	End   int
}

type Limits struct {
	MaxInputBytes     int
	MaxDepth          int
	MaxStringBytes    int
	MaxContainerItems int
}

func DefaultLimits() Limits {
	return Limits{
		MaxInputBytes:     16 << 20,
		MaxDepth:          64,
		MaxStringBytes:    8 << 20,
		MaxContainerItems: 100_000,
	}
}

func Decode(data []byte, limits Limits) (*Node, error) {
	if limits.MaxInputBytes <= 0 || limits.MaxDepth <= 0 || limits.MaxStringBytes <= 0 || limits.MaxContainerItems <= 0 {
		return nil, fmt.Errorf("%w: invalid decoder limits", ErrLimitExceeded)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty input", ErrInvalid)
	}
	if len(data) > limits.MaxInputBytes {
		return nil, fmt.Errorf("%w: input is %d bytes", ErrLimitExceeded, len(data))
	}

	d := decoder{data: data, limits: limits}
	node, err := d.value(0)
	if err != nil {
		return nil, err
	}
	if d.pos != len(data) {
		return nil, fmt.Errorf("%w at offset %d: trailing data", ErrInvalid, d.pos)
	}
	return node, nil
}

type decoder struct {
	data   []byte
	pos    int
	limits Limits
}

func (d *decoder) value(depth int) (*Node, error) {
	if depth > d.limits.MaxDepth {
		return nil, fmt.Errorf("%w at offset %d: maximum depth exceeded", ErrLimitExceeded, d.pos)
	}
	if d.pos >= len(d.data) {
		return nil, fmt.Errorf("%w at offset %d: unexpected end", ErrInvalid, d.pos)
	}

	switch b := d.data[d.pos]; {
	case b == 'i':
		return d.integer()
	case b == 'l':
		return d.list(depth)
	case b == 'd':
		return d.dictionary(depth)
	case b >= '0' && b <= '9':
		return d.byteString()
	default:
		return nil, fmt.Errorf("%w at offset %d: unexpected token", ErrInvalid, d.pos)
	}
}

func (d *decoder) integer() (*Node, error) {
	start := d.pos
	d.pos++ // i
	valueStart := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		d.pos++
	}
	if d.pos >= len(d.data) {
		return nil, fmt.Errorf("%w at offset %d: unterminated integer", ErrInvalid, start)
	}
	if d.pos == valueStart {
		return nil, fmt.Errorf("%w at offset %d: empty integer", ErrInvalid, start)
	}

	raw := d.data[valueStart:d.pos]
	if raw[0] == '+' {
		return nil, fmt.Errorf("%w at offset %d: non-canonical integer", ErrInvalid, start)
	}
	if raw[0] == '-' {
		if len(raw) == 1 || (len(raw) > 1 && raw[1] == '0') {
			return nil, fmt.Errorf("%w at offset %d: non-canonical integer", ErrInvalid, start)
		}
	} else if len(raw) > 1 && raw[0] == '0' {
		return nil, fmt.Errorf("%w at offset %d: non-canonical integer", ErrInvalid, start)
	}

	value, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w at offset %d: invalid integer", ErrInvalid, start)
	}
	d.pos++ // e
	return &Node{Kind: KindInteger, Int: value, Start: start, End: d.pos}, nil
}

func (d *decoder) byteString() (*Node, error) {
	start := d.pos
	lengthStart := d.pos
	for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
		d.pos++
	}
	if d.pos >= len(d.data) || d.data[d.pos] != ':' {
		return nil, fmt.Errorf("%w at offset %d: invalid byte-string length", ErrInvalid, start)
	}
	lengthRaw := d.data[lengthStart:d.pos]
	if len(lengthRaw) == 0 || (len(lengthRaw) > 1 && lengthRaw[0] == '0') {
		return nil, fmt.Errorf("%w at offset %d: non-canonical byte-string length", ErrInvalid, start)
	}
	length64, err := strconv.ParseInt(string(lengthRaw), 10, 64)
	if err != nil || length64 < 0 {
		return nil, fmt.Errorf("%w at offset %d: invalid byte-string length", ErrInvalid, start)
	}
	if length64 > int64(d.limits.MaxStringBytes) {
		return nil, fmt.Errorf("%w at offset %d: byte string is too large", ErrLimitExceeded, start)
	}
	length := int(length64)
	d.pos++ // :
	if length > len(d.data)-d.pos {
		return nil, fmt.Errorf("%w at offset %d: truncated byte string", ErrInvalid, start)
	}
	value := append([]byte(nil), d.data[d.pos:d.pos+length]...)
	d.pos += length
	return &Node{Kind: KindBytes, Bytes: value, Start: start, End: d.pos}, nil
}

func (d *decoder) list(depth int) (*Node, error) {
	start := d.pos
	d.pos++ // l
	items := make([]*Node, 0)
	for {
		if d.pos >= len(d.data) {
			return nil, fmt.Errorf("%w at offset %d: unterminated list", ErrInvalid, start)
		}
		if d.data[d.pos] == 'e' {
			d.pos++
			return &Node{Kind: KindList, List: items, Start: start, End: d.pos}, nil
		}
		if len(items) >= d.limits.MaxContainerItems {
			return nil, fmt.Errorf("%w at offset %d: too many list items", ErrLimitExceeded, start)
		}
		item, err := d.value(depth + 1)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
}

func (d *decoder) dictionary(depth int) (*Node, error) {
	start := d.pos
	d.pos++ // d
	entries := make(map[string]*Node)
	for {
		if d.pos >= len(d.data) {
			return nil, fmt.Errorf("%w at offset %d: unterminated dictionary", ErrInvalid, start)
		}
		if d.data[d.pos] == 'e' {
			d.pos++
			return &Node{Kind: KindDictionary, Dict: entries, Start: start, End: d.pos}, nil
		}
		if len(entries) >= d.limits.MaxContainerItems {
			return nil, fmt.Errorf("%w at offset %d: too many dictionary entries", ErrLimitExceeded, start)
		}
		keyNode, err := d.byteString()
		if err != nil {
			return nil, err
		}
		key := string(keyNode.Bytes)
		if _, exists := entries[key]; exists {
			return nil, fmt.Errorf("%w at offset %d: duplicate dictionary key", ErrInvalid, keyNode.Start)
		}
		value, err := d.value(depth + 1)
		if err != nil {
			return nil, err
		}
		entries[key] = value
	}
}
