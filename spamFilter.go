package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
	"unicode"
)

func buildCorpus(dirName string) (string, int) {
	corpusFiles, err := ioutil.ReadDir(dirName) // corpusFiles is a slice of file descriptions, of type FileInfo
	if err != nil {
		panic(err)
	}
	var dirLength int64
	for _, fileInfo := range corpusFiles {
		dirLength += fileInfo.Size()
	}
	catCorpus := make([]byte, dirLength)
	nCopiedBytes := 0
	for _, fileInfo := range corpusFiles {
		filePath := fmt.Sprintf("%s/%s", dirName, fileInfo.Name())
		email, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		bytes, err := ioutil.ReadAll(email)
		if err != nil {
			panic(err)
		}
		nCopiedBytes += copy(catCorpus[nCopiedBytes:], bytes)
		err = email.Close()
		if err != nil {
			panic(err)
		}
	}
	corpus := string(catCorpus)
	nMail := len(corpusFiles)
	return corpus, nMail
}

func fileToString(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(file)
	fileString := string(bytes)
	return fileString
}

func wordCount(fileString string) map[string]int {
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
	for word := range goodMap {
		probMap[word] = 0.0
	}
	for word := range badMap {
		probMap[word] = 0.0
	}
	for word := range probMap {
		goodCount, inGoodMap := goodMap[word]
		badCount, inBadMap := badMap[word]
		var flGoodCount, flBadCount, flNGoodMail, flNBadMail = float64(goodCount), float64(badCount), float64(nGoodMail), float64(nBadMail)
		flGoodCount = flGoodCount * 2
		if flGoodCount+flBadCount < 5 {
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

type wordProb struct {
	token       string
	probability float64
	interest    float64
}

type ByInterest []wordProb

func (a ByInterest) Len() int           { return len(a) }
func (a ByInterest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInterest) Less(i, j int) bool { return a[i].interest < a[j].interest }

func isSpam(mailString string, probMap map[string]float64) bool {

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	mailSlice := strings.FieldsFunc(mailString, f) //split string into slice
	newMailMap := make(map[string]float64)
	for _, mailWord := range mailSlice {
		mailWord = strings.ToLower(mailWord)
		prob, hasKnownProb := probMap[mailWord]
		if hasKnownProb == false { //fill map with words from mail and their probabilities
			newMailMap[mailWord] = 0.4
		} else {
			newMailMap[mailWord] = prob
		}
	}
	mapLength := 0
	for range newMailMap {
		mapLength++
	}
	wordProbs := make([]wordProb, mapLength) //make a slice and fill with words from mail and their probs
	i := 0
	for word, prob := range newMailMap {
		wordProbs[i].token = word
		wordProbs[i].probability = prob
		wordProbs[i].interest = math.Abs(0.5 - prob)
		i++
	}
	sort.Sort(ByInterest(wordProbs))
	//fmt.Println("the least to most interesting words are: ", wordProbs)
	var mostInteresting []wordProb
	if len(wordProbs) > 15 {
		mostInteresting = wordProbs[len(wordProbs)-15:]
	} else {
		mostInteresting = wordProbs
	}
	//fmt.Println("the most interesting words are: ", mostInteresting)
	var probProd, invProbProd float64 = 1.0, 1.0
	for _, v := range mostInteresting {
		probProd *= v.probability
		invProb := 1 - v.probability
		invProbProd *= invProb
	}
	return probProd/(probProd+invProbProd) >= 0.9
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
