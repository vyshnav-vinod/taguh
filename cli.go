package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

type CmdAdd struct {
	Args []string `arg:"positional" help:"Args should be of the form (filepath tag1 tag2 .. tagn)"`
	Tag  bool     `arg:"-t" help:"Use this flag to add a new tag. If this flag is used, the args should be of the form (tagname tagdescription)"`
}

var args struct {
	Add *CmdAdd `arg:"subcommand:add" help:"Add a file/tag to taguh"`
}

func Cli() {
	arg.MustParse(&args)

	switch {
	case args.Add != nil:
		addArgs := args.Add.Args
		if addArgs == nil {
			_subCommandUsage("add")
		}
		if args.Add.Tag {
			// Add a tag
			if len(addArgs) < 2 {
				_subCommandUsage("add")
			}
			tagName := addArgs[0]
			tagDesc := addArgs[1:]
			tags := getTags() // Load the tags  to memory

			tags[tagName] = TagDbSchema{
				Description: strings.Join(tagDesc, " "),
				CreatedOn:   time.Now().Format("2006-01-02 15:04:05"),
			} // Append the new tag to the in-memory tags and write to JSON once again
			WriteJsonToFile(TagsFileName, tags)

		} else {
			// Add a file to the DB
			if len(addArgs) < 2 {
				_subCommandUsage("add")
			}
			fileName := addArgs[0]
			tags := addArgs[1:]
			if !DataValidate(fileName, "file"){
				HandleError(errors.New("file does not exists"))
			}
			if !DataValidate(strings.Join(tags, ","), "tag"){
				// TODO: Validate when more than one tag is given
				// TODO: Specifiy which tag was not found 
				HandleError(errors.New("tag(s) not found"))
			}
			db := getDBVal() // Load the contents of the db to memory
			db[fileName] = FileData{
				Tags:      strings.Join(tags, ","),
				CreatedOn: time.Now().Format("2006-01-02 15:04:05"),
			} // Append the new file to the in-memory db contents and then write to json once again
			WriteJsonToFile(DbFileName, db)
		}

	}
}

func _subCommandUsage(cmd string) {
	fmt.Fprintf(os.Stderr, "Please provide arguments. For usage : taguh %s -h\n", cmd)
	os.Exit(1)
}

// TODO: First do the validations for the add commands
// TODO: Then write tests
// TODO: Then move to the next command