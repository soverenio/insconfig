package utils

import (
	"io"
	"regexp"
	"strings"
)

type Trimmer struct {
	buf    strings.Builder
	output io.Writer
}

var (
	trimmerRegex   = regexp.MustCompile("\\s+\n")
	trimmerReplace = "\n"
)

func (t Trimmer) Close() error {
	trimmed := trimmerRegex.ReplaceAllString(t.buf.String(), trimmerReplace)
	_, err := t.output.Write([]byte(trimmed))
	return err
}

func (t *Trimmer) Write(p []byte) (int, error) {
	return t.buf.Write(p)
}

func NewTrimmer(output io.Writer) *Trimmer {
	return &Trimmer{output: output}
}
