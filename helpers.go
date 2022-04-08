package sqlq

import (
	"fmt"
	"strings"
)

func cleanStringArray(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

type Map map[string]interface{}

func (sqlq *Sqlq) Select(columns ...interface{}) string {
	if len(columns) == 1 && fmt.Sprint(columns[0]) == "*" {
		return "SELECT *"
	}
	return "SELECT " + sqlq.I(columns, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) SelectDistinct(columns ...interface{}) string {
	if len(columns) == 1 && fmt.Sprint(columns[0]) == "*" {
		return "SELECT DISTINCT *"
	}
	return "SELECT DISTINCT " + sqlq.I(columns, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) Count(column interface{}) string {
	if fmt.Sprint(column) == "*" {
		return "SELECT COUNT(*)"
	}
	return "SELECT COUNT(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) CountDistinct(column interface{}) string {
	if fmt.Sprint(column) == "*" {
		return "SELECT COUNT(DISTINCT *)"
	}
	return "SELECT COUNT(DISTINCT " + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) Into(newTable string) string {
	return "INTO " + sqlq.I(newTable, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) From(tables ...interface{}) string {
	return "FROM " + sqlq.I(tables, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) Where(conditions ...string) string {
	return "WHERE " + sqlq.Conditions(conditions...)
}

func (Sqlq) Conditions(conditions ...string) string {
	conditions = cleanStringArray(conditions)
	if len(conditions) == 0 {
		return ""
	}
	if len(conditions) == 1 {
		return conditions[0]
	}
	return "(" + strings.Join(conditions, " AND ") + ")"
}

func (Sqlq) OrConditions(conditions ...string) string {
	conditions = cleanStringArray(conditions)
	if len(conditions) == 0 {
		return ""
	}
	if len(conditions) == 1 {
		return conditions[0]
	}
	return "(" + strings.Join(conditions, " OR ") + ")"
}

func (sqlq *Sqlq) Compare(column interface{}, operator string, value interface{}) string {
	if operator != "<" && operator != ">" && operator != "<=" && operator != ">=" && operator != "=" && operator != "<>" && operator != "!=" {
		operator = "="
	}
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " " + operator + " " + sqlq.L(value, !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) Between(column, value1, value2 interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " BETWEEN " + sqlq.L(value1, !sqlq.Escape, sqlq.skipEscapeOnce) + " AND " + sqlq.L(value2, !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) NotBetween(column, value1, value2 interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " NOT BETWEEN " + sqlq.L(value1, !sqlq.Escape, sqlq.skipEscapeOnce) + " AND " + sqlq.L(value2, !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) IsNull(column interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " IS NULL" + ")"
}

func (sqlq *Sqlq) IsNotNull(column interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " IS NOT NULL" + ")"
}

// Like : pattern "%?%" ? mark replaced by value
func (sqlq *Sqlq) Like(column interface{}, pattern string, value interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " LIKE " + sqlq.L(strings.Replace(pattern, "?", fmt.Sprint(value), 1), !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) NotLike(column interface{}, pattern string, value interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + "NOT LIKE " + sqlq.L(strings.Replace(pattern, "?", fmt.Sprint(value), 1), !sqlq.Escape, sqlq.skipEscapeOnce) + ")"
}

func (sqlq *Sqlq) In(column interface{}, values ...interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " IN " + "(" + sqlq.L(values, !sqlq.Escape, sqlq.skipEscapeOnce) + ") )"
}

func (sqlq *Sqlq) NotIn(column interface{}, values ...interface{}) string {
	return "(" + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " NOT IN " + "(" + sqlq.L(values, !sqlq.Escape, sqlq.skipEscapeOnce) + ") )"
}

func (sqlq *Sqlq) OrderBy(column interface{}, order string) string {
	o := sqlq.O(order, !sqlq.Escape, sqlq.skipEscapeOnce)
	if o == "" {
		return "ORDER BY " + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce)
	}
	return "ORDER BY " + sqlq.I(column, !sqlq.Escape, sqlq.skipEscapeOnce) + " " + o
}

func (sqlq *Sqlq) GroupBy(columns ...interface{}) string {
	return "GROUP BY " + sqlq.I(columns, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) Limit(value int) string {
	return "LIMIT " + sqlq.L(value, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) Offset(value int) string {
	return "OFFSET " + sqlq.L(value, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) Paginate(page, perPage int) string {
	return sqlq.Limit(perPage) + " " + sqlq.Offset((page-1)*perPage)
}

func (sqlq *Sqlq) Join(table, table1Column, table2Column interface{}) string {
	return "JOIN " + sqlq.I(table, !sqlq.Escape, sqlq.skipEscapeOnce) + " ON " + sqlq.I(table1Column, !sqlq.Escape, sqlq.skipEscapeOnce) + "=" + sqlq.I(table2Column, !sqlq.Escape, sqlq.skipEscapeOnce)
}

func (sqlq *Sqlq) InnerJoin(table, table1Column, table2Column interface{}) string {
	return "INNER " + sqlq.Join(table, table1Column, table2Column)
}

func (sqlq *Sqlq) LeftJoin(table, table1Column, table2Column interface{}) string {
	return "LEFT " + sqlq.Join(table, table1Column, table2Column)
}

func (sqlq *Sqlq) LeftOuterJoin(table, table1Column, table2Column interface{}) string {
	return "LEFT OUTER " + sqlq.Join(table, table1Column, table2Column)
}

func (sqlq *Sqlq) RightJoin(table, table1Column, table2Column interface{}) string {
	return "RIGHT " + sqlq.Join(table, table1Column, table2Column)
}

func (sqlq *Sqlq) RightOuterJoin(table, table1Column, table2Column interface{}) string {
	return "RIGHT OUTER " + sqlq.Join(table, table1Column, table2Column)
}

func (sqlq *Sqlq) FullJoin(table, table1Column, table2Column interface{}) string {
	return "FULL " + sqlq.Join(table, table1Column, table2Column)
}

func (sqlq *Sqlq) FullOuterJoin(table, table1Column, table2Column interface{}) string {
	return "FULL OUTER " + sqlq.Join(table, table1Column, table2Column)
}

func (Sqlq) Union(selectStatements ...string) string {
	return strings.Join(selectStatements, " UNION ")
}

func (Sqlq) UnionAll(selectStatements ...string) string {
	return strings.Join(selectStatements, " UNION ALL ")
}

func (sqlq *Sqlq) InsertInto(table interface{}, columns Map) string {
	cols := []string{}
	values := []string{}
	for key, value := range columns {
		cols = append(cols, sqlq.I(key, !sqlq.Escape, sqlq.skipEscapeOnce))
		values = append(values, sqlq.L(value, !sqlq.Escape, sqlq.skipEscapeOnce))
	}
	return "INSERT INTO " + sqlq.I(table, !sqlq.Escape, sqlq.skipEscapeOnce) + " (" + strings.Join(cols, ", ") + ")" + " VALUES " + "(" + strings.Join(values, ", ") + ")"
}

func (sqlq *Sqlq) Update(table interface{}, columns Map) string {
	cols := []string{}
	for key, value := range columns {
		cols = append(cols, sqlq.I(key, !sqlq.Escape, sqlq.skipEscapeOnce)+"="+sqlq.L(value, !sqlq.Escape, sqlq.skipEscapeOnce))
	}
	return "UPDATE " + sqlq.I(table, !sqlq.Escape, sqlq.skipEscapeOnce) + " SET " + strings.Join(cols, ", ")
}

func (sqlq *Sqlq) Delete(table interface{}) string {
	return "DELETE FROM " + sqlq.I(table, !sqlq.Escape, sqlq.skipEscapeOnce)
}
