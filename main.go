package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
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

func saveRowsAsJson(headers []*dat.Header, rows [][]dat.FieldValue, path string) {
	var buffer bytes.Buffer

	for _, row := range rows {
		rowMap := map[string]dat.FieldValue{}
		for i, header := range headers {
			name := header.Name
			if len(name) == 0 {
				name = fmt.Sprintf("Unknown %d", i)
			}
			rowMap[name] = row[i]
		}

		jsonBytes, err := json.MarshalIndent(rowMap, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		buffer.WriteString(string(jsonBytes))
	}

	os.WriteFile(path, buffer.Bytes(), 0644)
}

var datPath = flag.String("d", "", "the dat file path\n\texported json will saved as {dartPath}.json")
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

	headers := export.ImportHeaders(*tableName, datFile, &schemaFile)
	rows := export.ExportAllRows(headers, datFile)
	saveRowsAsJson(headers, rows, *datPath+".json")
}
