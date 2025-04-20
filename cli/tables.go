package cli

import (
	"log"

	"github.com/cn-poe-community/poedat/dat"
	"github.com/cn-poe-community/poedat/schema"
)

// https://github.com/SnosMe/poedat-viewer/blob/master/lib/src/cli/export-tables.ts

func ImportHeaders(name string, datFile *dat.DatFile, schemaFile *schema.SchemaFile, validFor int) []*dat.Header {
	var headers []*dat.Header
	var sch *schema.SchemaTable

	for _, table := range schemaFile.Tables {
		if table.Name == name && (validFor == int(table.ValidFor) || table.ValidFor == schema.ValidForCommon) {
			sch = table
			break
		}
	}

	if sch == nil {
		log.Fatalf("find no schema table matched %v", name)
		return nil
	}

	offset := 0

	for _, column := range sch.Columns {
		var header dat.Header
		var htype dat.HeaderType

		if column.Name == nil {
			header.Name = ""
		} else {
			header.Name = *column.Name
		}

		header.Offset = offset

		htype.Array = column.Array
		switch column.Type {
		case "u16":
			htype.Integer = &dat.IntergerType{
				Unsigned: true,
				Size:     2,
			}
		case "u32":
			htype.Integer = &dat.IntergerType{
				Unsigned: true,
				Size:     4,
			}
		case "i16":
			htype.Integer = &dat.IntergerType{
				Unsigned: false,
				Size:     2,
			}
		case "i32":
			htype.Integer = &dat.IntergerType{
				Unsigned: false,
				Size:     4,
			}
		case "enumrow":
			htype.Integer = &dat.IntergerType{
				Unsigned: false,
				Size:     4,
			}
		case "f32":
			htype.Decimal = &dat.DecimalType{
				Size: 4,
			}
		case "string":
			htype.String = true
		case "bool":
			htype.Boolean = true
		case "row":
			htype.Key = &dat.KeyType{}
		case "foreignrow":
			htype.Key = &dat.KeyType{
				Foreign: true,
			}
		default:
			log.Fatalf("unhandled header type in schema: %v", column.Type)
		}
		header.Type = htype

		headers = append(headers, &header)

		offset += header.HeaderLength(datFile.FieldSize)
	}

	return headers
}

func ExportAllRows(headers []*dat.Header, datFile *dat.DatFile) [][]dat.FieldValue {
	rows := [][]dat.FieldValue{}
	for i := 0; i < datFile.RowCount; i++ {
		rows = append(rows, dat.ReadRow(i, headers, datFile))
	}

	return rows
}
