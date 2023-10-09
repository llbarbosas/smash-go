package smash

import (
	"context"
	"errors"
	"time"

	"github.com/gobwas/glob"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	errBusNotLinked = errors.New("bus not linked")
)

type Bus interface {
	Emit(ctx context.Context, optsSetters ...func(*EmitOptions)) error
	RegisterHandler(handler *Handler) (string, error)
}

type LocalBus struct {
	topics      []string
	handlers    map[string]Handler
	remoteBuses map[string]Bus
}

func NewLocalBus() (*LocalBus, error) {
	return &LocalBus{
		topics:      []string{},
		handlers:    map[string]Handler{},
		remoteBuses: map[string]Bus{},
	}, nil
}

func (b *LocalBus) RegisterTopics(topics ...string) {
	b.topics = append(b.topics, topics...)
}

func (b *LocalBus) RegisterHandler(handler *Handler) (string, error) {
	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	b.handlers[id] = *handler

	return id, nil
}

func (b *LocalBus) Link(remoteBus Bus) (string, error) {
	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	handler, err := NewHandler("**", func(ctx context.Context, m Message) {
		b.Emit(ctx, WithDefaults(m.Topic, m.Payload), WithSource(m.Source), LocalOnly)
	})

	if err != nil {
		return "", err
	}

	_, err = remoteBus.RegisterHandler(handler)

	if err != nil {
		return "", err
	}

	b.remoteBuses[id] = remoteBus

	return id, nil
}

func (b *LocalBus) Unlink(id string) error {
	if _, ok := b.remoteBuses[id]; !ok {
		return errBusNotLinked
	}

	delete(b.remoteBuses, id)

	return nil
}

type EmitOptions struct {
	Message
	LocalOnly bool
}

func WithDefaults(topic string, payload any) func(*EmitOptions) {
	return func(opts *EmitOptions) {
		opts.Topic = topic
		opts.Payload = payload
	}
}

func WithSource(source string) func(*EmitOptions) {
	return func(opts *EmitOptions) {
		opts.Source = source
	}
}

func WithMessage(message Message) func(*EmitOptions) {
	return func(opts *EmitOptions) {
		opts.Message = message
	}
}

func LocalOnly(opts *EmitOptions) {
	opts.LocalOnly = true
}

func buildOpts(optsSetters ...func(*EmitOptions)) (*EmitOptions, error) {
	opts := &EmitOptions{}

	for _, setter := range optsSetters {
		setter(opts)
	}

	if opts.Message.ID == "" {
		id, err := gonanoid.New()

		if err != nil {
			return nil, err
		}

		opts.Message.ID = id
	}

	return opts, nil
}

func (b *LocalBus) Emit(ctx context.Context, optsSetters ...func(*EmitOptions)) error {
	opts, err := buildOpts(optsSetters...)

	if err != nil {
		return err
	}

	for _, handler := range b.handlers {
		if handler.Matcher.Match(opts.Topic) {
			handler.Handle(ctx, opts.Message)
		}
	}

	if !opts.LocalOnly {
		for _, bus := range b.remoteBuses {
			if err := bus.Emit(ctx, append(optsSetters, LocalOnly)...); err != nil {
				return err
			}
		}
	}

	return nil
}

type HandleFunc func(context.Context, Message)

type Handler struct {
	Handle  HandleFunc
	Matcher glob.Glob
}

func NewHandler(matcher string, handle HandleFunc) (*Handler, error) {
	g, err := glob.Compile(matcher, ':')

	if err != nil {
		return nil, err
	}

	return &Handler{
		Handle:  handle,
		Matcher: g,
	}, nil
}

type Message struct {
	ID         string    `json:"id"`
	Topic      string    `json:"topic"`
	Source     string    `json:"source"`
	OccorredAt time.Time `json:"occurred_at"`
	Payload    any       `json:"payload"`
}
