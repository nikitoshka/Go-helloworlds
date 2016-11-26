package counter

import (
	"bufio"
	"bytes"
	"fmt"
)

// ByteCounter is hjast an alias for int
type ByteCounter int

// WordCounter is hjast an alias for int
type WordCounter int

// LineCounter is hjast an alias for int
type LineCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))

	return len(p), nil
}

func (c *WordCounter) Write(p []byte) (int, error) {
	reader := bytes.NewReader(p)
	scanner := bufio.NewScanner(reader)
	count := 0

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error during word scanning: %v", err))
	}

	*c += WordCounter(count)

	return count, nil
}

func (c *LineCounter) Write(p []byte) (int, error) {
	reader := bytes.NewReader(p)
	scanner := bufio.NewScanner(reader)
	count := 0

	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error during line scanning: %v", err))
	}

	*c += LineCounter(count)

	return count, nil
}
