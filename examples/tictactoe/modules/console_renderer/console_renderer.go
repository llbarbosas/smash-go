package main

import (
	"context"
	"fmt"

	"github.com/llbarbosas/smash-tictactoe"
)

var (
	Name = "ttt.renderer.console"
)

type State interface {
	Board() []uint8
}

func Register(bus smash.Bus, scheduler *smash.Scheduler) error {
	stateUpdateHandler, err := smash.NewHandler("state:update", func(ctx context.Context, msg smash.Message) {
		state := msg.Payload.(State)
		fmt.Print("\033[H\033[2J")

		for i, place := range state.Board() {
			if i%3 == 0 {
				fmt.Println()
			}

			fmt.Print(place, "\t")
		}
	})

	if err != nil {
		return err
	}

	_, err = bus.RegisterHandler(stateUpdateHandler)

	if err != nil {
		return err
	}

	return nil
}
