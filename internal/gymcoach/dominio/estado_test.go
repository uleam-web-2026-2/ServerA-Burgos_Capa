package dominio

import "testing"

func TestPuedeTransicionarRespetaLaMaquinaDeEstados(t *testing.T) {
	casos := []struct {
		de, a string
		ok    bool
	}{
		{"pendiente", "en_progreso", true},
		{"en_progreso", "completada", true},
		{"pendiente", "completada", false},
		{"completada", "pendiente", false},
		{"desconocido", "pendiente", false},
	}
	for _, caso := range casos {
		if resultado := PuedeTransicionar(caso.de, caso.a); resultado != caso.ok {
			t.Errorf("PuedeTransicionar(%q, %q) = %v; se esperaba %v", caso.de, caso.a, resultado, caso.ok)
		}
	}
}
