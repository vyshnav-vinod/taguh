package main

import (
	"encoding/json"
	// "fmt"
	"os"
	"time"

	"github.com/buger/jsonparser"
)

const (
	TagsFileName = "data/tags.json"
	DbFileName   = "data/taguh.json"
)

var (
	Tags map[string]TagDbSchema
)

type FileData struct {

	/*
		Structure of taguh.json
		FileName1: {
			"tags": ...,
			"created_on": ...
		},
		FileName2: {
			.....
		},
		...
	*/

	Tags      []string  `json:"tags"`
	CreatedOn time.Time `json:"created_on"`
}

type TagDbSchema struct {

	/*
		Structure of the tags.json
		TagName1: {
			"description": ...,
			"created_on": ...
		},
		TagName2: {
			.....
		},
		...
	*/

	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
}

func main() {

	/*
		Load the list of tags to memory and close tags.json
		Then execute the command given by the CLI
		Load the DB only when called?
	*/
	/*
		To access elements use tags[name].Description...
		Example : To access Description of Starred, we use tags["Starred"].Description
		NOTE: To check if a key is not present in a map, just try to access it
		NOTE: Use for range to iterate through the map
	*/

	Tags = getTags()
	Cli()

}

func getTags() map[string]TagDbSchema {
	// TODO: Check if file exists and if not, then create a new tags file with base tags
	tags, err := os.ReadFile(TagsFileName)
	if err != nil {
		HandleError(err)
	}
	if len(tags) == 0 {
		// If content of the tags.json was somehow removed, create a new tags.json
		_createBaseTags()
		tags, err = os.ReadFile(TagsFileName)
		if err != nil {
			HandleError(err)
		}
	}

	// Below parses the tags.json into a map (tagMap)
	tagMap := make(map[string]TagDbSchema)
	handlerFunc := func(key []byte, value []byte, datatype jsonparser.ValueType, offset int) error {
		keyTmp := string(key)
		tmp := make(map[string]string)
		jsonparser.ObjectEach(value, func(key []byte, value []byte, datatype jsonparser.ValueType, offset int) error {
			tmp[string(key)] = string(value)
			return nil
		})
		tagMap[keyTmp] = TagDbSchema{
			Description: tmp["description"],
			CreatedOn:   tmp["created_on"],
		}
		return nil
	}
	// Goes through each JSON object and calls the handlerFunc
	err = jsonparser.ObjectEach(tags, handlerFunc)
	if err != nil {
		HandleError(err)
	}
	return tagMap
}

func _createBaseTags() {

	// Creates tags.json with the base tags such as Starred, Important and Archived

	baseTags := map[string]string{"Starred": "Tag for favourite files", "Important": "Tag for important files", "Archived": "Tag for archived files or files that may not be used anymore"}
	tagsMap := make(map[string]*TagDbSchema)

	for key, value := range baseTags {
		tagsMap[key] = &TagDbSchema{
			Description: value,
			CreatedOn:   time.Now().Format("2006-01-02 15:04:05"),
		}
	}
	tagsJson, err := json.MarshalIndent(tagsMap, "", "	")
	if err != nil {
		HandleError(err)
	}
	err = os.WriteFile(TagsFileName, tagsJson, 0666)
	if err != nil {
		HandleError(err)
	}
}
