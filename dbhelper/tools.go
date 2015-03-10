package main

import (
	"./models"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	schemaFilePath string
	driver         string
	dbUser         string
	dbPwd          string
	dbHost         string
	dbPort         int
	drop           = true
	operation      string
)

func init() {
	flag.StringVar(&schemaFilePath, "schema", "schema.xml", "The schema xml file path")
	flag.StringVar(&driver, "driver", "mysql", "The driver of the database")
	flag.StringVar(&dbUser, "user", "", "Database user name")
	flag.StringVar(&dbPwd, "pwd", "", "The password of the database user")
	flag.StringVar(&dbHost, "host", "127.0.0.1", "The database host name")
	flag.IntVar(&dbPort, "port", 3306, "The port for the database")
	flag.StringVar(&operation, "op", "show", "The operation to do, can be: show, populate, devsql, gostruct")
	flag.Parse()
}

func main() {
	dbSchema, err := models.ParseSchema(schemaFilePath)
	if err != nil {
		log.Println("Failed to parse schema", err)
		return
	}
	switch operation {
	case "populate":
		PopulateTables(dbSchema)
	case "show":
		ShowTableInfo(dbSchema)
	case "devsql":
		ShowDevSql(dbSchema)
	case "gostruct":
		ShowGoStruct(dbSchema)
	default:
		log.Println("Unknown operation:", operation)
	}
}

func ShowDevSql(dbSchema models.DBSchema) {
	for index, table := range dbSchema.Tables {
		if index > 0 {
			fmt.Println()
		}
		fmt.Println("========== " + table.Name + " ==========")
		sql, err := table.GetSelectSQL()
		if err != nil {
			log.Println("Failed to get select sql.", err)
		} else {
			fmt.Printf("Select SQL:\n\x1b[31;1m%s\x1b[0m\n", sql)
		}

		sql, err = table.GetInsertSQL()
		if err != nil {
			log.Println("Failed to get insert sql.", err)
		} else {
			fmt.Printf("Insert SQL:\n\x1b[31;1m%s\x1b[0m\n", sql)
		}
		sql, err = table.GetUpdateSQL()
		if err != nil {
			log.Println("Failed to get update sql.", err)
		} else {
			fmt.Printf("Update SQL:\n\x1b[31;1m%s\x1b[0m\n", sql)
		}

	}
}

func ShowGoStruct(dbSchema models.DBSchema) {
	for index, table := range dbSchema.Tables {
		if index > 0 {
			fmt.Println()
		}
		fmt.Println("// model for table " + table.Name)
		result, err := table.GetGoStruct()
		if err != nil {
			log.Println("Failed to get go struct.", err)
		} else {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", result)
		}
	}
}

func ShowTableInfo(dbSchema models.DBSchema) {
	for _, table := range dbSchema.Tables {
		sql, err := table.GetDropSQL()
		if err != nil {
			log.Println("Failed to generate drop sql.", err)
			continue
		}
		fmt.Printf("Drop sql:\n\x1b[31;1m%s\x1b[0m\n\n", sql)

		sql, err = table.GetCreateSQL()
		if err != nil {
			log.Println("Generate create sql failed.", err)
			continue
		}
		fmt.Printf("Create sql:\n\x1b[31;1m%s\x1b[0m\n\n", sql)
	}
}

func PopulateTables(dbSchema models.DBSchema) {
	if dbSchema.Name == "" || dbUser == "" || dbHost == "" || dbPwd == "" {
		log.Println("The db name, user, host or password is not specified.")
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPwd, dbHost, dbPort, dbSchema.Name)
	log.Println("The data source name is", dsn)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Println("Connect database failed.")
		return
	}
	defer db.Close()

	for _, table := range dbSchema.Tables {
		if drop {
			sql, err := table.GetDropSQL()
			if err != nil {
				log.Println("Failed to generate drop sql.", err)
				continue
			}
			log.Println("Execute sql:\n" + sql)
			_, err = db.Exec(sql)
			if err != nil {
				log.Println("Drop table failed.", err)
				log.Println("We will try to create table.")
			}
		}
		sql, err := table.GetCreateSQL()
		if err != nil {
			log.Println("Generate create sql failed.", err)
			continue
		}
		log.Println("Execute sql:\n" + sql)
		_, err = db.Exec(sql)
		if err != nil {
			log.Println("Execute sql "+sql+" Failed", err)
			return
		}
	}
}
