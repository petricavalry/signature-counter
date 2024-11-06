package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMissCountStart(t *testing.T) {
	code := `<!-- END LGBT-CN COUNT -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please put count block start first")
}

func TestMissCountEnd(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN COUNT -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please put count block end")
}

func TestDuplicateCountStart(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN COUNT -->
<!-- BEGIN LGBT-CN COUNT -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please don't use nested count block")
}

func TestMissSignatureStart(t *testing.T) {
	code := `<!-- END LGBT-CN SIGNATURE -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please put signature block start first")
}

func TestMissSignatureEnd(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN SIGNATURE -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please put signature block end")
}

func TestDuplicateSignatureStart(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN SIGNATURE -->
<!-- BEGIN LGBT-CN SIGNATURE -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please don't use nested signature block")
}

func TestMissCount(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN SIGNATURE -->
<!-- END LGBT-CN SIGNATURE -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	_, err := scan(scanner)
	assert.EqualError(t, err, "Please put count block")
}

// Count should be zero when signature block not exists.
func TestMissSignature(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN COUNT -->
<!-- END LGBT-CN COUNT -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	lines, err := scan(scanner)	
	assert.Nil(t, err)
	assert.Equal(t, lines, []string{
		"<!-- BEGIN LGBT-CN COUNT -->",
		"已有0人签署！",
		"<!-- END LGBT-CN COUNT -->",
	})
}

func TestSkipEmptyInSignature(t *testing.T) {
	code := `<!-- BEGIN LGBT-CN COUNT -->
<!-- END LGBT-CN COUNT -->

<!-- BEGIN LGBT-CN SIGNATURE -->
unrivaled scalded

reputable overripe
<!-- END LGBT-CN SIGNATURE -->`
	scanner := bufio.NewScanner(strings.NewReader(code))
	lines, err := scan(scanner)
	assert.Nil(t, err)
	assert.Equal(t, lines, []string{
		"<!-- BEGIN LGBT-CN COUNT -->",
		"已有2人签署！",
		"<!-- END LGBT-CN COUNT -->",
		"",
		"<!-- BEGIN LGBT-CN SIGNATURE -->",
		"unrivaled scalded",
		"",
		"reputable overripe",
		"<!-- END LGBT-CN SIGNATURE -->",
	})
}
