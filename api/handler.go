package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	api *Api
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.queryGet(w, r)
	case http.MethodPost:
		h.queryPost(w, r)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) context(r *http.Request) context.Context {
	fmt.Printf("r.Header: %+v\n", r.Header)

	ctx := r.Context()

	ctx = context.WithValue(ctx, "model", h.api.model)
	ctx = context.WithValue(ctx, "uid", "f8abd4df-f8b4-4cb9-9514-6395e7907f2b")

	return ctx
}

func (h *Handler) query(w http.ResponseWriter, r *http.Request, q Query) error {
	res := h.api.schema.Exec(h.context(r), q.Query, "", q.Variables)

	if len(res.Errors) > 0 {
		w.WriteHeader(403)
	}

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
