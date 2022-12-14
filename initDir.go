package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func initDir(default_recco bool) error {
	file, err := os.Create(".tflint.hcl")
	if err != nil {
		return errors.New(fmt.Sprintf(`{ "error" : "%s", "message" : "File creation error" }`, err))
	}
	_, errW := file.WriteString(fmt.Sprintf(`
	config {
		disabled_by_default = %t
	}
	plugin "developer"{
		enabled = true
		version = "1.2.0"
		source  = "github.com/trilogy-group/tflint-ruleset-developer"
	}
	rule "flag_reccomend"{
		enabled=true
	}
	`, default_recco))
	if errW != nil {
		return errors.New(fmt.Sprintf(`{ "error" : "%s", "message" : "File Writing error" }`, errW))
	}
	_, errT := exec.Command(tflint(), "--init").Output()
	if errT != nil {
		return errors.New(fmt.Sprintf(`{ "error" : "%s", "message" : "Tflint initialization error. Check the linter install file." }`, errT))
	}
	_, errY := exec.Command(yor(), "tag", "-d", ".", "--tag-groups", "code2cloud").Output()
	if errY != nil {
		return errors.New(fmt.Sprintf(`{ "error" : "%s", "message" : "Problem with yor tagging" }`, errY))
	}
	fmt.Println(`Directory initialised succesfully. Run "terraform apply" to apply tags.`)
	return nil
}
