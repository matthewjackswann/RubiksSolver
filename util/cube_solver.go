package util

import (
	"database/sql"
	"fmt"
	"github.com/davidminor/uint128"
	"github.com/matthewjackswann/rubiks/cube"
)

func loadSolution(id uint128.Uint128, preparedStmt *sql.Stmt) (string, bool) {
	//result, err := dbConnection.db.Query(fmt.Sprintf("SELECT solution FROM cubes WHERE cube_id_l = %d AND cube_id_h = %d;", int64(id.L), int64(id.H)))
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
	return dbConnection.lookupCube(cubeId, rotation, stmt)
}

func (dbConnection *DBConnection) lookupCube(cubeId uint128.Uint128, rotation string, stmt *sql.Stmt) (string, bool) {
	idSolution, success := loadSolution(cubeId, stmt)
	if !success {
		return "", false
	}
	return cube.RotateTransform(rotation, idSolution), true
}
