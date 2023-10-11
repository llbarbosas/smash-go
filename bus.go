package smash

import (
	"container/list"
	"context"
	"encoding/gob"
	"log"
	"net"
	"time"

	"github.com/gobwas/glob"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Message struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Source     string    `json:"source"`
	OccorredAt time.Time `json:"occurred_at"`
	Payload    any       `json:"payload"`
}

type HandleFunc func(context.Context, Message)

type Handler struct {
	Matcher glob.Glob
	Handle  HandleFunc
}

type RemoteBus struct {
	Addr string
}

func (b *RemoteBus) Emit(ctx context.Context, opts EmitOptions) ([]string, error) {
	conn, err := net.Dial("tcp", b.Addr)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if err := gob.NewEncoder(conn).Encode(opts.Message); err != nil {
		return nil, err
	}

	var ids []string

	if err := gob.NewDecoder(conn).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (b *RemoteBus) RegisterHandler(matcher string, handle HandleFunc) (string, error) {
	return "", nil
}

type Bus interface {
	Emit(context.Context, EmitOptions) ([]string, error)
	RegisterHandler(matcher string, handle HandleFunc) (string, error)
}

type LocalBus struct {
	addr        string
	handlers    map[string]*Handler
	linkedBuses map[string]Bus
}

func NewLocalBus(addr string) (*LocalBus, error) {
	return &LocalBus{
		addr:        addr,
		handlers:    map[string]*Handler{},
		linkedBuses: map[string]Bus{},
	}, nil
}

func (b *LocalBus) Link(bus Bus, propagate bool) (string, error) {
	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	if propagate {
		_, err = bus.Emit(context.Background(), EmitOptions{
			Message: Message{
				Type:    "__link",
				Payload: []string{b.addr},
			},
		})

		if err != nil {
			return "", err
		}
	}

	b.linkedBuses[id] = bus

	return id, nil
}

func (b *LocalBus) Serve() error {
	ln, _ := net.Listen("tcp", b.addr)

	for {
		func() {
			conn, err := ln.Accept()

			if err != nil {
				log.Println(err)
				return
			}

			defer conn.Close()

			var msg Message

			if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
				log.Println(err)
				return
			}

			if msg.Type == "__link" {
				ids := msg.Payload.([]string)
				addr := ids[0]

				id, err := b.Link(&RemoteBus{
					Addr: addr,
				}, false)

				if err != nil {
					log.Println(err)
					return
				}

				if err := gob.NewEncoder(conn).Encode([]string{id}); err != nil {
					log.Println(err)
					return
				}
			} else {
				ids, err := b.Emit(context.Background(), EmitOptions{
					LocalOnly: true,
					Message:   msg,
				})

				if err != nil {
					log.Println(err)
					return
				}

				if err := gob.NewEncoder(conn).Encode(ids); err != nil {
					log.Println(err)
					return
				}
			}

		}()
	}
}

func (b *LocalBus) RegisterHandler(matcher string, handle HandleFunc) (string, error) {
	g, err := glob.Compile(matcher, ':')

	if err != nil {
		return "", err
	}

	handler := Handler{
		Handle:  handle,
		Matcher: g,
	}

	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	b.handlers[id] = &handler

	return id, nil
}

type EmitOptions struct {
	Message
	LocalOnly bool
}

func (b *LocalBus) Emit(ctx context.Context, opts EmitOptions) ([]string, error) {
	handlerIdsList := list.New()

	for id, handler := range b.handlers {
		if handler.Matcher.Match(opts.Type) {
			handler.Handle(ctx, opts.Message)

			handlerIdsList.PushBack(id)
		}
	}

	if !opts.LocalOnly {
		opts.LocalOnly = true

		for _, bus := range b.linkedBuses {
			ids, err := bus.Emit(ctx, opts)

			if err != nil {
				return nil, err
			}

			for _, id := range ids {
				handlerIdsList.PushBack(id)
			}
		}
	}

	handlerIds := make([]string, 0, handlerIdsList.Len())

	for e := handlerIdsList.Front(); e != nil; e = e.Next() {
		handlerIds = append(handlerIds, e.Value.(string))
	}

	return handlerIds, nil
}
