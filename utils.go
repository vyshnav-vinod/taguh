/*
	Contains all the utility/helper functions for taguh
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func HandleError(e error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
	os.Exit(-1)
}

func WriteJsonToFile(file string, fileData any) {
	f, err := os.Create(file)
	if err != nil {
		HandleError(err)
	}
	tagsJson, err := json.MarshalIndent(fileData, "", "	")
	if err != nil {
		HandleError(err)
	}
	f.Write(tagsJson)
	defer f.Close()
}
