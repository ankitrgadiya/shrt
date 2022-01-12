package web

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"argc.in/shrt/internal/datastore"
	"argc.in/shrt/internal/model"
	"argc.in/shrt/internal/respond"
)

func RegisterRoutes(r *mux.Router, store datastore.RouteStore) {
	i := &impl{
		store: store,
		fs:    http.FileServer(http.FS(_staticFS)),
	}

	r.Path("/links/").Methods(http.MethodGet).HandlerFunc(i.AllLinks)
	r.PathPrefix("/s/").Methods(http.MethodGet).Handler(http.StripPrefix("/s/", http.HandlerFunc(i.Static)))
	r.PathPrefix("/edit/").Methods(http.MethodGet).HandlerFunc(i.Edit)
	r.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(i.Redirect)
}

type impl struct {
	store datastore.RouteStore
	fs    http.Handler
}

func (i *impl) Static(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = filepath.Join("static", r.URL.Path)
	r.URL.RawPath = filepath.Join("static", r.URL.RawPath)

	i.fs.ServeHTTP(w, r)
}

func (i *impl) Redirect(w http.ResponseWriter, r *http.Request) {
	route := &model.Route{Slug: getSlug(r)}
	if len(route.Slug) == 0 {
		respond.WithRedirect(w, r, "/edit/")
	}

	err := i.store.Query(r.Context(), route)
	if err != nil && !datastore.IsErrNotFound(err) {
		respond.WithError(w, r, http.StatusInternalServerError, err)
		return
	}

	if datastore.IsErrNotFound(err) {
		respond.WithRedirect(w, r, "/edit/"+route.Slug)
		return
	}

	respond.WithRedirect(w, r, route.URL)
}

func (i *impl) Edit(w http.ResponseWriter, r *http.Request) {
	f, err := _staticFS.Open(filepath.Join("static", "edit.html"))
	if err != nil {
		respond.WithError(w, r, http.StatusInternalServerError, err)
	}

	defer f.Close()

	io.Copy(w, f)
}

func (i *impl) AllLinks(w http.ResponseWriter, r *http.Request) {
	routes, err := i.store.QueryAll(r.Context())
	if err != nil {
		respond.WithError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.WithTemplate(w, r, _linksTmpl, routes)
}

func getSlug(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/")
}
