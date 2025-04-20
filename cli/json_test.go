package cli

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cn-poe-community/poedat/dat"
	"github.com/cn-poe-community/poedat/schema"
)

func TestDatToJSON(t *testing.T) {
	schemaBytes, err := os.ReadFile("../schema/schema.min.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaFile schema.SchemaFile
	err = json.Unmarshal(schemaBytes, &schemaFile)
	if err != nil {
		t.Fatal(err)
	}

	datBytes, err := os.ReadFile("./testfiles/baseitemtypes.dat64")
	if err != nil {
		t.Fatal(err)
	}
	datFile, err := dat.ReadDatFile(datBytes)
	if err != nil {
		t.Fatal(err)
	}

	DatToJSON(datFile, &schemaFile, "BaseItemTypes", 1)
}
