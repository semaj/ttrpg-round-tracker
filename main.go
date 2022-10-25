package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chzyer/readline"

	"github.com/semaj/sotdl-tracker/tracker"
)

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       "/tmp/readline.tmp",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()
	log.SetOutput(l.Stderr())

	state := tracker.NewState()
	dat, err := os.ReadFile("game.json")
	if err == nil {
		var s tracker.State
		if err := json.Unmarshal(dat, &s); err != nil {
			panic(err)
		}
		fmt.Println("read state from game.json")
		state = &s
	}

	for {
		fmt.Println(state.String())
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				fmt.Println("breaking")
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			fmt.Println("breaking")
			break
		}
		// text = text[:len(text)-1]
		messages, err := state.Input(line)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, message := range messages {
			fmt.Println(message)
		}
	}
	raw, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("game.json", raw, 0644); err != nil {
		panic(err)
	}
}
