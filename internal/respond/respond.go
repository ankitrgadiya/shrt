// Package respond is inspired by matryer/respond. The API is very similar but
// the package is customized for usage in this project.
package respond

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"argc.in/shrt/internal/model"
)

func WithStatus(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
	resp := new(Msg)

	if code >= 200 && code < 300 {
		resp.Ok = true
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		WithError(w, r, http.StatusInternalServerError, err)
	}
}

func With(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)

	var resp *Msg

	switch v := data.(type) {
	case *model.Route:
		resp = &Msg{Route: v}
	case []model.Route:
		resp = &Msg{Routes: v}
	}

	if code >= 200 && code < 300 {
		resp.Ok = true
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		WithError(w, r, http.StatusInternalServerError, err)
	}
}

func WithError(w http.ResponseWriter, _ *http.Request, code int, err error) {
	log.Printf("E: %+v", err)
	w.WriteHeader(code)

	resp := &Msg{Ok: false, Error: err.Error()}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("E: marshalling error response: %+v", err)
	}
}

func WithRedirect(w http.ResponseWriter, r *http.Request, dest string) {
	http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
}

func WithTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data interface{}) {
	tmpl.Execute(w, data)
}
