package xfmr

type InDB struct {
	Fixed   []Column `yaml:"fixed,omitempty"`
	Schemas []Schema `yaml:"schemas,omitempty"`
}

type Schema struct {
	Name   string  `yaml:"name,omitempty" default:"Schema"`
	Desc   string  `yaml:"description,omitempty"`
	Tables []Table `yaml:"tables,omitempty"`
}

type Table struct {
	Name       string      `yaml:"name,omitempty"`
	Title      string      `yaml:"title,omitempty"`
	Desc       string      `yaml:"desc,omitempty"`
	Columns    []Column    `yaml:"columns,omitempty"`
	OutColumns []OutColumn `yaml:"out_columns,omitempty"`
}

type Column struct {
	Name        string `yaml:"na,omitempty"`
	DataType    string `yaml:"ty,omitempty"`
	Identity    string `yaml:"id,omitempty"`
	NotNull     string `yaml:"nu,omitempty" default:"N"`
	Unique      string `yaml:"un,omitempty"`
	Value       string `yaml:"va,omitempty"`
	ForeignKey  string `yaml:"fk,omitempty"`
	Cardinality string `yaml:"cd,omitempty"`
	Title       string `yaml:"tt,omitempty"`
	Index       string `yaml:"in,omitempty"`
	Desc        string `yaml:"dc,omitempty"`
}

type TextArgs struct {
	InFile       string
	OutFile      string
	TemplateFile string
	Pattern      string
}

// ExcelArgs uses to save the cli arguments of Excel
type ExcelArgs struct {
	InFile  string
	OutFile string
	Simple  bool
}

// DiagramArgs uses to save the cli arguments of UML
type DiagramArgs struct {
	DigType   string
	InFile    string
	OutFile   string
	JarFile   string
	Schema    string
	Pattern   string
	Simple    bool // simple type, no column except PK & FK if true
	IncludeFK bool // include foreign key name in the arrow line
}

type Xfmr struct {
	Data *InDB
}

func NewXMFR() *Xfmr {
	return &Xfmr{}
}
