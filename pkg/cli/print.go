package cli

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func printTable(rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader(rows[0])

	for _, row := range rows[1:] {
		table.Append(row)
	}

	table.Render()
}
