package snowflake

import (
	"errors"
	"sync"
	"time"
)

type Snowflake struct {
	WorkerID      uint16
	Sequence      uint16
	LastTimestamp int64

	mu sync.Mutex
}

func New(workerID uint16) (*Snowflake, error) {
	if workerID > 1023 {
		return nil, errors.New("WORKER_ID must be between 0 and 1023")
	}
	return &Snowflake{
		WorkerID: uint16(workerID),
	}, nil
}

func encodeBase62(id int64) string {
	if id == 0 {
		return "a"
	}

	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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

func (g *Snowflake) Generate() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().UnixMilli()
	workerID := int64(g.WorkerID)
	if timestamp == g.LastTimestamp && g.Sequence < 4095 {
		g.Sequence++
	} else {
		g.Sequence = 0
		g.LastTimestamp = timestamp
	}

	// timestamp(41), workerID(10), sequence(12)
	id := (timestamp << 22) | (workerID << 12) | int64(g.Sequence)

	return encodeBase62(id)
}
