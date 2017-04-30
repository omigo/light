package mapper

import (
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/arstd/light/example/enum"
)

func initSQLMock(t *testing.T) {
	var err error
	var mock sqlmock.Sqlmock
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("mock error: '%s' ", err)
	}
	// defer db.Close()

	// TestModelMapperInsert
	mock.ExpectQuery(`insert into .+ returning id`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// TestModelMapperInsertTx
	mock.ExpectBegin()
	mock.ExpectQuery(`insert into .+ returning id`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// TestModelMapperBatchInsert
	mock.ExpectExec(`insert into .+ values\s*\(.+?\)(\s*,\s*\(.+?\))+`).
		WillReturnResult(sqlmock.NewResult(3, 3))

	rows := sqlmock.NewRows([]string{"id", "name", "flag", "score", "map",
		"time", "xarray", "slice", "status", "pointer", "struct_slice", "uint32"}).
		AddRow(1, "name", true, 1.23, `{"a": 1}`, time.Now().Add(3*time.Hour),
			`{1,2,3}`, `{Slice Elem 1,Slice Elem 2}`,
			enum.StatusNormal, `{"Name": "Pointer"}`,
			`[{"Name": "StructSlice"}]`, uint32(time.Now().Unix()))
	// TestModelMapperGet
	mock.ExpectQuery(`select`).WillReturnRows(rows)

	// TestModelMapperUpdate
	mock.ExpectExec(`update`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// TestModelMapperDelete
	mock.ExpectExec(`delete`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// TestModelMapperCount
	mock.ExpectQuery(`select count`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	rows2 := sqlmock.NewRows([]string{"id", "name", "flag", "score", "map",
		"time", "xarray", "slice", "status", "pointer", "struct_slice", "uint32"}).
		AddRow(1, "name", true, 1.23, `{"a": 1}`, time.Now().Add(3*time.Hour),
			`{1,2,3}`, `{Slice Elem 1,Slice Elem 2}`,
			enum.StatusNormal, `{"Name": "Pointer"}`,
			`[{"Name": "StructSlice"}]`, uint32(time.Now().Unix()))
	// TestModelMapperList
	mock.ExpectQuery(`select.+offset.+`).
		WillReturnRows(rows2)

	rows3 := sqlmock.NewRows([]string{"id", "name", "flag", "score", "map",
		"time", "xarray", "slice", "status", "pointer", "struct_slice", "uint32"}).
		AddRow(1, "name", true, 1.23, `{"a": 1}`, time.Now().Add(3*time.Hour),
			`{1,2,3}`, `{Slice Elem 1,Slice Elem 2}`,
			enum.StatusNormal, `{"Name": "Pointer"}`,
			`[{"Name": "StructSlice"}]`, uint32(time.Now().Unix()))
	// TestModelMapperPage
	mock.ExpectQuery(`select count`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(`select.+offset.+`).
		WillReturnRows(rows3)
}
