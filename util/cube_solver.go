package util

import (
	"database/sql"
	"fmt"
	"github.com/davidminor/uint128"
	"github.com/matthewjackswann/rubiks/cube"
)

func loadSolution(id uint128.Uint128, preparedStmt *sql.Stmt) (string, bool) {
	result, err := preparedStmt.Query(int64(id.L), int64(id.H))
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	if !result.Next() {
		return "", false
	}

	var encodedSolution uint64
	err = result.Scan(&encodedSolution)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	err = result.Close()
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	solution := ""
	for i := 0; i < 16; i++ {
		solutionPart := IntToTransform[encodedSolution&0xF]
		if solutionPart == 0 {
			break
		}
		solution += string(solutionPart) // Extract the 4 least significant bits
		encodedSolution >>= 4
	}

	return solution, true
}

// LookupCube is used to find the solution for a single cube if it exists in the database
func (dbConnection *DBConnection) LookupCube(cubeId uint128.Uint128, rotation string) (string, bool) {
	if cubeId.Equals(cube.SolvedCubeId) {
		return "", true
	}
	stmt, err := dbConnection.db.Prepare("SELECT solution FROM cubes WHERE cube_id_l = ? AND cube_id_h = ?;")
	if err != nil {
		fmt.Printf("Error creating prepared statement %s\n", err)
		return "", false
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println("Error closing statement")
			fmt.Println(err)
		}
	}(stmt)
	return dbConnection.lookupCube(cubeId, rotation, stmt)
}

func (dbConnection *DBConnection) lookupCube(cubeId uint128.Uint128, rotation string, stmt *sql.Stmt) (string, bool) {
	idSolution, success := loadSolution(cubeId, stmt)
	if !success {
		return "", false
	}
	return cube.RotateTransform(rotation, idSolution), true
}

func (dbConnection *DBConnection) SolveCubeBySearch(baseCube *cube.Cube, workers, maxDepth int) (string, bool) {
	solution, success := dbConnection.LookupCube(baseCube.EncodeCube())
	// not in lookup table, start brute forcing from cube direction
	if success {
		return solution, true
	}

	parallelLookup := CreateLookupWorkers(32, workers, dbConnection.path)
	cubeTransformChan := make(chan string, 32)
	cubeTransformWorkerStop := make(chan interface{}, workers)

	for i := 0; i < workers; i++ {
		go func() {
			for {
				transform := <-cubeTransformChan
				if transform == "" {
					cubeTransformWorkerStop <- nil
					return
				}
				c := cube.NewCube(baseCube.Layout)
				c.Transform(transform)
				parallelLookup.requestChan <- &lookupWorkerRequest{
					cube: c,
					data: transform,
				}
			}
		}()
	}

	baseRotations := baseCube.GetNonSymmetricalRotations()

	var generator cube.Generator
	if len(baseRotations) < 6 {
		generator = cube.CreateNewGenerator([]int{0}, 0, cube.ID_TRANSFORM_GRAPH)
	} else {
		generator = cube.CreateNewGenerator([]int{0}, 0, cube.TRANSFORM_GRAPH)
		baseRotations = []string{""} // no need to consider any other rotations. Just use the identity
	}

	resultsFound := false
	var lookupResult *lookupWorkerResponse

	transformsInBuffer := 0
	currentDepth := 0

	for {
		// send cubes
		if generator.GetCurrentDepth() != currentDepth && transformsInBuffer == 0 {
			currentDepth = generator.GetCurrentDepth()
			// if max depth reached
			if currentDepth == maxDepth+1 {
				for i := 0; i < workers; i++ {
					cubeTransformChan <- ""
				}
				for i := 0; i < workers; i++ {
					<-cubeTransformWorkerStop
				}
				parallelLookup.StopForcefully()
				return "", false
			}
		}
		if generator.GetCurrentDepth() == currentDepth && transformsInBuffer < 32-len(baseRotations) {
			baseTransform := generator.Next()
			for _, baseRotation := range baseRotations {
				cubeTransformChan <- baseRotation + baseTransform
				transformsInBuffer += 1
			}
		}

		// receive results, if found stop all
		select {
		case lookupResult = <-parallelLookup.resultsChan:
			resultsFound = lookupResult.success
			transformsInBuffer -= 1
		default:

		}

		if resultsFound {
			break
		}
	}

	for i := 0; i < workers; i++ {
		cubeTransformChan <- ""
	}
	for i := 0; i < workers; i++ {
		<-cubeTransformWorkerStop
	}
	parallelLookup.StopForcefully()
	return cube.RemoveRotationTransforms(lookupResult.data.(string) + lookupResult.solution), true
}
