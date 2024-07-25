package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		start := 0
		end := 0
		for i, c := range line {
			if c == '[' {
				start = i
			}
			if c == ']' {
				end = i
				break
			}
		}
		commaCount := 0
		for i := start; i < end; i++ {
			if line[i] == ',' {
				commaCount++
			}
		}
		//always on more than amount of commas
		commaCount++
		count += commaCount
	}
	fmt.Println("Found ", count, " untriggered select statements")
}
