package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

type CmdAdd struct {
	Args []string `arg:"positional" help:"Args should be of the form (filepath tag1 tag2 .. tagn)"`
	Tag  bool     `arg:"-t" help:"Use this flag to add a new tag. If this flag is used, the args should be of the form (tagname tagdescription)"`
}

type CmdList struct {
	Type string `arg:"positional" help:"Should be either files or tags"`
}

var args struct {
	Add  *CmdAdd  `arg:"subcommand:add" help:"Add a file/tag to taguh"`
	List *CmdList `arg:"subcommand:list" help:"List all files/tags added to taguh"`
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
			// TODO: Find a way for multiple files to be tagged at the same time (i.e using one command only)
			fileName := addArgs[0]
			tags := addArgs[1:]
			if !DataValidate(fileName, "file") {
				HandleError(errors.New("file does not exists"))
			}
			if !DataValidate(strings.Join(tags, ","), "tag") {
				HandleError(errors.New("tag(s) not found"))
			}
			db := getDBVal(DbFileName) // Load the contents of the db to memory

			var tagsFinal string
			// if there are no tags in the file(i.e file is new to db) -> do normal strings.Join() and remove the end comma
			// if there are tags in the file -> check for redundancy and update with only the new tag

			if len(db[fileName].Tags) == 0 {
				// No tags exists for the file (i.e. File was not added to taguh)
				tagsFinal = strings.Join(tags, ",")
			} else {
				// File already has tags
				tagsFinal = db[fileName].Tags
				existingTags := strings.Split(db[fileName].Tags, ",")
				for _, i := range tags {
					exists := false
					for _, j := range existingTags {
						if strings.EqualFold(i, j) {
							exists = true
						}
					}
					if exists {
						tags = SlicePop(tags, slices.Index(tags, i))
					}
				}
				if !(len(tags) == 0) {
					// Join to tagsFinal only if there is a new tag, else
					// the tags will be rewritten by itself
					tagsFinal = tagsFinal + "," + strings.Join(tags, ",")
				}
			}

			db[fileName] = FileData{
				Tags:      tagsFinal,
				CreatedOn: time.Now().Format("2006-01-02 15:04:05"),
			} // Append the new file to the in-memory db contents and then write to json once again
			WriteJsonToFile(DbFileName, db)
		}

	case args.List != nil:
		listType := args.List.Type
		if len(listType) == 0 {
			// No type was provided
			_subCommandUsage("list")
		}
		switch strings.ToLower(listType) {
		// Use colors/formats to better print the output
		case "files":
			db := getDBVal(DbFileName)
			fmt.Printf("The list of files added to taguh :\n\n")
			for name := range db {
				fmt.Println(name)
			}
			fmt.Println()
		case "tags":
			fmt.Printf("The list of tags added to taguh :\n\n")
			tags := getTags()
			for tag := range tags {
				fmt.Printf("%v: %v\n", tag, tags[tag].Description)
			}
			fmt.Println()
		default:
			fmt.Fprintln(os.Stderr, "Only accepted types are [files,tags]")
			os.Exit(1)
		}

	}
}

func _subCommandUsage(cmd string) {
	fmt.Fprintf(os.Stderr, "Please provide arguments. For usage : taguh %s -h\n", cmd)
	os.Exit(1)
}
