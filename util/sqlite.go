package util

import (
	"database/sql"
	"fmt"
	"github.com/davidminor/uint128"
	"github.com/matthewjackswann/rubiks/cube"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DBConnection struct {
	db        *sql.DB
	connected bool
}

func Create(path string) DBConnection {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	dbConnection := DBConnection{
		db:        db,
		connected: true,
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `cubes` (" +
		"`cube_id_l` BINARY(8) NOT NULL, " +
		"`cube_id_h` BINARY(8) NOT NULL, " +
		"`solution` BINARY(8), " +
		"PRIMARY KEY (cube_id_l, cube_id_h));")
	if err != nil {
		fmt.Println("Error creating cubes table")
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `next_transform` (" +
		"`id` INTEGER NOT NULL PRIMARY KEY, " +
		"`transform_no` INTEGER NOT NULL, " +
		"`stack` TEXT NOT NULL);")
	if err != nil {
		fmt.Println("Error creating next_transform table")
		panic(err)
	}

	return dbConnection
}

func (dbConnection *DBConnection) Close() {
	if !dbConnection.connected {
		panic("close called on disconnected DBConnection")
	}
	err := dbConnection.db.Close()
	if err != nil {
		panic(err)
	}
	dbConnection.connected = false
}

func (dbConnection *DBConnection) Save(results map[uint128.Uint128]uint64, transformNo int, stack string) bool {
	transaction, err := dbConnection.db.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}

	stmt, err := transaction.Prepare("INSERT OR IGNORE INTO cubes (cube_id_l, cube_id_h, solution) VALUES (?,?,?);")
	if err != nil {
		fmt.Println(err)
		return false
	}

	for cubeId, transform := range results {
		_, err := stmt.Exec(int64(cubeId.L), int64(cubeId.H), int64(transform))
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	stmt, err = transaction.Prepare("INSERT OR REPLACE INTO next_transform (id, transform_no, stack) VALUES (1, ?, ?);")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(transformNo, stack)
	if err != nil {
		fmt.Println("Couldn't update next transform")
		fmt.Println(err)
		return false
	}

	err = transaction.Commit()
	if err != nil {
		fmt.Println(fmt.Errorf("couldn't commit to database"))
		fmt.Println(err)
		return false
	}
	return true
}

type QueryResult struct {
	NextNum      int
	EncodedStack string
}

func (dbConnection *DBConnection) GetNextTransforms() QueryResult {
	rows, err := dbConnection.db.Query("SELECT transform_no, stack FROM next_transform;")
	if err != nil {
		log.Fatalln(err)
	}
	if rows.Next() {
		var nextNum int
		var stack string
		err := rows.Scan(&nextNum, &stack)
		if err != nil {
			fmt.Println("Error loading next transforms")
			panic(err)
		}
		rowErr := rows.Close()
		if rowErr != nil {
			fmt.Printf("Error closing the DB: %s\n", rowErr)
		}
		return QueryResult{
			NextNum:      nextNum,
			EncodedStack: stack,
		}
	}
	rowErr := rows.Close()
	if rowErr != nil {
		fmt.Printf("Error closing the DB: %s\n", rowErr)
	}
	return QueryResult{
		NextNum:      0,
		EncodedStack: "0",
	}
}

func (dbConnection *DBConnection) loadSolution(id uint128.Uint128) (string, bool) {
	result, err := dbConnection.db.Query(fmt.Sprintf("SELECT solution FROM cubes WHERE cube_id_l = %d AND cube_id_h = %d;", int64(id.L), int64(id.H)))
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

func (dbConnection *DBConnection) LookupCube(cubeId uint128.Uint128, rotation string) (string, bool) {
	if cubeId.Equals(cube.SolvedCubeId) {
		return "", true
	}
	idSolution, success := dbConnection.loadSolution(cubeId)
	if !success {
		return "", false
	}
	return cube.RotateTransform(rotation, idSolution), true
}

type lookupWorkerRequest struct {
	cube *cube.Cube
	data interface{}
}

type lookupWorkerResponse struct {
	cube     *cube.Cube
	success  bool
	solution string
	data     interface{}
}

type ParallelDatabaseLookup struct {
	requestChan chan *lookupWorkerRequest
	resultsChan chan *lookupWorkerResponse
	workerCount int
}

func (p ParallelDatabaseLookup) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.requestChan <- nil
	}
	for i := 0; i < p.workerCount; i++ {
		r := <-p.resultsChan
		if r != nil {
			fmt.Println("When closing lookup workers the results chan wasn't empty")
			i -= 1 // wait an extra iteration
		}
	}
}

// StartLookupWorkers takes requests from the requestChan, looks it up in the db and sends a result to the resultChan
func (dbConnection *DBConnection) StartLookupWorkers(bufferSize, workerCount int) ParallelDatabaseLookup {
	requestChan := make(chan *lookupWorkerRequest, bufferSize)
	resultsChan := make(chan *lookupWorkerResponse, bufferSize)
	for worker := 0; worker < workerCount; worker++ {
		go func() {
			for {
				job := <-requestChan
				if job == nil {
					resultsChan <- nil
					return
				}

				solution, success := dbConnection.LookupCube(job.cube.EncodeCube())
				resultsChan <- &lookupWorkerResponse{
					cube:     job.cube,
					success:  success,
					solution: solution,
					data:     job.data,
				}
			}
		}()
	}
	return ParallelDatabaseLookup{
		requestChan: requestChan,
		resultsChan: resultsChan,
		workerCount: workerCount,
	}
}
