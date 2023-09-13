package cube

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func newTestingCube() *Cube {
	return &Cube{
		Layout:         [54]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53},
		previousLayout: [54]int{},
	}
}

func TestNewSolvedCube(t *testing.T) {
	c := NewSolvedCube()
	if !reflect.DeepEqual(c.Layout, [54]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5}) {
		t.Errorf("Not a solved cube")
	}
}

func TestCube_Transform_Basic(t *testing.T) {
	identities := []string{
		"Ff", "Ll", "Rr", "Bb", "Uu", "Dd", "Xx", "Yy", "Zz", // transform then inverse
		"fF", "lL", "rR", "bB", "uU", "dD", "xX", "yY", "zZ", // inverse then transform
		"FFFF", "LLLL", "RRRR", "BBBB", "UUUU", "DDDD", "XXXX", "YYYY", "ZZZZ",
		"ffff", "llll", "rrrr", "bbbb", "uuuu", "dddd", "xxxx", "yyyy", "zzzz",
	}
	c := newTestingCube()
	targetCube := newTestingCube()
	for _, tr := range identities {
		c.Transform(tr)
		if !reflect.DeepEqual(c.Layout, targetCube.Layout) {
			t.Errorf("Cubes are not equal. Transform %v", tr)
		}
	}
}

func TestCube_Transform_Complex(t *testing.T) {
	identities := []string{
		"XRXRXRXR", "xRxRxRxR",
		"YUYUYUYU", "yUyUyUyU",
		"ZFZFZFZF", "zFzFzFzF",
		"lxlxlxlx", "lXlXlXlX",
		"dydydydy", "dYdYdYdY",
		"bzbzbzbz", "bZbZbZbZ",
		"FXux", "yLYb", "ZDzr",
	}
	c := newTestingCube()
	targetCube := newTestingCube()
	for _, tr := range identities {
		c.Transform(tr)
		if !reflect.DeepEqual(c.Layout, targetCube.Layout) {
			t.Errorf("Cubes are not equal. Transform %v\n%v\n%v", tr, c.Layout, targetCube.Layout)
		}
		c = newTestingCube()
	}
}

func TestCube_Rotations(t *testing.T) {
	tests := [][3]string{
		{"X", "FfLlRrBbUuDd", "UuLlRrDdBbFf"},
		{"x", "FfLlRrBbUuDd", "DdLlRrUuFfBb"},
		{"Y", "FfLlRrBbUuDd", "LlBbFfRrUuDd"},
		{"y", "FfLlRrBbUuDd", "RrFfBbLlUuDd"},
		{"Z", "FfLlRrBbUuDd", "FfUuDdBbRrLl"},
		{"z", "FfLlRrBbUuDd", "FfDdUuBbLlRr"},
	}
	for _, test := range tests {
		expected := test[2]
		actual := RotateTransform(test[0], test[1])
		if expected != actual {
			t.Errorf("%s + %s -> %s. Should be %s", test[0], test[1], actual, expected)
		}
	}
}

func cubeIdEquivalentCheck(t *testing.T, baseTransform string, equivalentTransforms []string) {
	c := NewSolvedCube()
	c.Transform(baseTransform)
	baseCubeId, _ := c.EncodeCube()

	for _, transform := range equivalentTransforms {
		d := NewSolvedCube()
		d.Transform(transform)
		cubeId, _ := d.EncodeCube()
		if baseCubeId != cubeId {
			t.Errorf("%s should have the same cubeId as %s", transform, baseTransform)
		}
	}
}

func TestCube_Ids(t *testing.T) {
	cubeIdEquivalentCheck(t, "f", []string{"l", "r", "u", "d", "b"})
	cubeIdEquivalentCheck(t, "F", []string{"L", "R", "U", "D", "B"})
	cubeIdEquivalentCheck(t, "LrUDBf", []string{"UdBFLr", "BfUDRl", "RlUDFb", "UdLRFb", "UdRLBf"})
	cubeIdEquivalentCheck(t, "UUrrDFbULR", []string{"DDffURlDBF", "RRffLUdRBF", "BBddFRlBUD", "RRddLFbRUD", "FFddBLrFUD"})
}

func TestCube_Transform_Commutative(t *testing.T) {
	c := NewSolvedCube()
	c.Transform("FBudRbLrUru")
	d := NewSolvedCube()
	for _, move := range strings.Split("FBudRbLrUru", "") {
		d.Transform(move)
	}
	cId, _ := c.EncodeCube()
	dId, _ := d.EncodeCube()
	if !cId.Equals(dId) {
		fmt.Println(cId)
		fmt.Println(dId)
		t.Errorf("cube ids cId and dId should be equivilent after equivilent transforms")
	}
}

func TestCube_Transform_Rotate_Commutative(t *testing.T) {
	c := NewSolvedCube()
	c.Transform("FzBudXRbLryUru")
	d := NewSolvedCube()
	for _, move := range strings.Split("FzBudXRbLryUru", "") {
		d.Transform(move)
	}
	cId, _ := c.EncodeCube()
	dId, _ := d.EncodeCube()
	if !cId.Equals(dId) {
		fmt.Println(cId)
		fmt.Println(dId)
		t.Errorf("cube ids cId and dId should be equivilent after equivilent transforms")
	}
}

func TestCube_GetNonSymmetricalRotations(t *testing.T) {
	c := NewSolvedCube()
	if len(c.GetNonSymmetricalRotations()) != 1 {
		t.Error("A solved cube should have all rotations symmetrical")
	}
	c.Transform("F")
	if len(c.GetNonSymmetricalRotations()) != 6 {
		t.Error("A cube with setup F should have 6 different rotations")
	}
	c.Transform("B")
	if len(c.GetNonSymmetricalRotations()) != 3 {
		t.Error("A cube with setup FB should have 3 different rotations")
	}
}

func TestRemoveRotationTransforms(t *testing.T) {
	tests := [][2]string{
		{"XFXFXFXF", "DBUF"},
		{"xFxFxFxF", "UBDF"},
		{"XXFXXF", "BF"},
		{"YYFYYF", "BF"},
		{"ZZFZZF", "FF"},
	}
	for _, test := range tests {
		before := test[0]
		after := test[1]
		if RemoveRotationTransforms(before) != after {
			t.Errorf("%s should be reduced to %s rather than %s", before, after, RemoveRotationTransforms(before))
		}
	}
}
