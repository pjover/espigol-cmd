package reports

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/pjover/espigol/internal/domain/ports"
)

// SubReport is a renderable section of a report.
type SubReport interface {
	GetTitle() string
	Render(m pdf.Maroto)
}

// pageBreakSubReport is a sentinel SubReport that forces a new page.
type pageBreakSubReport struct{}

// NewPageBreak returns a SubReport that inserts a page break when rendered.
func NewPageBreak() SubReport { return pageBreakSubReport{} }

func (pageBreakSubReport) GetTitle() string    { return "" }
func (pageBreakSubReport) Render(m pdf.Maroto) { m.AddPage() }

// sectionTitleSubReport renders a level-1 section heading.
type sectionTitleSubReport struct{ text string }

// NewSectionTitle returns a SubReport that renders a prominent section heading.
func NewSectionTitle(text string) SubReport { return sectionTitleSubReport{text: text} }

func (s sectionTitleSubReport) GetTitle() string { return "" }
func (s sectionTitleSubReport) Render(m pdf.Maroto) {
	m.Row(16, func() {
		m.Col(12, func() {
			m.Text(s.text, props.Text{
				Top:   4,
				Style: consts.Bold,
				Align: consts.Left,
				Color: color.Color{Red: 0, Green: 51, Blue: 51},
				Size:  16,
			})
		})
	})
}

// ReportDefinition describes the full report to be rendered.
type ReportDefinition struct {
	PageOrientation consts.Orientation
	Title           string
	Footer          string
	SubReports      []SubReport
}

// ReportPdf saves a ReportDefinition to a PDF file.
type ReportPdf struct {
	config ports.ConfigService
}

// NewReportPdf creates a new ReportPdf.
func NewReportPdf(config ports.ConfigService) *ReportPdf {
	return &ReportPdf{config: config}
}

// SaveToFile renders the report to the given file path.
func (r *ReportPdf) SaveToFile(def ReportDefinition, filePath string) error {
	m := pdf.NewMaroto(def.PageOrientation, consts.A4)
	m.SetPageMargins(15, 10, 15)

	r.header(m)
	r.footer(def.Footer, m)
	r.title(def.Title, m)
	for _, sub := range def.SubReports {
		r.subTitle(sub.GetTitle(), m)
		sub.Render(m)
	}

	dirPath, _ := path.Split(filePath)
	if dirPath != "" {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("cannot create output directory '%s': %w", dirPath, err)
		}
	}

	if err := m.OutputFileAndClose(filePath); err != nil {
		return fmt.Errorf("cannot save PDF to '%s': %w", filePath, err)
	}
	return nil
}

func (r *ReportPdf) header(m pdf.Maroto) {
	logoPath := r.expandTilde(r.config.GetString("files.logo"))
	businessName := r.config.GetString("business.name")
	_, logoErr := os.Stat(logoPath)
	hasLogo := logoErr == nil

	m.RegisterHeader(func() {
		m.Row(20, func() {
			if hasLogo {
				m.Col(3, func() {
					_ = m.FileImage(logoPath, props.Rect{
						Left:    2,
						Center:  true,
						Percent: 80,
					})
				})
				m.ColSpace(5)
				m.Col(4, func() {
					m.Text(businessName, props.Text{
						Style:       consts.BoldItalic,
						Size:        10,
						Align:       consts.Left,
						Extrapolate: false,
					})
				})
			} else {
				m.Col(12, func() {
					m.Text(businessName, props.Text{
						Style:       consts.BoldItalic,
						Size:        10,
						Align:       consts.Left,
						Extrapolate: false,
					})
				})
			}
		})
	})
}

func (r *ReportPdf) footer(footer string, m pdf.Maroto) {
	if footer == "" {
		return
	}
	m.RegisterFooter(func() {
		m.Row(4, func() {
			m.Col(12, func() {
				m.Text(footer, props.Text{
					Top:   4,
					Style: consts.Italic,
					Size:  8,
					Align: consts.Right,
				})
			})
		})
	})
}

func (r *ReportPdf) title(title string, m pdf.Maroto) {
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(title, props.Text{
				Top:   4,
				Style: consts.Bold,
				Align: consts.Center,
				Color: color.Color{Red: 0, Green: 51, Blue: 51},
				Size:  18,
			})
		})
	})
}

func (r *ReportPdf) subTitle(subTitle string, m pdf.Maroto) {
	if subTitle == "" {
		return
	}
	m.Row(14, func() {
		m.Col(12, func() {
			m.Text(subTitle, props.Text{
				Top:   6,
				Style: consts.Bold,
				Align: consts.Left,
				Color: color.Color{Red: 0, Green: 51, Blue: 51},
				Size:  13,
			})
		})
	})
}

// expandTilde replaces a leading ~ with the user home directory.
func (r *ReportPdf) expandTilde(p string) string {
	if !strings.HasPrefix(p, "~/") {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	return path.Join(home, p[2:])
}

// ExpandTilde is exported for use from the report service.
func ExpandTilde(p string) string {
	if !strings.HasPrefix(p, "~/") {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	return path.Join(home, p[2:])
}
