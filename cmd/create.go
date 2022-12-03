package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/shrt/internal/datastore"
	"argc.in/shrt/internal/model"
	"argc.in/shrt/internal/respond"
)

func createCommand() *cobra.Command {
	c := &cobra.Command{
		Use:          "create SLUG URL",
		Short:        "Creates a new Golink",
		RunE:         create,
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
	}

	c.Flags().BoolVar(&localOp, "local", false, "Run operation against local database")

	return c
}

func create(c *cobra.Command, args []string) error {
	r := &model.Route{Slug: args[0], URL: args[1]}

	if localOp {
		return createLocal(c.OutOrStdout(), r)
	}

	return createOnServer(c.OutOrStdout(), r)
}

func createOnServer(w io.Writer, r *model.Route) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		return err
	}

	ep, err := url.JoinPath(serverAddr, "api", "url", r.Slug)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ep, &buf)
	if err != nil {
		return err
	}

	if len(clientID) != 0 {
		req.Header.Set(_headerClientID, clientID)
		req.Header.Set(_headerClientSecret, clientSecret)
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var m respond.Msg

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return err
	}

	if !m.Ok {
		return errors.Errorf("something went wrong: %v", m)
	}

	_, err = fmt.Fprintf(w, "%s/%s\n", serverAddr, m.Route.Slug)
	return err
}

func createLocal(w io.Writer, r *model.Route) error {
	store, err := datastore.NewSQLiteStore(databasePath)
	if err != nil {
		return err
	}

	if err := store.Save(context.Background(), r); err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s/%s\n", serverAddr, r.Slug)
	return err
}
