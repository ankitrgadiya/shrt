package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"

	"argc.in/shrt/internal/datastore"
	"argc.in/shrt/internal/model"
	"argc.in/shrt/internal/respond"
)

func RegisterRoutes(r *mux.Router, store datastore.RouteStore) {
	i := &impl{store: store}
	r.Path("/api/urls/").Methods(http.MethodGet).HandlerFunc(i.APIGetURLs)
	r.PathPrefix("/api/url/").Methods(http.MethodGet).HandlerFunc(i.APIGetURL)
	r.PathPrefix("/api/url/").Methods(http.MethodPost).HandlerFunc(i.APIPostURL)
	r.PathPrefix("/api/url/").Methods(http.MethodDelete).HandlerFunc(i.APIDeleteURL)
}

type impl struct {
	store datastore.RouteStore
}

func (i *impl) APIGetURLs(w http.ResponseWriter, r *http.Request) {
	routes, err := i.store.QueryAll(r.Context())
	if err != nil {
		respond.WithError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, routes)
}

func (i *impl) APIGetURL(w http.ResponseWriter, r *http.Request) {
	route := &model.Route{Slug: getSlug(r)}

	err := i.store.Query(r.Context(), route)
	if err != nil {
		respond.WithError(w, r, http.StatusInternalServerError, err)
		return
	}

	if datastore.IsErrNotFound(err) {
		respond.WithStatus(w, r, http.StatusNotFound)
		return
	}

	respond.With(w, r, http.StatusOK, route)
}

func (i *impl) APIPostURL(w http.ResponseWriter, r *http.Request) {
	var route model.Route

	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		respond.WithError(w, r, http.StatusBadRequest, err)
		return
	}

	if _, err := url.Parse(route.URL); err != nil {
		respond.WithError(w, r, http.StatusBadRequest, err)
		return
	}

	route.Slug = getSlug(r)

	if err := i.store.Save(r.Context(), &route); err != nil {
		respond.WithError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusCreated, &route)
}

func (i *impl) APIDeleteURL(w http.ResponseWriter, r *http.Request) {
	route := &model.Route{Slug: getSlug(r)}

	if err := i.store.Delete(r.Context(), route); err != nil {
		respond.WithError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.WithStatus(w, r, http.StatusAccepted)
}

func getSlug(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/api/url/")
}
