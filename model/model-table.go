package model

type Table struct {
	Name       string       `yaml:"name,omitempty"`
	Title      string       `yaml:"title,omitempty"`
	Desc       string       `yaml:"desc,omitempty"`
	Version    bool         `yaml:"version,omitempty"`
	Columns    []*Column    `yaml:"columns,omitempty"`
	Indexes    []*Index     `yaml:"indexes,omitempty"`
	References []*Reference `yaml:"references,omitempty"`
}

// func (t *Table) MarshalYAML() (any, error) {
// 	content := getContent(*t, nil, nil)
// 	node := &yaml.Node{
// 		Kind:    yaml.MappingNode,
// 		Content: content,
// 	}
// 	return node, nil
// }

// PrimaryKeyColumns returns all columns that contain identity or are numeric, sorted by identity.
// func (s *Table) PrimaryKeyColumns() []*Column {
// 	// get all columns contains identity
// 	cols := lo.Filter(s.Columns, func(c *Column, _ int) bool {
// 		return c.Identity == "Y" || isNumeric(c.Identity)
// 	})
// 	// sort the the identity columns
// 	sort.Slice(cols, func(i, j int) bool {
// 		return toInt(cols[i].Identity) < toInt(cols[j].Identity)
// 	})
// 	return cols
// }
