package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/chzyer/readline"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
)

// OpenDatabase otwiera bazę danych, która znajduje się w pliku
// o nazwie `filename`
func OpenDatabase(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Execute wysyła polecenie do bazy danych `db`. To polecenie powstaje
// z szablonu `f`, w którym kolejne symbole zastępcze zaczynające się
// znakiem % są zastępowane kolejnymi argumentami `args`
func Execute(db *sql.DB, f string, args ...any) {
	_, err := db.Exec(fmt.Sprintf(f, args...))
	if err != nil {
		log.Fatal(err)
	}
}

// Query wysyła zapytanie do bazy danych `db`. To zapytanie powstaje z
// szablonu `f`, w którym kolejne symbole zastępcze zaczynające się
// znakiem % są zastępowane kolejnymi argumentami `args`
func Query(db *sql.DB, f string, args ...any) *sql.Rows {
	rows, err := db.Query(fmt.Sprintf(f, args...))
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

// ScanRow kopiuje kolejne pola bieżącego wiersza argumentu `rows` do
// tych wartości, na które wskazują kolejne elementy argumentu `args`
func ScanRow(rows *sql.Rows, args ...any) {
	err := rows.Scan(args...)
	if err != nil {
		log.Fatal(err)
	}
}

// BeginTransaction rozpoczyna transakcję w bazie danych `db`
func BeginTransaction(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

// CommitTransaction zatwierdza transakcję `tx`
func CommitTransaction(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// PrepareStatement przygotowuje polecenie z parametrami tak, żeby
// można było używać tego polecenia wewnątrz transakcji `tx`. To
// polecenie powstaje z szablonu `f`, w którym kolejne symbole
// zastępcze zaczynające się znakiem % są zastępowane kolejnymi
// argumentami `args`
func PrepareStatement(tx *sql.Tx, f string, args ...any) *sql.Stmt {
	stmt, err := tx.Prepare(fmt.Sprintf(f, args...))
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

// ExecuteStatement wykonuje przygotowane polecenie `stmt`. Kolejne
// parametry tego polecenia są zastępowane kolejnymi argumentami
// `args`
func ExecuteStatement(stmt *sql.Stmt, args ...any) int64 {
	res, err := stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
	}
	rowid, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return rowid
}

// OpenFile otwiera plik o nazwie `name` do odczytu
func OpenFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// ReadCsvRecord odczytuje 1 wiersz pliku tekstowego za pomocą
// czytnika `reader`. Jeśli w tym pliku nie ma więcej danych,
// ReadCsvRecord zwraca `nil, io.EOF`
func ReadCsvRecord(reader *csv.Reader) ([]string, error) {
	rec, err := reader.Read()
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		log.Fatal(err)
	}
	return rec, nil
}

// CreateReadline tworzy nową instancję edytora wiersza poleceń. Każdy
// wiersz wyświetlany przez tę instancję zaczyna się od łańcucha
// `prompt`
func CreateReadline(prompt string) *readline.Instance {
	r, err := readline.New(prompt)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

// GetLine wczytuje polecenie użytkownika z wiersza poleceń za pomocą
// instancji edytora wiersza poleceń `rl`. Jeśli użytkownik wprowadził
// znak końca pliku lub polecenie, które ma 0 znaków, GetLine zwraca
// `nil, io.EOF`
func GetLine(rl *readline.Instance) (string, error) {
	q, err := rl.Readline()
	if q == "" || err == io.EOF {
		return q, io.EOF
	}
	if err != nil {
		log.Fatal(err)
	}
	return q, nil
}
