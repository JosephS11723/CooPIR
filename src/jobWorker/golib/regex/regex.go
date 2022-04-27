package regex

// code adapted from: https://github.com/vertoforce/streamregex/blob/master/streamregex.go
// file created by github.com/brandonvessel for a personal project

import (
	"bufio"
	"context"
	"io"
	"regexp"
)

// SplitRegex takes a regex and returns a split function that will find that regex in a byte slice
func SplitRegex(re *regexp.Regexp, maxMatchLength int) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, io.EOF
		}
		if loc := re.FindIndex(data); loc != nil {
			return loc[1], data[loc[0]:loc[1]], nil
		}
		if atEOF {
			return 0, nil, io.EOF
		}
		if len(data) >= maxMatchLength {
			return len(data) - maxMatchLength, nil, nil
		}
		return 0, nil, nil
	}
}

// FindReader return channel of matched []byte from reader.
// This function will allocate maxMatchLength*2 bytes of memory
func FindReader(ctx context.Context, r *regexp.Regexp, maxMatchLength int, reader io.Reader) chan string {
	allMatches := make(chan string)

	buf := make([]byte, maxMatchLength*2)

	go func() {
		defer close(allMatches)

		scanner := bufio.NewScanner(reader)
		scanner.Buffer(buf, maxMatchLength)
		scanner.Split(SplitRegex(r, maxMatchLength))
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case allMatches <- scanner.Text():
			}
		}
	}()

	return allMatches
}

// FindByReader returns a list of all the information matched by the regex
func FindByReader(r *regexp.Regexp, reader io.Reader, maxMatchLength int) []string {
	var matches []string

	// create scanner from buffer
	scanner := bufio.NewScanner(reader)

	// maker buffer large enough to hold maxMatchLength
	buf := make([]byte, maxMatchLength*2)

	scanner.Buffer(buf, maxMatchLength)
	scanner.Split(SplitRegex(r, maxMatchLength))

	for scanner.Scan() {
		matches = append(matches, scanner.Text())
	}

	return matches
}
