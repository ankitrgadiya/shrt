package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/shrt/internal/datastore"
	"argc.in/shrt/internal/model"
	"argc.in/shrt/internal/respond"
)

func listCommand() *cobra.Command {
	c := &cobra.Command{
		Use:          "list",
		Short:        "Lists all Golinks",
		RunE:         list,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	c.Flags().BoolVar(&localOp, "local", false, "Run operation against local database")

	return c
}

func list(c *cobra.Command, args []string) error {
	if localOp {
		return listLocal(c.OutOrStdout())
	}

	return listOnServer(c.OutOrStdout())
}

func listOnServer(w io.Writer) error {
	ep, err := url.JoinPath(serverAddr, "api", "urls")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, ep, nil)
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

	return displayRoutes(w, m.Routes)
}

func listLocal(w io.Writer) error {
	store, err := datastore.NewSQLiteStore(databasePath)
	if err != nil {
		return err
	}

	routes, err := store.QueryAll(context.Background())
	if err != nil {
		return err
	}

	return displayRoutes(w, routes)
}

func displayRoutes(w io.Writer, routes []model.Route) error {
	tw := tabwriter.NewWriter(w, 3, 3, 3, ' ', 0)

	for _, r := range routes {
		fmt.Fprintf(tw, "%s\t%s\n", r.Slug, r.URL)
	}

	return tw.Flush()
}
