package utils

type Table struct {
	Header []string
	Rows   [][]string
}

func NewTable(header []string) *Table {
	return &Table{
		Header: header,
	}
}

func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
}

func (t *Table) AddRows(rows [][]string) {
	t.Rows = append(t.Rows, rows...)
}

func (t *Table) Render() string {
	out := ""
	out += "|"
	for _, h := range t.Header {
		out += h + "|"
	}
	out += "\n"

	for i := 0; i < len(t.Header); i++ {
		out += "|---"
	}
	out += "|\n"

	for _, row := range t.Rows {
		out += "|"
		for _, cell := range row {
			out += cell + "|"
		}
		out += "\n"
	}

	return out
}
