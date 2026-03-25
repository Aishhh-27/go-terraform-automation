package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run Terraform plan",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Terraform Plan...")
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
