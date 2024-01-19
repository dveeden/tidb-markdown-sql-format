package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/charset"
	"github.com/pingcap/tidb/pkg/parser/format"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

var logger *slog.Logger

func formatStmt(stmts string) string {
	var sb strings.Builder
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)
	p := parser.New()

	logger.Debug("parsing statements", "statements", stmts)

	if strings.HasPrefix(strings.TrimSpace(stmts), "Query OK") {
		logger.Warn("SQL code starts with resultset instead of statement", "statements", stmts)
		return strings.TrimRight(stmts, "\n")
	}

	splitStmts := strings.Split(stmts, ";")
	nrSplitStmts := len(splitStmts)
	for i, stmt := range splitStmts {
		if len(strings.TrimSpace(stmt)) == 0 {
			continue
		}

		logger.Debug("parsing statement", "statement", stmt)

		// Tabular resultset
		if strings.HasPrefix(stmt, "\n+") {
			sb.WriteString(stmt[1:])
			continue
		}

		if strings.HasPrefix(stmt, "mysql> ") {
			sb.WriteString("mysql> ")
			stmt = stmt[7:]
		}

		parsedStmt, err := p.ParseOneStmt(stmt+";", charset.CharsetUTF8MB4, charset.CollationUTF8MB4)
		if err != nil {
			logger.Warn("failed to parse statement", "statement", stmt, "statement number", i, "error", err)
			if i < nrSplitStmts-1 {
				sb.WriteString(stmt + ";")
			} else {
				sb.WriteString(stmt)
			}
		} else {
			err = parsedStmt.Restore(ctx)
			if err != nil {
				logger.Warn("failed to restore statement", "statement", stmt, "error", err)
				sb.WriteString(stmt)
			} else if i < nrSplitStmts {
				sb.WriteString(";\n")
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

func main() {
	logHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger = slog.New(logHandler)

	var markdownFileName = flag.String("filename", "", "markdown filename")
	flag.Parse()

	codeStart := regexp.MustCompile("^\\s*```sql")
	codeEnd := regexp.MustCompile("^\\s*```")
	insideCode := false

	if *markdownFileName == "" {
		logger.Error("Filename is required.")
		os.Exit(1)
	}
	logger.Info("parsing file", "filename", *markdownFileName)

	fileHandle, err := os.Open(*markdownFileName)
	if err != nil {
		logger.Error("failed to open markdown file", "error", err)
		os.Exit(1)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	sqlStmt := ""
	for scanner.Scan() {
		line := scanner.Text()
		if !insideCode && codeStart.FindString(line) != "" {
			insideCode = true
			fmt.Println(line)
		} else if insideCode && codeEnd.FindString(line) != "" {
			insideCode = false
			fmt.Println(formatStmt(sqlStmt))
			sqlStmt = ""
			fmt.Println(line)
		} else if insideCode {
			sqlStmt += line + "\n"
		} else {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("scanner failed", "error", err)
		os.Exit(1)
	}
}
