package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

func assert(t *testing.T, condition bool) {
	if !condition {
		t.Errorf("Test failed")
	}
}

func testString(t *testing.T, original string, expected string) {
	reader := strings.NewReader(original)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	minify(reader, writer)
	writer.Flush()

	assert(t, buffer.String() == expected)
}

func TestEmptyString(t *testing.T) {
	testString(t, "", "")
}

func TestSimple(t *testing.T) {
	original := `
	{ 
		"xxx     ": "yyy 123 ", 
		 "yyy",     "zzz z ",
		 "zzz" : [  " 123 ", 123,    321   ]
	}
 `
	expected := `{"xxx     ":"yyy 123 ","yyy","zzz z ","zzz":[" 123 ",123,321]}`
	testString(t, original, expected)
}

func TestEscaping(t *testing.T) {
	original := `
	{
		"escape2": "how are you?\\",
		"es ca\" ped": "he said \"hello\" and... ",
		"księgi wieczyste": "漢 \\字 "
	}
	`
	expected := `{"escape2":"how are you?\\","es ca\" ped":"he said \"hello\" and... ","księgi wieczyste":"漢 \\字 "}`
	testString(t, original, expected)
}

// here we read huge JSON file, then prettify it using Go's standard library,
// then minify it and prettify again to check if it 100% matches original prettified string
func TestLargeFile(t *testing.T) {
	fileBuffer, err := ioutil.ReadFile("large.json")
	if err != nil {
		panic(err)
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, fileBuffer, "", "\t")
	if err != nil {
		panic(err)
	}

	origJSON := prettyJSON.String()

	prettyReader := bufio.NewReader(&prettyJSON)
	var minifiedBuffer bytes.Buffer
	minifiedBufferWriter := bufio.NewWriter(&minifiedBuffer)
	minify(prettyReader, minifiedBufferWriter)
	minifiedBufferWriter.Flush()

	err = json.Indent(&prettyJSON, minifiedBuffer.Bytes(), "", "\t")
	if err != nil {
		panic(err)
	}

	assert(t, origJSON == prettyJSON.String())
}
