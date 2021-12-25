package main

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

func minify(reader io.RuneReader, writer *bufio.Writer) {
	inString := false
	escaping := false

	for {
		curRune, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		if curRune == '"' && !escaping {
			inString = !inString
		}

		if !unicode.IsSpace(curRune) || inString {
			writer.WriteRune(curRune)
		}

		escaping = inString && !escaping && curRune == '\\'
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	minify(reader, writer)

	writer.Flush()
}
