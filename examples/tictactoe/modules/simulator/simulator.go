package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"

	"github.com/llbarbosas/smash-go"
	"github.com/llbarbosas/smash-go/examples/tictactoe/modules/state"
)

var (
	Name = "ttt.simulator"
)

func Register(bus smash.Bus, scheduler *smash.Scheduler) error {
	gob.Register(state.State{})

	state := state.State{
		Board: []uint8{
			0, 0, 0,
			0, 0, 0,
			0, 0, 0,
		},
		NextPlayer: 1,
	}

	_, err := bus.RegisterHandler("input:read", func(ctx context.Context, msg smash.Message) {
		key := msg.Payload.(string)

		numberKey, err := strconv.Atoi(key)

		if err != nil {
			return
		}

		state.Update(uint8(numberKey))

		winner := state.HaveWinner()

		if winner != 0 {
			fmt.Println("Winner", winner)
			os.Exit(0)
		}

		bus.Emit(
			context.Background(),
			smash.EmitOptions{
				Message: smash.Message{
					Type:    "state:update",
					Payload: state,
					Source:  Name,
				},
			})
	})

	if err != nil {
		return err
	}

	_, err = bus.Emit(
		context.Background(),
		smash.EmitOptions{
			Message: smash.Message{
				Type:    "state:update",
				Payload: state,
				Source:  Name,
			},
		})

	if err != nil {
		return err
	}

	return nil
}
