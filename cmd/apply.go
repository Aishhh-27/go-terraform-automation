package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/cobra"
	"github.com/aishhh27/go-terraform-automation/internal/terraform"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Terraform for all environments (parallel)",
	Run: func(cmd *cobra.Command, args []string) {

		workspaces, err := terraform.ListWorkspaces("configs")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("🚀 Running Terraform in parallel...\n")

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 2) // limit concurrency

		for _, ws := range workspaces {
			wg.Add(1)

			go func(workspace string) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				fmt.Println("➡️ Starting:", workspace)

				err := terraform.PlanApply(workspace, true)
				if err != nil {
					fmt.Println("❌ Error:", workspace, err)
					return
				}

				fmt.Println("✅ Done:", workspace)

			}(ws)
		}

		wg.Wait()
		fmt.Println("\n🎉 All environments processed!")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
