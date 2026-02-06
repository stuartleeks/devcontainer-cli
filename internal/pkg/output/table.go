package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Formatter is a function that formats a cell value for a table.
// The first argument is the column index
// The second argument is the width of the column
// The third argument is the value to format
type Formatter func(int, int, interface{}) string
type WidthFunc func(interface{}) int

// These are the default properties for all Tables created from this package
// and can be modified.
var (
	// DefaultPadding specifies the number of spaces between columns in a table.
	DefaultPadding = 2

	// DefaultWriter specifies the output io.Writer for the Table.Print method.
	DefaultWriter io.Writer = os.Stdout

	// DefaultHeaderFormatter specifies the default Formatter for the table header.
	DefaultHeaderFormatter Formatter
)

func DefaultWidthFunc(value interface{}) int {
	return utf8.RuneCountInString(fmt.Sprint(value))
}
func DefaultValueFormatter(columnIndex int, width int, val interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", width), val)
}

type table struct {
	ValueFormatter  Formatter
	HeaderFormatter Formatter
	Padding         int
	Writer          io.Writer
	Width           WidthFunc

	header []string
	rows   [][]interface{}
	widths []int
}

type Table interface {
	WithHeaderFormatter(f Formatter) Table
	WithValueFormatter(f Formatter) Table
	WithPadding(p int) Table
	WithWriter(w io.Writer) Table
	WithWidthFunc(f WidthFunc) Table

	AddRow(vals ...interface{}) Table
	Print() error
	PrintRows() error
}

// New creates a Table instance with the specified header(s) provided. The number
// of columns is fixed at this point to len(columnHeaders) and the defined defaults
// are set on the instance.
func NewTable(columnHeaders ...interface{}) Table {
	t := table{header: make([]string, len(columnHeaders))}

	t.WithPadding(DefaultPadding)
	t.WithWriter(DefaultWriter)
	t.WithHeaderFormatter(DefaultHeaderFormatter)
	t.WithValueFormatter(DefaultValueFormatter)
	t.WithWidthFunc(DefaultWidthFunc)

	for i, col := range columnHeaders {
		t.header[i] = fmt.Sprint(col)
	}

	return &t
}

func (t *table) WithHeaderFormatter(f Formatter) Table {
	t.HeaderFormatter = f
	return t
}

func (t *table) WithValueFormatter(f Formatter) Table {
	t.ValueFormatter = f
	return t
}

func (t *table) WithPadding(p int) Table {
	if p < 0 {
		p = 0
	}

	t.Padding = p
	return t
}

func (t *table) WithWriter(w io.Writer) Table {
	t.Writer = w
	return t
}

func (t *table) WithWidthFunc(f WidthFunc) Table {
	t.Width = f
	return t
}

func (t *table) AddRow(vals ...interface{}) Table {
	t.rows = append(t.rows, vals)
	return t
}

func (t *table) Print() error {
	if t.Writer == nil {
		return fmt.Errorf("writer not set")
	}

	t.calculateWidths(true)
	if err := t.printHeader(); err != nil {
		return err
	}
	return t.printRows()
}
func (t *table) PrintRows() error {
	if t.Writer == nil {
		return nil
	}

	t.calculateWidths(false)
	return t.printRows()
}

func (t *table) calculateWidths(includeHeader bool) {
	t.widths = make([]int, len(t.header))
	if includeHeader {
		for i, header := range t.header {
			t.widths[i] = t.Width(header)
		}
	}

	for _, row := range t.rows {
		for i, cell := range row {
			width := t.Width(cell)
			if width > t.widths[i] {
				t.widths[i] = width
			}
		}
	}
}

func (t *table) printHeader() error {
	if t.HeaderFormatter == nil {
		t.HeaderFormatter = t.ValueFormatter
	}

	for i, header := range t.header {
		formatted := t.HeaderFormatter(i, t.widths[i], header)
		_, err := fmt.Fprint(t.Writer, formatted)
		if err != nil {
			return err
		}
		if i < len(t.header)-1 {
			_, err = fmt.Fprint(t.Writer, strings.Repeat(" ", t.Padding))
			if err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintln(t.Writer)
	return err
}
func (t *table) printRows() error {
	for _, row := range t.rows {
		for i, cell := range row {
			formatted := t.ValueFormatter(i, t.widths[i], cell)
			_, err := fmt.Fprint(t.Writer, formatted)
			if err != nil {
				return err
			}
			if i < len(row)-1 {
				_, err = fmt.Fprint(t.Writer, strings.Repeat(" ", t.Padding))
				if err != nil {
					return err
				}
			}
		}
		_, err := fmt.Fprintln(t.Writer)
		if err != nil {
			return err
		}
	}
	return nil
}
