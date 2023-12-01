package schema

// https://github.com/poe-tool-dev/dat-schema/blob/main/src/types.ts

var SchemaVersion = 3

type Ref struct {
	Table  string  `json:"table"`
	Column *string `json:"column"`
}

type TableColumn struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Array       bool    `json:"array"`
	Type        string  `json:"type"`
	Unique      bool    `json:"unique"`
	Localized   bool    `json:"localized"`
	Util        *string `json:"util"`
	References  *Ref    `json:"references"`
}

type SchemaTable struct {
	Name    string         `json:"name"`
	Columns []*TableColumn `json:"columns"`
	Tags    []string       `json:"tags"`
}

type SchemaEnumeration struct {
	Name        string    `json:"name"`
	Indexing    int       `json:"indexing"`
	Enumerators []*string `json:"enumerators"`
}

type SchemaMetadata struct {
	Version   int `json:"version"`
	CreatedAt int `json:"createdAt"`
}

type SchemaFile struct {
	Version      int                  `json:"version"`
	CreatedAt    int                  `json:"createdAt"`
	Tables       []*SchemaTable       `json:"tables"`
	Enumerations []*SchemaEnumeration `json:"enumerations"`
}
