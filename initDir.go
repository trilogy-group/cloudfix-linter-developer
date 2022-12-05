package main

import (
	"fmt"
	"os"
	"os/exec"
)

func initDir(default_recco bool) error {
	file, err := os.Create(".tflint.hcl")
	if err != nil {
		return err
	}
	_, errW := file.WriteString(fmt.Sprintf(`
	config {
		disabled_by_default = %t
	}
	plugin "developer"{
		enabled = true
		version = "0.3.5"
		source  = "github.com/trilogy-group/tflint-ruleset-developer"
	}
	rule "flag_reccomend"{
		enabled=true
	}
	`,default_recco))
	if errW != nil {
		return errW
	}
	_, errT := exec.Command(tflint(), "--init").Output()
	if errT != nil {
		return errT
	}
	_, errY := exec.Command(yor(), "tag", "-d", ".", "--tag-groups", "code2cloud").Output()
	if errY != nil {
		return err
	}
	fmt.Println(`Directory initialised succesfully. Run "terraform apply" to apply tags.`)
	return nil
}
