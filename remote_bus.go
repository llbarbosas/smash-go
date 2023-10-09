package smash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type RemoteBus struct {
	remoteAddr string
	handlers   map[string]Handler
	client     *http.Client
}

func NewRemoteBus(remoteAddr string) (*RemoteBus, error) {
	return &RemoteBus{
		remoteAddr: remoteAddr,
		handlers:   map[string]Handler{},
		client:     &http.Client{},
	}, nil
}

func (b *RemoteBus) Emit(ctx context.Context, optsSetters ...func(*EmitOptions)) error {
	opts, err := buildOpts(optsSetters...)

	if err != nil {
		return err
	}

	if !opts.LocalOnly {
		return fmt.Errorf("remote buses cannot emit remote messages")
	}

	req := BusEmitRequest{
		Message: opts.Message,
	}

	reqJson, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := http.Post(b.remoteAddr, "application/json", bytes.NewBuffer(reqJson))

	if err != nil || res.StatusCode != 200 {
		return fmt.Errorf("emit error")
	}

	for _, handler := range b.handlers {
		if handler.Matcher.Match(opts.Topic) {
			handler.Handle(ctx, opts.Message)
		}
	}

	return nil
}

func (b *RemoteBus) RegisterHandler(handler *Handler) (string, error) {
	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	b.handlers[id] = *handler

	return id, nil
}
