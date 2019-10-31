package cmd

import (
	"context"
	"fmt"
	"github.com/gorilla/mux" // need to use dep for package management
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewCmdRunClient() *cobra.Command  {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run client where probes will be executed",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running... client")
			return runClient()
		} ,
	}
	return cmd
}
func runClient() error {
	router := mux.NewRouter()
	router.HandleFunc("/", httpGETHandler).Methods("GET")
	router.HandleFunc("/success", httpGETHandler).Methods("GET")
	router.HandleFunc("/fail", httpGETHandler).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
	return nil
}

func httpGETHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("============== Received request")
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/success":
		fmt.Println("Request in path: /success")
		w.WriteHeader(http.StatusOK)
	case "/fail":
		fmt.Println("Request in path: /fail")
		w.WriteHeader(http.StatusForbidden)
	}
}