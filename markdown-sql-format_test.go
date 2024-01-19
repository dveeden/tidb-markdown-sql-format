package main

import (
	"log/slog"
	"os"
	"testing"

	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
	"github.com/stretchr/testify/require"
)

func Test_formatStmt(t *testing.T) {
	logHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger = slog.New(logHandler)

	tests := []struct {
		name  string
		stmts string
		want  string
	}{
		{
			"simple select",
			`select 1;`,
			`SELECT 1;`},
		{
			"mysql prompt",
			`mysql> select 2;`,
			`mysql> SELECT 2;`},
		{
			"multi statement",
			"select 1;\nselect 2;\nselect 3;\n",
			"SELECT 1;\nSELECT 2;\nSELECT 3;",
		},
		{
			"with resultset",
			`mysql> SelECT * FROM t;
+----+---+
| id | c |
+----+---+
| 1  | 1 |
| 2  | 2 |
| 3  | 3 |
| 4  | 4 |
| 5  | 5 |
+----+---+
5 rows in set (0.01 sec)`,
			`mysql> SELECT * FROM ` + "`t`" + `;
+----+---+
| id | c |
+----+---+
| 1  | 1 |
| 2  | 2 |
| 3  | 3 |
| 4  | 4 |
| 5  | 5 |
+----+---+
5 rows in set (0.01 sec)`,
		},
		{"not sql", "foobar", "foobar"},
		{"not sql with semicolon", "foo;\nbar", "foo;\nbar"},
		{"not sql with semicolon no newline", "foo;bar", "foo;bar"},
		{
			"start with resultset",
			`Query OK, 0 rows affected (0.11 sec)`,
			`Query OK, 0 rows affected (0.11 sec)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatStmt(tt.stmts); got != tt.want {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
