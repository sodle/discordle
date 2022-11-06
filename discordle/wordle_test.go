package discordle

import "testing"

func TestWordle_CheckGuess(t *testing.T) {
	type fields struct {
		correctWord string
		Guesses     []string
	}
	type args struct {
		guess string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantIsCorrect bool
		wantScore     string
		wantErr       bool
	}{
		{
			"correct word",
			fields{"GLUES", []string{}},
			args{"GLUES"},
			true,
			"GGGGG",
			false,
		},
		{
			"partly correct",
			fields{"GLUES", []string{}},
			args{"GOALS"},
			false,
			"GBBYG",
			false,
		},
		{
			"partly correct, double yellow letters",
			fields{"GLUES", []string{}},
			args{"GOLLY"},
			false,
			"GBYBB",
			false,
		},
		{
			"doubled a letter that needn't be doubled",
			fields{"PHOTO", []string{}},
			args{"POPPY"},
			false,
			"GYBBB",
			false,
		},
		{
			"wrong word length",
			fields{"GLUES", []string{}},
			args{"GRISTLE"},
			false,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Wordle{
				CorrectWord: tt.fields.correctWord,
				Guesses:     tt.fields.Guesses,
			}
			gotIsCorrect, gotScore, err := w.CheckGuess(tt.args.guess)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckGuess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsCorrect != tt.wantIsCorrect {
				t.Errorf("CheckGuess() gotIsCorrect = %v, want %v", gotIsCorrect, tt.wantIsCorrect)
			}
			if gotScore != tt.wantScore {
				t.Errorf("CheckGuess() gotScore = %v, want %v", gotScore, tt.wantScore)
			}
		})
	}
}

func TestWordle_InputGuess(t *testing.T) {
	type fields struct {
		correctWord string
		Guesses     []string
	}
	type args struct {
		guess string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantIsCorrect bool
		wantScore     string
		wantErr       bool
	}{
		{
			"first guess, incorrect",
			fields{"GLUES", []string{}},
			args{"GOALS"},
			false,
			"GBBYG",
			false,
		},
		{
			"guess after win",
			fields{"GLUES", []string{"GLUES"}},
			args{"GOALS"},
			false,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Wordle{
				CorrectWord: tt.fields.correctWord,
				Guesses:     tt.fields.Guesses,
			}
			gotIsCorrect, gotScore, err := w.InputGuess(tt.args.guess)
			if (err != nil) != tt.wantErr {
				t.Errorf("InputGuess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsCorrect != tt.wantIsCorrect {
				t.Errorf("InputGuess() gotIsCorrect = %v, want %v", gotIsCorrect, tt.wantIsCorrect)
			}
			if gotScore != tt.wantScore {
				t.Errorf("InputGuess() gotScore = %v, want %v", gotScore, tt.wantScore)
			}
		})
	}
}

func TestWordle_IsOver(t *testing.T) {
	type fields struct {
		correctWord string
		Guesses     []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"not over",
			fields{"GUESS", []string{}},
			false,
		},
		{
			"won",
			fields{"GUESS", []string{"GUESS"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Wordle{
				CorrectWord: tt.fields.correctWord,
				Guesses:     tt.fields.Guesses,
			}
			if got := w.IsOver(); got != tt.want {
				t.Errorf("IsOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
