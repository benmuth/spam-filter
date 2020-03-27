package main

import (
	//"fmt"
	//"io/ioutil"
	//"os"
	//"bufio"
	//"strings"
	//"unicode"
	//"math"
	//"sort"
)
/*
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

func wordCount (fileString string) map[string]int {
	wordCountMap := make(map[string]int)
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	slicedString := strings.FieldsFunc(fileString, f)
	for _, word := range slicedString {
		word = strings.ToLower(word)
		wordCountMap[word]++
	}
	return wordCountMap
}

*/
func probTable(goodMap map[string]int, badMap map[string]int, nGoodMail int, nBadMail int) map[string]float64 {
	probMap := make(map[string]float64)
	for word, _ := range goodMap {
		probMap[word] = 0.0
	}	
	for word, _ := range badMap {
		probMap[word] = 0.0
	}
	for word, _ := range probMap {
		goodCount, inGoodMap := goodMap[word]
		badCount, inBadMap := badMap[word]
		var flGoodCount, flBadCount, flNGoodMail, flNBadMail = float64(goodCount), float64(badCount), float64(nGoodMail), float64(nBadMail)
		flGoodCount = flGoodCount * 2
		if flGoodCount + flBadCount < 5 {
			delete(probMap, word)
			continue
		} else if inGoodMap == false && inBadMap == true {
			probMap[word] = 0.99
		} else if inGoodMap == true && inBadMap == false {
			probMap[word] = 0.01
		} else {
			probMap[word] = float64((flBadCount / flNBadMail) / ((flGoodCount / flNGoodMail) +(flBadCount / flNBadMail)))
		}
	}
	return probMap
}

func quickSort(wordProbSlice []wordProb, start int, end int) {
	if start >= end {
		pivot := wordProbSlice[end]
		pIndex := start
		temp := wordProb{}
		for i := start; i < end - 1; i++  {
			if wordProbSlice[i].interest <= pivot.interest {
				temp = pivot
				pivot = wordProbSlice[i]
				wordProbSlice[i] = temp
				pIndex++
			}
		}
		temp = pivot
		pivot = wordProbSlice[pIndex]
		wordProbSlice[pIndex] = temp	
		quickSort(wordProbSlice, start, pIndex - 1)
		quickSort(wordProbSlice, pIndex + 1, end)
	}
	
}

type wordProb struct {
	token string
	probability float64
	interest float64
}


func isSpam(newMail *os.File, probMap map[string]float64) bool {
	bytes, err := ioutil.ReadAll(newMail)				//Read mail file into string
	mailString := string(bytes)
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	mailSlice := strings.FieldsFunc(mailString, f)	//split string into slice
	newMailMap := make(map[string]float64)			
	for _, mailWord := range mailSlice {
		mailWord = strings.ToLower(mailWord)
		prob, isKnownProb := probMap[mailWord]			
		if isKnownProb == false {			//fill map with words from mail and their probabilities
			newMailMap[mailWord] = 0.4
		} else {
			newMailMap[mailWord] = prob
		}
	}
	mapLength := 0
	for _, _ = range newMailMap {
		mapLength++
	}
	wordProbSlice := make([]wordProb, mapLength)		//make a slice and fill with words from mail and their probs
	i := 0
	for word, prob := range newMailMap {
		wordProbSlice[i].token = word
		wordProbSlice[i].probability = prob
		wordProbSlice[i].interest = math.Abs(0.5 - prob)
		i++
	}
	quickSort(wordProbSlice, 0, len(wordProbSlice) - 1)
	
	
}

func main() {
	//fileString := fileToString()
	//fmt.Println(wordFreq(fileString))
}