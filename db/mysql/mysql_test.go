package mysql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bongnv/gwf"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func Test_WithMYSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.6")
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(rows)
	app := gwf.New(WithMYSQL(Config{
		Conn: db,
	}))
	component, err := app.Component("db")
	require.NoError(t, err)
	require.IsType(t, &gorm.DB{}, component)
}
