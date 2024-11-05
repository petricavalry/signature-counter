/*
 * Copyright (c) 2024. LGBT-CN & KevinZonda
 * This file is part of LGBT-CN Signature Counting.
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	SIGNATURE_START = "<!-- BEGIN LGBT-CN SIGNATURE -->"
	SIGNATURE_END   = "<!-- END LGBT-CN SIGNATURE -->"
	COUNT_START     = "<!-- BEGIN LGBT-CN COUNT -->"
	COUNT_END       = "<!-- END LGBT-CN COUNT -->"
)

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	if len(os.Args) < 2 {
		log.Fatal("Please specific file name in argument")
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inSignatureBlock := false
	inCountBlock := false
	lineCount := 0
	lineNumber := 0
	countStartLine := 0
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		slog.Debug(fmt.Sprintf("Parsing line: %d %s", lineNumber, line))

		if line == COUNT_START {
			slog.Debug("Found count block start")
			if inSignatureBlock {
				log.Fatal("Please put count block start outside signature block")
			}
			countStartLine = lineNumber
			inCountBlock = true
			goto exit
		}
		if line == COUNT_END {
			slog.Debug("Found count block end")
			if !inCountBlock {
				log.Fatal("Please put count block start first")
			}
			inCountBlock = false
			goto exit
		}
		if line == SIGNATURE_START {
			slog.Debug("Found signature block start")
			if inCountBlock {
				log.Fatal("Please put signature block end outside count block")
			}
			inSignatureBlock = true
			goto exit
		}
		if line == SIGNATURE_END {
			slog.Debug("Found signature block end")
			if !inSignatureBlock {
				log.Fatal("Please put signature block start first")
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
		fmt.Println("Please put count block end")
	}
	if inSignatureBlock {
		fmt.Println("Please put signature block end")
	}

	slog.Debug(fmt.Sprintf("%d lines in signature block", lineCount))
	slog.Debug(fmt.Sprintf("count block start at line %d", countStartLine))
	lines = slices.Insert(lines, countStartLine+1, fmt.Sprintf("已有%d人签署！", lineCount))

	for _, line := range lines {
		fmt.Println(line)
	}

	output := strings.Join(lines, "\n")

	err = os.WriteFile(fileName, []byte(output), 644)
	if err != nil {
		log.Fatal(err)
	}
}
