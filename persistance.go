package main

import (
	"fmt"
	"os"
)

type Persistance struct {
	// No data fields required for this class. Data shall be persisted using the file system
}

//Member functions for the class follow:

func (p *Persistance) store_reccos(reccosMap map[string]map[string][]string, fileNameForReccos string) error {
	file, err := os.Create(fileNameForReccos)
	if err != nil {
		//Add error log
		return err
	}
	for key, innerMap := range reccosMap {
		for innerKey, innerList := range innerMap {
			toWrite := fmt.Sprintf("%s:%s", key, innerKey)
			for _, innerValue := range innerList {
				toWrite = toWrite + fmt.Sprintf(":%s",innerValue)
			}
			toWrite = toWrite + "\n"
			_, err := file.WriteString(toWrite)
			if err != nil {
				//Add Error Log
				return err
			}
		}
	}
	//Add Info Log
	return nil
}

func (p *Persistance) store_tagMap(tagToIDMap map[string]map[string]string, fileNameForTagMap string) error {
	file, err := os.Create(fileNameForTagMap)
	if err != nil {
		//Add error log
		return err
	}
	for key, innerMap := range tagToIDMap {
		for innerKey, resourceId := range innerMap {
			toWrite := fmt.Sprintf("%s->%s->%s\n", key, innerKey, resourceId)
			_, err := file.WriteString(toWrite)
			if err != nil {
				//Add Error Log
				return err
			}
		} 
	}
	//Add Info Log
	return nil
}
