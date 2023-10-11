package smash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
)

type ManagmentAPI struct {
	http.HandlerFunc
}

type ManagmentAPIConfig struct {
	ModuleManager *ModuleManager
	Bus           *LocalBus
}

func NewManagmentAPI(cfg ManagmentAPIConfig) (*ManagmentAPI, error) {
	r := chi.NewRouter()

	/*logger :=*/
	httplog.NewLogger("httplog-example", httplog.Options{
		JSON: true,
	})

	// r.Use(httplog.RequestLogger(logger))

	r.Post("/modules/load", func(w http.ResponseWriter, r *http.Request) {
		var req ModuleLoadRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, map[string]interface{}{"error": "invalid-input-error", "message": err.Error()})
			return
		}

		module, err := cfg.ModuleManager.Load(req.Path)

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
		err := cfg.ModuleManager.RegisterModules()

		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, map[string]interface{}{"error": "modules-register-error", "message": err.Error()})
			return
		}

		render.Status(r, 200)
	})

	return &ManagmentAPI{
		HandlerFunc: r.ServeHTTP,
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
	Addr   string
	client *http.Client
}

func NewManagmentAPIClient(addr string) (*ManagmentAPIClient, error) {
	return &ManagmentAPIClient{
		Addr:   addr,
		client: &http.Client{},
	}, nil
}

func (c ManagmentAPIClient) LoadModule(req ModuleLoadRequest) (*ModuleLoadResponse, error) {
	reqJson, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := http.Post(fmt.Sprintf("%s/modules/load", c.Addr), "application/json", bytes.NewBuffer(reqJson))

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
	res, err := http.Post(fmt.Sprintf("%s/modules/register", c.Addr), "application/json", new(bytes.Buffer))

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
