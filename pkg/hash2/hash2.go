package hash2

import (
	"bufio"
	"hash/crc32"
	"io"
	"strings"
)

type H struct {
	state [32]int
}

func (h *H) Add(x uint32) {
	for i := 0; i < 32; i++ {
		b := int(^x & 1)
		h.state[i] += 1 - (b << 1)
		x >>= 1
	}
}

func (h *H) Sum32() uint32 {
	x := uint32(0)
	for i := 31; i >= 0; i-- {
		x <<= 1
		if h.state[i] > 0 {
			x |= 1
		}
	}
	return x
}

func Trigrams(r io.Reader) (map[string]struct{}, error) {
	br := bufio.NewReader(r)
	m := make(map[string]struct{})
	for {
		line, err := br.ReadString('\n')
		if line != "" || err != nil {
			line = strings.TrimRight(line, "\r\n")
			for i := 0; i < len(line)-2; i++ {
				trigram := line[i : i+3]
				m[trigram] = struct{}{}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func Hash2(r io.Reader) (uint32, error) {
	br := bufio.NewReader(r)
	var h H
	for {
		line, err := br.ReadString('\n')
		if line != "" || err != nil {
			line = strings.TrimRight(line, "\r\n")
			for i := 0; i < len(line)-2; i++ {
				trigram := line[i : i+3]
				h.Add(crc32.ChecksumIEEE([]byte(trigram)))
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return h.Sum32(), nil
}

func Hash3(r io.Reader) (int, error) {
	trigrams, err := Trigrams(r)
	if err != nil {
		return 0, err
	}
	return len(trigrams), nil
}
