package fileOps

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FileParser struct {
	FilePath   string
	ParsedFile [][]string
	DidParse   bool
}

func (f *FileParser) Parse() error {
	if f.DidParse {
		return nil
	}
	file, err := os.Open(f.FilePath)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", f.FilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f.ParsedFile = append(f.ParsedFile, strings.Fields(scanner.Text()))
	}
	if scanner.Err() != nil {
		return fmt.Errorf("could not scan %s %v", f.FilePath, err)
	}
	f.DidParse = true
	return nil
}

func (f *FileParser) FindTag(tag string) (int, error) {
	if !f.DidParse {
		if err := f.Parse(); err != nil {
			return 0, fmt.Errorf("failed to parse %s: %v", f.FilePath, err)
		}
	}

	for index, line := range f.ParsedFile {
		if line[0] == tag {
			return index, nil
		}
	}
	return 0, fmt.Errorf("tag \"%s\" was not found", tag)
}
