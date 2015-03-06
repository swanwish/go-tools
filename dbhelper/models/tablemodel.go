package models

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	indent = "    "
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

type DBSchema struct {
	Name   string  `xml:"name,attr"`
	Tables []Table `xml:"table"`
}

func (table Table) GetDropSQL() (string, error) {
	if table.Name == "" {
		log.Println("Table name not specified.")
		return "", errors.New("Table name not specified.")
	}
	return "DROP TABLE IF EXISTS " + table.Name, nil
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
	sql := "CREATE TABLE " + table.Name + " (\n"

	pks := ""
	pkCount := 0
	for index, column := range table.Columns {
		defineSql, err := column.GetDefineSQL()
		if err != nil {
			log.Println("Failed to get column define sql.", err)
			return "", err
		}
		if index > 0 {
			sql += ",\n"
		}
		sql += indent + defineSql
		if column.PK == 1 {
			if pkCount > 0 {
				pks += ", "
			}
			pks += column.Name
			pkCount++
		}
	}
	if pkCount > 0 {
		sql += ",\n" + indent + "PRIMARY KEY (" + pks + ")"
	}
	sql += "\n)"

	if table.Description != "" {
		sql += " COMMENT='" + table.Description + "'"
	}

	return sql, nil
}

func (table Table) GetSelectSQL() (string, error) {
	if table.Name == "" {
		log.Println("Table name not specified.")
		return "", errors.New("Table name not specified.")
	}

	if len(table.Columns) == 0 {
		log.Println("No column defined for table ", table.Name)
		return "", errors.New("No column defined for table " + table.Name)
	}
	sql := "SELECT "

	columns := make([]string, 0)
	for _, column := range table.Columns {
		columns = append(columns, column.Name)
	}
	sql += strings.Join(columns, ", ")
	sql += " FROM " + table.Name

	return sql, nil
}

func (table Table) GetUpdateSQL() (string, error) {
	if table.Name == "" {
		log.Println("Table name not specified.")
		return "", errors.New("Table name not specified.")
	}

	if len(table.Columns) == 0 {
		log.Println("No column defined for table ", table.Name)
		return "", errors.New("No column defined for table " + table.Name)
	}
	sql := "UPDATE " + table.Name + " SET "

	columns := make([]string, 0)
	wheres := make([]string, 0)
	for _, column := range table.Columns {
		if column.PK == 1 {
			wheres = append(wheres, column.Name+"=?")
		} else {
			columns = append(columns, column.Name+"=?")
		}
	}
	sql += strings.Join(columns, ", ")
	sql += " WHERE " + strings.Join(wheres, " AND ")

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

func ParseSchema(schemaFilePath string) (DBSchema, error) {
	xmlFile, err := os.Open(schemaFilePath)
	if err != nil {
		log.Println("Open file failed.", err)
		return DBSchema{}, err
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Read file failed.", err)
		return DBSchema{}, err
	}

	var dbSchema = DBSchema{}
	err = xml.Unmarshal(b, &dbSchema)

	if err != nil {
		log.Println("Failed unmarshal db schema", err)
		return DBSchema{}, err
	}

	return dbSchema, nil
}
