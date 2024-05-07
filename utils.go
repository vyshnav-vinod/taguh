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
	/*
		Validate the file/tag
		s -> The content to validate
		t -> The type of content (file or tag)
	*/
	// TODO: Return what tag was invalid
	if strings.ToLower(t) == "file" {
		if checkIfExists(s) {
			return true
		} else {
			return false
		}
	} else if strings.ToLower(t) == "tag" {
		tags := getTags()
		if strings.Contains(s, ",") {
			// More than one tag was provided by the user
			tagsList := strings.Split(s, ",")
			var count = 0
			for _, tag := range tagsList {
				for name := range tags {
					if strings.EqualFold(tag, name) {
						count++
					}
				}
			}
			if len(tagsList) == count { // Check if all the provided tags are valid
				return true
			} else {
				return false
			}
		}

		for name := range tags {
			if strings.EqualFold(s, name) {
				return true
			}
		}
		return false
	} else {
		panic(fmt.Sprintf("Type %s is not found. Please report!!", s))
	}
}

func SlicePop(slice []string, index int) []string {
	// Remove an element of index from the slice
	// TODO: Make this better, use the append based method??
	var tmp []string
	for i, j := range slice {
		if !(i == index) {
			tmp = append(tmp, j)
		}
	}
	return tmp
}

func PrintOutputFind(path string, tags string, added string) {
	fmt.Printf("\nFile path: %s\nTags: %s\nAdded on: %s\n", path, tags, added)
}
