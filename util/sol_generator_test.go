package util

import (
	"fmt"
	"github.com/matthewjackswann/rubiks/cube"
	"math/rand"
	"strings"
	"testing"
)

func getLastFullLayer(stack string) int {
	stackString := strings.Split(stack, ",")
	return len(stackString) - 1
}

func checkTransformInverseExists(t *testing.T, moves string, db DBConnection, maxLength int) {
	c := cube.NewSolvedCube()
	c.Transform(moves)

	solution, solFound := db.SolveCube(c.EncodeCube())
	if !solFound {
		t.Errorf("Cube with setup %s should have a solution in the DB", moves)
	}
	if len(solution) > maxLength {
		t.Errorf("Cube with setup %s has a solution (%s) which is larger than the test layer %d", moves, solution, maxLength)
	}
	c.Transform(solution)
	if !c.IsSolved() {
		t.Errorf("Provided solution %s for cube with scramble %s doesn't solve the cube", solution, moves)
	}
}

func TestOneMoveCubes(t *testing.T) {
	db := Create("/media/swanny/Lexar/rubiks.db")

	if getLastFullLayer(db.GetNextTransforms().EncodedStack) < 1 {
		t.Skip("Layer 1 is not in the database")
	}

	for _, m := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
		checkTransformInverseExists(t, m, db, 1)
	}

	db.Close()
}

func TestTwoMoveCubes(t *testing.T) {
	db := Create("/media/swanny/Lexar/rubiks.db")

	if getLastFullLayer(db.GetNextTransforms().EncodedStack) < 2 {
		t.Skip("Layer 2 is not in the database")
	}

	for _, m0 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
		for _, m1 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
			checkTransformInverseExists(t, m0+m1, db, 2)
		}
	}

	db.Close()
}

func TestThreeMoveCubes(t *testing.T) {
	db := Create("/media/swanny/Lexar/rubiks.db")

	if getLastFullLayer(db.GetNextTransforms().EncodedStack) < 3 {
		t.Skip("Layer 3 is not in the database")
	}

	for _, m0 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
		for _, m1 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
			for _, m2 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
				checkTransformInverseExists(t, m0+m1+m2, db, 3)
			}
		}
	}

	db.Close()
}

func TestFourMoveCubes(t *testing.T) {
	db := Create("/media/swanny/Lexar/rubiks.db")

	if getLastFullLayer(db.GetNextTransforms().EncodedStack) < 4 {
		t.Skip("Layer 4 is not in the database")
	}

	for _, m0 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
		for _, m1 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
			for _, m2 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
				for _, m3 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
					checkTransformInverseExists(t, m0+m1+m2+m3, db, 4)
				}
			}
		}
	}

	db.Close()
}

func TestFiveMoveCubes(t *testing.T) {
	db := Create("/media/swanny/Lexar/rubiks.db")

	if getLastFullLayer(db.GetNextTransforms().EncodedStack) < 5 {
		t.Skip("Layer 5 is not in the database")
	}

	for _, m0 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
		for _, m1 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
			for _, m2 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
				for _, m3 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
					for _, m4 := range []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"} {
						checkTransformInverseExists(t, m0+m1+m2+m3+m4, db, 5)
					}
				}
			}
		}
	}

	db.Close()
}

func TestOneMillionCubes(t *testing.T) {
	db := Create("/media/swanny/Lexar/rubiks.db")
	rand.Seed(0)

	stringLength := getLastFullLayer(db.GetNextTransforms().EncodedStack)

	for i := 0; i < 1000000; i++ {
		if i%100000 == 0 {
			fmt.Printf("Finished iteration %d/1000000\n", i)
		}
		cubeSetup := ""
		for j := 0; j < stringLength; j++ {
			cubeSetup += []string{"f", "F", "u", "U", "l", "L", "r", "R", "b", "B", "d", "D"}[rand.Intn(12)]
		}
		checkTransformInverseExists(t, cubeSetup, db, stringLength)
	}

	db.Close()
}
