package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/config"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach"
)

func main() {
	cfg, err := config.Cargar()
	if err != nil {
		log.Fatal("configuración: ", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar: ", err)
	}

	if err := db.AutoMigrate(&gymcoach.Entrenador{}); err != nil {
		log.Fatal("no se pudo migrar la tabla de entrenadores: ", err)
	}

	registrosActualizados, err := gymcoach.MigrarDatosAntiguos(db)
	if err != nil {
		log.Fatal("no se pudieron adaptar los datos existentes: ", err)
	}
	if registrosActualizados > 0 {
		log.Printf("se adaptaron %d valores de clientes y rutinas existentes; revise los datos temporales .invalid", registrosActualizados)
	}

	if err := db.AutoMigrate(&gymcoach.Cliente{}); err != nil {
		log.Fatal("no se pudo migrar la tabla de clientes: ", err)
	}

	if err := db.AutoMigrate(&gymcoach.Rutina{}); err != nil {
		log.Fatal("no se pudo migrar la tabla de rutinas: ", err)
	}

	if err := gymcoach.Sembrar(db); err != nil {
		log.Fatal("no se pudieron cargar los datos de ejemplo: ", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.TiempoEspera))
	(&gymcoach.Manejador{DB: db}).Rutas(r)

	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		Handler:      r,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
	}
	log.Println("escuchando en el puerto", cfg.Puerto)
	log.Fatal(servidor.ListenAndServe())
}
