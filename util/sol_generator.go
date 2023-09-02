package util

import (
	"bufio"
	"fmt"
	"github.com/davidminor/uint128"
	"github.com/matthewjackswann/rubiks/cube"
	"golang.org/x/sys/unix"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

type cubeResult struct {
	id        uint128.Uint128
	transform uint64
}

func StartSolutionGenerator(db DBConnection, init []int, i, maximumDepth int) {
	// setup generator
	generator := cube.CreateNewGenerator(init, i, "cube/graph.csv") // todo include and increment i

	stop := make(chan struct{})
	// func for receiving signal to start stopping the generator
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		fmt.Println("Stopping...")
		stop <- struct{}{}
	}()

	//batchSize := 1000000
	batchSize := 1000

	transformsSent := 0
	cubeTransforms := make(chan string, batchSize)
	idsReceived := 0
	cubeIds := make(chan cubeResult, batchSize)

	wg := new(sync.WaitGroup)
	workerStopChannel := make(chan interface{})
	wg.Add(1)
	go cubeWorker(cubeTransforms, cubeIds, workerStopChannel, wg)

	fmt.Println("Made cube workers")

	dbSaveChan := make(chan batchResults)
	dbSaveChanResult := make(chan bool, 1)
	dbSaveChanResult <- true // skips over first save as successful
	wg.Add(1)
	go saveWorker(db, dbSaveChan, dbSaveChanResult, wg)

	generatingCubes := true

	currentDepth := generator.GetCurrentDepth()

	fmt.Println(currentDepth)

	for generatingCubes { // while generating or ids haven't been processed yet

		for i := 0; i < batchSize && currentDepth == generator.GetCurrentDepth(); i++ {
			cubeTransforms <- generator.Next()
			transformsSent += 1
		}
		successfulSave := <-dbSaveChanResult
		if !successfulSave {
			fmt.Print("Error saving last batch. Quitting\n\n")
			break
		}

		resultMap := make(map[uint128.Uint128]uint64, batchSize)

		for transformsSent != idsReceived {
			cr := <-cubeIds
			idsReceived += 1
			_, keyExists := resultMap[cr.id]
			if !keyExists {
				resultMap[cr.id] = cr.transform
			}
		}

		dbSaveChan <- batchResults{
			results:       resultMap,
			transformNo:   generator.GetCurrentTransformNum(),
			lastTransform: generator.TransformStack,
		}

		currentDepth = generator.GetCurrentDepth()
		if currentDepth > maximumDepth {
			generatingCubes = false
		}

		select {
		case <-stop:
			generatingCubes = false
		default:
		}
	}

	fmt.Println("Stopping worker")
	workerStopChannel <- new(interface{})
	fmt.Println("Stopping db goroutine and closing db connection")
	dbSaveChan <- batchResults{results: nil}

	wg.Wait()
	fmt.Println("Stopped")
}

var TransformToInt = map[rune]uint64{ // not transform is represented by a 0
	'F': 1,
	'f': 2,
	'L': 3,
	'l': 4,
	'U': 5,
	'u': 6,
	'B': 7,
	'b': 8,
	'R': 9,
	'r': 10,
	'D': 11,
	'd': 12,
}

var IntToTransform = map[uint64]rune{
	1:  'F',
	2:  'f',
	3:  'L',
	4:  'l',
	5:  'U',
	6:  'u',
	7:  'B',
	8:  'b',
	9:  'R',
	10: 'r',
	11: 'D',
	12: 'd',
}

func cubeWorker(generatorResult <-chan string, resultChan chan<- cubeResult, stop <-chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stop:
			return
		case generatorResult := <-generatorResult:
			c := cube.NewSolvedCube()

			c.Transform(generatorResult)

			id, rotationTransform := c.EncodeCube()

			// encode the reverse of the transform
			transform := cube.RotateTransform(cube.ReverseTransform(rotationTransform), cube.ReverseTransform(generatorResult))
			encodedTransform := uint64(0)
			for i := range transform {
				encodedTransform = encodedTransform << 4
				c := rune(transform[len(transform)-1-i])
				encodedTransform += TransformToInt[c]
			}

			resultChan <- cubeResult{
				id:        id,
				transform: encodedTransform,
			}
		}
	}
}

type batchResults struct {
	results       map[uint128.Uint128]uint64
	transformNo   int
	lastTransform []int // it's a waste of time encoding this value, just use "," separated list
}

// closes db connection when stopping
func saveWorker(db DBConnection, dbSaveChan chan batchResults, dbSaveChanResult chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	lastStackSize := 0

	for {
		toSave := <-dbSaveChan

		if toSave.results == nil {
			db.Close()
			return
		}

		fmt.Print("\033[1A\x1b[2K")

		b := make([]string, len(toSave.lastTransform))
		for i, stackElement := range toSave.lastTransform {
			b[i] = strconv.Itoa(stackElement)
		}
		encodedStack := strings.Join(b, ",")

		fmt.Printf("Saving: %d, Stack: %s Estimated Layer:%.2f%%\n", toSave.transformNo, encodedStack, estimateStackPercentage(toSave.lastTransform))

		if len(toSave.lastTransform) != lastStackSize {
			lastStackSize = len(toSave.lastTransform)
			fmt.Print("\033[1A\x1b[2K")
			fmt.Printf("Stack %d: %f%%\n\n", lastStackSize, getDiskUsePercentage())
		}

		success := db.Save(toSave.results, toSave.transformNo, encodedStack)
		//if getDiskUsePercentage() > 99.9 {
		if getDiskUsePercentage() > 0.1 {
			fmt.Println("Disk low on space")
			dbSaveChanResult <- false
		} else {
			dbSaveChanResult <- success
		}
	}
}

func getDiskUsePercentage() float64 {
	var stat unix.Statfs_t
	err := unix.Statfs("/media/swanny/Lexar", &stat)
	if err != nil {
		fmt.Printf("Error getting disk info: %s\n\n", err)
		return -1
	} else {
		return 100 * (1 - float64(stat.Bavail)/float64(stat.Blocks))
	}
}

func estimateStackPercentage(stack []int) float64 {
	if len(stack) == 0 {
		return 0
	}

	res := float64(stack[0]) * 0.5

	for i := 1; i < len(stack); i++ {
		res += 0.5 * float64(stack[i]) / math.Pow(11, float64(i))
	}

	return res * 100
}
