package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/handler"
	rp "github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/repository"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/service"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/config"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/logger"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/middleware"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/repository"
)

func main() {
	// initial logger
	lg := logger.New(os.Stdout)

	// Initial Config
	cfg := config.MustLoadConfig()

	//  Initial Database
	db, err := repository.InitDB(cfg.DB)
	if err != nil {
		lg.Fatal("failed to connect to database: %v", err)
	}

	// Initial repository, service, handler
	repo := rp.NewPostgresProductRepository(db)
	svc := service.NewProductService(repo)
	hdl := handler.NewProductHandler(svc, lg)

	// init middlerware
	mw := middleware.New(lg)

	// router
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			hdl.Create(w, r)
		case http.MethodGet:
			hdl.GetAll(w, r)
		default:
			lg.Warn("invalid method: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/product/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			hdl.Update(w, r)
		case http.MethodDelete:
			hdl.Delete(w, r)
		case http.MethodGet:
			hdl.GetById(w, r)
		default:
			lg.Warn("invalid method: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Server
	server := &http.Server{
		Addr: ":8080",
		Handler: mw.Recovery(
			mw.Logging(mux),
		),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			lg.Error("failed to start server: %v", err)
		}
	}()

	// shutdown
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGTERM)

	<-sg

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		lg.Fatal("failed to shutdown server: %v", err)
	}

	func(db *sql.DB) {
		if err := db.Close(); err != nil {
			lg.Fatal("failed to close database: %v", err)
		}
	}(db)

	lg.Info("Server exiting")
}
