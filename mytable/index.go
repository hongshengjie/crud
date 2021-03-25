package mytable

import (
	"database/sql"

	"github.com/hongshengjie/crud/snaker"
)

// Index Index
type Index struct {
	IndexName    string    // index_name
	GoIndexName  string    // go index name
	IsUnique     bool      // is_unique
	IndexFields  []string  // index colomn name
	IndexColumns []*Column // index columns
}

// MyTableIndexes  MyTableIndexes
func MyTableIndexes(db *sql.DB, schema string, table string) ([]*Index, error) {
	var err error

	const sqlstr = `SELECT ` +
		`DISTINCT index_name, ` +
		`NOT non_unique AS is_unique ` +
		`FROM information_schema.statistics ` +
		`WHERE index_schema = ? AND table_name = ?`

	q, err := db.Query(sqlstr, schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	indexes := []*Index{}
	for q.Next() {
		i := Index{}
		err = q.Scan(&i.IndexName, &i.IsUnique)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, &i)
	}
	// column
	for _, v := range indexes {
		indexColumns, err := MyIndexColumns(db, schema, table, v.IndexName)
		if err != nil {
			return nil, err
		}
		v.IndexFields = indexColumns
		var indexGoName string
		for _, i := range indexColumns {
			indexGoName = indexGoName + snaker.SnakeToCamelIdentifier(i)
		}
		v.GoIndexName = indexGoName

	}

	return indexes, nil
}

// MyIndexColumns MyIndexColumns
func MyIndexColumns(db *sql.DB, schema string, table string, index string) ([]string, error) {
	var err error

	const sqlstr = `SELECT ` +
		`column_name ` +
		`FROM information_schema.statistics ` +
		`WHERE index_schema = ? AND table_name = ? AND index_name = ? ` +
		`ORDER BY seq_in_index`

	q, err := db.Query(sqlstr, schema, table, index)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	res := []string{}
	for q.Next() {
		var ic string
		err = q.Scan(&ic)
		if err != nil {
			return nil, err
		}

		res = append(res, ic)
	}

	return res, nil
}
