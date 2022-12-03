package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/shrt/internal/datastore"
	"argc.in/shrt/internal/model"
	"argc.in/shrt/internal/respond"
)

func deleteCommand() *cobra.Command {
	c := &cobra.Command{
		Use:          "delete SLUG",
		Short:        "Deletes the Golink",
		RunE:         delete,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	}

	c.Flags().BoolVar(&localOp, "local", false, "Run operation against local database")

	return c
}

func delete(c *cobra.Command, args []string) error {
	r := &model.Route{Slug: args[0]}

	if localOp {
		return deleteLocal(r)
	}

	return deleteOnServer(r)
}

func deleteOnServer(r *model.Route) error {
	ep, err := url.JoinPath(serverAddr, "api", "url", r.Slug)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, ep, nil)
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

	return nil
}

func deleteLocal(r *model.Route) error {
	store, err := datastore.NewSQLiteStore(databasePath)
	if err != nil {
		return err
	}

	if err := store.Delete(context.Background(), r); err != nil {
		return err
	}

	return nil
}
