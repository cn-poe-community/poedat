package cli

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cn-poe-community/poedat/dat"
	"github.com/cn-poe-community/poedat/schema"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func DatToJSON(datFile *dat.DatFile, shecmaFile *schema.SchemaFile, tableName string, gameVersion int) []byte {
	headers := ImportHeaders(tableName, datFile, shecmaFile, gameVersion)
	rows := ExportAllRows(headers, datFile)

	rowMaps := []*linkedhashmap.Map{}

	for rid, row := range rows {
		rowMap := linkedhashmap.New()
		rowMap.Put("_rid", rid)
		for i, header := range headers {
			name := header.Name
			if len(name) == 0 {
				name = fmt.Sprintf("Unknown %d", i+1)
			}
			rowMap.Put(name, row[i])
		}

		rowMaps = append(rowMaps, rowMap)
	}

	data, err := json.MarshalIndent(rowMaps, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	return data
}
