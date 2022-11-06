package main

import (
	"discordle/discordle"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func scoreToDiscord(score string) (discordScore string) {
	for _, r := range score {
		switch r {
		case 'G':
			discordScore += ":green_square:"
			break
		case 'Y':
			discordScore += ":yellow_square:"
			break
		case 'B':
			discordScore += ":black_square_button:"
			break
		}
	}
	return
}

//func (app *App) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
//	if m.Author.ID == s.State.User.ID {
//		return
//	}
//
//	if m.Type == discordgo.MessageTypeReply {
//		if m.ReferencedMessage.Author.ID != s.State.User.ID {
//			// not a reply to the bot
//			return
//		}
//	} else if m.Type == discordgo.MessageTypeThreadStarterMessage {
//		return
//	} else {
//		mentionedBot := false
//		for _, mention := range m.Mentions {
//			if mention.ID == s.State.User.ID {
//				mentionedBot = true
//				break
//			}
//		}
//		if !mentionedBot {
//			// doesn't mention the bot
//			return
//		}
//	}
//
//	game := app.gameForUser(m.Author.ID)
//
//	guess := nonLetterRegex.ReplaceAllString(m.Content, "")
//	correct, score, err := game.InputGuess(strings.ToUpper(guess))
//	if err != nil {
//		msg := fmt.Sprintf("error: %s", err)
//		_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
//		if err != nil {
//			panic(err)
//		}
//	} else if correct {
//		msg := fmt.Sprintf("%d/6 %s :tada:", len(game.Guesses)+1, scoreToDiscord(score))
//		_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
//		if err != nil {
//			panic(err)
//		}
//	} else if len(game.Guesses) == 6 {
//		msg := fmt.Sprintf("6/6 %s :x:", scoreToDiscord(score))
//		_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
//		if err != nil {
//			panic(err)
//		}
//	} else {
//		msg := fmt.Sprintf("%d/6 %s", len(game.Guesses), scoreToDiscord(score))
//		_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
//		if err != nil {
//			panic(err)
//		}
//	}
//}

func formatWin(w *discordle.Wordle, guess string) string {
	guessCount := len(w.Guesses)
	return fmt.Sprintf("%s\n%d %s :tada:", guess, guessCount, scoreToDiscord("GGGGG"))
}

func formatScore(w *discordle.Wordle, guess string, score string) string {
	guessCount := len(w.Guesses)
	return fmt.Sprintf("%s\n%d %s", guess, guessCount, scoreToDiscord(score))
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, t *discordle.ThreadManager) {
	if m.Message.Author.ID == s.State.User.ID {
		return
	}

	ch, _ := s.State.Channel(m.Message.ChannelID)
	if !ch.IsThread() {
		return
	}

	if ch.OwnerID == s.State.User.ID {
		game, err := t.GameForThread(ch.ID)
		if err != nil {
			fmt.Println("error getting game for thread,", err)
			return
		}
		guess := strings.ToUpper(strings.TrimSpace(m.Content))
		isCorrect, score, err := game.InputGuess(guess)
		if err != nil {
			_, err := s.ChannelMessageSend(ch.ID, fmt.Sprintf("Error: %e", err))
			if err != nil {
				fmt.Println("error sending response,", err)
			}
		} else if isCorrect {
			t.GamesWon++
			_, err := s.ChannelMessageSend(ch.ID, formatWin(game, guess))
			if err != nil {
				fmt.Println("error sending response,", err)
			}
		} else {
			_, err := s.ChannelMessageSend(ch.ID, formatScore(game, guess, score))
			if err != nil {
				fmt.Println("error sending response,", err)
			}
		}
	}
}

func handleStartGame(s *discordgo.Session, m *discordgo.InteractionCreate, t *discordle.ThreadManager) {
	err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Game started! Please place guesses in this thread.",
		},
	})
	if err != nil {
		fmt.Println("error handling interaction,", err)
		return
	}

	response, err := s.InteractionResponse(m.Interaction)
	if err != nil {
		fmt.Println("error handling retrieving interaction response,", err)
		return
	}

	threadStart, err := s.MessageThreadStart(response.ChannelID, response.ID, "Guess Thread", 7*24*60)
	if err != nil {
		fmt.Println("error starting thread,", err)
		return
	}

	_, err = t.GameForThread(threadStart.ID)
	if err != nil {
		fmt.Println("error starting game,", err)
		return
	}
}

func main() {
	discordToken := os.Getenv("DISCORD_TOKEN")

	// Create bot
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating client,", err)
		return
	}

	// Print to console after connection success
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", discord.State.User.Username, discord.State.User.Discriminator)
	})

	// Open connection to Discord
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	wordBank, err := discordle.LoadWordBank("word_bank/candidate_words.txt", "word_bank/valid_guesses.txt")
	if err != nil {
		fmt.Println("error loading word bank,", err)
		return
	}

	threader := &discordle.ThreadManager{
		WordBank: wordBank,
	}

	// Implement custom slash command
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "wordle" {
			handleStartGame(s, i, threader)
		}
	})

	// Implement message receive handler
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handleMessage(s, m, threader)
	})

	// Register custom slash command
	command, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "wordle",
		Description: "Start a collaborative Wordle inside a thread.",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Remove slash command
	err = discord.ApplicationCommandDelete(discord.State.User.ID, "", command.ID)
	if err != nil {
		fmt.Println("error removing command", err)
	}

	// Cleanly close down the Discord session.
	err = discord.Close()
	if err != nil {
		fmt.Println("error closing connection", err)
		return
	}
}
