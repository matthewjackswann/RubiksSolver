package cube

import (
	"fmt"
	"github.com/davidminor/uint128"
	"golang.org/x/exp/slices"
	"strings"
	"unicode"
)

type Cube struct {
	Layout         [54]int
	previousLayout [54]int
}

var SolvedCubeId, _ = NewSolvedCube().EncodeCube()

var fRotationMap = map[int]int{
	6:  15,
	15: 47,
	47: 35,
	35: 6,
	7:  27,
	27: 46,
	46: 23,
	23: 7,
	8:  39,
	39: 45,
	45: 11,
	11: 8,
	12: 14,
	14: 38,
	38: 36,
	36: 12,
	13: 26,
	26: 37,
	37: 24,
	24: 13,
}

var fiRotationMap = inverseTransformMap(fRotationMap)

var lRotationMap = map[int]int{
	0:  12,
	12: 45,
	45: 44,
	44: 0,
	3:  24,
	24: 48,
	48: 32,
	32: 3,
	6:  36,
	36: 51,
	51: 20,
	20: 6,
	9:  11,
	11: 35,
	35: 33,
	33: 9,
	10: 23,
	23: 34,
	34: 21,
	21: 10,
}

var liRotationMap = inverseTransformMap(lRotationMap)

var rRotationMap = map[int]int{
	2:  42,
	42: 47,
	47: 14,
	14: 2,
	5:  30,
	30: 50,
	50: 26,
	26: 5,
	8:  18,
	18: 53,
	53: 38,
	38: 8,
	15: 17,
	17: 41,
	41: 39,
	39: 15,
	16: 29,
	29: 40,
	40: 27,
	27: 16,
}

var riRotationMap = inverseTransformMap(rRotationMap)

var bRotationMap = map[int]int{
	0:  33,
	33: 53,
	53: 17,
	17: 0,
	1:  21,
	21: 52,
	52: 29,
	29: 1,
	2:  9,
	9:  51,
	51: 41,
	41: 2,
	18: 20,
	20: 44,
	44: 42,
	42: 18,
	19: 32,
	32: 43,
	43: 30,
	30: 19,
}

var biRotationMap = inverseTransformMap(bRotationMap)

var uRotationMap = map[int]int{
	0:  2,
	2:  8,
	8:  6,
	6:  0,
	1:  5,
	5:  7,
	7:  3,
	3:  1,
	9:  18,
	18: 15,
	15: 12,
	12: 9,
	10: 19,
	19: 16,
	16: 13,
	13: 10,
	11: 20,
	20: 17,
	17: 14,
	14: 11,
}

var uiRotationMap = inverseTransformMap(uRotationMap)

var dRotationMap = map[int]int{
	33: 36,
	36: 39,
	39: 42,
	42: 33,
	34: 37,
	37: 40,
	40: 43,
	43: 34,
	35: 38,
	38: 41,
	41: 44,
	44: 35,
	45: 47,
	47: 53,
	53: 51,
	51: 45,
	46: 50,
	50: 52,
	52: 48,
	48: 46,
}

var diRotationMap = inverseTransformMap(dRotationMap)

var xRotationMap = map[int]int{
	0:  44,
	44: 45,
	45: 12,
	12: 0,
	1:  43,
	43: 46,
	46: 13,
	13: 1,
	2:  42,
	42: 47,
	47: 14,
	14: 2,
	3:  32,
	32: 48,
	48: 24,
	24: 3,
	4:  31,
	31: 49,
	49: 25,
	25: 4,
	5:  30,
	30: 50,
	50: 26,
	26: 5,
	6:  20,
	20: 51,
	51: 36,
	36: 6,
	7:  19,
	19: 52,
	52: 37,
	37: 7,
	8:  18,
	18: 53,
	53: 38,
	38: 8,
	9:  33,
	33: 35,
	35: 11,
	11: 9,
	10: 21,
	21: 34,
	34: 23,
	23: 10,
	15: 17,
	17: 41,
	41: 39,
	39: 15,
	16: 29,
	29: 40,
	40: 27,
	27: 16,
}

var xiRotationMap = inverseTransformMap(xRotationMap)

var yRotationMap = map[int]int{
	0:  2,
	2:  8,
	8:  6,
	6:  0,
	1:  5,
	5:  7,
	7:  3,
	3:  1,
	9:  18,
	18: 15,
	15: 12,
	12: 9,
	10: 19,
	19: 16,
	16: 13,
	13: 10,
	11: 20,
	20: 17,
	17: 14,
	14: 11,
	21: 30,
	30: 27,
	27: 24,
	24: 21,
	22: 31,
	31: 28,
	28: 25,
	25: 22,
	23: 32,
	32: 29,
	29: 26,
	26: 23,
	33: 42,
	42: 39,
	39: 36,
	36: 33,
	34: 43,
	43: 40,
	40: 37,
	37: 34,
	35: 44,
	44: 41,
	41: 38,
	38: 35,
	45: 51,
	51: 53,
	53: 47,
	47: 45,
	46: 48,
	48: 52,
	52: 50,
	50: 46,
}

var yiRotationMap = inverseTransformMap(yRotationMap)

var zRotationMap = map[int]int{
	0:  17,
	17: 53,
	53: 33,
	33: 0,
	1:  29,
	29: 52,
	52: 21,
	21: 1,
	2:  41,
	41: 51,
	51: 9,
	9:  2,
	3:  16,
	16: 50,
	50: 34,
	34: 3,
	4:  28,
	28: 49,
	49: 22,
	22: 4,
	5:  40,
	40: 48,
	48: 10,
	10: 5,
	6:  15,
	15: 47,
	47: 35,
	35: 6,
	7:  27,
	27: 46,
	46: 23,
	23: 7,
	8:  39,
	39: 45,
	45: 11,
	11: 8,
	12: 14,
	14: 38,
	38: 36,
	36: 12,
	13: 26,
	26: 37,
	37: 24,
	24: 13,
	18: 42,
	42: 44,
	44: 20,
	20: 18,
	19: 30,
	30: 43,
	43: 32,
	32: 19,
}

var ziRotationMap = inverseTransformMap(zRotationMap)

func inverseTransformMap(m map[int]int) map[int]int {
	n := make(map[int]int, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

func NewSolvedCube() *Cube {
	return &Cube{
		Layout:         [54]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5},
		previousLayout: [54]int{},
	}
}

func NewCube(layout [54]int) *Cube {
	return &Cube{
		Layout:         layout,
		previousLayout: [54]int{},
	}
}

func (cube *Cube) applyTransformMap(transformMap map[int]int) {
	for i := 0; i < 54; i++ {
		newIndex, exists := transformMap[i]
		if exists {
			cube.previousLayout[newIndex] = cube.Layout[i] // loads into old state
		} else {
			cube.previousLayout[i] = cube.Layout[i] // loads into old state
		}
	}
	cube.previousLayout, cube.Layout = cube.Layout, cube.previousLayout // swaps old state and new state
}

func (cube *Cube) transform(t string) {
	switch t {
	case "F":
		cube.applyTransformMap(fRotationMap)
	case "f":
		cube.applyTransformMap(fiRotationMap)
	case "L":
		cube.applyTransformMap(lRotationMap)
	case "l":
		cube.applyTransformMap(liRotationMap)
	case "R":
		cube.applyTransformMap(rRotationMap)
	case "r":
		cube.applyTransformMap(riRotationMap)
	case "B":
		cube.applyTransformMap(bRotationMap)
	case "b":
		cube.applyTransformMap(biRotationMap)
	case "U":
		cube.applyTransformMap(uRotationMap)
	case "u":
		cube.applyTransformMap(uiRotationMap)
	case "D":
		cube.applyTransformMap(dRotationMap)
	case "d":
		cube.applyTransformMap(diRotationMap)
	case "X":
		cube.applyTransformMap(xRotationMap)
	case "x":
		cube.applyTransformMap(xiRotationMap)
	case "Y":
		cube.applyTransformMap(yRotationMap)
	case "y":
		cube.applyTransformMap(yiRotationMap)
	case "Z":
		cube.applyTransformMap(zRotationMap)
	case "z":
		cube.applyTransformMap(ziRotationMap)
	default:
		fmt.Printf("Invalid Transform :%v\n", t)
	}
}

func (cube *Cube) Transform(t string) {
	for _, t := range strings.Split(t, "") {
		cube.transform(t)
	}
}

var idTranslations = [24][54]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53},
	{33, 21, 9, 34, 22, 10, 35, 23, 11, 51, 48, 45, 36, 24, 12, 6, 3, 0, 20, 32, 44, 52, 49, 46, 37, 25, 13, 7, 4, 1, 19, 31, 43, 53, 50, 47, 38, 26, 14, 8, 5, 2, 18, 30, 42, 39, 27, 15, 40, 28, 16, 41, 29, 17},
	{53, 52, 51, 50, 49, 48, 47, 46, 45, 41, 40, 39, 38, 37, 36, 35, 34, 33, 44, 43, 42, 29, 28, 27, 26, 25, 24, 23, 22, 21, 32, 31, 30, 17, 16, 15, 14, 13, 12, 11, 10, 9, 20, 19, 18, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	{17, 29, 41, 16, 28, 40, 15, 27, 39, 2, 5, 8, 14, 26, 38, 47, 50, 53, 42, 30, 18, 1, 4, 7, 13, 25, 37, 46, 49, 52, 43, 31, 19, 0, 3, 6, 12, 24, 36, 45, 48, 51, 44, 32, 20, 11, 23, 35, 10, 22, 34, 9, 21, 33},
	{6, 3, 0, 7, 4, 1, 8, 5, 2, 12, 13, 14, 15, 16, 17, 18, 19, 20, 9, 10, 11, 24, 25, 26, 27, 28, 29, 30, 31, 32, 21, 22, 23, 36, 37, 38, 39, 40, 41, 42, 43, 44, 33, 34, 35, 47, 50, 53, 46, 49, 52, 45, 48, 51},
	{36, 24, 12, 37, 25, 13, 38, 26, 14, 45, 46, 47, 39, 27, 15, 8, 7, 6, 11, 23, 35, 48, 49, 50, 40, 28, 16, 5, 4, 3, 10, 22, 34, 51, 52, 53, 41, 29, 17, 2, 1, 0, 9, 21, 33, 42, 30, 18, 43, 31, 19, 44, 32, 20},
	{51, 48, 45, 52, 49, 46, 53, 50, 47, 44, 43, 42, 41, 40, 39, 38, 37, 36, 35, 34, 33, 32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 2, 5, 8, 1, 4, 7, 0, 3, 6},
	{20, 32, 44, 19, 31, 43, 18, 30, 42, 0, 1, 2, 17, 29, 41, 53, 52, 51, 33, 21, 9, 3, 4, 5, 16, 28, 40, 50, 49, 48, 34, 22, 10, 6, 7, 8, 15, 27, 39, 47, 46, 45, 35, 23, 11, 14, 26, 38, 13, 25, 37, 12, 24, 36},
	{8, 7, 6, 5, 4, 3, 2, 1, 0, 15, 16, 17, 18, 19, 20, 9, 10, 11, 12, 13, 14, 27, 28, 29, 30, 31, 32, 21, 22, 23, 24, 25, 26, 39, 40, 41, 42, 43, 44, 33, 34, 35, 36, 37, 38, 53, 52, 51, 50, 49, 48, 47, 46, 45},
	{39, 27, 15, 40, 28, 16, 41, 29, 17, 47, 50, 53, 42, 30, 18, 2, 5, 8, 14, 26, 38, 46, 49, 52, 43, 31, 19, 1, 4, 7, 13, 25, 37, 45, 48, 51, 44, 32, 20, 0, 3, 6, 12, 24, 36, 33, 21, 9, 34, 22, 10, 35, 23, 11},
	{45, 46, 47, 48, 49, 50, 51, 52, 53, 35, 34, 33, 44, 43, 42, 41, 40, 39, 38, 37, 36, 23, 22, 21, 32, 31, 30, 29, 28, 27, 26, 25, 24, 11, 10, 9, 20, 19, 18, 17, 16, 15, 14, 13, 12, 0, 1, 2, 3, 4, 5, 6, 7, 8},
	{11, 23, 35, 10, 22, 34, 9, 21, 33, 6, 3, 0, 20, 32, 44, 51, 48, 45, 36, 24, 12, 7, 4, 1, 19, 31, 43, 52, 49, 46, 37, 25, 13, 8, 5, 2, 18, 30, 42, 53, 50, 47, 38, 26, 14, 17, 29, 41, 16, 28, 40, 15, 27, 39},
	{2, 5, 8, 1, 4, 7, 0, 3, 6, 18, 19, 20, 9, 10, 11, 12, 13, 14, 15, 16, 17, 30, 31, 32, 21, 22, 23, 24, 25, 26, 27, 28, 29, 42, 43, 44, 33, 34, 35, 36, 37, 38, 39, 40, 41, 51, 48, 45, 52, 49, 46, 53, 50, 47},
	{42, 30, 18, 43, 31, 19, 44, 32, 20, 53, 52, 51, 33, 21, 9, 0, 1, 2, 17, 29, 41, 50, 49, 48, 34, 22, 10, 3, 4, 5, 16, 28, 40, 47, 46, 45, 35, 23, 11, 6, 7, 8, 15, 27, 39, 36, 24, 12, 37, 25, 13, 38, 26, 14},
	{47, 50, 53, 46, 49, 52, 45, 48, 51, 38, 37, 36, 35, 34, 33, 44, 43, 42, 41, 40, 39, 26, 25, 24, 23, 22, 21, 32, 31, 30, 29, 28, 27, 14, 13, 12, 11, 10, 9, 20, 19, 18, 17, 16, 15, 6, 3, 0, 7, 4, 1, 8, 5, 2},
	{14, 26, 38, 13, 25, 37, 12, 24, 36, 8, 7, 6, 11, 23, 35, 45, 46, 47, 39, 27, 15, 5, 4, 3, 10, 22, 34, 48, 49, 50, 40, 28, 16, 2, 1, 0, 9, 21, 33, 51, 52, 53, 41, 29, 17, 20, 32, 44, 19, 31, 43, 18, 30, 42},
	{12, 13, 14, 24, 25, 26, 36, 37, 38, 11, 23, 35, 45, 46, 47, 39, 27, 15, 8, 7, 6, 10, 22, 34, 48, 49, 50, 40, 28, 16, 5, 4, 3, 9, 21, 33, 51, 52, 53, 41, 29, 17, 2, 1, 0, 44, 43, 42, 32, 31, 30, 20, 19, 18},
	{9, 10, 11, 21, 22, 23, 33, 34, 35, 20, 32, 44, 51, 48, 45, 36, 24, 12, 6, 3, 0, 19, 31, 43, 52, 49, 46, 37, 25, 13, 7, 4, 1, 18, 30, 42, 53, 50, 47, 38, 26, 14, 8, 5, 2, 41, 40, 39, 29, 28, 27, 17, 16, 15},
	{18, 19, 20, 30, 31, 32, 42, 43, 44, 17, 29, 41, 53, 52, 51, 33, 21, 9, 0, 1, 2, 16, 28, 40, 50, 49, 48, 34, 22, 10, 3, 4, 5, 15, 27, 39, 47, 46, 45, 35, 23, 11, 6, 7, 8, 38, 37, 36, 26, 25, 24, 14, 13, 12},
	{15, 16, 17, 27, 28, 29, 39, 40, 41, 14, 26, 38, 47, 50, 53, 42, 30, 18, 2, 5, 8, 13, 25, 37, 46, 49, 52, 43, 31, 19, 1, 4, 7, 12, 24, 36, 45, 48, 51, 44, 32, 20, 0, 3, 6, 35, 34, 33, 23, 22, 21, 11, 10, 9},
	{44, 43, 42, 32, 31, 30, 20, 19, 18, 33, 21, 9, 0, 1, 2, 17, 29, 41, 53, 52, 51, 34, 22, 10, 3, 4, 5, 16, 28, 40, 50, 49, 48, 35, 23, 11, 6, 7, 8, 15, 27, 39, 47, 46, 45, 12, 13, 14, 24, 25, 26, 36, 37, 38},
	{35, 34, 33, 23, 22, 21, 11, 10, 9, 36, 24, 12, 6, 3, 0, 20, 32, 44, 51, 48, 45, 37, 25, 13, 7, 4, 1, 19, 31, 43, 52, 49, 46, 38, 26, 14, 8, 5, 2, 18, 30, 42, 53, 50, 47, 15, 16, 17, 27, 28, 29, 39, 40, 41},
	{38, 37, 36, 26, 25, 24, 14, 13, 12, 39, 27, 15, 8, 7, 6, 11, 23, 35, 45, 46, 47, 40, 28, 16, 5, 4, 3, 10, 22, 34, 48, 49, 50, 41, 29, 17, 2, 1, 0, 9, 21, 33, 51, 52, 53, 18, 19, 20, 30, 31, 32, 42, 43, 44},
	{41, 40, 39, 29, 28, 27, 17, 16, 15, 42, 30, 18, 2, 5, 8, 14, 26, 38, 47, 50, 53, 43, 31, 19, 1, 4, 7, 13, 25, 37, 46, 49, 52, 44, 32, 20, 0, 3, 6, 12, 24, 36, 45, 48, 51, 9, 10, 11, 21, 22, 23, 33, 34, 35},
}

var idTranslationTransforms = [24]string{
	"", "Z", "ZZ", "ZZZ", "Y", "YZ", "YZZ", "YZZZ", "YY", "YYZ", "YYZZ", "YYZZZ", "YYY", "YYYZ", "YYYZZ", "YYYZZZ", "X", "XZ", "XZZ", "XZZZ", "XXX", "XXXZ", "XXXZZ", "XXXZZZ",
}

var faceCenters = []int{4, 22, 25, 28, 31, 49}

func (cube *Cube) EncodeCube() (uint128.Uint128, string) {
	lowestId := uint128.Uint128{H: ^uint64(0), L: ^uint64(0)}
	lowestIdRotation := 0
	for i, translation := range idTranslations {
		thisId := uint128.Uint128{}
		thisMap := map[int]int{
			cube.Layout[translation[4]]:  0,
			cube.Layout[translation[22]]: 1,
			cube.Layout[translation[25]]: 2,
			cube.Layout[translation[28]]: 3,
			cube.Layout[translation[31]]: 4,
			cube.Layout[translation[49]]: 5,
		}
		for i, t := range translation {
			if slices.Contains(faceCenters, i) {
				continue
			}
			colour := cube.Layout[t]
			mappedColour, success := thisMap[colour]
			if !success {
				_ = fmt.Errorf("can't map colour %d, %d should be < 6", colour, colour)
			}
			thisId = thisId.Mult(uint128.Uint128{L: 6})
			thisId = thisId.Add(uint128.Uint128{L: uint64(mappedColour)})
		}
		if thisId.H < lowestId.H || (thisId.H == lowestId.H && thisId.L < lowestId.L) {
			lowestId = thisId
			lowestIdRotation = i
		}
	}
	return lowestId, strings.ToLower(idTranslationTransforms[lowestIdRotation])
}

func (cube *Cube) IsSolved() bool {
	cubeId, _ := cube.EncodeCube()
	return cubeId.Equals(SolvedCubeId)
}

func (cube *Cube) GetNonSymmetricalRotations() []string {
	var rotations []string
	var ids []uint128.Uint128
	for i, translation := range idTranslations {
		thisId := uint128.Uint128{}
		thisMap := map[int]int{
			cube.Layout[translation[4]]:  0,
			cube.Layout[translation[22]]: 1,
			cube.Layout[translation[25]]: 2,
			cube.Layout[translation[28]]: 3,
			cube.Layout[translation[31]]: 4,
			cube.Layout[translation[49]]: 5,
		}
		for i, t := range translation {
			if slices.Contains(faceCenters, i) {
				continue
			}
			colour := cube.Layout[t]
			mappedColour, success := thisMap[colour]
			if !success {
				_ = fmt.Errorf("can't map colour %d, %d should be < 6", colour, colour)
			}
			thisId = thisId.Mult(uint128.Uint128{L: 6})
			thisId = thisId.Add(uint128.Uint128{L: uint64(mappedColour)})
		}
		unseenId := true
		for _, id := range ids {
			if thisId.Equals(id) {
				unseenId = false
				break
			}
		}
		if unseenId {
			rotations = append(rotations, idTranslationTransforms[i])
			ids = append(ids, thisId)
		}
	}
	return rotations
}

var xRotationTransform = map[rune]rune{
	'L': 'L',
	'R': 'R',
	'F': 'U',
	'U': 'B',
	'B': 'D',
	'D': 'F',
}

var xiRotationTransform = inverseRotationTransform(xRotationTransform)

var yRotationTransform = map[rune]rune{
	'U': 'U',
	'D': 'D',
	'F': 'L',
	'L': 'B',
	'B': 'R',
	'R': 'F',
}

var yiRotationTransform = inverseRotationTransform(yRotationTransform)

var zRotationTransform = map[rune]rune{
	'F': 'F',
	'B': 'B',
	'L': 'U',
	'U': 'R',
	'R': 'D',
	'D': 'L',
}

var ziRotationTransform = inverseRotationTransform(zRotationTransform)

func inverseRotationTransform(m map[rune]rune) map[rune]rune {
	n := make(map[rune]rune, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

var rotationMap = map[rune]map[rune]rune{
	'X': xRotationTransform,
	'x': xiRotationTransform,
	'Y': yRotationTransform,
	'y': yiRotationTransform,
	'Z': zRotationTransform,
	'z': ziRotationTransform,
}

func RotateTransform(rotation, transform string) string {
	faceMapA := map[rune]rune{
		'F': 'F',
		'L': 'L',
		'R': 'R',
		'B': 'B',
		'U': 'U',
		'D': 'D',
	}
	faceMapB := make(map[rune]rune, 6)
	for _, char := range []rune(rotation) {
		m, validRotation := rotationMap[char]
		if !validRotation {
			continue
		}
		for k, v := range m {
			faceMapB[k] = faceMapA[v]
		}
		faceMapA, faceMapB = faceMapB, faceMapA
	}
	sb := strings.Builder{}
	for _, char := range []rune(transform) {
		if unicode.IsUpper(char) {
			sb.WriteRune(faceMapA[char])
		} else {
			sb.WriteRune(unicode.ToLower(faceMapA[unicode.ToUpper(char)]))
		}
	}
	return sb.String()
}

func ReverseTransform(transform string) string {
	res := strings.Builder{}
	transformRunes := []rune(transform)
	for i := len(transform) - 1; i >= 0; i-- {
		c := transformRunes[i]
		if unicode.IsUpper(c) {
			res.WriteRune(unicode.ToLower(c))
		} else {
			res.WriteRune(unicode.ToUpper(c))
		}
	}
	return res.String()
}

func RemoveRotationTransforms(transform string) string {
	result := strings.Builder{}
	faceMapA := map[rune]rune{
		'F': 'F',
		'L': 'L',
		'R': 'R',
		'B': 'B',
		'U': 'U',
		'D': 'D',
	}
	faceMapB := make(map[rune]rune, 6)
	for _, char := range []rune(transform) {
		_, validRotation := rotationMap[char]
		if validRotation {
			var m map[rune]rune
			if unicode.IsUpper(char) {
				m = rotationMap[unicode.ToLower(char)]
			} else {
				m = rotationMap[unicode.ToUpper(char)]
			}
			for k, v := range m {
				faceMapB[k] = faceMapA[v]
			}
			faceMapA, faceMapB = faceMapB, faceMapA
		} else { // must be a transform
			if unicode.IsUpper(char) {
				result.WriteRune(faceMapA[char])
			} else {
				result.WriteRune(unicode.ToLower(faceMapA[unicode.ToUpper(char)]))
			}
		}

	}
	return result.String()
}
