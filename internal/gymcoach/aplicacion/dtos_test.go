package aplicacion

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
)

type repositorioDtoPrueba struct{}

func (*repositorioDtoPrueba) ListarRutinas(context.Context, string) ([]dominio.Rutina, error) {
	return nil, nil
}

func (*repositorioDtoPrueba) CrearRutina(context.Context, *dominio.Rutina) error {
	return nil
}

func (*repositorioDtoPrueba) ObtenerRutina(context.Context, uint) (dominio.Rutina, error) {
	return dominio.Rutina{}, nil
}

func (*repositorioDtoPrueba) ObtenerRutinaConCliente(context.Context, uint) (dominio.Rutina, error) {
	return dominio.Rutina{
		BaseModel:   dominio.BaseModel{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now().Add(1 * time.Second)},
		Titulo:      "Fuerza",
		Descripcion: "Entrenamiento de fuerza",
		DiasSemana:  4,
		Estado:      "pendiente",
		ClienteID:   12,
		Cliente: dominio.Cliente{
			BaseModel:    dominio.BaseModel{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now().Add(2 * time.Second)},
			Nombre:       "Ana",
			Correo:       "ana@example.com",
			Objetivo:     "Ganar masa muscular",
			EntrenadorID: 3,
		},
	}, nil
}

func (*repositorioDtoPrueba) ActualizarRutina(context.Context, uint, dominio.Rutina) (dominio.Rutina, error) {
	return dominio.Rutina{}, nil
}

func (*repositorioDtoPrueba) EliminarRutina(context.Context, uint) (int64, error) { return 0, nil }

func (*repositorioDtoPrueba) ActualizarEstado(context.Context, uint, string, string) (int64, error) {
	return 1, nil
}

func TestRutinaSalidaDesdeDominioExponeSoloCamposDeLaAPI(t *testing.T) {
	rutina := dominio.Rutina{
		BaseModel:   dominio.BaseModel{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now().Add(time.Second)},
		Titulo:      "Fuerza",
		Descripcion: "Rutina semanal",
		DiasSemana:  4,
		Estado:      "pendiente",
		ClienteID:   12,
		Cliente: dominio.Cliente{
			BaseModel:    dominio.BaseModel{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now().Add(2 * time.Second)},
			Nombre:       "Ana",
			Correo:       "ana@example.com",
			Objetivo:     "Ganar masa muscular",
			EntrenadorID: 3,
		},
	}

	salida := RutinaSalidaDesdeDominio(rutina)
	serializado, err := json.Marshal(salida)
	if err != nil {
		t.Fatalf("no se pudo serializar la rutina DTO: %v", err)
	}
	jsonStr := string(serializado)
	for _, campo := range []string{"created_at", "updated_at", "deleted_at", "CreatedAt", "UpdatedAt"} {
		if strings.Contains(jsonStr, campo) {
			t.Fatalf("el DTO no debe exponer campos del modelo: %s -> %s", campo, jsonStr)
		}
	}
	if !strings.Contains(jsonStr, "\"cliente_id\":12") || !strings.Contains(jsonStr, "\"estado\":\"pendiente\"") {
		t.Fatalf("la salida debería conservar los campos públicos de la API: %s", jsonStr)
	}
}

func TestClienteListaSalidaIncluyeRelacionesSinCamposInternos(t *testing.T) {
	cliente := dominio.Cliente{
		BaseModel:    dominio.BaseModel{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Nombre:       "Ana",
		Correo:       "ana@example.com",
		Objetivo:     "Ganar masa muscular",
		EntrenadorID: 3,
		Entrenador: dominio.Entrenador{
			BaseModel: dominio.BaseModel{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Nombre:    "Luis",
			Correo:    "luis@example.com",
		},
		Rutinas: []dominio.Rutina{{
			BaseModel:   dominio.BaseModel{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Titulo:      "Fuerza",
			Descripcion: "Rutina semanal",
			DiasSemana:  4,
			Estado:      "pendiente",
			ClienteID:   12,
		}},
	}

	salida := ClienteListaSalidaDesdeDominio(cliente)
	serializado, err := json.Marshal(salida)
	if err != nil {
		t.Fatalf("no se pudo serializar el DTO de lista de clientes: %v", err)
	}
	jsonStr := string(serializado)
	for _, campo := range []string{"created_at", "updated_at", "deleted_at", "CreatedAt", "UpdatedAt"} {
		if strings.Contains(jsonStr, campo) {
			t.Fatalf("el DTO no debe exponer campos del modelo: %s -> %s", campo, jsonStr)
		}
	}
	for _, campo := range []string{`"entrenador":{"id":3`, `"rutinas":[{"id":8`, `"cliente_id":12`} {
		if !strings.Contains(jsonStr, campo) {
			t.Fatalf("el DTO de lista debe incluir %s: %s", campo, jsonStr)
		}
	}
}

func TestServicioObtenerConvierteDominioEnDTO(t *testing.T) {
	servicio := NuevoServicioRutinas(&repositorioDtoPrueba{})
	rutina, err := servicio.Obtener(context.Background(), 8)
	if err != nil {
		t.Fatalf("se esperaba éxito al convertir la rutina: %v", err)
	}
	if rutina.ID != 8 || rutina.Cliente == nil || rutina.Cliente.Nombre != "Ana" {
		t.Fatalf("el DTO no se convirtió correctamente: %+v", rutina)
	}
}
