package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

	c.Flags().StringVar(&accessToken, "access", "", "Cloudflare Access Token (optional)")
	c.Flags().StringVar(&serverAddr, "server", "https://argv.in", "Address for the web server (optional)")
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
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/url/%s", serverAddr, r.Slug), nil)
	if err != nil {
		return err
	}

	if len(accessToken) != 0 {
		req.Header.Set(_cloudflareAccessHeader, accessToken)
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
