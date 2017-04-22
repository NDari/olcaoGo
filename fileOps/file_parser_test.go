package fileOps

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	p := FileParser{
		FilePath: "data/output.txt",
	}

	err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 2, len(p.ParsedFile), "must be two lines")
	assert.Equal(t, 4, len(p.ParsedFile[0]), "must have 4 words")
	assert.Equal(t, 5, len(p.ParsedFile[1]), "must have 5 words")
}

func TestFindTag(t *testing.T) {
	p := FileParser{
		FilePath: "data/output.txt",
	}

	err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}
	tagLine, err := p.FindTag("some")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 0, tagLine)

	tagLine, err = p.FindTag("here")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 1, tagLine)

	tagLine, err = p.FindTag("fake_tag")
	assert.Equal(t, fmt.Errorf("tag \"%s\" was not found", "fake_tag"), err)
}
