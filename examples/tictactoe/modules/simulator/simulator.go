package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/llbarbosas/smash-tictactoe"
)

var (
	Name = "ttt.simulator"
)

type State struct {
	board      []uint8
	nextPlayer uint8
}

func (s State) Board() []uint8 {
	return s.board
}

func (s *State) Update(move uint8) {
	if s.board[move] == 0 {
		s.board[move] = s.nextPlayer
	}

	if s.nextPlayer == 1 {
		s.nextPlayer = 2
	} else {
		s.nextPlayer = 1
	}
}

func (s *State) HaveWinner() uint8 {
	diagonalWin := s.board[4] != 0 && ((s.board[0] == s.board[4] && s.board[4] == s.board[8]) || (s.board[2] == s.board[4] && s.board[4] == s.board[6]))

	if diagonalWin {
		return s.board[4]
	}

	for i := 0; i < 3; i++ {
		if s.board[i*3] != 0 {
			rowWin := s.board[i*3] == s.board[i*3+1] && s.board[i*3+1] == s.board[i*3+2]
			colWin := s.board[i*3] == s.board[(i+1)*3] && s.board[(i+1)*3] == s.board[(i+2)*3]

			if rowWin || colWin {
				return s.board[i*3]
			}
		}
	}

	return 0
}

func Register(bus smash.Bus, scheduler *smash.Scheduler) error {
	state := State{
		board: []uint8{
			0, 0, 0,
			0, 0, 0,
			0, 0, 0,
		},
		nextPlayer: 1,
	}

	inputHandler, err := smash.NewHandler("input:read", func(ctx context.Context, msg smash.Message) {
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
			smash.WithDefaults("state:update", state),
			smash.WithSource(Name),
		)
	})

	if err != nil {
		return err
	}

	_, err = bus.RegisterHandler(inputHandler)

	if err != nil {
		return err
	}

	bus.Emit(
		context.Background(),
		smash.WithDefaults("state:update", state),
		smash.WithSource(Name),
	)

	return nil
}
