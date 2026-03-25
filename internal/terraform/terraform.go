package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// ==============================
// 🔹 STRUCTS
// ==============================

type Change struct {
	Resource string
	Action   string
}

// ==============================
// 🔹 WORKSPACE HANDLING
// ==============================

func ListWorkspaces(baseDir string) ([]string, error) {
	dirs, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var workspaces []string
	for _, d := range dirs {
		if d.IsDir() {
			workspaces = append(workspaces, filepath.Join(baseDir, d.Name()))
		}
	}

	return workspaces, nil
}

// ==============================
// 🔹 TERRAFORM INIT
// ==============================

func InitTerraform(dir string) (*tfexec.Terraform, error) {

	tfPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, fmt.Errorf("terraform not found in PATH")
	}

	tf, err := tfexec.NewTerraform(dir, tfPath)
	if err != nil {
		return nil, err
	}

	err = tf.Init(context.Background())
	if err != nil {
		return nil, err
	}

	return tf, nil
}

// ==============================
// 🔹 APPLY FLOW
// ==============================

func PlanApply(dir string, apply bool) error {

	tf, err := InitTerraform(dir)
	if err != nil {
		return err
	}

	hasChanges, err := tf.Plan(context.Background())
	if err != nil {
		return err
	}

	if hasChanges {
		fmt.Println("🔍 Changes detected in:", dir)
	} else {
		fmt.Println("✅ No changes in:", dir)
	}

	if apply && hasChanges {
		fmt.Println("⚙️ Applying changes in:", dir)

		err = tf.Apply(context.Background())
		if err != nil {
			return err
		}

		fmt.Println("✅ Apply complete:", dir)
	}

	return nil
}

// ==============================
// 🔹 AUDIT FLOW (PLAN JSON)
// ==============================

func GeneratePlanJSON(dir string) ([]byte, error) {

	tf, err := InitTerraform(dir)
	if err != nil {
		return nil, err
	}

	planFile := "tfplan"

	_, err = tf.Plan(context.Background(), tfexec.Out(planFile))
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("terraform", "show", "-json", planFile)
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}

// ==============================
// 🔹 PARSE TERRAFORM JSON
// ==============================

func ParseChanges(jsonData []byte) ([]Change, error) {

	var raw map[string]interface{}

	err := json.Unmarshal(jsonData, &raw)
	if err != nil {
		return nil, err
	}

	var changes []Change

	resourceChanges, ok := raw["resource_changes"].([]interface{})
	if !ok {
		return changes, nil
	}

	for _, r := range resourceChanges {

		res := r.(map[string]interface{})

		address := res["address"].(string)

		changeBlock := res["change"].(map[string]interface{})
		actions := changeBlock["actions"].([]interface{})

		action := actions[0].(string)

		changes = append(changes, Change{
			Resource: address,
			Action:   action,
		})
	}

	return changes, nil
}

// ==============================
// 🔥 DRIFT DETECTION
// ==============================

func DetectDrift(dir string) (bool, error) {

	tf, err := InitTerraform(dir)
	if err != nil {
		return false, err
	}

	hasChanges, err := tf.Plan(
		context.Background(),
		tfexec.RefreshOnly(true),
	)

	if err != nil {
		return false, err
	}

	return hasChanges, nil
}
