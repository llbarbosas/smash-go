package smash

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ManagmentAPI struct {
	handler http.HandlerFunc
}

func NewManagmentAPI(moduleManager *ModuleManager, bus *LocalBus) (*ManagmentAPI, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/modules/load", func(w http.ResponseWriter, r *http.Request) {
		var req ModuleLoadRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, map[string]interface{}{"error": "invalid-input-error", "message": err.Error()})
			return
		}

		module, err := moduleManager.Load(req.Path)

		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, map[string]interface{}{"error": "module-load-error", "message": err.Error()})
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, ModuleLoadResponse{
			Name: module.Name,
		})
	})

	r.Post("/modules/register", func(w http.ResponseWriter, r *http.Request) {
		err := moduleManager.RegisterModules()

		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, map[string]interface{}{"error": "modules-register-error", "message": err.Error()})
			return
		}

		render.Status(r, 200)
	})

	r.Get("/bus/emit", func(w http.ResponseWriter, r *http.Request) {
		var req BusEmitRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, map[string]interface{}{"error": "invalid-input-error", "message": err.Error()})
			return
		}

		ctx := context.Background()

		err := bus.Emit(ctx, WithMessage(req.Message), LocalOnly)

		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, map[string]interface{}{"error": "bus-emit-error", "message": err.Error()})
			return
		}

		render.Status(r, 200)
	})

	return &ManagmentAPI{
		handler: r.ServeHTTP,
	}, nil
}

type ModuleLoadRequest struct {
	Path string `json:"path"`
}

type ModuleLoadResponse struct {
	Name string `json:"name"`
}

type BusEmitRequest struct {
	Message Message `json:"message"`
}

type ManagmentAPIClient struct {
	client *http.Client
}

func NewManagmentAPIClient(addr string) (*ManagmentAPIClient, error) {
	return &ManagmentAPIClient{
		client: &http.Client{},
	}, nil
}

func (c ManagmentAPIClient) LoadModule(req ModuleLoadRequest) (*ModuleLoadResponse, error) {
	reqJson, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := http.Post("http://127.0.0.1:3000/modules/load", "application/json", bytes.NewBuffer(reqJson))

	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var resJson ModuleLoadResponse

	if json.Unmarshal(b, &resJson); err != nil {
		return nil, err
	}

	return &resJson, nil
}

func (c ManagmentAPIClient) RegisterModules() error {
	res, err := http.Post("http://127.0.0.1:3000/modules/register", "application/json", new(bytes.Buffer))

	if err != nil {
		return err
	}

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var resJson ModuleLoadResponse

	if json.Unmarshal(b, &resJson); err != nil {
		return err
	}

	return nil
}
