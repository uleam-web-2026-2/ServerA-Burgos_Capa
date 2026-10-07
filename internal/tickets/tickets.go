package tickets

import "sync"

// Ticket representa un caso de soporte o solicitud del usuario.
type Ticket struct {
	ID        int    `json:"id"`
	Titulo    string `json:"titulo"`
	Prioridad string `json:"prioridad"`
	Estado    string `json:"estado"`
}

// Almacen guarda los tickets en memoria con acceso concurrente seguro.
type Almacen struct {
	mu     sync.Mutex
	ultimo int
	items  map[int]Ticket
}

func NuevoAlmacen() *Almacen {
	return &Almacen{items: make(map[int]Ticket)}
}

// Crear crea un ticket nuevo y devuelve la instancia generada.
func (a *Almacen) Crear(titulo, prioridad string) Ticket {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ultimo++
	t := Ticket{
		ID:        a.ultimo,
		Titulo:    titulo,
		Prioridad: prioridad,
		Estado:    "abierto",
	}
	a.items[t.ID] = t
	return t
}

// crear mantiene compatibilidad con implementaciones previas en minúsculas.
func (a *Almacen) crear(titulo, prioridad string) Ticket {
	return a.Crear(titulo, prioridad)
}

// Obtener devuelve el ticket por ID.
func (a *Almacen) Obtener(id int) (Ticket, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	t, ok := a.items[id]
	return t, ok
}

// obtener mantiene compatibilidad con implementaciones previas en minúsculas.
func (a *Almacen) obtener(id int) (Ticket, bool) {
	return a.Obtener(id)
}

// Listar devuelve todos los tickets en orden ascendente por ID.
func (a *Almacen) Listar() []Ticket {
	a.mu.Lock()
	defer a.mu.Unlock()

	lista := make([]Ticket, 0, len(a.items))
	for i := 1; i <= a.ultimo; i++ {
		if t, ok := a.items[i]; ok {
			lista = append(lista, t)
		}
	}
	return lista
}

// listar mantiene compatibilidad con implementaciones previas en minúsculas.
func (a *Almacen) listar() []Ticket {
	return a.Listar()
}
