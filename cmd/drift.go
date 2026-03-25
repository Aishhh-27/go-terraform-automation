package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/aishhh27/go-terraform-automation/internal/terraform"
)

var driftCmd = &cobra.Command{
	Use:   "drift",
	Short: "Detect infrastructure drift",
	Run: func(cmd *cobra.Command, args []string) {

		workspaces, err := terraform.ListWorkspaces("configs")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("🔍 Detecting infrastructure drift...\n")

		for _, ws := range workspaces {

			fmt.Println("Checking:", ws)

			drift, err := terraform.DetectDrift(ws)
			if err != nil {
				fmt.Println("❌ Error:", err)
				continue
			}

			if drift {
				fmt.Println("⚠️ Drift detected in:", ws)
			} else {
				fmt.Println("✅ No drift in:", ws)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(driftCmd)
}
