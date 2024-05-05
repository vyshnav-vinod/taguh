/*
	Contains all the utility/helper functions for taguh
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

func HandleError(e error) {
	_, f, l, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
	fmt.Fprintf(os.Stderr, "[LINE %d] %s\n", l, f)
	os.Exit(-1)
}

func WriteJsonToFile(file string, fileData any) {

	// CAUTION : This function rewrites the file contents fully

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

func CheckDataFiles() {
	if !checkIFExists("data") {
		err := os.Mkdir("data", 0755)
		if err != nil {
			HandleError(err)
		}
	}
	if !checkIFExists(TagsFileName) {
		_, err := os.Create(TagsFileName)
		if err != nil {
			HandleError(err)
		}
	}
	if !checkIFExists(DbFileName) {
		_, err := os.Create(DbFileName)
		if err != nil {
			HandleError(err)
		}
	}
}

func checkIFExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		HandleError(err)
		return false
	} else {
		return true
	}
}
