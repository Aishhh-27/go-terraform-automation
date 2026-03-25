package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aishhh27/go-terraform-automation/internal/report"
	"github.com/aishhh27/go-terraform-automation/internal/terraform"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate Terraform audit report",
	Run: func(cmd *cobra.Command, args []string) {

		workspaces, err := terraform.ListWorkspaces("configs")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("📊 Generating REAL audit report...\n")

		for _, ws := range workspaces {

			fmt.Println("🔍 Analyzing:", ws)

			jsonData, err := terraform.GeneratePlanJSON(ws)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			changes, err := terraform.ParseChanges(jsonData)
			if err != nil {
				fmt.Println("Parse error:", err)
				continue
			}

			report.GenerateAuditReport(ws, changes)
		}

		fmt.Println("\n✅ Audit report saved to audit-report.txt")
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
