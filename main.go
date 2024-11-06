/*
 * Copyright (c) 2024. LGBT-CN & KevinZonda
 * This file is part of LGBT-CN Signature Counting.
 */

package main

import (
	"bufio"
	"log"
	"os"
	"strings"
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
	lines, err := scan(scanner)
	if err != nil {
		log.Fatal(err)
	}
	output := strings.Join(lines, "\n")

	err = os.WriteFile(fileName, []byte(output), 644)
	if err != nil {
		log.Fatal(err)
	}
}
