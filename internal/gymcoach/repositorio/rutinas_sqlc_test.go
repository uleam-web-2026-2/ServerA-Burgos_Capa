package repositorio

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	consultassqlc "github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/repositorio/sqlc"
)

type sinFilasConnector struct{}

func (sinFilasConnector) Connect(context.Context) (driver.Conn, error) {
	return sinFilasConn{}, nil
}

func (sinFilasConnector) Driver() driver.Driver {
	return sinFilasDriver{}
}

type sinFilasDriver struct{}

func (sinFilasDriver) Open(string) (driver.Conn, error) {
	return sinFilasConn{}, nil
}

type sinFilasConn struct{}

func (sinFilasConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("no soportado")
}

func (sinFilasConn) Close() error {
	return nil
}

func (sinFilasConn) Begin() (driver.Tx, error) {
	return nil, errors.New("no soportado")
}

func (sinFilasConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return sinFilasRows{}, nil
}

type sinFilasRows struct{}

func (sinFilasRows) Columns() []string {
	return []string{"id"}
}

func (sinFilasRows) Close() error {
	return nil
}

func (sinFilasRows) Next([]driver.Value) error {
	return io.EOF
}

func TestRutinasSqlcTraduceErrNoRows(t *testing.T) {
	db := sql.OpenDB(sinFilasConnector{})
	defer db.Close()

	repo := &RutinasSqlc{consultas: consultassqlc.New(db)}
	_, err := repo.ObtenerRutina(context.Background(), 123)
	if !errors.Is(err, dominio.ErrRutinaNoEncontrada) {
		t.Fatalf("se esperaba ErrRutinaNoEncontrada; llegó %v", err)
	}
}
