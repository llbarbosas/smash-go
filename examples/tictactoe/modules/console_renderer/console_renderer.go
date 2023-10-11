package main

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/llbarbosas/smash-go"
	"github.com/llbarbosas/smash-go/examples/tictactoe/modules/state"
)

var (
	Name = "ttt.renderer.console"
)

func Register(bus smash.Bus, scheduler *smash.Scheduler) error {
	gob.Register(state.State{})

	_, err := bus.RegisterHandler("state:update", func(ctx context.Context, msg smash.Message) {
		state := msg.Payload.(state.State)
		fmt.Print("\033[H\033[2J")

		for i, place := range state.Board {
			if i%3 == 0 {
				fmt.Println()
			}

			fmt.Print(RenderPlace(place), "\t")
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func RenderPlace(place uint8) string {
	if place == 1 {
		return "X"
	}

	if place == 2 {
		return "O"
	}

	return "-"
}

func GetState(payload any) state.State {
	state := payload.(state.State)

	return state
}
