package main

import (
	"encoding/json"
	"os"

	"github.com/trilogy-group/cloudfix-linter-developer/cloudfixIntegration"
)

type Persistance struct {
	// No data fields required for this class. Data shall be persisted using the file system
}

//Member functions for the class follow:
func (p *Persistance) store_reccos(reccosMap map[string]cloudfixIntegration.Recommendation, fileNameForReccos string) error {
	file, err := os.Create(fileNameForReccos)
	defer file.Close()
	if err != nil {
		//Add error log
		return err
	}
	recommendationJSON, e := json.Marshal(reccosMap)
	if e != nil {
		return e
	}
	file.Write(recommendationJSON)
	return nil
}

func (p *Persistance) store_tagMap(tagToIDMap map[string]map[string][]string, fileNameForTagMap string) error {
	file, err := os.Create(fileNameForTagMap)
	if err != nil {
		//Add error log
		return err
	}
	recommendationJSON, e := json.Marshal(tagToIDMap)
	if e != nil {
		return e
	}
	file.Write(recommendationJSON)
	return nil
}
