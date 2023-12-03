package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cn-poe-community/poedat/dat"
	"github.com/cn-poe-community/poedat/schema"
)

func DatToJSON(datFile *dat.DatFile, shecmaFile *schema.SchemaFile, tableName string) []byte {
	headers := ImportHeaders(tableName, datFile, shecmaFile)
	rows := ExportAllRows(headers, datFile)

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

	return buffer.Bytes()
}
