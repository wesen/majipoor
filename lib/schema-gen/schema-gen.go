package schema_gen

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
)

type ProtoColumn interface {
	GetName() string
	Instance() *Column
}

type Column struct {
	Name              string
	HasLength         bool
	HasPrecisionScale bool

	Nullable     bool
	Length       int64
	DatabaseType string
	Precision    int64
	Scale        int64

	DefaultValue interface{}
	EnumValues   []string

	ProtoColumn ProtoColumn
}

type Index struct {
	Name        string
	ColumnNames []string
	IsUnique    bool
}

type Table struct {
	Name    string
	Columns []*Column
	Indexes []*Index

	PrimaryKey string
}

var tableNames = []string{
	"accounts",
	"customers",
	"items",
	"values",
	"orders", "stores", "order_items",
	"widgets", "categories", "keys", "objects", "tags",
	"roles", "permissions", "posts", "post_metadata",
	"metadata", "entries", "logs", "log_entries",
	"capabilities", "jobs", "queues", "job_logs",
	"queue_items",
}

type numericColumn struct {
	name string
}

func (c *numericColumn) GetName() string {
	return c.name
}

func randomString(values []string) string {
	return values[rand.Intn(len(values))]
}

func (c *numericColumn) Instance() *Column {
	isNullable := rand.Intn(2) == 0
	hasDefaultValue := rand.Intn(3) == 0

	switch rand.Intn(3) {
	case 0:
		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: randomString([]string{"int", "bigint", "smallint", "tinyint"}),
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			col.DefaultValue = rand.Intn(100)
		}
		return col

	case 1:
		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: randomString([]string{"decimal", "numeric"}),
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			col.DefaultValue = rand.Float64()
		}
		return col

	case 2:
		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: randomString([]string{"float", "double", "real"}),
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			col.DefaultValue = rand.Float64()
		}
		return col
	}

	return nil
}

type dateColumn struct {
	name string
}

func (c *dateColumn) GetName() string {
	return c.name
}

func randInRange(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (c *dateColumn) Instance() *Column {
	defaultYear := fmt.Sprintf("%04d-%02d-%02d", randInRange(1900, 2100), randInRange(1, 12), randInRange(1, 28))
	defaultTime := fmt.Sprintf("%02d:%02d:%02d", randInRange(0, 23), randInRange(0, 59), randInRange(0, 59))
	defaultFraction := fmt.Sprintf("%04d", randInRange(0, 9999))
	isNullable := rand.Intn(2) == 0
	hasDefaultValue := rand.Intn(3) == 0

	switch rand.Intn(4) {
	case 0:
		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: "date",
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			col.DefaultValue = defaultYear
		}
		return col

	case 1:
		isDefaultCurrent := rand.Intn(2) == 0
		isUpdateCurrent := rand.Intn(2) == 0

		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: randomString([]string{"timestamp", "datetime"}),
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			if isDefaultCurrent {
				col.DefaultValue = "current_timestamp"
				if isUpdateCurrent {
					col.DefaultValue = "current_timestamp on update current_timestamp"
				}
			} else {
				col.DefaultValue = fmt.Sprintf("%s %s.%s", defaultYear, defaultTime, defaultFraction)
			}
		}
		return col

	case 2:
		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: "time",
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			col.DefaultValue = fmt.Sprintf("%03d:%02d:%02d", randInRange(-800, 800), randInRange(0, 59), randInRange(0, 59))
		}
		return col

	case 3:
		col := &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: "year",
			ProtoColumn:  c,
		}
		if hasDefaultValue {
			col.DefaultValue = randInRange(1900, 2100)
		}
		return col
	}

	return nil
}

type textColumn struct {
	name     string
	valueGen func() string
}

func newTextColumn(name string) *textColumn {
	return &textColumn{
		name:     name,
		valueGen: gofakeit.Noun,
	}
}

func (c *textColumn) GetName() string {
	return c.name
}

func (c *textColumn) Instance() *Column {
	isNullable := rand.Intn(2) == 0
	hasDefaultValue := rand.Intn(3) == 0
	useCharacterSet := rand.Intn(3) == 0
	characterSet := randomString([]string{"utf8", "utf8mb4", "latin1", "ascii"})

	defaultValue := c.valueGen()

	var col *Column

	switch rand.Intn(2) {
	case 0:
		col = &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: "text",
			ProtoColumn:  c,
		}
		hasDefaultValue = false

	case 1:
		precision := randInRange(5, 30)
		col = &Column{
			Name:     c.name,
			Nullable: isNullable,
			DatabaseType: fmt.Sprintf("%s(%d)",
				randomString([]string{"char", "varchar", "binary", "varbinary"}), precision),
			ProtoColumn: c,
		}
	}

	if hasDefaultValue {
		col.DefaultValue = defaultValue
	}
	if useCharacterSet {
		col.DatabaseType += " character set " + characterSet
	}

	return col
}

func makeParagraph() string {
	return gofakeit.Paragraph(randInRange(1, 4), randInRange(2, 10), randInRange(5, 20), "")
}

type blobColumn struct {
	name string
}

func (b *blobColumn) GetName() string {
	return b.name
}

func (b *blobColumn) Instance() *Column {
	return &Column{
		Name:         b.name,
		Nullable:     rand.Intn(2) == 0,
		DatabaseType: "blob",
		ProtoColumn:  b,
	}
}

var protoColumns = []ProtoColumn{
	&numericColumn{"count"},
	&numericColumn{"amount"},
	&numericColumn{"quantity"},
	&numericColumn{"price"},
	&numericColumn{"value"},
	&numericColumn{"age"},

	&dateColumn{"created_at"},
	&dateColumn{"updated_at"},
	&dateColumn{"deleted_at"},
	&dateColumn{"start_time"},
	&dateColumn{"end_time"},

	newTextColumn("text"),
	newTextColumn("description"),
	newTextColumn("title"),
	newTextColumn("paragraph"),
	newTextColumn("content"),
	newTextColumn("body"),

	&textColumn{"name", gofakeit.Name},
	&textColumn{"description", makeParagraph},
	&textColumn{"title", gofakeit.SentenceSimple},
	&textColumn{"paragraph", makeParagraph},
	&textColumn{"address", gofakeit.Street},
	&textColumn{"body", makeParagraph},
	&textColumn{"phone", gofakeit.Phone},

	&blobColumn{"data"},
	&blobColumn{"binary_data"},
	&blobColumn{"binary_contents"},
}

func GenerateTable() *Table {
	hasId := rand.Intn(2) == 0

	usedColumns := make(map[string]bool)
	for _, column := range protoColumns {
		usedColumns[column.GetName()] = false
	}

	nColumns := randInRange(1, len(protoColumns))
	columns := []*Column{}
	if hasId {
		columns = append(columns, &Column{
			Name:         "id",
			Nullable:     false,
			DatabaseType: "int",
			DefaultValue: "auto_increment",
			ProtoColumn:  &numericColumn{"id"},
		})
		usedColumns["id"] = true

	}

	for i := 0; i < nColumns; i++ {
		column := protoColumns[rand.Intn(len(protoColumns))]
		if usedColumns[column.GetName()] {
			i--
			continue
		}
		instance := column.Instance()
		if instance != nil {
			columns = append(columns, column.Instance())
			usedColumns[column.GetName()] = true
		}
	}

	return &Table{
		Name:       randomString(tableNames),
		Columns:    columns,
		PrimaryKey: "id",
	}
}
