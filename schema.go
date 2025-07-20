package lunodbgo

import (
	nodepb "github.com/cloudproud/lunodb.api/proto/node"
	typespb "github.com/cloudproud/lunodb.api/proto/types"
)

type Catalog struct {
	UID         uint64
	Namespace   string
	Name        string
	Description string
	Labels      []string
	Tables      Tables
	Hidden      bool
}

func (catalog Catalog) Proto() *nodepb.Catalog {
	return &nodepb.Catalog{
		UID:         catalog.UID,
		Namespace:   catalog.Namespace,
		Name:        catalog.Name,
		Description: catalog.Description,
		Labels:      catalog.Labels,
		Tables:      catalog.Tables.Proto(),
		Hidden:      catalog.Hidden,
	}
}

type Tables []Table

func (tables Tables) Proto() []*nodepb.Table {
	result := make([]*nodepb.Table, len(tables))
	for index, table := range tables {
		result[index] = table.Proto()
	}
	return result
}

type Table struct {
	Name       string
	Schema     string
	Catalog    string
	Schemaless bool
	Columns    Columns
	Operators  Operators
}

func (t Table) Proto() *nodepb.Table {
	return &nodepb.Table{
		Name:       t.Name,
		Schema:     t.Schema,
		Catalog:    t.Catalog,
		Schemaless: t.Schemaless,
		Columns:    t.Columns.Proto(),
		Operators:  t.Operators.Proto(),
	}
}

type Columns []Column

func (columns Columns) Proto() []*nodepb.Column {
	result := make([]*nodepb.Column, len(columns))
	for index, column := range columns {
		result[index] = &nodepb.Column{
			Name:      column.Name,
			Type:      column.Type,
			Indexed:   column.Indexed,
			Nullable:  column.Nullable,
			Operators: column.Operators.Proto(),
		}
	}
	return result
}

type Column struct {
	Name      string
	Type      *typespb.Type
	Required  bool
	Indexed   bool
	Nullable  bool
	Operators Operators
}

type Operators []Operator

func (operators Operators) Proto() []*nodepb.Operator {
	result := make([]*nodepb.Operator, len(operators))
	for index, operator := range operators {
		result[index] = &nodepb.Operator{
			Statement:       nodepb.OperatorStatement(operator.Statement),
			ComparisonTypes: operator.ComparisonTypes.Proto(),
			Required:        operator.Required,
		}
	}
	return result
}

type Operator struct {
	Statement       Statement
	ComparisonTypes ComparisonTypes
	Required        bool
}

type ComparisonTypes []ComparisonType

func (types ComparisonTypes) Proto() []nodepb.ComparisonType {
	result := make([]nodepb.ComparisonType, len(types))
	for index, typ := range types {
		result[index] = nodepb.ComparisonType(typ)
	}
	return result
}

type ComparisonType int32

const (
	VariableConstant ComparisonType = 1
	VariableVariable ComparisonType = 2
)

type Statement int32

const (
	StatementWhere              Statement = 1
	StatementLimit              Statement = 2
	StatementOffset             Statement = 3
	StatementOrder              Statement = 4
	StatementLeftJoin           Statement = 5
	StatementRightJoin          Statement = 6
	StatementInnerJoin          Statement = 7
	StatementOuterJoin          Statement = 8
	StatementEqual              Statement = 9
	StatementNotEqual           Statement = 10
	StatementIn                 Statement = 11
	StatementNotIn              Statement = 12
	StatementGreaterThan        Statement = 13
	StatementGreaterOrEqualThan Statement = 14
	StatementLessThan           Statement = 15
	StatementLessOrEqualThan    Statement = 16
	StatementLike               Statement = 17
	StatementNotLike            Statement = 18
	StatementILike              Statement = 19
	StatementNotILike           Statement = 20
	StatementRegMatch           Statement = 21
	StatementNotRegMatch        Statement = 22
	StatementRegIMatch          Statement = 23
	StatementNotRegIMatch       Statement = 24
	StatementIsDistinctFrom     Statement = 25
	StatementIsNotDistinctFrom  Statement = 26
	StatementAny                Statement = 27
	StatementAll                Statement = 28
	StatementBinary             Statement = 29
)
