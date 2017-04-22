package fileOps

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSklTiTle(t *testing.T) {
	s, err := NewSklParser("data/olcao.skl")
	if err != nil {
		log.Fatal(err)
	}

	title, err := s.Title()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "some title\nsecond title line", title)
}

func TestSklCellInfo(t *testing.T) {
	s, err := NewSklParser("data/olcao.skl")
	if err != nil {
		log.Fatal(err)
	}

	cellInfo, err := s.CellInfo()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 10.1, cellInfo[0])
	assert.Equal(t, 12.122, cellInfo[1])
	assert.Equal(t, 1.0, cellInfo[2])
	assert.Equal(t, 90.0, cellInfo[3])
	assert.Equal(t, 24.56, cellInfo[4])
	assert.Equal(t, 123.10123, cellInfo[5])
}

func TestNumAtoms(t *testing.T) {
	s, err := NewSklParser("data/olcao.skl")
	if err != nil {
		log.Fatal(err)
	}

	numAtoms, err := s.NumAtoms()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 2, numAtoms)
}
