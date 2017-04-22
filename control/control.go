package control

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NDari/gocrunch/mat"
)

// This type is made private, such that it MUST be created using
// the NewStructure function below.
type Structure struct {
	Title      string
	CellInfo   []float64
	CoordType  string
	NumAtoms   int
	AtomCoors  [][]float64
	AtomNames  []string
	SpaceGroup string
	Rlm        [][]float64
	Mlr        [][]float64
}

// Struture takes a file name (currently only olcao.skl file is
// supported) and returns a pointer to a structure data type whos fields
// are initialized based on the content of that file.
func New(fileName string) *Structure {
	s := new(Structure)
	i := 1
	ext := filepath.Ext(fileName)
	if ext == ".skl" {
		skl, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer skl.Close()
		var words [][]string
		scanner := bufio.NewScanner(skl)
		for scanner.Scan() {
			line := scanner.Text()
			words = append(words, strings.Fields(line))
		}
		for i < len(words) {
			if words[i][0] != "end" {
				s.Title += strings.Join(words[i], " ")
				s.Title += "\n"
				i++
			} else { // remove extra newline at the end.
				s.Title = strings.TrimRight(s.Title, "\n")
				i++
				break
			}
		}
		i++
		s.CellInfo = make([]float64, 6)
		for j := 0; j < len(words[i]); j++ {
			info, err := strconv.ParseFloat(words[i][j], 64)
			if err != nil {
				panic(err)
			}
			s.CellInfo[j] = info
		}
		i++
		if strings.Contains(words[i][0], "frac") {
			s.CoordType = "F"
		} else if strings.Contains(words[i][0], "cart") {
			s.CoordType = "C"
		} else {
			panic("Coordinate type does not contain frac or cart")
		}
		s.NumAtoms, err = strconv.Atoi(words[i][1])
		if err != nil {
			panic(err)
		}
		i++
		s.AtomNames = make([]string, s.NumAtoms)
		s.AtomCoors = mat.New(s.NumAtoms, 3)
		for j := 0; j < s.NumAtoms; j++ {
			s.AtomNames[j] = words[i][0]
			val, err := strconv.ParseFloat(words[i][1], 64)
			if err != nil {
				panic(err)
			}
			s.AtomCoors[j][0] = val
			val, err = strconv.ParseFloat(words[i][2], 64)
			if err != nil {
				panic(err)
			}
			s.AtomCoors[j][1] = val
			val, err = strconv.ParseFloat(words[i][3], 64)
			if err != nil {
				panic(err)
			}
			s.AtomCoors[j][2] = val
			i++
		}
		s.SpaceGroup = words[i][1]
		i++

	}
	a := s.CellInfo[0]
	b := s.CellInfo[1]
	c := s.CellInfo[2]
	alf := s.CellInfo[3]
	bet := s.CellInfo[4]
	gam := s.CellInfo[5]

	// convert the angles to radians
	angToRad := math.Pi / 180.0
	alf *= angToRad
	bet *= angToRad
	gam *= angToRad

	s.Rlm = mat.New(3)

	// assume a and x are colinear.
	s.Rlm[0][0] = a
	s.Rlm[0][1] = 0.0
	s.Rlm[0][2] = 0.0

	// assume b to be in the xy-plane
	s.Rlm[1][0] = b * math.Cos(gam)
	s.Rlm[1][1] = b * math.Sin(gam)
	s.Rlm[1][2] = 0.0

	// c is then a mix of all three (x, y, z)
	s.Rlm[2][0] = c * math.Cos(bet)
	s.Rlm[2][1] = c * (math.Cos(alf) - math.Cos(gam)*math.Cos(bet)) / math.Sin(gam)
	s.Rlm[2][2] = c * math.Sqrt(1.0-math.Pow(math.Cos(bet), 2.0)-
		math.Pow((s.Rlm[2][1]/c), 2.0))

	// zero out small rlm values. this helps with numerical errors.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(s.Rlm[i][j]) < 0.000000001 {
				s.Rlm[i][j] = 0.0
			}
		}
	}
	// the matrix calculated here is lifted from the code in pyMol: look in
	// www.pymolewiki.org/index.php/Cart_to_frac

	// calculate the volume of the scaled parallelpiped
	v := math.Sqrt(1.0 - (math.Cos(alf) * math.Cos(alf)) -
		(math.Cos(bet) * math.Cos(bet)) - (math.Cos(gam) * math.Cos(gam)) +
		2.0*math.Cos(alf)*math.Cos(bet)*math.Cos(gam))

	s.Mlr = mat.New(3)
	// assume a and x are colinear.
	s.Mlr[0][0] = 1.0 / a
	s.Mlr[0][1] = 0.0
	s.Mlr[0][2] = 0.0

	// assume b in the xy-plane.
	s.Mlr[1][0] = -math.Cos(gam) / (a * math.Sin(gam))
	s.Mlr[1][1] = 1.0 / (b * math.Sin(gam))
	s.Mlr[1][2] = 0.0

	// c is then a mix of all three axes.
	s.Mlr[2][0] = (math.Cos(alf)*math.Cos(gam) - math.Cos(bet)) /
		(a * v * math.Sin(gam))
	s.Mlr[2][1] = (math.Cos(bet)*math.Cos(gam) - math.Cos(alf)) /
		(b * v * math.Sin(gam))
	s.Mlr[2][2] = (math.Sin(gam)) / (c * v)

	// zero out mlr enteries that are too small.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(s.Mlr[i][j]) < 0.00000001 {
				s.Mlr[i][j] = 0.0
			}
		}
	}
	return s
}

func (s *Structure) Clone() *Structure {
	r := new(Structure)
	r.Title = s.Title
	copy(r.CellInfo, s.CellInfo)
	r.CoordType = s.CoordType
	r.NumAtoms = s.NumAtoms
	r.AtomCoors = mat.Copy(s.AtomCoors)
	copy(r.AtomNames, s.AtomNames)
	r.SpaceGroup = s.SpaceGroup
	r.Rlm = mat.Copy(s.Rlm)
	r.Mlr = mat.Copy(s.Mlr)
	return r
}

func (s *Structure) ToFrac() *Structure {
	if s.CoordType == "F" {
		return s
	}
	s.AtomCoors = mat.Dot(s.AtomCoors, s.Mlr)
	s.CoordType = "F"
	return s
}

func (s *Structure) ToCart() *Structure {
	if s.CoordType == "C" {
		return s
	}
	s.AtomCoors = mat.Dot(s.AtomCoors, s.Rlm)
	s.CoordType = "C"
	return s
}

func (s *Structure) MinDistMat() [][]float64 {
	mdm := mat.New(s.NumAtoms)
	mat.Set(mdm, 1000000000.0)
	if s.CoordType == "C" {
		s.ToFrac()
		defer s.ToCart()
	}
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			for z := -1; z < 2; z++ {
				c := s.Clone()
				mat.Add(c.AtomCoors, []float64{float64(x), float64(y), float64(z)})
				c.ToCart()

				for i := 0; i < s.NumAtoms-1; i++ {
					for j := i + 1; j < s.NumAtoms; j++ {
						r := math.Sqrt((s.AtomCoors[j][0]-c.AtomCoors[i][0])*
							(s.AtomCoors[j][0]-c.AtomCoors[i][0]) +
							(s.AtomCoors[j][1]-c.AtomCoors[i][1])*
								(s.AtomCoors[j][1]-c.AtomCoors[i][1]) +
							(s.AtomCoors[j][2]-c.AtomCoors[i][2])*
								(s.AtomCoors[j][2]-c.AtomCoors[i][2]))

						if r < mdm[i][j] {
							mdm[i][j] = r
							mdm[j][i] = r
						}
					}

				}
			}
		}
	}
	// min distance between an atom and itself is zero.
	for i := range mdm {
		mdm[i][i] = 0.0
	}
	return mdm
}

// WriteSkl creates a generic olcao.skl file from some info
// which is passed to it. We assume that the space group is 1_a and
// that the supercell is "1 1 1".
func (s *Structure) WriteSkl(fileName string) {

	// start our string.
	str := ""

	// put the olcao.skl header
	str += "title\n"
	str += s.Title
	str += "\nend\ncell\n"

	// now convert the cellInfo
	for i := range s.CellInfo {
		// convert string to float with no decimals ("f"), with 14 digits
		// after the period, from a float64 number.
		str += strconv.FormatFloat(s.CellInfo[i], 'f', 4, 64)
		str += " "
	}
	str += "\n"

	// add in the line for coordinates.
	if s.CoordType == "C" {
		str += "cart "
	} else if s.CoordType == "F" {
		str += "frac "
	}
	str += strconv.Itoa(s.NumAtoms)
	str += "\n"
	// now lets add in the atomic info.
	for i := range s.AtomCoors {
		// first the atom name
		str += s.AtomNames[i]
		str += " "
		// then the 3 coordinates. add spaces between, and a newline at end.
		for j := range s.AtomCoors[i] {
			str += strconv.FormatFloat(s.AtomCoors[i][j], 'f', 4, 64)
			str += " "
		}
		str += "\n"
	}

	// add the olcao.skl footer
	str += "space "
	str += s.SpaceGroup
	str += "\nsupercell 1 1 1\nfull\n"

	// create and open the file. Check for errors. Then write and close.
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(str)
	if err != nil {
		panic(err)
	}
}

// Mutate moves the atoms in a structure by a given distance, mag, in
// a random direction. The probablity that any given atom will be moves
// is given by the argument 'prob', where 0.0 means 0% chance that  the atom
// will not move, and 1.0 means a 100% chance that a given atom will move.
func (s *Structure) Mutate(mag, prob float64) {
	// first, make sure that if the atoms in the structure are in fractional,
	// that we convert them to cartesian, and back again after the move is
	// done.
	if s.CoordType == "F" {
		s.ToCart()
		defer s.ToFrac()
	}
	// for each atomic position, we roll the die. If its less than the prob,
	// we generate a theta and phi randomly from 0 to 2Pi, then move the
	// atom in that direction.
	for i := 0; i < len(s.AtomCoors); i++ {
		if rand.Float64() < prob {
			theta := rand.Float64() * math.Pi
			phi := rand.Float64() * 2.0 * math.Pi
			xAdded := mag * math.Sin(theta) * math.Cos(phi)
			yAdded := mag * math.Sin(theta) * math.Sin(phi)
			zAdded := mag * math.Cos(theta)
			s.AtomCoors[i][0] += xAdded
			s.AtomCoors[i][1] += yAdded
			s.AtomCoors[i][2] += zAdded
		}
	}
	// since we mutated this structure, the space group no longer
	// applies. We make sure that the it is appropriately set to
	// reflect this: it is set to '1_a'.
	s.SpaceGroup = "1_a"
	s.ApplyPBC()
}

func (s *Structure) ApplyPBC() {
	if s.CoordType == "C" {
		s.ToFrac()
		defer s.ToCart()
	}
	for i := range s.AtomCoors {
		for s.AtomCoors[i][0] < 0 {
			s.AtomCoors[i][0] += 1
		}
		for s.AtomCoors[i][1] < 0 {
			s.AtomCoors[i][1] += 1
		}
		for s.AtomCoors[i][2] < 0 {
			s.AtomCoors[i][2] += 1
		}
		for s.AtomCoors[i][0] > 1 {
			s.AtomCoors[i][0] -= 1
		}
		for s.AtomCoors[i][1] > 1 {
			s.AtomCoors[i][1] -= 1
		}
		for s.AtomCoors[i][2] > 1 {
			s.AtomCoors[i][2] -= 1
		}
	}
}

func (s *Structure) WriteXYZ(fileName, comment string) {
	// the xyz file must contain cartesian coordinate. So, if the
	// coordinates of the structure are fractional, we must convert
	// them.
	if s.CoordType == "F" {
		s.ToCart()
		defer s.ToFrac()
	}

	// start our string.
	str := ""

	// set the number of atoms.
	str += strconv.Itoa(s.NumAtoms)
	str += "\n"

	// next we need the comment line. make sure no extra newlines are in the
	// string
	comment = strings.TrimSpace(comment)
	str += comment
	str += "\n"
	// now lets add in the atomic info.
	for i := range s.AtomCoors {

		// first the atom name
		str += s.AtomNames[i]
		str += " "

		// then the 3 coordinates. Add spaces between, and a newline at end.
		for j := range s.AtomCoors[i] {
			str += strconv.FormatFloat(s.AtomCoors[i][j], 'f', 4, 64)
			str += " "
		}
		str += "\n"
	}

	// create and open the file. Check for errors. Then write and close.
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(str)
	if err != nil {
		panic(err)
	}
}

func (s *Structure) MinVecs() [][][]float64 {
	if s.CoordType == "C" {
		s.ToFrac()
		defer s.ToCart()
	}
	// create the space needed to hold the min distance matrix..
	mdm := mat.New(s.NumAtoms)
	mat.Set(mdm, 1000000000.0)

	// create the space needed for the min vectors.
	mv := make([][][]float64, s.NumAtoms)
	for i := range mv {
		mv[i] = mat.New(s.NumAtoms, 3)
	}

	// now lets begin calculating the closest distance between two atoms. to
	// do calculate the distance between two atoms while considering the PBC
	// we need to calculate the distance of a particular atom, i, to another
	// atom, j, and compare that to the distance between i and the reflections
	// of j across the boundaries (the ab-plane, the ac-plane, etc). there is
	// 27 "boxes" that we must consider, including the box containing the
	// original system.
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			for z := -1; z < 2; z++ {

				// make a copy
				c := s.Clone()

				// move the atoms to the box we wish to check.
				mat.Add(c.AtomCoors, []float64{float64(x), float64(y), float64(z)})
				// and convert the new position to cartesian.
				c.ToCart()

				// now calculate the distance from atoms in the original "box"
				// to this "box".
				for i := 0; i < s.NumAtoms-1; i++ {
					for j := i + 1; j < s.NumAtoms; j++ {
						r := math.Sqrt((s.AtomCoors[j][0]-c.AtomCoors[i][0])*
							(s.AtomCoors[j][0]-c.AtomCoors[i][0]) +
							(s.AtomCoors[j][1]-c.AtomCoors[i][1])*
								(s.AtomCoors[j][1]-c.AtomCoors[i][1]) +
							(s.AtomCoors[j][2]-c.AtomCoors[i][2])*
								(s.AtomCoors[j][2]-c.AtomCoors[i][2]))
						if r < mdm[i][j] {
							mdm[i][j] = r
							mdm[j][i] = r
							for l := 0; l < 3; l++ {
								mv[i][j][l] = s.AtomCoors[j][l] - c.AtomCoors[i][l]
								mv[j][i][l] = mv[i][j][l]
							}
						}
					}
				}
			}
		}
	}

	// finally, set the diagonal elements (vector from an atom to itself)
	// to 0.0
	for i := 0; i < s.NumAtoms; i++ {
		for l := 0; l < 3; l++ {
			mv[i][i][l] = 0.0
		}
	}
	return mv
}

func (s *Structure) GetSymFns() [][]float64 {
	if s.CoordType == "C" {
		s.ToFrac()
		defer s.ToCart()
	}

	allSym := make([][]float64, s.NumAtoms)
	var set []float64

	mdm := s.MinDistMat()
	mv := s.MinVecs()

	set = GenSymFn1(mdm, 1.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 1.2)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 1.4)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 1.6)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 1.8)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 2.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 2.2)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 2.4)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 2.6)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 2.8)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 3.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 3.2)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 3.4)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 3.6)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 3.8)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 4.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 4.2)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 4.4)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 4.6)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 4.8)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 5.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 5.2)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 5.4)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 5.6)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 5.8)
	mat.AppendCol(allSym, set)
	set = GenSymFn1(mdm, 6.0)
	mat.AppendCol(allSym, set)

	set = GenSymFn2(mdm, 5.0, 2.0, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 2.2, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 2.4, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 2.6, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 2.8, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 3.0, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 3.2, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 3.4, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 3.6, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 3.8, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 4.0, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 4.2, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 4.4, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 4.6, 12.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn2(mdm, 5.0, 4.8, 12.0)
	mat.AppendCol(allSym, set)

	set = GenSymFn3(mdm, 5.0, 1.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 1.5)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 2.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 2.5)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 3.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 3.5)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 4.0)
	mat.AppendCol(allSym, set)
	set = GenSymFn3(mdm, 5.0, 4.5)
	mat.AppendCol(allSym, set)

	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 1.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 2.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 4.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 16.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 64.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 1.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 2.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 4.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 16.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 64.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 1.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 2.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 4.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 16.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn4(mdm, s.AtomCoors, mv, 5.0, 1.0, 64.0, 0.001)
	mat.AppendCol(allSym, set)

	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 1.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 2.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 4.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 16.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 64.0, 0.1)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 1.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 2.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 4.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 16.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 64.0, 0.01)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 1.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 2.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 4.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 16.0, 0.001)
	mat.AppendCol(allSym, set)
	set = GenSymFn5(mdm, s.AtomCoors, mv, 5.0, 1.0, 64.0, 0.001)
	mat.AppendCol(allSym, set)

	return allSym
}

// CutoffFn returns the value of the cutoff function defined in:
//
// "Atom Centered Symmetry Fucntions for Contructing High-Dimentional Neural
// Network Potentials", by Jorg Behler, J. Chem. Phys. 134, 074106 (2011)
func CutoffFn(dist, cutoffRad float64) float64 {
	if dist > cutoffRad {
		return 0.0
	}
	return (0.5 * (math.Cos((math.Pi*dist)/cutoffRad) + 1.0))
}

// GenSymFn1 returns the value of the first symmetry function defined in:
//
// "Atom Centered Symmetry Fucntions for Contructing High-Dimentional Neural
// Network Potentials", by Jorg Behler, J. Chem. Phys. 134, 074106 (2011)
func GenSymFn1(minDistMat [][]float64, cutoff float64) []float64 {

	symFn1 := make([]float64, len(minDistMat))

	for i := 0; i < len(minDistMat); i++ {
		for j := 0; j < len(minDistMat[i]); j++ {
			symFn1[i] += CutoffFn(minDistMat[i][j], cutoff)
		}
	}

	return symFn1
}

// GenSymFn2 returns the value of the second symmetry function defined in:
//
// "Atom Centered Symmetry Fucntions for Contructing High-Dimentional Neural
// Network Potentials", by Jorg Behler, J. Chem. Phys. 134, 074106 (2011)
func GenSymFn2(minDistMat [][]float64, cutoff, rs, eta float64) []float64 {

	symFn2 := make([]float64, len(minDistMat))

	for i := 0; i < len(minDistMat); i++ {
		for j := 0; j < len(minDistMat[i]); j++ {
			val := CutoffFn(minDistMat[i][j], cutoff)
			if val != 0.0 {
				symFn2[i] += (math.Exp(-eta*
					math.Pow((minDistMat[i][j]-rs), 2.0)) * val)
			}
		}
	}

	return symFn2
}

// GenSymFn3 returns the value of the third symmetry function defined in:
//
// "Atom Centered Symmetry Fucntions for Contructing High-Dimentional Neural
// Network Potentials", by Jorg Behler, J. Chem. Phys. 134, 074106 (2011)
func GenSymFn3(minDistMat [][]float64, cutoff, kappa float64) []float64 {

	symFn3 := make([]float64, len(minDistMat))

	for i := 0; i < len(minDistMat); i++ {
		for j := 0; j < len(minDistMat[i]); j++ {
			val := CutoffFn(minDistMat[i][j], cutoff)
			if val != 0.0 {
				symFn3[i] += (math.Cos(kappa * minDistMat[i][j] * val))
			}
		}
	}

	return symFn3
}

// GenSymFn4 returns the value of the symmetry function 4 defined in:
//
// "Atom Centered Symmetry Fucntions for Contructing High-Dimentional Neural
// Network Potentials", by Jorg Behler, J. Chem. Phys. 134, 074106 (2011)
func GenSymFn4(minDistMat, cartCoors [][]float64, mv [][][]float64,
	cutoff, lambda, zeta, eta float64) []float64 {

	// lets make the space for the symmetry function.
	symFn4 := make([]float64, len(minDistMat))

	// lets make space for a vector that is from atom i to atom j.
	Rij := make([]float64, 3)

	// and a vector from atom i to atom k.
	Rik := make([]float64, 3)

	for i := 0; i < len(minDistMat); i++ {
		for j := 0; j < len(minDistMat); j++ {
			// skip this loop if i and j are the same atom.
			if j == i {
				continue
			}
			FcRij := CutoffFn(minDistMat[i][j], cutoff)
			// skip this loop if the cutoff function value for atom j is zero.
			// this means that the atom is outside the cutoff radius.
			if FcRij == 0.0 {
				continue
			}
			// lets get the vector from atom i to atom j.
			Rij = mv[i][j]
			for k := 0; k < len(minDistMat); k++ {
				// skip this loop if i and k are the same atom.
				if i == k {
					continue
				}
				FcRik := CutoffFn(minDistMat[i][k], cutoff)
				// skip this loop if the cutoff function value for k is zero.
				if FcRik == 0.0 {
					continue
				}
				// also skip this loop if the cutoff function between j and
				// k is 0.
				FcRjk := CutoffFn(minDistMat[j][k], cutoff)
				if FcRjk == 0.0 {
					continue
				}
				// lets get the vector from atom i to atom k.
				Rik = mv[i][k]

				// calculate the cosine of the angle between the vectors from
				// i to j, and from i to k. this is the dot product of the two
				// vectors devided by product of their magnitudes.
				cosThetaijk := (Rij[0]*Rik[0] + Rij[1]*Rik[1] + Rij[2]*Rik[2])
				cosThetaijk /= (minDistMat[i][j] * minDistMat[i][k])

				// calculate the symmetry function for the ijk triple, and add
				// it to the symmetry function for this atom.
				symFn4[i] += (math.Pow(1.0+lambda*cosThetaijk, zeta) *
					math.Exp(-eta*(minDistMat[i][j]*minDistMat[i][j]+
						minDistMat[i][k]*minDistMat[i][k]+
						minDistMat[j][k]*minDistMat[j][k])) *
					FcRij * FcRik * FcRjk)
			}
		}
		symFn4[i] *= math.Pow(2.0, 1.0-zeta)
	}
	return symFn4
}

// GenSymFn5 returns the value of the symmetry function 5 defined in:
//
// "Atom Centered Symmetry Fucntions for Contructing High-Dimentional Neural
// Network Potentials", by Jorg Behler, J. Chem. Phys. 134, 074106 (2011)
func GenSymFn5(minDistMat, cartCoors [][]float64, mv [][][]float64,
	cutoff, lambda, zeta, eta float64) []float64 {

	// lets make the space for the symmetry function.
	symFn5 := make([]float64, len(minDistMat))

	// lets make space for a vector that is from atom i to atom j.
	Rij := make([]float64, 3)

	// and a vector from atom i to atom k.
	Rik := make([]float64, 3)

	for i := 0; i < len(minDistMat); i++ {
		for j := 0; j < len(minDistMat); j++ {
			// skip this loop if i and j are the same atom.
			if j == i {
				continue
			}
			FcRij := CutoffFn(minDistMat[i][j], cutoff)
			// skip this loop if the cutoff function value for atom j is zero.
			// this means that the atom is outside the cutoff radius.
			if FcRij == 0.0 {
				continue
			}
			// lets get the vector from atom i to atom j.
			Rij = mv[i][j]
			for k := 0; k < len(minDistMat); k++ {
				// skip this loop if i and k are the same atom.
				if i == k {
					continue
				}
				FcRik := CutoffFn(minDistMat[i][k], cutoff)
				// skip this loop if the cutoff function value for k is zero.
				if FcRik == 0.0 {
					continue
				}
				// lets calculate the vector from atom i to atom k.
				Rik = mv[i][k]

				// calculate the cosine of the angle between the vectors from
				// i to j, and from i to k. this is the dot product of the two
				// vectors devided by product of their magnitudes.
				cosThetaijk := (Rij[0]*Rik[0] + Rij[1]*Rik[1] + Rij[2]*Rik[2])
				cosThetaijk /= (minDistMat[i][j] * minDistMat[i][k])

				// calculate the symmetry function for the ijk triple, and add
				// it to the symmetry function for this atom.
				symFn5[i] += (math.Pow(1.0+lambda*cosThetaijk, zeta) *
					math.Exp(-eta*(minDistMat[i][j]*minDistMat[i][j]+
						minDistMat[i][k]*minDistMat[i][k])) * FcRij * FcRik)
			}
		}
		symFn5[i] *= math.Pow(2.0, 1.0-zeta)
	}
	return symFn5
}
