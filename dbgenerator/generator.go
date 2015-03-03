package main

import (
	// "database/sql"
	"encoding/xml"
	"errors"
	"flag"
	// "fmt"
	// _ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
)

var (
	schemaFilePath string
	driver         string
	dbUser         string
	dbPwd          string
	dbHost         string
	dbPort         int
	indent         = "    "
	drop           = true
)

func init() {
	flag.StringVar(&schemaFilePath, "schema", "schema.xml", "The schema xml file path")
	flag.StringVar(&driver, "driver", "mysql", "The driver of the database")
	flag.StringVar(&dbUser, "user", "", "Database user name")
	flag.StringVar(&dbPwd, "pwd", "", "The password of the database user")
	flag.StringVar(&dbHost, "host", "127.0.0.1", "The database host name")
	flag.IntVar(&dbPort, "port", 3306, "The port for the database")
}

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
	flag.Parse()
	xmlFile, err := os.Open(schemaFilePath)
	if err != nil {
		log.Println("Open file failed.", err)
		return
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Read file failed.", err)
	}

	var dbSchema = Db{}
	xml.Unmarshal(b, &dbSchema)

	// if dbSchema.Name == "" {
	// 	log.Println("The db name is not specified.")
	// 	return
	// }

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPwd, dbHost, dbPort, dbSchema.Name)
	// log.Println("The data source name is", dsn)
	// //	return
	// db, err := sql.Open(driver, dsn)
	// if err != nil {
	// 	log.Println("Connect database failed.")
	// 	return
	// }
	// defer db.Close()

	for _, table := range dbSchema.Tables {
		if drop {
			sql, err := table.GetDropSQL()
			if err != nil {
				log.Println("Failed to generate drop sql.", err)
				continue
			}
			log.Println("Execute sql:\n" + sql)
			// _, err = db.Exec(sql)
			// if err != nil {
			// 	log.Println("Drop table failed.", err)
			// 	log.Println("We will try to create table.")
			// }
		}
		sql, err := table.GetCreateSQL()
		if err != nil {
			log.Println("Generate create sql failed.", err)
			continue
		}
		log.Println("Execute sql:\n" + sql)
		// _, err = db.Exec(sql)
		// if err != nil {
		// 	log.Println("Execute sql "+sql+" Failed", err)
		// 	return
		// }
	}
}
