package util

import (
	"fmt"
	"strings"
)

type Condition struct {
	Terms []ConditionTerm
}

type ConditionTerm struct {
	Type  string // AND | OR
	Term  any
	Param any
}

func NewCondition() *Condition {
	return &Condition{}
}

func (c *Condition) And(term any) *Condition {
	c.Terms = append(c.Terms, ConditionTerm{
		Type: "AND",
		Term: term,
	})
	return c
}

func (c *Condition) AndWithParam(term any, param any) *Condition {
	c.Terms = append(c.Terms, ConditionTerm{
		Type:  "AND",
		Term:  term,
		Param: param,
	})
	return c
}

func (c *Condition) Or(term any) *Condition {
	c.Terms = append(c.Terms, ConditionTerm{
		Type: "OR",
		Term: term,
	})
	return c
}

func (c *Condition) OrWithParam(term any, param any) *Condition {
	c.Terms = append(c.Terms, ConditionTerm{
		Type:  "OR",
		Term:  term,
		Param: param,
	})
	return c
}

func (c *Condition) HasCondition() bool {
	return len(c.Terms) != 0
}

func (c *Condition) Build() (string, []any) {
	stat := strings.Builder{}
	parm := []any{}

	for _, term := range c.Terms {
		switch v := term.Term.(type) {
		case string:
			if stat.Len() > 0 {
				stat.WriteString(fmt.Sprintf(" %s ", term.Type))
			}
			stat.WriteString(v)
			if term.Param != nil {
				parm = append(parm, term.Param)
			}
		case *Condition:
			if stat.Len() > 0 {
				stat.WriteString(fmt.Sprintf(" %s ", term.Type))
			}
			stat.WriteString("(")
			childStmt, childParm := v.Build()
			stat.WriteString(childStmt)
			parm = append(parm, childParm...)
			stat.WriteString(")")
		default:
			panic("unknown condition term type")
		}
	}
	return stat.String(), parm
}
