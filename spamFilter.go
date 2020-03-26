package main

import (
	//"fmt"
	//"io/ioutil"
	//"os"
	//"bufio"
	//"strings"
	//"unicode"
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

func wordFreq (fileString string) map[string]int {
	wordCount := make(map[string]int)
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	slicedString := strings.FieldsFunc(fileString, f)
	for _, v := range slicedString {
		v = strings.ToLower(v)
		wordCount[v]++
	}
	return wordCount
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
		var flGoodCount float64 = float64(goodCount)
		var flBadCount float64 = float64(badCount)
		var flNGoodMail float64 = float64(nGoodMail)
		var flNBadMail float64 = float64(nBadMail)
		flGoodCount = flGoodCount * 2
		if flGoodCount + flBadCount < 5 {
			delete(probMap, word)
			continue
		} else if inGoodMap == false && inBadMap == true {
			probMap[word] = 0.99
		} else if inGoodMap == true && inBadMap == false {
			probMap[word] = 0.01
		} else {
			probMap[word] = float64((flBadCount / flNBadMail) / ((flGoodCount / flNGoodMail) + (flBadCount / flNBadMail)))
		}
	}
	return probMap
}

func main() {
	//fileString := fileToString()
	//fmt.Println(wordFreq(fileString))
}