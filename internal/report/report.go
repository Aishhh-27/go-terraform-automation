package report

import (
	"fmt"
	"os"

	"github.com/aishhh27/go-terraform-automation/internal/terraform"
)

func GenerateAuditReport(workspace string, changes []terraform.Change) error {

	file, err := os.OpenFile("audit-report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "\nWorkspace: %s\n", workspace)
	fmt.Fprintf(file, "----------------------------------\n")

	for _, c := range changes {
		line := fmt.Sprintf("%s → %s\n", c.Resource, c.Action)
		fmt.Print(line)
		file.WriteString(line)
	}

	return nil
}
