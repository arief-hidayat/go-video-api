package models

import (
	"database/sql"
	"log"

	"github.com/arief-hidayat/go-video-api/query"

	_ "github.com/lib/pq"
)

// This time the global variable is unexported.
var db *sql.DB

// https://www.alexedwards.net/blog/organising-database-access
// InitDB sets up setting up the connection pool global variable.
func InitDB(dataSourceName string, maxOpenConnections int) error {
	var err error

	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxOpenConnections)
	return db.Ping()
}

func CloseDB() {
    db.Close()
}

func GetVideos(queryStr string) ([]interface{}, error) {
	rows, err := db.Query(query.SqlSearch, queryStr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return makeStructJSON(rows)
}

func makeStructJSON(rows *sql.Rows) ([]interface{}, error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	count := len(columnTypes)
	finalRows := []interface{}{}
	for rows.Next() {
		scanArgs := make([]interface{}, count)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		masterData := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}
			masterData[v.Name()] = scanArgs[i]
		}
		finalRows = append(finalRows, masterData)
	}
	return finalRows, nil
}
