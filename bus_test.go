package smash_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/llbarbosas/smash-go"
)

func TestBus(t *testing.T) {
	log.SetFlags(log.Llongfile)

	bus1, err := smash.NewLocalBus("127.0.0.1:8000")

	if err != nil {
		log.Fatal(err)
	}

	bus2, err := smash.NewLocalBus("127.0.0.1:8001")

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := bus1.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := bus2.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Second * 2)

	_, err = bus1.Link(&smash.RemoteBus{
		Addr: "127.0.0.1:8001",
	}, true)

	if err != nil {
		log.Fatal(err)
	}

	_, err = bus1.RegisterHandler("component1:message", func(ctx context.Context, m smash.Message) {
		log.Println("This is component 1", m.Payload, time.Since(m.OccorredAt))
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = bus2.RegisterHandler("component2:message", func(ctx context.Context, m smash.Message) {
		log.Println("This is component 2", m.Payload, time.Since(m.OccorredAt))
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = bus1.Emit(context.Background(), smash.EmitOptions{
		Message: smash.Message{
			Type:    "component1:message",
			Payload: "hey!",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = bus1.Emit(context.Background(), smash.EmitOptions{
		Message: smash.Message{
			Type:    "component2:message",
			Payload: "hey!",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = bus2.Emit(context.Background(), smash.EmitOptions{
		Message: smash.Message{
			Type:    "component1:message",
			Payload: "hey!",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = bus2.Emit(context.Background(), smash.EmitOptions{
		Message: smash.Message{
			Type:    "component2:message",
			Payload: "hey!",
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
