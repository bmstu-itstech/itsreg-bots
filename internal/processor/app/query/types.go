package query

import "github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"

type Table struct {
	Head []string
	Body [][]string
}

func mapTableFromDomain(table *bots.Table) Table {
	return Table{
		Head: table.Head,
		Body: table.Body,
	}
}
