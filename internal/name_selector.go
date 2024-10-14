package server

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GetRandomName() string {
	file, err := os.Open("names.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	if len(names) == 0 {
		fmt.Println("No names found in the file.")
	} // Generate a random number between 1 and 200

	rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNumber := rand.Intn(200)
	randomName := names[randomNumber]

	return randomName
}
