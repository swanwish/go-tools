package main

import (
	// "database/sql"
	"encoding/xml"
	"errors"
	// _ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
)

const (
	schemaFilePath = "schema.xml"
	driver         = "mysql"
	dataSourceName = "username:password@tcp(host:3306)/dbname"
)

type Column struct {
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr"`
	Length      string `xml:"length,attr"`
	PK          int64  `xml:"pk,attr"`
	NotNull     int64  `xml:"notnull,attr"`
	Description string `xml:"description,attr"`
}

type Table struct {
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description,attr"`
	Columns     []Column `xml:"column"`
}

type Db struct {
	Name   string  `xml:"name,attr"`
	Tables []Table `xml:"table"`
}

func (table Table) GetCreateSQL() (string, error) {
	if table.Name == "" {
		log.Println("Table name not specified.")
		return "", errors.New("Table name not specified.")
	}

	if len(table.Columns) == 0 {
		log.Println("No column defined for table ", table.Name)
		return "", errors.New("No column defined for table " + table.Name)
	}
	sql := "CREATE TABLE " + table.Name + " ("

	pks := ""
	pkCount := 0
	for index, column := range table.Columns {
		defineSql, err := column.GetDefineSQL()
		if err != nil {
			log.Println("Failed to get column define sql.", err)
			return "", err
		}
		if index > 0 {
			sql += ", "
		}
		sql += defineSql
		if column.PK == 1 {
			if pkCount > 0 {
				pks += ", "
			}
			pks += column.Name
			pkCount++
		}
	}
	if pkCount > 0 {
		sql += ", PRIMARY KEY (" + pks + ")"
	}
	sql += ")"

	return sql, nil
}

func (column Column) GetDefineSQL() (string, error) {
	if column.Name == "" {
		return "", errors.New("Column has no name defined.")
	}
	if column.Type == "" {
		return "", errors.New("Column has no type defined.")
	}
	defineSql := column.Name + " " + column.Type
	if column.Length != "" {
		defineSql += "(" + column.Length + ")"
	}
	notNull := 0
	if column.NotNull == 1 {
		notNull = 1
	}
	if column.PK == 1 {
		notNull = 1
	}
	if notNull == 1 {
		defineSql += " NOT NULL"
	}
	if column.Description != "" {
		defineSql += " COMMENT '" + column.Description + "'"
	}
	return defineSql, nil
}

func main() {
	xmlFile, err := os.Open("schema.xml")
	if err != nil {
		log.Println("Open file failed.", err)
		return
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Read file failed.", err)
	}

	var schema = Db{}
	xml.Unmarshal(b, &schema)

	// db, err := sql.Open(driver, dataSourceName)
	// if err != nil {
	// 	log.Println("Connect database failed.")
	// 	return
	// }
	// defer db.Close()

	for _, table := range schema.Tables {
		sql, err := table.GetCreateSQL()
		if err != nil {
			log.Println("Generate create sql failed.", err)
			continue
		}
		log.Println("Execute sql:" + sql)
		// _, err = db.Exec(sql)
		// if err != nil {
		// 	log.Println("Execute sql "+sql+" Failed", err)
		// 	return
		// }
	}
}
