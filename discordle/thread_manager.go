package discordle

import "fmt"

type ThreadManager struct {
	gameThreads map[string]*Wordle
	WordBank    *WordBank
	GamesWon    int
}

func (t *ThreadManager) GameForThread(threadId string) (game *Wordle, err error) {
	if t.gameThreads == nil {
		t.gameThreads = map[string]*Wordle{}
	}

	game = t.gameThreads[threadId]
	if game != nil {
		if game.IsOver() {
			return nil, fmt.Errorf("this game is over")
		} else {
			return game, nil
		}
	}

	game = &Wordle{
		CorrectWord: t.WordBank.GetWord(),
		Guesses:     []string{},
	}
	t.gameThreads[threadId] = game
	return game, nil
}
