package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	// Nazwa pliku tekstowego, którego każdy wiersz zawiera dane
	// jednego pracownika AGH, rozdzielone przecinkami
	CSVFilename = "AGH.csv"
	// Nazwa pliku, który zawiera bazę danych SQLite3
	DbFilename = "AGH.sqlite3"
)

// Nazwy kolumn w pliku tekstowym `CSVFilename` i w tabeli Pracownicy
// w bazie danych `DBFilename`
type ColumnName string

const (
	Osoba      ColumnName = "osoba"
	Stanowisko ColumnName = "stanowisko"
	Jednostka  ColumnName = "jednostka"
	Budynek    ColumnName = "budynek"
	Piętro     ColumnName = "piętro"
	Pokój      ColumnName = "pokój"
	Telefon    ColumnName = "telefon"
	Adres      ColumnName = "adres"
)

// Nazwy kolejnych kolumn tabeli Pracownicy
var ColumnNames = []ColumnName{
	Osoba,
	Stanowisko,
	Jednostka,
	Budynek,
	Piętro,
	Pokój,
	Telefon,
	Adres,
}

// main pobiera z wiersza poleceń kolejne pytania użytkownika,
// wyrażone po polsku, przetwarza te pytania na zapytania do bazy
// danych SQLite3 i wypisuje na standardowym wyjściu wyniki tych
// zapytań. Jeśli funkcja main wczytała pusty łańcuch zamiast pytania,
// program się kończy. Jeśli program został uruchomiony z flagą -init,
// main najpierw tworzy nową bazę danych w pliku `DbFilename` i
// zapisuje w niej dane z pliku tekstowego `CSVFilename`, a potem
// działa tak, jak opisano powyżej
func main() {
	var db *sql.DB
	if len(os.Args) > 1 && os.Args[1] == "-init" {
		db = CreateDatabase(DbFilename)
		FillDatabase(CSVFilename, db)
	} else {
		db = OpenDatabase(DbFilename)
	}
	defer db.Close()
	colsOfStems := GetColumnsOfStems(db)

	rl := CreateReadline("AGH> ")
	defer rl.Close()
	for {
		s, err := GetLine(rl)
		if err == io.EOF {
			break
		}
		s = TransformPhoneNumbers(s)
		as := ToASCIIString(s)
		match, cols, err := ParseQuestion(as, colsOfStems)
		if err != nil {
			fmt.Println(err)
			continue
		}
		query := MakeQuery(match, cols)
		res := ExecuteQuery(query, cols, db)
		DisplayResult(res)
	}
}

// CreateDatabase tworzy nową bazę danych w pliku o nazwie `filename`
//
// Ta baza danych zawiera 2 tabele:
// + tabelę Pracownicy
// + tabelę PracownicyFTS
//
// Każdy wiersz tabeli Pracownicy odpowiada 1 osobie, która jest
// pracownikiem AGH i składa się z kolumn o nazwach wymienionych
// w zmiennej globalnej `ColumnNames`
//
// Każdy wiersz tabeli PracownicyFTS ma 2 pola:
// + pole rowid
// + pole dane
//
// Pole rowid wiersza tabeli PracownicyFTS ma tę samą wartość, co pole
// rowid odpowiedniego wiersza tabeli Pracownicy
//
// Pole dane wiersza tabeli PracownicyFTS zawiera tematy tych wyrazów,
// które występują w odpowiednim wierszu tabeli Pracownicy
//
// Przykład:
//
// Wiersz tabeli Pracownicy
// rowid:      14
// osoba:      inż. Anna Kot
// stanowisko: specjalista
// jednostka:  Wydział Informatyki
// budynek:    D-17
// piętro:     V
// pokój:      6.11
// telefon:    12-328-99-99
// adres:      ul. Kawiory 21
//
// Wiersz tabeli PracownicyFTS
// rowid: 14
// dane:  inz ann kot specjalist wydzial informatyk d-17 v 6.11 12-328-99-99 ul kawior 21
func CreateDatabase(dbFilename string) *sql.DB {
	_ = os.Remove(dbFilename)
	db := OpenDatabase(dbFilename)
	Execute(db, `CREATE TABLE Pracownicy(
docid INTEGER PRIMARY KEY
, %s TEXT
, %s TEXT
, %s TEXT
, %s TEXT
, %s TEXT
, %s TEXT
, %s TEXT
, %s TEXT)`, ToAnySlice(ColumnNames)...)
	Execute(db, `CREATE VIRTUAL TABLE PracownicyFTS USING fts5(dane)`)
	return db
}

// FillDatabase zapisuje w bazie danych `db` dane z pliku tekstowego o
// nazwie `csvFilename`. Każdy wiersz tego pliku zawiera te same pola,
// co tabela Pracownicy, rozdzielone przecinkami
func FillDatabase(csvFilename string, db *sql.DB) {
	tx := BeginTransaction(db)
	insertStmt := PrepareStatement(
		tx, `INSERT INTO Pracownicy(%s,%s,%s,%s,%s,%s,%s,%s)
VALUES (?,?,?,?,?,?,?,?)`, ToAnySlice(ColumnNames)...)
	insertFTSStmt := PrepareStatement(
		tx, `INSERT INTO PracownicyFTS(rowid,dane) VALUES (?,?)`)
	file := OpenFile(csvFilename)
	defer file.Close()

	csvFile := csv.NewReader(file)
	firstRecord := true
	header := []ColumnName{}
	for {
		rec, err := ReadCsvRecord(csvFile)
		if err == io.EOF {
			break
		}
		if firstRecord {
			for _, c := range rec {
				header = append(header, ColumnName(c))
			}
			firstRecord = false
			continue
		}
		row := map[ColumnName]string{}
		for i, c := range header {
			row[c] = rec[i]
		}
		rowid := ExecuteStatement(
			insertStmt,
			row[Osoba], row[Stanowisko], row[Jednostka],
			row[Budynek], row[Piętro], row[Pokój],
			TransformPhoneNumbers(row[Telefon]),
			row[Adres])
		stems := []ASCIIStem{}
		for _, c := range ColumnNames {
			as := ToASCIIString(row[c])
			// Nie usuwaj piętra I
			ss := ASCIIStringToASCIIStemSlice(as, c != Piętro)
			stems = append(stems, ss...)
		}
		// Zadanie 6.
		//
		// Jeśli row[Jednostka] zaczyna się od łańcucha "Wydział ", to:
		// 1. Podziel row[Jednostka] na części rozdzielone przecinkami,
		//    po których następuje wyraz nie zakończony na -i ani na -y
		// 2. Zamień zerową otrzymaną część row[Jednostka] na wartość
		//    typu ASCIIString
		// 3. Zamień otrzymaną wartość typu ASCIIString na wycinek
		//    zawierający wartości typu ASCIIStem, nie usuwając żadnych
		//    tematów wyrazów i połącz pierwsze litery wartości z tego
		//    wycinka w łańcuch
		// 4. Zamień otrzymaną wartość typu ASCIIString na wycinek
		//    zawierający wartości typu ASCIIStem, usuwając takie
		//    tematy wyrazów, które można pomylić z innymi wyrazami
		//    z ich tematami, i połącz pierwsze litery wartości z tego
		//    wycinka w łańcuch
		// 5. Dodaj do wycinka stems łańcuchy otrzymane w punktach
		//    3 i 4
		//
		// Przykład:
		//
		// Jeśli row[Jednostka] ==
		// "Wydział Geologii, Geofizyki i Ochrony Środowiska, Dziekanat"
		// to dodaj do wycinka stems łańcuchy "wggios" i "wggos"
		//
		// Wskazówki:
		//
		// + Proszę odkomentować funkcję TestAbbreviateFacultyName
		//   w pliku transform_test.go
		// + Polecenie "go doc strings" jest państwa przyjacielem
		ExecuteStatement(insertFTSStmt, rowid, JoinASCIIStems(stems))
	}
	CommitTransaction(tx)
}

// GetColumnsOfStems zwraca mapę tematów wyrazów pochodzących z tabeli
// Pracownicy na zbiory tych kolumn, w których występują te wyrazy
//
// Przykład:
//
// Tabela Pracownicy
// osoba: Anna Kot, stanowisko: specjalista, jednostka: Wydział Informatyki,...
// osoba: Anna Nowy, stanowisko: specjalista ds. kotów, jednostka: Wydział Biologii,...
//
// Wynik funkcji GetColumnsOfStems:
// ann:        {osoba}
// kot:        {osoba, stanowisko}
// now:        {osoba}
// specjalist: {stanowisko}
// ds:         {stanowisko}
// wydzial:    {jednostka}
// informatyk: {jednostka}
// biolog:     {jednostka}
func GetColumnsOfStems(db *sql.DB) map[ASCIIStem]map[ColumnName]bool {
	rec, args := MakeStringSliceAndAnySlice(len(ColumnNames))
	ret := map[ASCIIStem]map[ColumnName]bool{}
	rows := Query(db, `SELECT %s,%s,%s,%s,%s,%s,%s,%s FROM Pracownicy`,
		ToAnySlice(ColumnNames)...)
	for rows.Next() {
		ScanRow(rows, args...)
		for i, c := range ColumnNames {
			as := ToASCIIString(rec[i])
			// Nie usuwaj piętra I
			stems := ASCIIStringToASCIIStemSlice(as, c != Piętro)
			for _, stem := range stems {
				if ret[stem] == nil {
					ret[stem] = map[ColumnName]bool{}
				}
				ret[stem][c] = true
			}
		}
	}
	return ret
}

// Wyrazy, które zmieniają
var Negations = map[ASCIIWord]bool{
	"nie":    true,
	"ani":    true,
	"bez":    true,
	"oprocz": true,
	"poza":   true,
}

var Replacements = map[ASCIIStem]ASCIIStem{
	"licencjat":   "lic",
	"inzynier":    "inz",
	"magister":    "mgr",
	"magistr":     "mgr",
	"doktor":      "dr",
	"habilitowan": "hab",
	"profesor":    "prof",
	"zwyczajn":    "zw",
	"nadzwyczajn": "nadzw",
	"gmach":       "a-0",
}

var StemsToColumnNames = map[ASCIIStem][]ColumnName{
	"stanowisk": []ColumnName{Stanowisko},
	"funkcj":    []ColumnName{Stanowisko},
	"gd":        []ColumnName{Stanowisko, Jednostka},
	"jednostk":  []ColumnName{Stanowisko, Jednostka},
	"pokoj":     []ColumnName{Budynek, Piętro, Pokój},
	"gabinet":   []ColumnName{Budynek, Piętro, Pokój},
	"sal":       []ColumnName{Budynek, Piętro, Pokój},
	"pietr":     []ColumnName{Budynek, Piętro},
	"budynk":    []ColumnName{Budynek},
	"numer":     []ColumnName{Telefon},
	"telefon":   []ColumnName{Telefon},
	"adres":     []ColumnName{Adres},
	"ulic":      []ColumnName{Adres},
}

type Conjunction string

const (
	And Conjunction = "AND"
	Not Conjunction = "NOT"
)

// ParseQuestion przetwarza pytanie `as`, wyrażone po polsku, na
// łańcuch i na nazwy kolumn. Łańcuch opisuje te wartości pól tabeli
// Pracownicy, które zna użytkownik. Nazwy kolumn nazywają te kolumny
// tabeli Pracownicy, których zawartość chce poznać użytkownik
func ParseQuestion(
	as ASCIIString,
	colsOfStems map[ASCIIStem]map[ColumnName]bool) (string, []ColumnName, error) {
	conj := And
	cols := map[ColumnName]bool{Osoba: true}
	parts := map[Conjunction][]string{}
	for _, s := range SplitASCIIString(as) {
		w := ToASCIIWord(s)
		if Negations[w] {
			conj = Not
			continue
		}
		stems := []ASCIIStem{}
		for _, stem := range ToASCIIStems(w) {
			if repl := Replacements[stem]; repl != "" {
				stem = repl
			}
			if colNames := StemsToColumnNames[stem]; colNames != nil {
				for _, c := range colNames {
					cols[c] = true
				}
			} else if colMap := colsOfStems[stem]; colMap != nil {
				for c, _ := range colMap {
					cols[c] = true
				}
				stems = append(stems, stem)
			}
		}
		if len(stems) > 0 {
			parts[conj] = append(
				parts[conj], JoinQuotedStems(stems, ` OR `))
			conj = And
		}
	}
	if len(parts[And]) == 0 {
		return "", nil, errors.New("W Twoim pytaniu brak konkretów")
	}
	retMatch := strings.Join(parts[And], " AND ")
	if len(parts[Not]) > 0 {
		retMatch += " NOT " + strings.Join(parts[Not], " NOT ")
	}
	retCols := []ColumnName{}
	for _, c := range ColumnNames {
		if cols[c] {
			retCols = append(retCols, c)
		}
	}
	return retMatch, retCols, nil
}

// MakeQuery tworzy zapytanie w języku SQL z łańcucha `match` i z nazw
// kolumn `cols`. Łańcuch `match` opisuje te wartości pól tabeli
// Pracownicy, które zna użytkownik. Nazwy kolumn `cols` nazywają te
// kolumny tabeli Pracownicy, których zawartość chce poznać użytkownik
func MakeQuery(match string, cols []ColumnName) string {
	return fmt.Sprintf(`SELECT %s
FROM Pracownicy JOIN PracownicyFTS
ON Pracownicy.rowid = PracownicyFTS.rowid
WHERE dane MATCH '%s'`, strings.Join(ToStringSlice(cols), ", "), match)
}

// ExecuteQuery zwraca wynik zapytania `q` do bazy danych `db`. Wynik
// funkcji ExecuteQuery to wycinek, który składa się z wycinków
// złożonych z łańcuchów. Pierwszy z tych wycinków zawiera łańcuchy,
// które są zapisem kolejnych liczb porządkowych. Następne wycinki
// odpowiadają kolumnom o nazwach `cols`. Każdy z tych wycinków
// zawiera wartość odpowiedniej kolumny w kolejnych wierszach wyniku
// zapytania `q`. Pierwszy element każdego wycinka złożonego z
// łańcuchów to jego nagłówek
func ExecuteQuery(q string, cols []ColumnName, db *sql.DB) [][]string {
	ret := [][]string{[]string{"lp"}}
	for _, c := range cols {
		ret = append(ret, []string{string(c)})
	}
	row, arg := MakeStringSliceAndAnySlice(len(cols))
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	for n := 1; rows.Next(); n++ {
		err := rows.Scan(arg...)
		if err != nil {
			log.Fatal(err)
		}
		ret[0] = append(ret[0], fmt.Sprintf("%d", n))
		for i, c := range row {
			ret[i+1] = append(ret[i+1], c)
		}
	}
	return ret
}

// DisplayResult wypisuje na standardowym wyjściu `result`, czyli
// wynik zapytania do bazy danych. Każda kolumna wyniku ma stałą
// szerokość. Pierwsza kolumna, która zawiera zapis liczby
// porządkowej, jest wyrównana do prawej strony. Kolejne kolumny są
// wyrównane do lewej strony. Jeśli `result` składa się tylko z
// nagłówków kolumn, DisplayResult wypisuje zamiast tych nagłówków
// komunikat "Nie znam takich osób"
func DisplayResult(res [][]string) {
	if len(res) == 0 {
		return
	}
	if len(res[0]) <= 1 {
		fmt.Println("Nie znam takich osób")
		return
	}
	lens := make([]int, len(res))
	for x, col := range res {
		lens[x] = 0
		for _, s := range col {
			lens[x] = max(lens[x], len(s))
		}
	}
	for y := range res[0] {
		for x, col := range res {
			if x == 0 {
				fmt.Printf("%*s ", lens[x], col[y])
			} else {
				fmt.Printf("%-*s ", lens[x], col[y])
			}
		}
		fmt.Println()
	}
}
