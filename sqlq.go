package sqlq

import "strings"

type Sqlq struct {
	Dialect        string
	Escape         bool
	skipEscapeOnce bool
}

func New(dialect string, escape bool) *Sqlq {
	if d := strings.ToLower(dialect); d != "mysql" && d != "postgres" && d != "sqlite3" {
		dialect = "mysql"
	}
	return &Sqlq{Dialect: dialect, Escape: escape, skipEscapeOnce: false}
}

func (sqlq *Sqlq) StopEscape() {
	sqlq.Escape = false
}

func (sqlq *Sqlq) StartEscape() {
	sqlq.Escape = true
}

func (sqlq *Sqlq) SkipEscapeOnce() {
	sqlq.skipEscapeOnce = true
}

func (sqlq *Sqlq) Raw() *Sqlq {
	sqlq.skipEscapeOnce = true
	return sqlq
}

type query struct {
	queryA []string
}

func (q *query) String() string {
	return strings.Join(q.queryA, " ")
}

func (q *query) Append(queryA ...*query) *query {
	for _, p := range queryA {
		q.queryA = append(q.queryA, p.queryA...)
	}
	return q
}

func (q *query) Query(qParts ...string) *query {
	q.queryA = append(q.queryA, qParts...)
	return q
}

type Condition struct {
	res bool
	t   string
}

func (sqlq *Sqlq) If(condition bool, s ...string) *Condition {
	return &Condition{res: condition, t: sqlq.Query(s...).String()}
}

func (c *Condition) ElseIf(condition bool, s ...string) *Condition {
	if !c.res {
		c.res = condition
		c.t = Sqlq{}.Query(s...).String()
	}
	return c
}

func (c *Condition) Else(s ...string) string {
	if c.res {
		return c.t
	}
	return Sqlq{}.Query(s...).String()
}

func (Sqlq) Query(qParts ...string) *query {
	var q query
	q.Query(qParts...)
	return &q
}

func (Sqlq) Concat(queryA ...*query) *query {
	var newquery query
	for _, query := range queryA {
		newquery.queryA = append(newquery.queryA, query.queryA...)
	}
	return &newquery
}
