package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraform(t *testing.T) {
    opts := &terraform.Options{
        TerraformDir: "../configs/dev",
    }

    terraform.InitAndApply(t, opts)
}
