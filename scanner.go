package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
)

const (
	SIGNATURE_START = "<!-- BEGIN LGBT-CN SIGNATURE -->"
	SIGNATURE_END   = "<!-- END LGBT-CN SIGNATURE -->"
	COUNT_START     = "<!-- BEGIN LGBT-CN COUNT -->"
	COUNT_END       = "<!-- END LGBT-CN COUNT -->"
)

func scan(scanner *bufio.Scanner) ([]string, error) {
	inSignatureBlock := false
	inCountBlock := false
	lineCount := 0
	lineNumber := 0
	// count block does not exists if count equal to zero
	countStartLine := 0
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == COUNT_START {
			if inCountBlock {
				return nil, errors.New("Please don't use nested count block")
			}
			if inSignatureBlock {
				return nil, errors.New("Please put count block start outside signature block")
			}
			// insert count in next line
			countStartLine = lineNumber + 1
			inCountBlock = true
			goto exit
		}
		if line == COUNT_END {
			if !inCountBlock {
				return nil, errors.New("Please put count block start first")
			}
			inCountBlock = false
			goto exit
		}
		if line == SIGNATURE_START {
			if inSignatureBlock {
				return nil, errors.New("Please don't use nested signature block")
			}
			if inCountBlock {
				return nil, errors.New("Please put signature block end outside count block")
			}
			inSignatureBlock = true
			goto exit
		}
		if line == SIGNATURE_END {
			if !inSignatureBlock {
				return nil, errors.New("Please put signature block start first")
			}
			inSignatureBlock = false
			goto exit
		}
		if inSignatureBlock && line != "" {
			lineCount++
		}
		// skip content in count block
		if inCountBlock {
			continue
		}
	exit:
		lines = append(lines, line)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if inCountBlock {
		return nil, errors.New("Please put count block end")
	}
	if inSignatureBlock {
		return nil, errors.New("Please put signature block end")
	}

	if countStartLine == 0 {
		return nil, errors.New("Please put count block")
	}
	lines = slices.Insert(lines, countStartLine, fmt.Sprintf("已有%d人签署！", lineCount))

	return lines, nil
}
