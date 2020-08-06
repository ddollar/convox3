package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/convox/console/api/resolver"
)

type Handler struct {
	api *Api
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.api.schema.Exec(h.context(r), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		h.queryGet(w, r)
// 	case http.MethodPost:
// 		h.queryPost(w, r)
// 	default:
// 		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
// 	}
// }

func (h *Handler) context(r *http.Request) context.Context {
	ctx := r.Context()

	ctx = context.WithValue(ctx, resolver.ContextModel, h.api.model)

	if parts := strings.Fields(r.Header.Get("Authorization")); len(parts) == 2 && parts[0] == "Bearer" {
		ctx = context.WithValue(ctx, resolver.ContextToken, parts[1])
	}

	return ctx
}

func (h *Handler) query(w http.ResponseWriter, r *http.Request, q Query) error {
	res := h.api.schema.Exec(h.context(r), q.Query, "", q.Variables)

	// if len(res.Errors) > 0 {
	// 	w.WriteHeader(403)
	// }

	data, err := json.Marshal(res)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func (h *Handler) queryGet(w http.ResponseWriter, r *http.Request) {
	qv := r.URL.Query()

	q := Query{
		Query: qv.Get("query"),
	}

	if vars := qv.Get("variables"); vars != "" {
		if err := json.Unmarshal([]byte(vars), &q.Variables); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	if err := h.query(w, r, q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *Handler) queryPost(w http.ResponseWriter, r *http.Request) {
	var q Query

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := json.Unmarshal(data, &q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.query(w, r, q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
