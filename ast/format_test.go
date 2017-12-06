package ast_test

import (
	"bytes"
	"fmt"

	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/parser"
)

var _ = Suite(&testAstFormatSuite{})

type testAstFormatSuite struct {
}

func getDefaultCharsetAndCollate() (string, string) {
	return "utf8", "utf8_bin"
}

func (ts *testAstFormatSuite) TestAstFormat(c *C) {
	var testcases = []struct {
		input  string
		output string
	}{
		// Literals.
		{`350`, `350`},
		{`345.678`, `345.678`},
		{`1e-12`, `1e-12`},
		{`null`, `NULL`},
		{`"Hello, world"`, `"Hello, world"`},
		{`'Hello, world'`, `"Hello, world"`},
		{`'Hello, "world"'`, `"Hello, \"world\""`},

		// Expressions.
		{`f between 30 and 50`, "`f` BETWEEN 30 AND 50"},
		{`345 + "  hello  "`, `345 + "  hello  "`},
		{`"hello world"    >=    'hello world'`, `"hello world" >= "hello world"`},
		{`case 3 when 1 then false else true end`, `CASE 3 WHEN 1 THEN FALSE ELSE TRUE END`},
		{`database.table.column`, "`database`.`table`.`column`"}, // ColumnNameExpr
		{`3 is null`, `3 IS NULL`},
		{`3 is not null`, `3 IS NOT NULL`},
		{`3 is true`, `3 IS TRUE`},
		{`3 is not true`, `3 IS NOT TRUE`},
		{`3 is false`, `3 IS FALSE`},
		{`  ( x is false  )`, "(`x` IS FALSE)"},
		{`3 in ( a,b,"h",6 )`, "3 IN (`a`,`b`,\"h\",6)"},
		{`"abc" like '%b%'`, `"abc" LIKE "%b%"`},
		{`"abc" like '%b%' escape '_'`, `"abc" LIKE "%b%" ESCAPE '_'`},
		{`"abc" regexp '.*bc?'`, `"abc" REGEXP ".*bc?"`},
		{`"abc" not regexp '.*bc?'`, `"abc" NOT REGEXP ".*bc?"`},
		{`-  4`, `-4`},
		{`- ( - 4 ) `, `-(-4)`},
		// Functions.
		{` json_extract ( a,'$.b',"$.\"c d\"" ) `, "json_extract(`a`, \"$.b\", \"$.\\\"c d\\\"\")"},
		{` length ( a )`, "length(`a`)"},
		// Cast, Convert and Binary.
		{` cast ( a as signed ) `, "CAST(`a` AS SIGNED)"},
		{` cast ( a as unsigned integer) `, "CAST(`a` AS UNSIGNED)"},
		{` cast ( a as char(3) binary) `, "CAST(`a` AS CHAR(3) BINARY)"},
		{` cast ( a as decimal ) `, "CAST(`a` AS DECIMAL(11))"},
		{` cast ( a as decimal (3) ) `, "CAST(`a` AS DECIMAL(3))"},
		{` cast ( a as decimal (3,3) ) `, "CAST(`a` AS DECIMAL(3, 3))"},
		{` convert (a, signed) `, "CONVERT(`a`, SIGNED)"},
		{` binary "hello"`, `BINARY "hello"`},
	}
	for _, tt := range testcases {
		expr := fmt.Sprintf("select %s", tt.input)
		charset, collation := getDefaultCharsetAndCollate()
		stmts, err := parser.New().Parse(expr, charset, collation)
		node := stmts[0].(*ast.SelectStmt).Fields.Fields[0].Expr
		c.Assert(err, IsNil)

		writer := bytes.NewBufferString("")
		node.Format(writer)
		c.Assert(writer.String(), Equals, tt.output)
	}
}
