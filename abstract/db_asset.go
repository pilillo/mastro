package abstract

// types used in parsing

// DBInfo ... Name and description for a database
type DBInfo struct {
	Name    string
	Comment string
}

// ColumnInfo ... Type and description for a table column
type ColumnInfo struct {
	Type    string
	Comment string
}

// TableInfo ... Name, schema and description for a table
type TableInfo struct {
	Name    string
	Schema  map[string]ColumnInfo
	Comment string
}

// structs shall be unexported by default
type databaseBuilder struct{ asset Asset }

// NewDatabaseBuilder ... builder for a database asset type
func NewDatabaseBuilder() *databaseBuilder {
	return &databaseBuilder{}
}

func (b *databaseBuilder) SetName(name string) *databaseBuilder {
	b.asset.Name = name
	return b
}

func (b *databaseBuilder) SetDescription(description string) *databaseBuilder {
	b.asset.Description = description
	return b
}

func (b *databaseBuilder) Build() (*Asset, error) {
	if err := b.asset.Validate(); err != nil {
		return nil, err
	}
	return &b.asset, nil
}

func (db *DBInfo) BuildAsset() (*Asset, error) {
	return NewDatabaseBuilder().SetName(db.Name).SetDescription(db.Comment).Build()
}

type tableBuilder struct{ asset Asset }

// NewTableBuilder ... table builder
func NewTableBuilder() *tableBuilder {
	return &tableBuilder{}
}

func (b *tableBuilder) SetName(name string) *tableBuilder {
	b.asset.Name = name
	return b
}

func (b *tableBuilder) SetDescription(description string) *tableBuilder {
	b.asset.Description = description
	return b
}

func (b *tableBuilder) SetSchema(schema map[string]ColumnInfo) *tableBuilder {
	b.asset.Labels[L_SCHEMA] = schema
	return b
}

func (b *tableBuilder) Build() (*Asset, error) {
	if err := b.asset.Validate(); err != nil {
		return nil, err
	}
	return &b.asset, nil
}

func (tb *TableInfo) BuildAsset() (*Asset, error) {
	return NewTableBuilder().SetName(tb.Name).SetDescription(tb.Comment).SetSchema(tb.Schema).Build()
}