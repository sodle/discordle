package discordle

import (
	"fmt"
	"strings"
)

type Wordle struct {
	CorrectWord string
	Guesses     []string
}

func (w *Wordle) CheckGuess(guess string) (isCorrect bool, score string, err error) {
	if guess == w.CorrectWord {
		return true, "GGGGG", nil
	}

	if len(guess) != 5 {
		return false, "", fmt.Errorf("guesses must be 5 characters")
	}

	guessRunes := []rune(guess)
	correctRunes := []rune(w.CorrectWord)

	scoreRunes := []string{"", "", "", "", ""}

	// Count occurrences of each letter
	letterCounts := map[rune]int64{}
	for _, letter := range correctRunes {
		letterCounts[letter]++
	}

	// Find all the green letters
	for idx, guessLetter := range guessRunes {
		if guessLetter == correctRunes[idx] {
			scoreRunes[idx] = "G"
			letterCounts[guessLetter]--
		}
	}

	// Find the yellows/blacks from the remaining letters
	for idx, scoreRune := range scoreRunes {
		if scoreRune == "G" {
			continue
		}
		guessRune := guessRunes[idx]
		if letterCounts[guessRune] > 0 {
			scoreRunes[idx] = "Y"
			// Wordle never paints the same letter yellow twice in the same guess
			// Even if it appears twice in the correct word
			// This prevents leaking the fact that there is a double letter
			letterCounts[guessRune] = 0
		} else {
			scoreRunes[idx] = "B"
		}
	}

	return false, strings.Join(scoreRunes, ""), nil
}

func (w *Wordle) InputGuess(guess string) (isCorrect bool, score string, err error) {
	if w.IsOver() {
		return false, "", fmt.Errorf("game is over")
	}

	isCorrect, score, err = w.CheckGuess(guess)
	if err == nil {
		w.Guesses = append(w.Guesses, guess)
	}
	return isCorrect, score, err
}

func (w *Wordle) IsOver() bool {
	if len(w.Guesses) > 0 {
		if w.Guesses[len(w.Guesses)-1] == w.CorrectWord {
			return true
		}
	}
	return false
}
