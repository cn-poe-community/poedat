package schema_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cn-poe-community/poedat/schema"
)

func TestSchemaTypes(t *testing.T) {
	data, err := os.ReadFile("schema.min.json")
	if err != nil {
		t.Error(err)
	}

	var schemaFile schema.SchemaFile

	err = json.Unmarshal(data, &schemaFile)
	if err != nil {
		t.Error(err)
	}

	if len(schemaFile.Tables) == 0 {
		t.Error("schemaFile.Tables unmarshaled failed")
	}

	if len(schemaFile.Enumerations) == 0 {
		t.Error("schemaFile.Enumerations unmarshaled failed")
	}
}
