package main

import (
	"context"
	"os"
	"time"

	"github.com/llbarbosas/smash-tictactoe"
)

var (
	Name = "ttt.input_manager"
)

func Register(bus smash.Bus, scheduler *smash.Scheduler) error {
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	//
	// if err != nil {
	// 	panic(err)
	// }

	scheduler.RunEvery(time.Millisecond, func() {
		b := make([]byte, 1)
		os.Stdin.Read(b)

		key := string(b)

		bus.Emit(
			context.Background(),
			smash.WithDefaults("input:read", key),
			smash.WithSource(Name),
		)

		// if key == "q" {
		// 	term.Restore(int(os.Stdin.Fd()), oldState)
		// 	os.Exit(0)
		// }
	})

	return nil
}
