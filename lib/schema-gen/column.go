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

// TODO(manuel) Add code to serialize the schema generator data to disk, so it can be reloaded for more generation
type Column struct {
	Name              string
	HasLength         bool
	HasPrecisionScale bool

	Nullable     bool
	Length       int64
	DatabaseType string
	Precision    int64
	Scale        int64

	DefaultValue string
	EnumValues   []string

	ProtoColumn ProtoColumn
}

func (c *Column) ColumnDefinition() string {
	s := fmt.Sprintf("%s %s", c.Name, c.DatabaseType)

	if !c.Nullable {
		s += " NOT NULL"
	}
	if c.DefaultValue != "" {
		s += " " + c.DefaultValue
	}

	return s
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
			col.DefaultValue = fmt.Sprintf("DEFAULT %d", rand.Intn(100))
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
			col.DefaultValue = fmt.Sprintf("DEFAULT %f", rand.Float64())
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
			col.DefaultValue = fmt.Sprintf("DEFAULT %f", rand.Float64())
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
	defaultYear := fmt.Sprintf("%04d-%02d-%02d", randInRange(1900, 2037), randInRange(1, 12), randInRange(1, 28))
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
			col.DefaultValue = fmt.Sprintf("DEFAULT '%s'", defaultYear)
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
				col.DefaultValue = fmt.Sprintf("DEFAULT '%s %s.%s'", defaultYear, defaultTime, defaultFraction)
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
			col.DefaultValue = fmt.Sprintf("DEFAULT '%03d:%02d:%02d'", randInRange(-800, 800), randInRange(0, 59), randInRange(0, 59))
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
			col.DefaultValue = fmt.Sprintf("DEFAULT %d", randInRange(1900, 2100))
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

	switch rand.Intn(3) {
	case 0:
		col = &Column{
			Name:         c.name,
			Nullable:     isNullable,
			DatabaseType: "text",
			ProtoColumn:  c,
		}
		hasDefaultValue = false

	case 1:
		length := randInRange(20, 100)
		col = &Column{
			Name:     c.name,
			Nullable: isNullable,
			Length:   int64(length),
			DatabaseType: fmt.Sprintf("%s(%d)",
				randomString([]string{"char", "varchar"}), length),
			ProtoColumn: c,
		}
		if len(defaultValue) > length {
			defaultValue = defaultValue[:length]
		}

	case 2:
		length := randInRange(5, 30)
		col = &Column{
			Name:     c.name,
			Nullable: isNullable,
			Length:   int64(length),
			DatabaseType: fmt.Sprintf("%s(%d)",
				randomString([]string{"binary", "varbinary"}), length),
			ProtoColumn: c,
		}
		useCharacterSet = false
		if len(defaultValue) > length {
			defaultValue = defaultValue[:length]
		}
	}

	if hasDefaultValue {
		col.DefaultValue = fmt.Sprintf("DEFAULT '%s'", defaultValue)
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
