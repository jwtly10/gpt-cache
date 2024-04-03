package main

// Parser is an interface that defines the behavior of a parser.
type Parser interface {
	// Parse takes a request body as input and returns a slice of strings and an error.
	// It parses the request body and extracts relevant information, into a list of strings ready for indexing.
	Parse(reqBody string) ([]string, error)
}
