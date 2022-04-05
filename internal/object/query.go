package object

import "fmt"

type Query struct {
	Schema map[string]string // all generated queries
}

type QuerySettings struct {
	Query       string        // prepared query for db
	QueryName   string        // key in constants
	QueryFields []interface{} // for fmt-symbols like '%s' in string

	Fields []interface{} // only for single queryRow
}

func (ss *QuerySettings) MakeQuery(q *Query) *QuerySettings {
	if len(ss.QueryFields) == 0 {
		ss.Query = q.Schema[ss.QueryName]
	} else {
		ss.Query = fmt.Sprintf(q.Schema[ss.QueryName], ss.QueryFields...)
	}
	return ss
}
