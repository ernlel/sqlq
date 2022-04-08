package sqlq

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func quoteIdent(s, quote string) string {
	part := strings.SplitN(s, ".", 2)
	if len(part) == 2 {
		if part[1] == "*" {
			return quoteIdent(part[0], quote) + "." + part[1]
		} else {
			return quoteIdent(part[0], quote) + "." + quoteIdent(part[1], quote)
		}
	}
	return quote + s + quote
}

func interfaceArryToString(i interface{}) string {
	if v, ok := i.([]interface{}); ok {
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, fmt.Sprint(vi))
		}
		return strings.Join(sArr, ", ")
	}
	return fmt.Sprint(i)
}

// L any type Literal.
func (sqlq *Sqlq) L(i interface{}, skipEscape ...bool) string {
	if len(skipEscape) > 0 && skipEscape[0] {
		return interfaceArryToString(i)
	} else if len(skipEscape) > 1 && skipEscape[1] {
		sqlq.skipEscapeOnce = false
		return interfaceArryToString(i)
	}
	switch sqlq.Dialect {
	case "mysql":
		return mysqlL(i)
	case "postgres":
		return postgresL(i)
	case "sqlite3":
		return sqlite3L(i)
	}
	return mysqlL(i)
}

// ------------ MYSQL -----------

// mysqlL any type Literal.
func mysqlL(i interface{}) string {
	switch v := i.(type) {
	case string:
		return mysqlStringEscape(v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case []interface{}:
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, mysqlL(vi))
		}
		return strings.Join(sArr, ", ")
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return fmt.Sprint(v)
	case time.Time:
		return `'` + v.UTC().Format("2006-01-02 15:04:05.000000") + `'`
	case []byte:
		return fmt.Sprintf(`0x%x`, v)
	case nil:
		return "NULL"
	}
	return mysqlStringEscape(fmt.Sprint(i))
}

// mysql String Escape
func mysqlStringEscape(s string) string {
	buf := new(bytes.Buffer)

	buf.WriteRune('\'')
	// https://dev.mysql.com/doc/refman/5.7/en/string-literals.html
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 0:
			buf.WriteString(`\0`)
		case '\'':
			buf.WriteString(`\'`)
		case '"':
			buf.WriteString(`\"`)
		case '\b':
			buf.WriteString(`\b`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case 26:
			buf.WriteString(`\Z`)
		case '\\':
			buf.WriteString(`\\`)
		default:
			buf.WriteByte(s[i])
		}
	}
	buf.WriteRune('\'')
	return buf.String()
}

// ------------ POSTGRES -----------

// postgresL any type Literal.
func postgresL(i interface{}) string {
	switch v := i.(type) {
	case string:
		return postgresStringEscape(v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case []interface{}:
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, postgresL(vi))
		}
		return strings.Join(sArr, ", ")
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return fmt.Sprint(v)
	case time.Time:
		return `'` + v.UTC().Format("2006-01-02 15:04:05.000000") + `'`
	case []byte:
		return fmt.Sprintf(`E'\\x%x'`, v)
	case nil:
		return "NULL"
	}
	return postgresStringEscape(fmt.Sprint(i))
}

// postgres String Escape
func postgresStringEscape(s string) string {
	// http://www.postgresql.org/docs/9.2/static/sql-syntax-lexical.html
	return `'` + strings.Replace(s, `'`, `''`, -1) + `'`
}

// ------------ SQLITE3 -----------

// sqlite3L any type Literal.
func sqlite3L(i interface{}) string {
	switch v := i.(type) {
	case string:
		return sqlite3StringEscape(v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case []interface{}:
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, sqlite3L(vi))
		}
		return strings.Join(sArr, ", ")
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return fmt.Sprint(v)
	case time.Time:
		// https://www.sqlite.org/lang_datefunc.html
		return `'` + v.UTC().Format("2006-01-02 15:04:05.000000") + `'`
	case []byte:
		// https://www.sqlite.org/lang_expr.html
		return fmt.Sprintf(`X'%x'`, v)
	case nil:
		return "NULL"
	}
	return sqlite3StringEscape(fmt.Sprint(i))
}

// sqlite3 String Escape
func sqlite3StringEscape(s string) string {
	// https://www.sqlite.org/faq.html
	return `'` + strings.Replace(s, `'`, `''`, -1) + `'`
}

// I any type identifier
func (sqlq *Sqlq) I(i interface{}, skipEscape ...bool) string {
	if len(skipEscape) > 0 && skipEscape[0] {
		return interfaceArryToString(i)
	} else if len(skipEscape) > 1 && skipEscape[1] {
		sqlq.skipEscapeOnce = false
		return interfaceArryToString(i)
	}
	switch sqlq.Dialect {
	case "mysql":
		return mysqlI(i)
	case "postgres":
		return postgresI(i)
	case "sqlite3":
		return sqlite3I(i)
	}

	return mysqlI(i)
}

// ------------ MYSQL -----------

// mysqlI any type identifier
func mysqlI(i interface{}) string {
	if v, ok := i.([]interface{}); ok {
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, mysqlIdentEscape(fmt.Sprint(vi)))
		}
		return strings.Join(sArr, ", ")
	}
	return mysqlIdentEscape(fmt.Sprint(i))
}

// mysql Ident Escape
func mysqlIdentEscape(s string) string {
	as := regexp.MustCompile(`(?i)(.*) AS (.*)`)
	if as.MatchString(s) {
		asa := as.FindStringSubmatch(s)
		return quoteIdent(asa[1], "`") + " AS " + quoteIdent(asa[2], "`")
	}
	return quoteIdent(s, "`")
}

// ------------ POSTGRES -----------

// postgresI any type identifier
func postgresI(i interface{}) string {
	if v, ok := i.([]interface{}); ok {
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, postgresIdentEscape(fmt.Sprint(vi)))
		}
		return strings.Join(sArr, ", ")
	}
	return postgresIdentEscape(fmt.Sprint(i))
}

// postgres Ident Escape
func postgresIdentEscape(s string) string {
	as := regexp.MustCompile(`(?i)(.*) AS (.*)`)
	if as.MatchString(s) {
		asa := as.FindStringSubmatch(s)
		return quoteIdent(asa[1], `"`) + " AS " + quoteIdent(asa[2], `"`)
	}
	return quoteIdent(s, `"`)
}

// ------------ SQLITE3 -----------

// sqlite3I any type identifier
func sqlite3I(i interface{}) string {
	if v, ok := i.([]interface{}); ok {
		var sArr []string
		for _, vi := range v {
			sArr = append(sArr, sqlite3IdentEscape(fmt.Sprint(vi)))
		}
		return strings.Join(sArr, ", ")
	}
	return sqlite3IdentEscape(fmt.Sprint(i))
}

// sqlite3 Ident Escape
func sqlite3IdentEscape(s string) string {
	as := regexp.MustCompile(`(?i)(.*) AS (.*)`)
	if as.MatchString(s) {
		asa := as.FindStringSubmatch(s)
		return quoteIdent(asa[1], `"`) + " AS " + quoteIdent(asa[2], `"`)
	}
	return quoteIdent(s, `"`)
}

// O order asc or desc.
func (sqlq *Sqlq) O(s string, skipEscape ...bool) string {
	if len(skipEscape) > 0 && skipEscape[0] {
		return s
	} else if len(skipEscape) > 1 && skipEscape[1] {
		sqlq.skipEscapeOnce = false
		return s
	}

	const (
		asc  = "ASC"
		desc = "DESC"
	)
	switch strings.ToUpper(s) {
	case asc:
		return asc
	case desc:
		return desc
	}
	return asc
}
