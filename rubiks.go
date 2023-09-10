package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/matthewjackswann/rubiks/cube"
	"github.com/matthewjackswann/rubiks/util"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CubeData struct {
	CubeLayout     [54]int
	Transformation string
}

func main() {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverPort := serverFlags.Int("port", 3000, "Port the server will be hosted on")

	//generateFlags := flag.NewFlagSet("generate", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments\nExpected 'server' or 'generate' subcommand")
		return
	}

	switch os.Args[1] {

	case "server":
		startServer(*serverPort)

	case "generate":
		db := util.CreateDBConnection("/media/swanny/Lexar/rubiks.db")
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

func startServer(port int) {
	fmt.Printf("Starting Server at localhost:%d \nUse ^C to stop\n", port)
	http.Handle("/", http.FileServer(http.Dir("./frontEnd/build")))
	http.HandleFunc("/cube", fulfillCubeTransformRequest)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Couldn't start server")
		fmt.Println(err.Error())
		return
	}
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

// generator stuff

func startGenerator(db util.DBConnection, init []int, i, maximumDepth int) {
	util.StartSolutionGenerator(db, init, i, maximumDepth)
}
