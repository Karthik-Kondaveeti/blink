package snowflake

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type Snowflake struct {
	WorkerID      uint16
	Sequence      uint16
	LastTimestamp uint64

	mu sync.Mutex
}

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func New(workerID uint16) (*Snowflake, error) {
	if workerID > 1023 {
		return nil, errors.New("WORKER_ID must be between 0 and 1023")
	}
	return &Snowflake{
		WorkerID: uint16(workerID),
	}, nil
}

func (g *Snowflake) Encode(id uint64) string {
	if id == 0 {
		return "a"
	}

	var result []byte
	for id > 0 {
		result = append(result, charSet[id%62])
		id /= 62
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func (g *Snowflake) Decode(code string) (uint64, error) {
	var id uint64
	for _, ch := range code {
		value := strings.IndexRune(charSet, ch)
		if value == -1 {
			return 0, errors.New("cannot decode given string!")
		}
		id = id*62 + uint64(value)
	}
	return id, nil
}

func (g *Snowflake) Generate() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	workerID := uint64(g.WorkerID)
	timestamp := uint64(time.Now().UnixMilli())
	if timestamp == g.LastTimestamp && g.Sequence < 4095 {
		g.Sequence++
	} else {
		g.Sequence = 0
		g.LastTimestamp = timestamp
	}

	// timestamp(41), workerID(10), sequence(12)
	id := (timestamp << 22) | (workerID << 12) | uint64(g.Sequence)

	return id
}
