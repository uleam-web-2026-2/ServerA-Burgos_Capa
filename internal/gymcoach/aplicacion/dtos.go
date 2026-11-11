package aplicacion

import "github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"

type ClienteEntrada struct {
	Nombre       string `json:"nombre"`
	Correo       string `json:"correo"`
	Objetivo     string `json:"objetivo"`
	EntrenadorID uint   `json:"entrenador_id"`
}

type ClienteSalida struct {
	ID           uint   `json:"id"`
	Nombre       string `json:"nombre"`
	Correo       string `json:"correo"`
	Objetivo     string `json:"objetivo"`
	EntrenadorID uint   `json:"entrenador_id"`
}

type EntrenadorSalida struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Correo string `json:"correo"`
}

type RutinaDeClienteSalida struct {
	ID          uint   `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	DiasSemana  int    `json:"dias_semana"`
	Estado      string `json:"estado"`
	ClienteID   uint   `json:"cliente_id"`
}

type ClienteListaSalida struct {
	ClienteSalida
	Entrenador *EntrenadorSalida       `json:"entrenador,omitempty"`
	Rutinas    []RutinaDeClienteSalida `json:"rutinas"`
}

func ClienteSalidaDesdeDominio(cliente dominio.Cliente) ClienteSalida {
	return ClienteSalida{
		ID:           cliente.ID,
		Nombre:       cliente.Nombre,
		Correo:       cliente.Correo,
		Objetivo:     cliente.Objetivo,
		EntrenadorID: cliente.EntrenadorID,
	}
}

func ClienteListaSalidaDesdeDominio(cliente dominio.Cliente) ClienteListaSalida {
	salida := ClienteListaSalida{
		ClienteSalida: ClienteSalidaDesdeDominio(cliente),
		Rutinas:       make([]RutinaDeClienteSalida, 0, len(cliente.Rutinas)),
	}
	if cliente.Entrenador.ID != 0 {
		salida.Entrenador = &EntrenadorSalida{
			ID:     cliente.Entrenador.ID,
			Nombre: cliente.Entrenador.Nombre,
			Correo: cliente.Entrenador.Correo,
		}
	}
	for _, rutina := range cliente.Rutinas {
		salida.Rutinas = append(salida.Rutinas, RutinaDeClienteSalida{
			ID:          rutina.ID,
			Titulo:      rutina.Titulo,
			Descripcion: rutina.Descripcion,
			DiasSemana:  rutina.DiasSemana,
			Estado:      rutina.Estado,
			ClienteID:   rutina.ClienteID,
		})
	}
	return salida
}

func (c ClienteEntrada) AModelo() dominio.Cliente {
	return dominio.Cliente{
		Nombre:       c.Nombre,
		Correo:       c.Correo,
		Objetivo:     c.Objetivo,
		EntrenadorID: c.EntrenadorID,
	}
}

type RutinaEntrada struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	DiasSemana  int    `json:"dias_semana"`
	ClienteID   uint   `json:"cliente_id"`
}

type RutinaActualizacion struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	DiasSemana  int    `json:"dias_semana"`
	ClienteID   uint   `json:"cliente_id"`
}

type RutinaSalida struct {
	ID          uint           `json:"id"`
	Titulo      string         `json:"titulo"`
	Descripcion string         `json:"descripcion"`
	DiasSemana  int            `json:"dias_semana"`
	Estado      string         `json:"estado"`
	ClienteID   uint           `json:"cliente_id"`
	Cliente     *ClienteSalida `json:"cliente,omitempty"`
}

func RutinaSalidaDesdeDominio(rutina dominio.Rutina) RutinaSalida {
	salida := RutinaSalida{
		ID:          rutina.ID,
		Titulo:      rutina.Titulo,
		Descripcion: rutina.Descripcion,
		DiasSemana:  rutina.DiasSemana,
		Estado:      rutina.Estado,
		ClienteID:   rutina.ClienteID,
	}
	if rutina.Cliente.ID != 0 || rutina.Cliente.Nombre != "" || rutina.Cliente.Correo != "" || rutina.Cliente.Objetivo != "" {
		cliente := ClienteSalidaDesdeDominio(rutina.Cliente)
		salida.Cliente = &cliente
	}
	return salida
}

func (r RutinaEntrada) AModelo() dominio.Rutina {
	return dominio.Rutina{
		Titulo:      r.Titulo,
		Descripcion: r.Descripcion,
		DiasSemana:  r.DiasSemana,
		ClienteID:   r.ClienteID,
	}
}

func (r RutinaActualizacion) AModelo() dominio.Rutina {
	return dominio.Rutina{
		Titulo:      r.Titulo,
		Descripcion: r.Descripcion,
		DiasSemana:  r.DiasSemana,
		ClienteID:   r.ClienteID,
	}
}
