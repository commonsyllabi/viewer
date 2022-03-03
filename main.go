package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFile() {
	file, err := os.Open("manifest.imscc")
	defer file.Close()

	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func main() {
	fmt.Println("cc viewer 0.1")

	readFile()
}
