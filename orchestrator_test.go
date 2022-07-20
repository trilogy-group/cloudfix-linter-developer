package main

import (
	"fmt"
	"reflect"
	"testing"
)

type ReturnByModuleParser struct {
	modules []string
}

type TestModuleParser struct {
	moduleJson []byte
	expected   ReturnByModuleParser
}

var addTestsModules = []TestModuleParser{
	{ // Test 1: Everything normal
		[]byte(`{
			"issues": [
				{
					"rule": {
						"name": "module_source",
						"severity": "info",
						"link": ""
					},
					"message": ".//module-1",
					"range": {
						"filename": "main.tf",
						"start": {
							"line": 84,
							"column": 23
						},
						"end": {
							"line": 84,
							"column": 36
						}
					},
					"callers": []
				},
				{
					"rule": {
						"name": "module_source",
						"severity": "info",
						"link": ""
					},
					"message": ".//module-2",
					"range": {
						"filename": "main.tf",
						"start": {
							"line": 89,
							"column": 12
						},
						"end": {
							"line": 89,
							"column": 25
						}
					},
					"callers": []
				}
			],
			"errors": []
		}`),
		ReturnByModuleParser{
			[]string{".//module-1", ".//module-2"},
		},
	},
	{ // Test 2: Empty JSON returned
		[]byte(``),
		ReturnByModuleParser{
			[]string{},
		},
	},
}

func TestExtractModules(t *testing.T) {
	var orches Orchestrator
	for _, test := range addTestsModules {
		expected := test.expected
		modules, _ := orches.extractModulePaths(test.moduleJson)
		if len(expected.modules) == 0 && len(modules) == 0 {
			continue
		}
		eq := reflect.DeepEqual(ReturnByModuleParser{modules}, expected)
		if !eq {
			fmt.Println(expected)
			fmt.Println(ReturnByModuleParser{modules})
			t.Errorf("Test failed!")
		}
	}
}
