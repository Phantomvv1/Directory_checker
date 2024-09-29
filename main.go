package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	directoryName := "test"
	content, fileNames := fileOperations(directoryName)
	fmt.Printf("Contents: %v\nFilenames: %v", content, fileNames)
	for i, name := range fileNames {
		updateFile(content[i], name, i)
	}
}

func fileOperations(directoryName string) ([]string, []string) {
	entries, err := os.ReadDir(directoryName)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}

	_, err = os.Stat("comparisonFile.txt")
	if os.IsNotExist(err) {
		_, err = os.Create("comparisonFile.txt")
		if err != nil {
			log.Fatal("Error creating comp file: ", err)
		}
	}

	var fileNames []string
	var fileContents []string

	for i, fileName := range entries {
		if !fileName.IsDir() {
			file, err := os.Open("./" + directoryName + "/" + entries[i].Name())
			if err != nil {
				log.Fatal("Error opening file:", err)
			}
			defer file.Close()
			if err != nil {
				log.Fatal("Error with reading the file", err)
			}
			content, err := compareContentsOfFiles(entries[i].Name(), directoryName)
			if err != nil {
				fmt.Println(err)
			}
			fileNames = append(fileNames, fileName.Name())
			fileContents = append(fileContents, string(content))
		} else {
			subDirContents, subDirFileNames := fileOperations("./" + directoryName + "/" + fileName.Name() + "/")
			fileContents = append(fileContents, subDirContents...)
			fileNames = append(fileNames, subDirFileNames...)
		}
	}
	return fileContents, fileNames
}

func compareContentsOfFiles(fileName string, directoryName string) ([]byte, error) {
	fileFound, fileDifference := false, false

	file, err := os.Open("comparisonFile.txt")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	compFile, err := os.Open("./" + directoryName + "/" + fileName)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer compFile.Close()

	scanner := bufio.NewScanner(file)
	compFileScanner := bufio.NewScanner(compFile)

	for scanner.Scan() {
		line := scanner.Text()
		if line == fileName+":" {
			fileFound = true
			continue
		}
		if fileFound {
			compFileScanner.Scan()
			if line != compFileScanner.Text() {
				fileDifference = true
				break
			}
		}
	}

	errorMessage := "Contents of the file " + fileName + " don't match"

	content, err := io.ReadAll(compFile)
	if err != nil {
		log.Fatal(err)
	}

	if fileDifference {
		return content, errors.New(errorMessage)
	}

	return content, nil
}

func updateFile(content string, fileName string, append int) {
	if append > 0 {
		file, err := os.OpenFile("comparisonFile.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(fileName + ":\n" + content + "\n")
	} else {
		file, err := os.OpenFile("comparisonFile.txt", os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(fileName + ":\n" + content + "\n")
	}
}
