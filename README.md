# SQLQ - SQL query builder from composable parts.

## Example 
```go
q := sqlq.New("sqlite", true)
query := q.Select("*") + q.From("a") + q.Where(q.Compare("b","=","c")) + q.Limit(10)
// query = "SELECT * FROM a WHERE b = c LIMIT 10"
```
## Available methods:

```go
Select(columns ...interface{}) string 

SelectDistinct(columns ...interface{}) string 

Count(column interface{}) string 

CountDistinct(column interface{}) string 

Into(newTable string) string 

From(tables ...interface{}) string 

Where(conditions ...string) string 

fConditions(conditions ...string) string 

fOrConditions(conditions ...string) string 

Compare(column interface{}, operator string, value interface{}) string 

Between(column, value1, value2 interface{}) 

NotBetween(column, value1, value2 interface{}) string 

IsNull(column interface{}) string 

IsNotNull(column interface{}) string 

// Like : pattern "%?%" ? mark replaced by value
Like(column interface{}, pattern string, value interface{}) string 

NotLike(column interface{}, pattern string, value interface{}) string 

In(column interface{}, values ...interface{}) string 

NotIn(column interface{}, values ...interface{}) string 

OrderBy(column interface{}, order string) string 

GroupBy(columns ...interface{}) string 

Limit(value int) string 

Offset(value int) string 

Paginate(page, perPage int) string 

Join(table, table1Column, table2Column interface{}) string 

InnerJoin(table, table1Column, table2Column interface{}) string 

LeftJoin(table, table1Column, table2Column interface{}) string 

LeftOuterJoin(table, table1Column, table2Column interface{}) string 

RightJoin(table, table1Column, table2Column interface{}) string 

RightOuterJoin(table, table1Column, table2Column interface{}) string 

FullJoin(table, table1Column, table2Column interface{}) string 

FullOuterJoin(table, table1Column, table2Column interface{}) string 

fUnion(selectStatements ...string) string 

fUnionAll(selectStatements ...string) string 

InsertInto(table interface{}, columns Map) string 

Update(table interface{}, columns Map) string 

Delete(table interface{}) string 
```

