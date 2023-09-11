package cube

import (
	"regexp"
	"testing"
)

// check if a string fails the regex then it's not in the generator output
func TestGenerator_Accuracy(t *testing.T) {
	g := CreateNewGenerator([]int{0}, 0, ID_TRANSFORM_GRAPH)

	// map of all transforms. Being used as a set implementation
	transforms := map[string]struct{}{}
	for true {
		s := g.Next()
		if len(s) > 4 {
			break
		}
		transforms[s] = struct{}{}
	}

	transform := []string{"F", "f", "L", "l", "U", "u", "B", "b", "R", "r", "D", "d", ""}

	// This is a blacklist, accepted transforms shouldn't match
	testingRegex := regexp.MustCompile("fF|ff|Ff|lL|ll|Ll|uU|uu|Uu|bB|bb|Bb|rR|rr|Rr|dD|dd|Dd|FFF|LLL|UUU|BBB|RRR|DDD|(B|b|(B2))(F|f|(F2))|(R|r|(R2))(L|l|(L2))|(D|d|(D2))(U|u|(U2))")

	// There are 12 possible transforms (+1 for no transform). Check all accepted combinations are covered
	maxTransformVal := 2 * 13 * 13 * 13
	for i := 0; i < maxTransformVal; i++ {
		transformString := transform[i%2] + transform[(i/2)%13] + transform[(i*2*13)%13] + transform[(i*2*13*13)%13]
		passesRegex := testingRegex.FindString(transformString) == ""
		if passesRegex {
			_, ok := transforms[transformString]
			if !ok {
				t.Errorf("Transform passes regex but is not in generator: %s", transformString)
			}
		} else {
			_, ok := transforms[transformString]
			if ok {
				t.Errorf("Transform fails regex but is in generator: %s", transformString)
			}
		}
	}
}

func TestGenerator_Valid(t *testing.T) {
	// This is a blacklist, accepted transforms shouldn't match
	testingRegex := regexp.MustCompile("fF|ff|Ff|lL|ll|Ll|uU|uu|Uu|bB|bb|Bb|rR|rr|Rr|dD|dd|Dd|(B|b|(B2))(F|f|(F2))|(R|r|(R2))(L|l|(L2))|(D|d|(D2))(U|u|(U2))")
	g := CreateNewGenerator([]int{0}, 0, ID_TRANSFORM_GRAPH)
	for i := 0; i < 10000; i++ {
		s := g.Next()
		failed := testingRegex.FindString(s) != "" // this doesn't work
		if failed {
			t.Errorf("Produced banned transform combination %s", s)
		}
	}
}
