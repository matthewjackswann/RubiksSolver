package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/matthewjackswann/rubiks/cube"
	"github.com/matthewjackswann/rubiks/util"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverPort := serverFlags.Int("port", 3000, "Port the server will be hosted on")
	dbPathServer := serverFlags.String("db", "", "Path to sqlite database")

	generateFlags := flag.NewFlagSet("generate", flag.ExitOnError)
	dbPathGenerator := generateFlags.String("db", "", "Path to sqlite database")

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments\nExpected 'server' or 'generate' subcommand")
		return
	}

	switch os.Args[1] {

	case "server":
		if err := serverFlags.Parse(os.Args[2:]); err != nil {
			fmt.Println("error processing server args")
			return
		}
		if *dbPathServer == "" {
			fmt.Println("Please provide a path to the database to use for cube lookups")
			return
		}
		if _, err := os.Stat(*dbPathServer); errors.Is(err, os.ErrNotExist) {
			fmt.Println("couldn't resolve file at location", *dbPathServer)
			return
		}
		startServer(*serverPort, *dbPathServer)

	case "generate":
		if err := generateFlags.Parse(os.Args[2:]); err != nil {
			fmt.Println("error processing generate args")
			return
		}
		if *dbPathGenerator == "" {
			fmt.Println("Please provide a path to the database to save the generated cubes to")
			return
		}
		db := util.CreateDBConnection(*dbPathGenerator)
		nextInfo := db.GetNextTransforms()
		stackString := strings.Split(nextInfo.EncodedStack, ",")
		initStack := make([]int, len(stackString))
		for i, s := range stackString {
			si, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("Error loading stack from db")
				panic(err)
			}
			initStack[i] = si
		}
		startGenerator(db, initStack, nextInfo.NextNum, 16)

	default:
		fmt.Println("Expected 'server' or 'generate' subcommand")
	}
}

// server stuff

func startServer(port int, dbPath string) {
	fmt.Printf("Starting Server at localhost:%d \nUse ^C to stop\n", port)
	http.Handle("/", http.FileServer(http.Dir("./frontEnd/build")))
	http.HandleFunc("/cube", fulfillCubeTransformRequest)
	http.HandleFunc("/cubeMinimalSol", fulfillCubeMinimalSolveRequest(dbPath))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Couldn't start server")
		fmt.Println(err.Error())
		return
	}
}

type CubeData struct {
	CubeLayout     [54]int
	Transformation string
}

func fulfillCubeTransformRequest(w http.ResponseWriter, r *http.Request) {
	data := new(CubeData)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := cube.NewCube(data.CubeLayout)
	c.Transform(data.Transformation)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(c.Layout)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type CubeDescription struct {
	CubeLayout [54]int
}

type CubeSolution struct {
	Success   bool   `json:"success"`
	Transform string `json:"transform"`
}

func fulfillCubeMinimalSolveRequest(dbPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(CubeDescription)
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		db := util.CreateDBConnection(dbPath)
		c := cube.NewCube(data.CubeLayout)
		solution, success := db.SolveCubeBySearch(c, 6, 10)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(CubeSolution{
			Success:   success,
			Transform: solution,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// generator stuff

func startGenerator(db util.DBConnection, init []int, i, maximumDepth int) {
	util.StartSolutionGenerator(db, init, i, maximumDepth)
}
