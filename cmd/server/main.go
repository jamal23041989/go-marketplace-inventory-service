package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/handler"
	rp "github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/repository"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/service"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/config"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/repository"
)

func main() {
	// Initial Config
	cfg := config.MustLoadConfig()

	//  Initial Database
	db, err := repository.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	// Initial repository, service, handler
	repo := rp.NewInMemoryProductRepository()
	svc := service.NewProductService(repo)
	hdl := handler.NewProductHandler(svc)

	// router
	mux := http.NewServeMux()

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			hdl.Create(w, r)
		case http.MethodGet:
			hdl.GetAll(w, r)
		default:
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Server
	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// shutdown
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGTERM)

	<-sg

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)
	log.Println("Server exiting")
}
