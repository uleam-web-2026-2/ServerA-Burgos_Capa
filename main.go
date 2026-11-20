package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/config"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/aplicacion"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/repositorio"
	consultassqlc "github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/repositorio/sqlc"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
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
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("no se pudo acceder a la conexión SQL: ", err)
	}

	if cfg.EjecutarMigraciones {
		if err := ejecutarMigraciones(db); err != nil {
			log.Fatal(err)
		}
	}
	if cfg.MigracionSolo {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("no se pudo cerrar la conexión de migración: ", err)
		}
		log.Println("migración y siembra completadas")
		return
	}

	r := chi.NewRouter()
	r.Use(configurarCORS())
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.TiempoEspera))

	tokens := sesion.NuevosTokens(cfg.ClaveFirma)
	login := &sesion.Manejador{DB: db, Tokens: tokens}
	r.With(limitadorLogin()).Post("/login", login.Login) // pública: limitada a 5 intentos por IP y minuto
	repositorioRutinas := repositorio.NuevoRutinasSqlc(
		&repositorio.Rutinas{DB: db},
		consultassqlc.New(sqlDB),
	)
	servicioRutinas := aplicacion.NuevoServicioRutinas(repositorioRutinas)
	(&gymcoach.Manejador{DB: db, ServicioRutinas: servicioRutinas}).Rutas(r, tokens)

	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		Handler:      r,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
	}
	log.Println("escuchando en el puerto", cfg.Puerto)
	log.Fatal(servidor.ListenAndServe())
}

func configurarCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

func limitadorLogin() func(http.Handler) http.Handler {
	return httprate.LimitBy(5, time.Minute, porIP)
}

func porIP(r *http.Request) (string, error) {
	ip := middleware.GetClientIP(r.Context())
	return httprate.CanonicalizeIP(ip), nil
}

func ejecutarMigraciones(db *gorm.DB) error {
	if err := db.AutoMigrate(&gymcoach.Entrenador{}); err != nil {
		return fmt.Errorf("no se pudo migrar la tabla de entrenadores: %w", err)
	}
	registrosActualizados, err := gymcoach.MigrarDatosAntiguos(db)
	if err != nil {
		return fmt.Errorf("no se pudieron adaptar los datos existentes: %w", err)
	}
	if registrosActualizados > 0 {
		log.Printf("se adaptaron %d valores de clientes y rutinas existentes; revise los datos temporales .invalid", registrosActualizados)
	}
	if err := db.AutoMigrate(&gymcoach.Cliente{}); err != nil {
		return fmt.Errorf("no se pudo migrar la tabla de clientes: %w", err)
	}
	if err := db.AutoMigrate(&gymcoach.Rutina{}); err != nil {
		return fmt.Errorf("no se pudo migrar la tabla de rutinas: %w", err)
	}
	if err := gymcoach.Sembrar(db); err != nil {
		return fmt.Errorf("no se pudieron cargar los datos de ejemplo: %w", err)
	}
	if err := db.AutoMigrate(&sesion.Usuario{}); err != nil {
		return fmt.Errorf("no se pudo migrar la tabla de usuarios: %w", err)
	}
	if err := sesion.Sembrar(db); err != nil {
		return fmt.Errorf("no se pudieron sembrar los usuarios: %w", err)
	}
	return nil
}
