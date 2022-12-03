package main

import (
	"fmt"
	"os"
	"os/exec"
)

func initDir() error {
	file, err := os.Create(".tflint.hcl")
	if err != nil {
		return err
	}
	_, errW := file.WriteString(`plugin "template"{
		enabled = true
		version = "0.3.0"
		source  = "github.com/prasheel-ti/tflint-ruleset-template"
}`)
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
