package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cn-poe-community/poedat/cli"
	"github.com/cn-poe-community/poedat/dat"
	"github.com/cn-poe-community/poedat/schema"
)

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func toJSONPath(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, name+".json")
}

var datPath = flag.String("d", "", "the dat file path")
var schemaPath = flag.String("s", "", "the schema file path")
var tableName = flag.String("t", "", "the table name")
var gameVersion = flag.Int("g", 1, "the game version, poe1 is 1, poe2 is 2")

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

	jsonBytes := cli.DatToJSON(datFile, &schemaFile, *tableName, *gameVersion)
	err = os.WriteFile(toJSONPath(*datPath), jsonBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
