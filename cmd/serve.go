package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"argc.in/shrt/internal/api"
	"argc.in/shrt/internal/datastore"
	"argc.in/shrt/internal/web"
)

func serveCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Run the server",
		RunE:  serve,
	}

	return c
}

func serve(c *cobra.Command, _ []string) error {
	store, err := datastore.NewSQLiteStore(databasePath)
	if err != nil {
		return err
	}

	r := mux.NewRouter().StrictSlash(true)
	api.RegisterRoutes(r, store)
	web.RegisterRoutes(r, store)

	s := http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("starting server %s\n", serverAddr)

		if err := s.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("stopping server")
				return
			}

			log.Println("E: server stopped: ", err)
		}
	}()

	<-ctx.Done()

	if err := s.Shutdown(context.Background()); err != nil {
		log.Println("E: stopping server: ", err)
	}

	return nil
}
