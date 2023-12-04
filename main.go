package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/cn-poe-community/poedat/dat"
	"github.com/cn-poe-community/poedat/export"
	"github.com/cn-poe-community/poedat/schema"
)

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

var datPath = flag.String("d", "", "the dat file path\n\texported json will be saved as {dartPath}.json")
var schemaPath = flag.String("s", "", "the schema file path")
var tableName = flag.String("t", "", "the table name")

func main() {
	flag.Parse()
	if *datPath == "" || *schemaPath == "" || *tableName == "" {
		flag.Usage()
		os.Exit(1)
	}

	schemaBytes := readFile(*schemaPath)
	var schemaFile schema.SchemaFile

	err := json.Unmarshal(schemaBytes, &schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	datBytes := readFile(*datPath)
	datFile, err := dat.ReadDatFile(datBytes)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes := export.DatToJSON(datFile, &schemaFile, *tableName)
	err = os.WriteFile(*datPath+".json", jsonBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
