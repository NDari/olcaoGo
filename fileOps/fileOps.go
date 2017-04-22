// fileOps package contains many functions which are related to reading
// and processing the various files in OLCAO. It contains functions to extract
// various needed data from various files used in a typical OLCAO run.
package fileOps

import (

	//"fmt"

	"os"
	"strconv"
	"strings"
	"utilFns"
)

// ReadFloats reads a file and converts all entries into float64

// PrintFloats prints a file from a [][]float64 that was passed to it.
func PrintFloats(data [][]float64, outFile string) {

	// record the name of this function for error handling.
	name := "structOps.PrintFloats"
	var str string
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			str += strconv.FormatFloat(data[i][j], 'f', 14, 64)
			str += " "
		}
		if i+1 <= len(data) {
			str += "\n"
		}
	}

	f, err := os.Create(outFile)
	utilFns.Check(err, name)
	defer f.Close()
	_, err = f.WriteString(str)
	utilFns.Check(err, name)
}

func SklTitle(skl [][]string) string {
	var title string

	for i := 1; i < len(skl); i++ {
		if skl[i][0] != "end" {
			title += strings.Join(skl[i], " ")
			title += "\n"
		} else {
			title = strings.TrimRight(title, "\n")
			break
		}
	}

	return title
}

// // SklCoors returns the coordinates of the atoms contained in an olcao.skl
// // file. the coordinates are returned in a 2D float64 array, where the first
// // dimension is the atom #, and the second dimension is the x, y, and z
// // coordinates.
// func SklCoors(skl [][]string) [][]float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SklCoors"

// 	// get the number of atoms.
// 	numAtoms := SklNumAtoms(skl)

// 	// make the coordinate array
// 	coors := make([][]float64, numAtoms)
// 	for i := 0; i < numAtoms; i++ {
// 		coors[i] = make([]float64, 3)
// 	}

// 	// find the line where the the coordinates are listed, and
// 	// assign that line number to k.
// 	var k int
// 	for i := 1; i < len(skl); i++ {
// 		if skl[i][0] == "cell" {
// 			k = i + 3
// 			break
// 		}
// 	}

// 	// fill up the coordinates
// 	for i := 0; i < numAtoms; i++ {
// 		for j := 0; j < 3; j++ {
// 			c, err := strconv.ParseFloat(skl[k+i][1+j], 64)
// 			utilFns.Check(err, name)
// 			coors[i][j] = c
// 		}
// 	}
// 	return coors
// }

// // SklAtomNames returns an array of strings which contain the names of the
// // atoms in an olcao.skl file.
// func SklAtomNames(skl [][]string) []string {
// 	// get the number of atoms.
// 	numAtoms := SklNumAtoms(skl)

// 	// make the atom names array
// 	aNames := make([]string, numAtoms)

// 	// find the line where the the coordinates are listed, and assign
// 	// that line number to k.
// 	var k int
// 	for i := 1; i < len(skl); i++ {
// 		if skl[i][0] == "cell" {
// 			k = i + 3
// 			break
// 		}
// 	}

// 	// get the element names, converting to lower case. The case conversion
// 	// is done to keep with the general rules used in the olcao package.
// 	for i := 0; i < numAtoms; i++ {
// 		aNames[i] = strings.ToLower(skl[i+k][0])
// 	}
// 	return aNames
// }

// // SklSpaceGroup return the space group of the structure in the olcao.skl
// // file.
// func SklSpaceGroup(skl [][]string) string {
// 	return skl[len(skl)-3][1]
// }

// // WriteSkl creates a generic olcao.skl file from some info
// // which is passed to it. We assume that the space group is 1_a and
// // that the supercell is "1 1 1".
// func SklWrite(cellInfo []float64, aNames []string, cartCoors [][]float64, fileName string) {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SklWrite"

// 	// start our string.
// 	str := ""

// 	// put the olcao.skl header
// 	str += "title\nGeneric candidate title\nend\ncell\n"

// 	// now convert the cellInfo
// 	for i := 0; i < 6; i++ {

// 		// convert string to float with no decimals ("f"), with 4 digits
// 		// after the period, from a float64 number.
// 		str += strconv.FormatFloat(cellInfo[i], 'f', 4, 64)

// 		// add space between entries.
// 		str += " "
// 	}

// 	// add newline.
// 	str += "\n"

// 	// add in the line for coordinates.
// 	str += "cart "
// 	str += strconv.Itoa(len(aNames))
// 	str += "\n"
// 	// now lets add in the atomic info.
// 	for i := 0; i < len(cartCoors); i++ {

// 		// first the atom name
// 		str += aNames[i]
// 		str += " "

// 		// then the 3 coordinates. add sapces between, and a newline at end.
// 		for j := 0; j < 3; j++ {
// 			str += strconv.FormatFloat(cartCoors[i][j], 'f', 4, 64)
// 			str += " "
// 		}
// 		str += "\n"
// 	}

// 	// add the olcao.skl footer
// 	str += "space 1_a\nsupercell 1 1 1\nfull"

// 	// create and open the file. check for errors.
// 	f, err := os.Create(fileName)
// 	utilFns.Check(err, name)
// 	defer f.Close()
// 	f.WriteString(str)
// 	return
// }

// // The next functions are concerned with extracting information from the
// // olcao.dat file

// // OdatTitle returns the title string from an olcao.dat file.
// func OdatTitle(odat [][]string) string {
// 	var title string

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] != "END_TITLE" {
// 			title += strings.Join(odat[i], " ")
// 			title += "\n"
// 		} else {
// 			title = strings.TrimRight(title, "\n")
// 			break
// 		}
// 	}

// 	return title
// }

// // OdatNumTypes returns the number of atom types in a system from an
// // olcao.dat file.
// func OdatNumAtomTypes(odat [][]string) int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatNumAtomTypes"
// 	var numTypes int

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_ATOM_TYPES" {
// 			num, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)

// 			// This is stupid. While i request a 32 bit int in the strconv
// 			// above, I still have to convert it to 32 bit int here to make
// 			// it work. Not sure where the error is.
// 			numTypes = int(num)
// 			break
// 		}
// 	}

// 	return numTypes
// }

// func OdatAtomElementIDs(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatAtomElementIDs"

// 	var elemIds []int

// 	j := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "ATOM_TYPE_ID__SEQUENTIAL_NUMBER" {
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			elemIds = append(elemIds, int(dummy))
// 			j += 1
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if j != int(dummy) {
// 				log.Fatalf("Types numbered out of order.")
// 			}
// 		}
// 	}

// 	return elemIds
// }

// func OdatAtomSpeciesIDs(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatAtomSpeciesTypeIDs"

// 	var speciesIds []int

// 	j := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "ATOM_TYPE_ID__SEQUENTIAL_NUMBER" {
// 			dummy, err := strconv.ParseInt(odat[i+1][1], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			speciesIds = append(speciesIds, int(dummy))
// 			j += 1
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if j != int(dummy) {
// 				log.Fatalf("Types numbered out of order.")
// 			}
// 		}
// 	}

// 	return speciesIds
// }

// func OdatAtomTypeIDs(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatAtomTypeIDs"

// 	var typeIds []int

// 	j := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "ATOM_TYPE_ID__SEQUENTIAL_NUMBER" {
// 			dummy, err := strconv.ParseInt(odat[i+1][2], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			typeIds = append(typeIds, int(dummy))
// 			j += 1
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if j != int(dummy) {
// 				log.Fatalf("Types numbered out of order.")
// 			}
// 		}
// 	}

// 	return typeIds
// }

// func OdatAtomTypeLabels(odat [][]string) []string {

// 	var typeLabels []string

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "ATOM_TYPE_LABEL" {
// 			typeLabels = append(typeLabels, odat[i+1][0])
// 		}
// 	}

// 	return typeLabels
// }

// func OdatNumAtomAlphas(odat [][]string) [][]int {

// 	var alphaSet []int
// 	var numAlphas [][]int

// 	// define name of the function. this is done for error handling.
// 	name := "fileOps.OdatNumAtomAlphas"

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_ALPHA_S_P_D_F" {
// 			// reset the alphaSet to be empty
// 			alphaSet = alphaSet[:0]
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			alphaSet = append(alphaSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][1], 10, 32)
// 			utilFns.Check(err, name)
// 			if alphaSet[0] < int(dummy) {
// 				log.Fatalf("Num of alphas not monotonically decreasing")
// 			}
// 			alphaSet = append(alphaSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][2], 10, 32)
// 			utilFns.Check(err, name)
// 			if alphaSet[1] < int(dummy) {
// 				log.Fatalf("Num of alphas not monotonically decreasing")
// 			}
// 			alphaSet = append(alphaSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if alphaSet[2] < int(dummy) {
// 				log.Fatalf("Num of alphas not monotonically decreasing")
// 			}
// 			alphaSet = append(alphaSet, int(dummy))
// 			numAlphas = append(numAlphas, alphaSet)
// 		}
// 	}

// 	return numAlphas
// }

// func OdatAtomAlphas(odat [][]string) [][]float64 {

// 	// define the name of the function. this is for error handling
// 	name := "fileOps.OdatAlphas"

// 	var alphas [][]float64
// 	var alphaSet []float64

// 	// first, get the number of alphas for all types in the system.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// set a counter, so that we know which type we are working on.
// 	typeNum := 0

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "ALPHAS" {
// 			// this is a bit tricky. so, once we found the "ALPHAS"
// 			// tag, we know that the actual alphas are on the next line.
// 			// there are 4 alphas per line, so we will read each item,
// 			// until 4 alphas are read, then we go to the next line.
// 			//
// 			// another note, we use the number of alphas in the S orbital
// 			// for this type (numAlphas[type][orbital]) as it will always
// 			// be the max number of alphas this type has.
// 			// start with line = 1,
// 			line := 1
// 			item := 0
// 			// reset the alphaSet to store a fresh set of alphas.
// 			alphaSet = alphaSet[:0]
// 			for alpha := 0; alpha < numAlphas[typeNum][0]; alpha++ {
// 				dummy, err := strconv.ParseFloat(odat[i+line][item], 64)
// 				utilFns.Check(err, name)
// 				alphaSet = append(alphaSet, dummy)
// 				item += 1
// 				if item == 4 {
// 					item = 0
// 					line += 1
// 				}
// 			}
// 			alphas = append(alphas, alphaSet)
// 			// increment the type number
// 			typeNum += 1

// 			// if we have found all the types, exit the loop.
// 			if typeNum == len(numAlphas) {
// 				break
// 			}
// 		}
// 	}
// 	return alphas
// }

// func OdatNumCoreRadFns(odat [][]string) [][]int {

// 	var radFnSet []int
// 	var numCoreRadFns [][]int

// 	// define name of the function. this is done for error handling.
// 	name := "fileOps.OdatNumAlphas"

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_CORE_RADIAL_FNS" {
// 			// reset the radFnSet to be empty
// 			radFnSet = radFnSet[:0]
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			radFnSet = append(radFnSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][1], 10, 32)
// 			utilFns.Check(err, name)
// 			radFnSet = append(radFnSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][2], 10, 32)
// 			utilFns.Check(err, name)
// 			radFnSet = append(radFnSet, int(dummy))
// 			numCoreRadFns = append(numCoreRadFns, radFnSet)
// 		}
// 	}

// 	return numCoreRadFns
// }

// func OdatNumValeRadFns(odat [][]string) [][]int {

// 	var radFnSet []int
// 	var numValeRadFns [][]int

// 	// define name of the function. this is done for error handling.
// 	name := "fileOps.OdatNumAlphas"

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_VALE_RADIAL_FNS" {
// 			// reset the radFnSet to be empty
// 			radFnSet = radFnSet[:0]
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			radFnSet = append(radFnSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][1], 10, 32)
// 			utilFns.Check(err, name)
// 			radFnSet = append(radFnSet, int(dummy))
// 			dummy, err = strconv.ParseInt(odat[i+1][2], 10, 32)
// 			utilFns.Check(err, name)
// 			radFnSet = append(radFnSet, int(dummy))
// 			numValeRadFns = append(numValeRadFns, radFnSet)
// 		}
// 	}

// 	return numValeRadFns
// }

// func OdatCoreQNs(odat [][]string) [][][]int {

// 	// define the name of the fnction. this is used for error handling.
// 	name := "fileOps.OdatCoreQNs"

// 	// lets get the number of types in this file.
// 	numTypes := OdatNumAtomTypes(odat)

// 	// lets get the number of alphas for the types. we need this to determine
// 	// how many lines to skip after reading each QN set.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// lets also get the number of core radial function this type has.
// 	// each radial function will have its own set of quantum numbers.
// 	numCoreRadFns := OdatNumCoreRadFns(odat)

// 	// now we can allocate the space
// 	coreQNs := make([][][]int, numTypes)
// 	for i := 0; i < numTypes; i++ {
// 		// we allocate enough space to hold the highest number of core radial
// 		// functions that this type can have, i.e. the number of core radial
// 		// functions in the extended basis. thats indicated by the [2] below.
// 		coreQNs[i] = make([][]int, numCoreRadFns[i][2])
// 		for j := 0; j < numCoreRadFns[i][2]; j++ {
// 			// each type has 5 QNs that we want: n, l, 2j, number of states
// 			// in component, and component index.
// 			coreQNs[i][j] = make([]int, 5)
// 		}
// 	}

// 	// now that space has been allocated, we are ready to fill it in.
// 	// lets make a counter so that we know which type we are working on.
// 	typeNum := 0

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_CORE_RADIAL_FNS" {
// 			// set a marker for the line number where the flag was found.
// 			line := i

// 			// we will skip the next two lines.
// 			line += 2

// 			// after each set of QNs, we need to skip over some lines
// 			// on which the radial functions are listed, incase this
// 			// type had more than 1 set of radial functions, and
// 			// therefore more than one set of QNs to be read. so lets
// 			// determine how many lines to skip. the radial functions
// 			// coefficients are written 4 per line. so the number of
// 			// lines is:
// 			numSkipLines := int(math.Floor(float64(numAlphas[typeNum][0] / 4)))
// 			// note that we used the number of alphas in the s obrital
// 			// above.
// 			// we also need to add one more line, if there are any left
// 			// over radial function coefficients that do not fill a
// 			// whole line of 4.
// 			if math.Mod(float64(numAlphas[typeNum][0]), 4.0) != 0.0 {
// 				numSkipLines += 1
// 			}
// 			// now we will read each set of radial functions in turn.
// 			for j := 0; j < numCoreRadFns[typeNum][2]; j++ {
// 				// skip 2 lines. this puts us at the line where the QNs
// 				// are listed
// 				line += 2

// 				// now we will read the QNs. first, the n.
// 				dummy, err := strconv.ParseInt(odat[line][0], 10, 32)
// 				utilFns.Check(err, name)
// 				coreQNs[typeNum][j][0] = int(dummy)
// 				// now the l.
// 				dummy, err = strconv.ParseInt(odat[line][1], 10, 32)
// 				utilFns.Check(err, name)
// 				coreQNs[typeNum][j][1] = int(dummy)
// 				// now the 2j.
// 				dummy, err = strconv.ParseInt(odat[line][2], 10, 32)
// 				utilFns.Check(err, name)
// 				coreQNs[typeNum][j][2] = int(dummy)
// 				// now the number of states in this component.
// 				dummy, err = strconv.ParseInt(odat[line][3], 10, 32)
// 				utilFns.Check(err, name)
// 				coreQNs[typeNum][j][3] = int(dummy)
// 				// finally, the component index.
// 				dummy, err = strconv.ParseInt(odat[line][4], 10, 32)
// 				utilFns.Check(err, name)
// 				coreQNs[typeNum][j][4] = int(dummy)
// 				line += numSkipLines
// 			}
// 			typeNum += 1
// 			if typeNum == numTypes {
// 				break
// 			}
// 		}
// 	}

// 	return coreQNs
// }

// func OdatValeQNs(odat [][]string) [][][]int {

// 	// define the name of the fnction. this is used for error handling.
// 	name := "fileOps.OdatValeQNs"

// 	// lets get the number of types in this file.
// 	numTypes := OdatNumAtomTypes(odat)

// 	// lets also get the number of alphas.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// lets also get the number of vale radial function this type has.
// 	// each radial function will have its own set of quantum numbers.
// 	numValeRadFns := OdatNumValeRadFns(odat)

// 	// now we can allocate the space
// 	valeQNs := make([][][]int, numTypes)
// 	for i := 0; i < numTypes; i++ {
// 		// we allocate enough space to hold the highest number of vale radial
// 		// functions that this type can have, i.e. the number of vale radial
// 		// functions in the extended basis. thats indicated by the [2] below.
// 		valeQNs[i] = make([][]int, numValeRadFns[i][2])
// 		for j := 0; j < numValeRadFns[i][2]; j++ {
// 			// each type has 5 QNs that we want: n, l, 2j, number of states
// 			// in component, and component index.
// 			valeQNs[i][j] = make([]int, 5)
// 		}
// 	}

// 	// now that space has been allocated, we are ready to fill it in.
// 	// lets make a counter so that we know which type we are working on.
// 	typeNum := 0

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_VALE_RADIAL_FNS" {
// 			// set a marker for the line number where the flag was found.
// 			line := i

// 			// we will skip the next two lines.
// 			line += 2

// 			// after each set of QNs, we need to skip over some lines
// 			// on which the radial functions are listed, incase this
// 			// type had more than 1 set of radial functions, and
// 			// therefore more than one set of QNs to be read. so lets
// 			// determine how many lines to skip. the radial functions
// 			// coefficients are written 4 per line. so the number of
// 			// lines is:
// 			numSkipLines := int(math.Floor(float64(numAlphas[typeNum][0] / 4)))
// 			// note that we used the number of alphas in the s obrital
// 			// above.
// 			// we also need to add one more line, if there are any left
// 			// over radial function coefficients that do not fill a
// 			// whole line of 4.
// 			if math.Mod(float64(numAlphas[typeNum][0]), 4.0) != 0.0 {
// 				numSkipLines += 1
// 			}
// 			// now we will read each set of radial functions in turn.
// 			for j := 0; j < numValeRadFns[typeNum][2]; j++ {
// 				// skip 2 lines. this puts us at the line where the QNs
// 				// are listed
// 				line += 2

// 				// now we will read the QNs. first, the n.
// 				dummy, err := strconv.ParseInt(odat[line][0], 10, 32)
// 				utilFns.Check(err, name)
// 				valeQNs[typeNum][j][0] = int(dummy)
// 				// now the l.
// 				dummy, err = strconv.ParseInt(odat[line][1], 10, 32)
// 				utilFns.Check(err, name)
// 				valeQNs[typeNum][j][1] = int(dummy)
// 				// now the 2j.
// 				dummy, err = strconv.ParseInt(odat[line][2], 10, 32)
// 				utilFns.Check(err, name)
// 				valeQNs[typeNum][j][2] = int(dummy)
// 				// now the number of states in this component.
// 				dummy, err = strconv.ParseInt(odat[line][3], 10, 32)
// 				utilFns.Check(err, name)
// 				valeQNs[typeNum][j][3] = int(dummy)
// 				// finally, the component index.
// 				dummy, err = strconv.ParseInt(odat[line][4], 10, 32)
// 				utilFns.Check(err, name)
// 				valeQNs[typeNum][j][4] = int(dummy)
// 				// now lets skip numLines ahead for the next set
// 				line += numSkipLines
// 			}
// 			typeNum += 1
// 			if typeNum == numTypes {
// 				break
// 			}
// 		}
// 	}

// 	return valeQNs
// }

// func OdatCoreRadFns(odat [][]string) [][][]float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatCoreRadFns"

// 	// lets get the dimesnions of the core radial functions in this system.
// 	// the first dimension is the number of types.
// 	numTypes := OdatNumAtomTypes(odat)

// 	// next isthe number of core radial functions for this type. in
// 	// particular, the number of core radial functions in the extended
// 	// basis, as the minimal and full bases are a subset of the extended
// 	// basis.
// 	numCoreRadFns := OdatNumCoreRadFns(odat)

// 	// the third dimension is the number of s-type alphas for this type. the
// 	// reason we choose the s-type alphas, is because p, d, and f type alphas
// 	// are a subset of s-type alphas.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// now lets allocate the sapce
// 	coreRadFns := make([][][]float64, numTypes)
// 	for i := 0; i < numTypes; i++ {
// 		coreRadFns[i] = make([][]float64, numCoreRadFns[i][2])
// 		for j := 0; j < numCoreRadFns[i][2]; j++ {
// 			coreRadFns[i][j] = make([]float64, numAlphas[i][0])
// 		}
// 	}

// 	// lets set a counter for the types.
// 	typeNum := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_CORE_RADIAL_FNS" {
// 			// set a marker for the line number where the flag was found.
// 			line := i

// 			// we will skip the next two lines.
// 			line += 2
// 			// now we will read each set of radial functions in turn.
// 			for j := 0; j < numCoreRadFns[typeNum][2]; j++ {
// 				// skip 3 lines. this puts us at the line where the coeffs
// 				// are listed
// 				line += 3

// 				// now we will read the radial function coefficients.
// 				// for an explaination about what is going on here, look at
// 				// the "OdatAlphas" function above.
// 				item := 0
// 				for alpha := 0; alpha < numAlphas[typeNum][0]; alpha++ {
// 					dummy, err := strconv.ParseFloat(odat[line][item], 64)
// 					utilFns.Check(err, name)
// 					coreRadFns[typeNum][j][alpha] = dummy
// 					item += 1
// 					if item == 4 {
// 						item = 0
// 						if alpha != (numAlphas[typeNum][0] - 1) {
// 							line += 1
// 						}
// 					}
// 				}
// 			}
// 			typeNum += 1
// 			if typeNum == numTypes {
// 				break
// 			}
// 		}
// 	}

// 	return coreRadFns
// }

// func OdatValeRadFns(odat [][]string) [][][]float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatValeRadFns"

// 	// lets get the dimesnions of the vale radial functions in this system.
// 	// the first dimension is the number of types.
// 	numTypes := OdatNumAtomTypes(odat)

// 	// next isthe number of vale radial functions for this type. in
// 	// particular, the number of vale radial functions in the extended
// 	// basis, as the minimal and full bases are a subset of the extended
// 	// basis.
// 	numValeRadFns := OdatNumValeRadFns(odat)

// 	// the third dimension is the number of s-type alphas for this type. the
// 	// reason we choose the s-type alphas, is because p, d, and f type alphas
// 	// are a subset of s-type alphas.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// now lets allocate the sapce
// 	valeRadFns := make([][][]float64, numTypes)
// 	for i := 0; i < numTypes; i++ {
// 		valeRadFns[i] = make([][]float64, numValeRadFns[i][2])
// 		for j := 0; j < numValeRadFns[i][2]; j++ {
// 			valeRadFns[i][j] = make([]float64, numAlphas[i][0])
// 		}
// 	}

// 	// lets set a counter for the types.
// 	typeNum := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_VALE_RADIAL_FNS" {
// 			// set a marker for the line number where the flag was found.
// 			line := i

// 			// we will skip the next two lines.
// 			line += 2
// 			// now we will read each set of radial functions in turn.
// 			for j := 0; j < numValeRadFns[typeNum][2]; j++ {
// 				// skip 3 lines. this puts us at the line where the coeffs
// 				// are listed
// 				line += 3

// 				// now we will read the radial function coefficients.
// 				// for an explaination about what is going on here, look at
// 				// the "OdatAlphas" function above.
// 				item := 0
// 				for alpha := 0; alpha < numAlphas[typeNum][0]; alpha++ {
// 					dummy, err := strconv.ParseFloat(odat[line][item], 64)
// 					utilFns.Check(err, name)
// 					valeRadFns[typeNum][j][alpha] = dummy
// 					item += 1
// 					if item == 4 {
// 						item = 0
// 						if alpha != (numAlphas[typeNum][0] - 1) {
// 							line += 1
// 						}
// 					}
// 				}
// 			}
// 			typeNum += 1
// 			if typeNum == numTypes {
// 				break
// 			}
// 		}
// 	}

// 	return valeRadFns
// }

// func OdatCoreBasisCodes(odat [][]string) [][]int {

// 	// define the name of the fnction. this is used for error handling.
// 	name := "fileOps.OdatCoreBasisCodes"

// 	// lets get the number of types in this file.
// 	numTypes := OdatNumAtomTypes(odat)

// 	// lets get the number of alphas for the types. we need this to determine
// 	// how many lines to skip after reading each basis code set.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// lets also get the number of core radial function this type has.
// 	// each radial function will have its own set of basis codes.
// 	numCoreRadFns := OdatNumCoreRadFns(odat)

// 	// now we can allocate the space
// 	coreBasisCodes := make([][]int, numTypes)
// 	for i := 0; i < numTypes; i++ {
// 		// we allocate enough space to hold the highest number of core radial
// 		// functions that this type can have, i.e. the number of core radial
// 		// functions in the extended basis. thats indicated by the [2] below.
// 		coreBasisCodes[i] = make([]int, numCoreRadFns[i][2])
// 	}

// 	// now that space has been allocated, we are ready to fill it in.
// 	// lets make a counter so that we know which type we are working on.
// 	typeNum := 0

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_CORE_RADIAL_FNS" {
// 			// set a marker for the line number where the flag was found.
// 			line := i

// 			// we will skip the next two lines.
// 			line += 2

// 			// after each set of QNs, we need to skip over some lines
// 			// on which the radial functions are listed, incase this
// 			// type had more than 1 set of radial functions, and
// 			// therefore more than one set of QNs to be read. so lets
// 			// determine how many lines to skip. the radial functions
// 			// coefficients are written 4 per line. so the number of
// 			// lines is:
// 			numSkipLines := int(math.Floor(float64(numAlphas[typeNum][0] / 4)))
// 			// note that we used the number of alphas in the s obrital
// 			// above.
// 			// we also need to add one more line, if there are any left
// 			// over radial function coefficients that do not fill a
// 			// whole line of 4.
// 			if math.Mod(float64(numAlphas[typeNum][0]), 4.0) != 0.0 {
// 				numSkipLines += 1
// 			}
// 			// now we will read each set of radial functions in turn.
// 			for j := 0; j < numCoreRadFns[typeNum][2]; j++ {
// 				// skip 1 line. this puts us at the line where the basis code
// 				// is listed
// 				line += 1

// 				// now we will read the QNs. first, the n.
// 				dummy, err := strconv.ParseInt(odat[line][1], 10, 32)
// 				utilFns.Check(err, name)
// 				coreBasisCodes[typeNum][j] = int(dummy)

// 				// skip a line to be put on the coefficient lines.
// 				line += 1
// 				// skip the coefficients
// 				line += numSkipLines
// 			}
// 			typeNum += 1
// 			if typeNum == numTypes {
// 				break
// 			}
// 		}
// 	}

// 	return coreBasisCodes
// }

// func OdatValeBasisCodes(odat [][]string) [][]int {

// 	// define the name of the fnction. this is used for error handling.
// 	name := "fileOps.OdatValeBasisCodes"

// 	// lets get the number of types in this file.
// 	numTypes := OdatNumAtomTypes(odat)

// 	// lets get the number of alphas for the types. we need this to determine
// 	// how many lines to skip after reading each basis code set.
// 	numAlphas := OdatNumAtomAlphas(odat)

// 	// lets also get the number of vale radial function this type has.
// 	// each radial function will have its own set of basis codes.
// 	numValeRadFns := OdatNumValeRadFns(odat)

// 	// now we can allocate the space
// 	valeBasisCodes := make([][]int, numTypes)
// 	for i := 0; i < numTypes; i++ {
// 		// we allocate enough space to hold the highest number of vale radial
// 		// functions that this type can have, i.e. the number of vale radial
// 		// functions in the extended basis. thats indicated by the [2] below.
// 		valeBasisCodes[i] = make([]int, numValeRadFns[i][2])
// 	}

// 	// now that space has been allocated, we are ready to fill it in.
// 	// lets make a counter so that we know which type we are working on.
// 	typeNum := 0

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_VALE_RADIAL_FNS" {
// 			// set a marker for the line number where the flag was found.
// 			line := i

// 			// we will skip the next two lines.
// 			line += 2

// 			// after each set of QNs, we need to skip over some lines
// 			// on which the radial functions are listed, incase this
// 			// type had more than 1 set of radial functions, and
// 			// therefore more than one set of QNs to be read. so lets
// 			// determine how many lines to skip. the radial functions
// 			// coefficients are written 4 per line. so the number of
// 			// lines is:
// 			numSkipLines := int(math.Floor(float64(numAlphas[typeNum][0] / 4)))
// 			// note that we used the number of alphas in the s obrital
// 			// above.
// 			// we also need to add one more line, if there are any left
// 			// over radial function coefficients that do not fill a
// 			// whole line of 4.
// 			if math.Mod(float64(numAlphas[typeNum][0]), 4.0) != 0.0 {
// 				numSkipLines += 1
// 			}
// 			// now we will read each set of radial functions in turn.
// 			for j := 0; j < numValeRadFns[typeNum][2]; j++ {
// 				// skip 1 line. this puts us at the line where the basis code
// 				// is listed
// 				line += 1

// 				// now we will read the QNs. first, the n.
// 				dummy, err := strconv.ParseInt(odat[line][1], 10, 32)
// 				utilFns.Check(err, name)
// 				valeBasisCodes[typeNum][j] = int(dummy)

// 				// skip a line to be put on the coefficient lines.
// 				line += 1
// 				// skip the coefficients
// 				line += numSkipLines
// 			}
// 			typeNum += 1
// 			if typeNum == numTypes {
// 				break
// 			}
// 		}
// 	}

// 	return valeBasisCodes
// }

// func OdatNumMeshPts(odat [][]string) [3]int {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatNumMeshPts"

// 	// define the array which will hold the number of mesh points
// 	var numMeshPts [3]int

// 	// find the flag, read in the 3 numbers.
// 	for i := 0; i < len(odat); i++ {
// 		if odat[i][0] == "WAVE_INPUT_DATA" {
// 			// the mesh points are on the line right after the above flag
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			numMeshPts[0] = int(dummy)
// 			dummy, err = strconv.ParseInt(odat[i+1][1], 10, 32)
// 			utilFns.Check(err, name)
// 			numMeshPts[1] = int(dummy)
// 			dummy, err = strconv.ParseInt(odat[i+1][2], 10, 32)
// 			utilFns.Check(err, name)
// 			numMeshPts[2] = int(dummy)
// 		}
// 	}

// 	return numMeshPts
// }

// func OdatNumPotTypes(odat [][]string) int {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatNumPotTypes"

// 	var numPotTypes int

// 	for i := 0; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_POTENTIAL_TYPES" {
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			numPotTypes = int(dummy)
// 			break
// 		}
// 	}
// 	return numPotTypes
// }

// func OdatPotElementIDs(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatPotElementIDs"

// 	var elemIds []int

// 	j := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "POTENTIAL_TYPE_ID__SEQUENTIAL_NUMBER" {
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			elemIds = append(elemIds, int(dummy))
// 			j += 1
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if j != int(dummy) {
// 				log.Fatalf("Types numbered out of order.")
// 			}
// 		}
// 	}

// 	return elemIds
// }

// func OdatPotSpeciesIDs(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatPotSpeciesIDs"

// 	var speciesIds []int

// 	j := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "POTENTIAL_TYPE_ID__SEQUENTIAL_NUMBER" {
// 			dummy, err := strconv.ParseInt(odat[i+1][1], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			speciesIds = append(speciesIds, int(dummy))
// 			j += 1
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if j != int(dummy) {
// 				log.Fatalf("Types numbered out of order.")
// 			}
// 		}
// 	}

// 	return speciesIds
// }

// func OdatPotTypeIDs(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatPotTypeIDs"

// 	var typeIds []int

// 	j := 0
// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "POTENTIAL_TYPE_ID__SEQUENTIAL_NUMBER" {
// 			dummy, err := strconv.ParseInt(odat[i+1][2], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			typeIds = append(typeIds, int(dummy))
// 			j += 1
// 			dummy, err = strconv.ParseInt(odat[i+1][3], 10, 32)
// 			utilFns.Check(err, name)
// 			if j != int(dummy) {
// 				log.Fatalf("Types numbered out of order.")
// 			}
// 		}
// 	}

// 	return typeIds
// }

// func OdatPotTypeLabels(odat [][]string) []string {

// 	var typeLabels []string

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "POTENTIAL_TYPE_LABEL" {
// 			typeLabels = append(typeLabels, odat[i+1][0])
// 		}
// 	}

// 	return typeLabels
// }

// func OdatNucCharges(odat [][]string) []float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatNucCharge"

// 	var nucCharges []float64

// 	for i := 0; i < len(odat); i++ {
// 		if odat[i][0] == "NUCLEAR_CHARGE__ALPHA" {
// 			dummy, err := strconv.ParseFloat(odat[i+1][0], 64)
// 			utilFns.Check(err, name)
// 			nucCharges = append(nucCharges, dummy)
// 		}
// 	}

// 	return nucCharges
// }

// func OdatNucAlphas(odat [][]string) []float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatNucCharge"

// 	var nucAlphas []float64

// 	for i := 0; i < len(odat); i++ {
// 		if odat[i][0] == "NUCLEAR_CHARGE__ALPHA" {
// 			dummy, err := strconv.ParseFloat(odat[i+1][1], 64)
// 			utilFns.Check(err, name)
// 			nucAlphas = append(nucAlphas, dummy)
// 		}
// 	}

// 	return nucAlphas
// }

// func OdatCovalRadii(odat [][]string) []float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatCovalRadii"

// 	var covalRadii []float64

// 	for i := 0; i < len(odat); i++ {
// 		if odat[i][0] == "COVALENT_RADIUS" {
// 			dummy, err := strconv.ParseFloat(odat[i+1][0], 64)
// 			utilFns.Check(err, name)
// 			covalRadii = append(covalRadii, dummy)
// 		}
// 	}

// 	return covalRadii
// }

// func OdatNumPotAlphas(odat [][]string) []int {

// 	// set the name of the function. this is used for error handling.
// 	name := "fileOps.OdatPotElementIDs"

// 	var numAlphas []int

// 	for i := 1; i < len(odat); i++ {
// 		if odat[i][0] == "NUM_ALPHAS" {
// 			dummy, err := strconv.ParseInt(odat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			// for some reason i have to do this...
// 			numAlphas = append(numAlphas, int(dummy))
// 		}
// 	}

// 	return numAlphas
// }

// func OdatPotAlphas(odat [][]string) [][]float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.OdatPotAlphas"

// 	var potAlphas [][]float64
// 	var alphaSet []float64

// 	for i := 0; i < len(odat); i++ {
// 		if odat[i][0] == "ALPHAS" {
// 			// reset the alpha set
// 			alphaSet = alphaSet[:0]
// 			// red the min alpha for this type
// 			dummy, err := strconv.ParseFloat(odat[i+1][0], 64)
// 			utilFns.Check(err, name)
// 			alphaSet = append(alphaSet, dummy)
// 			// read the max alpha for this type
// 			dummy, err = strconv.ParseFloat(odat[i+1][1], 64)
// 			utilFns.Check(err, name)
// 			alphaSet = append(alphaSet, dummy)

// 			potAlphas = append(potAlphas, alphaSet)
// 		}
// 	}

// 	return potAlphas
// }

// // the next functions are concerned with extracting information from the
// // structure.dat file.

// func SdatCellVecs(sdat [][]string) [3][3]float64 {

// 	// set the name of the function. this is used in error handling
// 	name := "fileOps.SdatCellVecs"

// 	// lets assign the array (not slice) that will contain the 3x3 real
// 	// vectors:
// 	//				ax ay az
// 	//				bx by bz
// 	//				cx cy cz
// 	//
// 	// where a, b, and c are the lattice vectors.
// 	var realVectors [3][3]float64

// 	// now lets start the conversion.
// 	// the information on on the second line (so line 1), hence why
// 	// i starts at 1.
// 	for i := 1; i < 4; i++ {
// 		dummy, err := strconv.ParseFloat(sdat[i][0], 64)
// 		utilFns.Check(err, name)
// 		realVectors[i-1][0] = dummy
// 		dummy, err = strconv.ParseFloat(sdat[i][1], 64)
// 		utilFns.Check(err, name)
// 		realVectors[i-1][1] = dummy
// 		dummy, err = strconv.ParseFloat(sdat[i][2], 64)
// 		utilFns.Check(err, name)
// 		realVectors[i-1][2] = dummy
// 	}

// 	return realVectors
// }

// func SdatNumAtomSites(sdat [][]string) int {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SdatNumAtomSites"

// 	var numAtomSites int

// 	for i := 0; i < len(sdat); i++ {
// 		if sdat[i][0] == "NUM_ATOM_SITES" {
// 			dummy, err := strconv.ParseInt(sdat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			numAtomSites = int(dummy)
// 			break
// 		}
// 	}

// 	return numAtomSites
// }

// func SdatAtomSites(sdat [][]string) [][]float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SdatAtomSites"

// 	// first get the number of atomic sites.
// 	numAtomSites := SdatNumAtomSites(sdat)

// 	// allocate the space for the cartesian coordinates of the sites
// 	atomSites := make([][]float64, numAtomSites)
// 	for i := 0; i < numAtomSites; i++ {
// 		atomSites[i] = make([]float64, 3)
// 	}

// 	for i := 0; i < len(sdat); i++ {
// 		if sdat[i][0] == "NUM_ATOM_SITES" {
// 			// the atom site info start 3 lines after the above
// 			// flag. start reading.
// 			line := i + 3
// 			for j := 0; j < numAtomSites; j++ {
// 				dumint, err := strconv.ParseInt(sdat[line+j][0], 10, 32)
// 				utilFns.Check(err, name)
// 				if j+1 != int(dumint) {
// 					log.Fatalf("atomic site list out of order")
// 				}
// 				dummy, err := strconv.ParseFloat(sdat[line+j][2], 64)
// 				utilFns.Check(err, name)
// 				atomSites[j][0] = dummy
// 				dummy, err = strconv.ParseFloat(sdat[line+j][3], 64)
// 				utilFns.Check(err, name)
// 				atomSites[j][1] = dummy
// 				dummy, err = strconv.ParseFloat(sdat[line+j][4], 64)
// 				utilFns.Check(err, name)
// 				atomSites[j][2] = dummy
// 			}
// 			break
// 		}
// 	}
// 	return atomSites
// }

// func SdatNumPotSites(sdat [][]string) int {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SdatNumPotSites"

// 	var numPotSites int

// 	for i := 0; i < len(sdat); i++ {
// 		if sdat[i][0] == "NUM_POTENTIAL_SITES" {
// 			dummy, err := strconv.ParseInt(sdat[i+1][0], 10, 32)
// 			utilFns.Check(err, name)
// 			numPotSites = int(dummy)
// 			break
// 		}
// 	}

// 	return numPotSites
// }

// func SdatPotTypeAssn(sdat [][]string) []int {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SdatPotTypeAssn"

// 	// first, get the number of potential sites in the system.
// 	numPotSites := SdatNumPotSites(sdat)

// 	var potTypeAssn []int

// 	for i := 0; i < len(sdat); i++ {
// 		if sdat[i][0] == "NUM_POTENTIAL_SITES" {
// 			// the potential type assignments start 3 lines after the above
// 			// flag. start reading.
// 			line := i + 3
// 			for j := 0; j < numPotSites; j++ {
// 				dummy, err := strconv.ParseInt(sdat[line+j][0], 10, 32)
// 				utilFns.Check(err, name)
// 				if j+1 != int(dummy) {
// 					log.Fatalf("potential site list out of order")
// 				}
// 				dummy, err = strconv.ParseInt(sdat[line+j][1], 10, 32)
// 				utilFns.Check(err, name)
// 				potTypeAssn = append(potTypeAssn, int(dummy))
// 			}
// 			break
// 		}
// 	}

// 	return potTypeAssn
// }

// func SdatPotSites(sdat [][]string) [][]float64 {

// 	// set the name of the function. this is for error handling.
// 	name := "fileOps.SdatPotSites"

// 	// first get the number of potential sites.
// 	numPotSites := SdatNumPotSites(sdat)

// 	// allocate the space for the cartesian coordinates of the sites
// 	potSites := make([][]float64, numPotSites)
// 	for i := 0; i < numPotSites; i++ {
// 		potSites[i] = make([]float64, 3)
// 	}

// 	for i := 0; i < len(sdat); i++ {
// 		if sdat[i][0] == "NUM_POTENTIAL_SITES" {
// 			// the potential site info start 3 lines after the above
// 			// flag. start reading.
// 			line := i + 3
// 			for j := 0; j < numPotSites; j++ {
// 				dumint, err := strconv.ParseInt(sdat[line+j][0], 10, 32)
// 				utilFns.Check(err, name)
// 				if j+1 != int(dumint) {
// 					log.Fatalf("potential site list out of order")
// 				}
// 				dummy, err := strconv.ParseFloat(sdat[line+j][2], 64)
// 				utilFns.Check(err, name)
// 				potSites[j][0] = dummy
// 				dummy, err = strconv.ParseFloat(sdat[line+j][3], 64)
// 				utilFns.Check(err, name)
// 				potSites[j][1] = dummy
// 				dummy, err = strconv.ParseFloat(sdat[line+j][4], 64)
// 				utilFns.Check(err, name)
// 				potSites[j][2] = dummy
// 			}
// 			break
// 		}
// 	}

// 	return potSites
// }
