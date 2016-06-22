package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func shouldNot(e error) {
	if e != nil {
		panic(e)
	}
}

const MAX_USER_COUNT = 15000000

func main() {
	var inputFileDir string
	flag.StringVar(&inputFileDir, "input", "input", "input file directory")

	flag.Parse()
	fmt.Println("Start with input:", inputFileDir)

	files, err := ioutil.ReadDir(inputFileDir)
	shouldNot(err)

	users := make([]int, MAX_USER_COUNT)
	for _, file := range files {
		count, members := parse(inputFileDir, file.Name())
		for i := 0; i < count; i += 1 {
			users[members[i]] += 1
		}
	}

	counts := make([]int, 100)
	for _, count := range users {
		for i := 0; i < 100; i += 1 {
			if i < count {
				counts[i] += 1
			}
		}
	}
	for i := 0; i < 100; i += 1 {
		fmt.Printf("%d %d\n", i+1, counts[i])
	}
}

func parse(inputDir string, fileName string) (count int, result []int) {
	path := path.Join(inputDir, fileName)

	input, err := os.Open(path)
	shouldNot(err)
	defer input.Close()

	reader := bufio.NewReader(input)
	scanner := bufio.NewScanner(reader)
	result = make([]int, MAX_USER_COUNT)
	count = 0
	i := 0
	for scanner.Scan() {
		shouldNot(scanner.Err())
		id, err := strconv.Atoi(scanner.Text())
		shouldNot(err)
		result[i] = id
		i += 1
		count += 1
	}
	return
}
