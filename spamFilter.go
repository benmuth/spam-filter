package main

import (
	"fmt"
	"io/ioutil"
	"os"
	//"bufio"
	"strings"
	"unicode"
	"math"
	"sort"
)

func buildCorpus(dirName string) (string, int) {
	fileSlice, err := ioutil.ReadDir(dirName)
	if err != nil {
		panic(err)
	}
	nMail := len(fileSlice)
	var dirLength int64
	for _, file := range fileSlice {
		dirLength += file.Size()
	}
	copiedLength := 0
	corpusSlice := make([]byte, dirLength)
	for _, file := range fileSlice {
		filePath := fmt.Sprintf("%s/%s", dirName, file.Name())
		pFile, err := os.Open(filePath)
		if err != nil {
        		panic(err)
   		}	
		bytes, err := ioutil.ReadAll(pFile)
		if err != nil {
			panic(err)
		}
		copiedLength += copy(corpusSlice[copiedLength:], bytes)
		err = pFile.Close()
		if err != nil {
			panic(err)
		}
	}
	corpus := string(corpusSlice)
	return corpus, nMail
}


func fileToString(fileName string) string {
	//reader := bufio.NewReader(os.Stdin)
   	//fmt.Print("Enter the name of the file you want to read: ")
	//fileName, _ := reader.ReadString('\n')
	//fileName = strings.TrimSpace(fileName)
	
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

func probCalc(goodMap map[string]int, badMap map[string]int, nGoodMail int, nBadMail int) map[string]float64 {
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

/*
func quickSort(wordProbSlice []wordProb, start int, end int) {
	j := 0
	var pivot wordProb{}
	var pIndex int
	if start < end {
		j++
		pivot = wordProbSlice[end]
		pIndex = start
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
	} else {
		fmt.Println("sorted ", j, " times!")
	}
}
*/

type wordProb struct {
	token string
	probability float64
	interest float64
}

type ByInterest []wordProb

func (a ByInterest) Len() int           { return len(a) }
func (a ByInterest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInterest) Less(i, j int) bool { return a[i].interest < a[j].interest }



func isSpam(mailString string, probMap map[string]float64) bool {
	
	//bytes, err := ioutil.ReadAll(newMail)			
	//if err != nil {
	//	panic(err)
	//}
	//mailString := string(bytes)
	
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
	//quickSort(wordProbSlice, 0, len(wordProbSlice) - 1)
	sort.Sort(ByInterest(wordProbSlice))
	//fmt.Println("the least to most interesting words are: ", wordProbSlice)
	var mostInteresting []wordProb
	if len(wordProbSlice) > 15 {
		mostInteresting = wordProbSlice[len(wordProbSlice) - 15:]
	} else {
		mostInteresting = wordProbSlice 
	}
	//fmt.Println("the most interesting words are: ", mostInteresting)
	var probProd, invProbProd float64 = 1.0, 1.0
	for _, v := range mostInteresting {
		probProd *= v.probability
		invProb := 1 - v.probability
		invProbProd *= invProb
	}	
	combProb := probProd/(probProd + invProbProd)
	//fmt.Println("the combined probability is: ", combProb)
	var mailIsSpam bool
	if combProb >= 0.9 {
		mailIsSpam = true
	} else {
		mailIsSpam = false
	}
	return mailIsSpam
}


func main() {
	goodCorpus, nGoodMail := buildCorpus("/Users/moose1/Documents/SpamFilter/lingspam_public/NotSpam")
	badCorpus, nBadMail := buildCorpus("/Users/moose1/Documents/SpamFilter/lingspam_public/Spam")
	goodMap := wordCount(goodCorpus)
	badMap := wordCount(badCorpus)
	probMap := probCalc(goodMap, badMap, nGoodMail, nBadMail)
	spamMailString := fileToString("/Users/moose1/Documents/SpamFilter/lingspam_public/Spam/spmsga1.txt")
	notSpamMailString := fileToString("/Users/moose1/Documents/SpamFilter/lingspam_public/NotSpam/3-1msg1.txt")
	isSpamSpam := isSpam(spamMailString, probMap)
	isNotSpamSpam := isSpam(notSpamMailString, probMap)
	fmt.Println("Is spam spam? ", isSpamSpam)
	fmt.Println("Is not spam spam? ", isNotSpamSpam)
}