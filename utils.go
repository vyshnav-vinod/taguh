/*
	Contains all the utility/helper functions for taguh
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
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
	if !checkIfExists("data") {
		err := os.Mkdir("data", 0755)
		if err != nil {
			HandleError(err)
		}
	}
	if !checkIfExists(TagsFileName) {
		_, err := os.Create(TagsFileName)
		if err != nil {
			HandleError(err)
		}
	}
	if !checkIfExists(DbFileName) {
		_, err := os.Create(DbFileName)
		if err != nil {
			HandleError(err)
		}
	}
}

func checkIfExists(f string) bool {
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

func DataValidate(s string, t string) bool {
	// s -> The content to validate
	// t -> The type of content (file or tag)
	if strings.ToLower(t) == "file" {
		if checkIfExists(s) {
			return true
		} else {
			return false
		}
	} else if strings.ToLower(t) == "tag" {
		tags := getTags()
		for name := range tags {
			if s == name {
				return true
			}
		}
		return false
	} else {
		panic(fmt.Sprint("Type %s is not found\n", s))
	}
}
