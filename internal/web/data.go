package web

import (
	"embed"
	"text/template"
)

var (
	//go:embed static/*
	_staticFS embed.FS

	//go:embed static/links.html
	_linksHTML string
	_linksTmpl = template.Must(template.New("links.html").Parse(_linksHTML))
)
