package discordle

import (
	"bufio"
	"golang.org/x/exp/slices"
	"log"
	"math/rand"
	"os"
	"strings"
)

type WordBank struct {
	candidateWords []string
	validGuesses   []string
}

func LoadWordBank(candidatePath string, guessPath string) (*WordBank, error) {
	wordList, err := loadWordList(candidatePath)
	if err != nil {
		return nil, err
	}

	guessList, err := loadWordList(guessPath)
	if err != nil {
		return nil, err
	}

	return &WordBank{
		candidateWords: wordList,
		validGuesses:   guessList,
	}, nil
}

func (b *WordBank) GetWord() string {
	wordCount := len(b.candidateWords)
	wordIdx := rand.Intn(wordCount)
	return b.candidateWords[wordIdx]
}

func (b *WordBank) ValidateGuess(guess string) bool {
	return slices.Contains(append(b.validGuesses, b.candidateWords...), guess)
}

func loadWordList(filename string) (wordList []string, err error) {
	file, err := os.Open(filename)
	defer func(wordFile *os.File) {
		err := wordFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w := scanner.Text()
		if !strings.Contains(w, "#") {
			w = strings.TrimSpace(w)
			if len(w) > 0 {
				wordList = append(wordList, strings.ToUpper(w))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wordList, nil
}
