package fileOps

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type SklParser struct {
	FileParser
}

func NewSklParser(path string) (*SklParser, error) {
	ext := filepath.Ext(path)
	if ext != ".skl" {
		return nil, fmt.Errorf("expected a file with .skl extension, but got %s", ext)
	}

	s := &SklParser{
		FileParser{
			FilePath: path,
		},
	}

	s.Parse()
	return s, nil
}

// Title returns the title string from an olcao.skl file.
func (s *SklParser) Title() (string, error) {
	var title string
	endTag, err := s.FindTag("end")
	if err != nil {
		return "", fmt.Errorf("failed to find the 'end' tag: %v", err)
	}

	for i := 1; i < endTag; i++ {
		title += strings.Join(s.ParsedFile[i], " ")
		title += "\n"
	}
	return strings.TrimSpace(title), nil
}

// CellInfo returns the a, b, c cell vectors and the alpha, beta, and gamma
// cell angles from an olcao.skl file.
func (s *SklParser) CellInfo() ([]float64, error) {
	cellInfo := make([]float64, 6)
	endTag, err := s.FindTag("cell")
	if err != nil {
		return nil, fmt.Errorf("failed to find the 'cell' tag: %v", err)
	}

	for i := 0; i < 6; i++ {
		info, err := strconv.ParseFloat(s.ParsedFile[endTag+1][i], 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse float: %v", err)
		}
		cellInfo[i] = info
	}
	return cellInfo, nil
}

// NumAtoms returns the number of atoms in a olcao.skl file.
func (s *SklParser) NumAtoms() (int, error) {
	cellTag, err := s.FindTag("cell")
	if err != nil {
		return 0, fmt.Errorf("failed to find the 'cell' tag: %v", err)
	}

	num, err := strconv.Atoi(s.ParsedFile[cellTag+2][1])
	if err != nil {
		return 0, fmt.Errorf("could not parse int: %v", err)
	}
	return num, nil
}

// CoordType checks the coordinates used in an olcao.skl file and returns
// "F" for fractional, and "C" for cartesian.
func (s *SklParser) CoordType() (string, error) {
	cellTag, err := s.FindTag("cell")
	if err != nil {
		return "", fmt.Errorf("failed to find the 'cell' tag: %v", err)
	}

	if strings.Contains(s.ParsedFile[cellTag+1][0], "frac") {
		return "F", nil
	}
	if strings.Contains(s.ParsedFile[cellTag+2][0], "cart") {
		return "C", nil
	}
	return "", fmt.Errorf("unknown coordinate type: %s", s.ParsedFile[cellTag+2][0])
}
