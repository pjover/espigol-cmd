package generate

import (
	"fmt"
	"time"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
)

type expenseForecastReportCmd struct {
	cmd             *cobra.Command
	generateService ports.GenerateService
	year            int
}

// NewExpenseForecastReportCmd creates the "expense-forecast-report" subcommand.
func NewExpenseForecastReportCmd(generateService ports.GenerateService) cli.Cmd {
	c := &expenseForecastReportCmd{
		generateService: generateService,
	}
	cmd := &cobra.Command{
		Use:   "expense-forecast-report",
		Short: "Generate expense forecast PDF report",
		Long:  "Generates a PDF report with expense forecasts grouped by category and scope for the given year.",
		RunE:  c.run,
	}
	cmd.Flags().IntVarP(&c.year, "year", "y", time.Now().Year(), "Year to generate the report for")
	c.cmd = cmd
	return c
}

func (c *expenseForecastReportCmd) Cmd() *cobra.Command {
	return c.cmd
}

func (c *expenseForecastReportCmd) run(cmd *cobra.Command, args []string) error {
	hasNegativeRemainder, msg, err := c.generateService.ExpenseForecastReport(c.year)
	if err != nil {
		return err
	}
	fmt.Println(msg)
	if hasNegativeRemainder {
		return fmt.Errorf("el remanent és negatiu per a alguna categoria; consulteu el PDF per als detalls de l'ajust")
	}
	return nil
}
