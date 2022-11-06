package discordle

import (
	"testing"
)

func TestLoadWordBank(t *testing.T) {
	expectedCandidateWordCount := 2315
	expectedValidGuessCount := 10657

	b, _ := LoadWordBank("../word_bank/candidate_words.txt", "../word_bank/valid_guesses.txt")
	candidateWordCount := len(b.candidateWords)
	validGuessCount := len(b.validGuesses)

	if candidateWordCount != expectedCandidateWordCount {
		t.Errorf(
			"Expected candidate word list to contain %d words, but it actually contains %d.",
			expectedCandidateWordCount, candidateWordCount,
		)
	}

	if validGuessCount != expectedValidGuessCount {
		t.Errorf(
			"Expected valid guess list to contain %d words, but it actually contains %d.",
			expectedValidGuessCount, validGuessCount,
		)
	}
}

func TestWordBank_GetWord(t *testing.T) {
	b, _ := LoadWordBank("../word_bank/candidate_words.txt", "../word_bank/valid_guesses.txt")
	word := b.GetWord()
	wordLength := len(word)
	if wordLength != 5 {
		t.Errorf("Word is not 5 letters as expected! %s, len %d.", word, wordLength)
	}
}

func TestWordBank_ValidateGuess(t *testing.T) {
	testWord := "TWEET"

	b, _ := LoadWordBank("../word_bank/candidate_words.txt", "../word_bank/valid_guesses.txt")
	if !b.ValidateGuess(testWord) {
		t.Errorf("Expected word %s not recognized as valid.", testWord)
	}
}
