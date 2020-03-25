package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"bufio"
	"strings"
	"unicode"
)

func fileToString() string {
	reader := bufio.NewReader(os.Stdin)
   	fmt.Print("Enter the name of the file you want to read: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)
	file, err := os.Open(fileName)
	if err != nil {
        	panic(err)
   	}	
	bytes, err := ioutil.ReadAll(file)
	fileString := string(bytes)
	return fileString
}

func wordFreq (fileString string) map[string]int {
	wordCount := make(map[string]int)
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	slicedString := strings.FieldsFunc(fileString, f)
	for _, v := range slicedString {
		wordCount[v]++
	}
	return wordCount
}



func main() {
	fileString := fileToString()
	fmt.Println(wordFreq(fileString))
}