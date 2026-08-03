package generator

type Generator interface {
	// Function: New(workerId int64) (*Generator, error)
	Generate() string
}
