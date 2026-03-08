package reports

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// RowDef represents a single row in a CustomTableSubReport.
type RowDef struct {
	// Cells are the text values per column.
	Cells []string
	// Bold renders all cells in bold.
	Bold bool
	// Header renders the row as a column header (bold + dark background).
	Header bool
}

// CustomTableSubReport renders a table with per-row bold/header styling.
// Column widths must sum to 12 (maroto grid units).
type CustomTableSubReport struct {
	Title   string
	Widths  []uint
	Headers []string
	Rows    []RowDef
}

func (c CustomTableSubReport) GetTitle() string {
	return c.Title
}

var headerBg = &color.Color{Red: 80, Green: 80, Blue: 100}
var sectionBg = &color.Color{Red: 200, Green: 220, Blue: 220}

func (c CustomTableSubReport) Render(m pdf.Maroto) {
	rowHeight := float64(6)

	// header row
	m.Row(rowHeight, func() {
		for i, h := range c.Headers {
			w := c.Widths[i]
			m.Col(w, func() {
				m.Text(h, props.Text{
					Top:   1,
					Style: consts.Bold,
					Size:  9,
					Align: consts.Left,
					Color: color.Color{Red: 255, Green: 255, Blue: 255},
				})
			})
			_ = headerBg
		}
	})

	// data rows
	for _, row := range c.Rows {
		rowCopy := row
		m.Row(rowHeight, func() {
			style := consts.Normal
			if rowCopy.Bold || rowCopy.Header {
				style = consts.Bold
			}
			for i, cell := range rowCopy.Cells {
				w := c.Widths[i]
				cellText := cell
				align := consts.Left
				if i == len(rowCopy.Cells)-1 {
					align = consts.Right
				}
				m.Col(w, func() {
					m.Text(cellText, props.Text{
						Top:   1,
						Style: style,
						Size:  9,
						Align: align,
					})
				})
			}
		})
		_ = sectionBg
	}
}
