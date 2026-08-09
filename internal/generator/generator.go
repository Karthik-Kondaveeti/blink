package generator

type Generator interface {
	// Function: New(workerId int64) (*Generator, error)
	Generate() uint64
	Encode(id uint64) string
	Decode(code string) (uint64, error)
}
