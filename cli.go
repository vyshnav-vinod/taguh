package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

type CmdFind struct {
	Type    string `arg:"positional" help:"Should be either file or tag"`
	Arg     string `arg:"positional" help:"The filename or the tag name to find"`
	Options string `arg:"positional" help:"If the type is tag, options are latest, oldest, asc and desc"`
}

var args struct {
	Add  *CmdAdd  `arg:"subcommand:add" help:"Add a file/tag to taguh"`
	List *CmdList `arg:"subcommand:list" help:"List all files/tags added to taguh"`
	Find *CmdFind `arg:"subcommand:find" help:"If a tag is given, list all the files associated with that tag. Else list the path of the filename and its tag(s)"`
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
				CreatedOn:   time.Now().Format(DateParseLayout),
			} // Append the new tag to the in-memory tags and write to JSON once again
			WriteJsonToFile(TagsFileName, tags)

		} else {
			// Add a file to the DB
			if len(addArgs) < 2 {
				_subCommandUsage("add")
			}
			// TODO: Find a way for multiple files to be tagged at the same time (i.e using one command only)
			fileName, err := filepath.Abs(addArgs[0])
			if err != nil {
				HandleError(err)
			}
			tags := addArgs[1:]
			if !DataValidate(fileName, "file") {
				HandleError(errors.New("file does not exists"))
			}
			if !DataValidate(strings.Join(tags, ","), "tag") {
				HandleError(errors.New("tag(s) not found"))
			}
			db := getDBVal(DbFileName) // Load the contents of the db to memory

			var tagsFinal string

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
				} else {
					fmt.Printf("%s is already added to taguh!!!\n", fileName)
					os.Exit(0)
				}
			}

			db[fileName] = FileData{
				Tags:      tagsFinal,
				CreatedOn: time.Now().Format(DateParseLayout),
			} // Append the new file to the in-memory db contents and then write to json once again
			WriteJsonToFile(DbFileName, db)
		}

	case args.List != nil:
		// TODO: Add options here as well (sorted, oldest, latest, etc)
		listType := args.List.Type
		if len(listType) == 0 {
			// No type was provided
			_subCommandUsage("list")
		}
		switch strings.ToLower(listType) {
		// TODO: Use colors/formats to better print the output
		// TODO: If no files/tags are found, give a better output
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

	case args.Find != nil:
		findType := strings.ToLower(args.Find.Type)
		findArg := args.Find.Arg
		// findOpt := strings.ToLower(args.Find.Options)
		if len(findType) == 0 || len(findArg) == 0 {
			_subCommandUsage("find")
		}
		db := getDBVal(DbFileName)
		switch findType {
		case "file":
			var resultFiles []string
			for fpath := range db {
				if strings.Contains(filepath.Base(fpath), findArg) {
					// Maybe more than one files with same name (but different
					// folders) may be found, so list them all to the user to choose
					resultFiles = append(resultFiles, fpath)
				}
			}
			if len(resultFiles) == 0 {
				// TODO: Suggestions for similar file [PRIORITY]
				fmt.Println("No such file found in taguh")
				os.Exit(1)
			} else {
				// TODO: Make the output better
				if len(resultFiles) == 1 {
					// Only 1 file is found
					PrintOutput(resultFiles[0], db[resultFiles[0]].Tags, db[resultFiles[0]].CreatedOn)
				} else {
					fmt.Printf("Found %d matching files\n", len(resultFiles))
					for i := range resultFiles {
						PrintOutput(resultFiles[i], db[resultFiles[i]].Tags, db[resultFiles[i]].CreatedOn)
					}
				}
			}
		case "tag":
			arg := args.Find.Arg
			options := args.Find.Options
			var result []string

			if !DataValidate(arg, "tag") {
				HandleError(errors.New(fmt.Sprintf("Tag: %s not found", arg)))
				// TODO: Suggest similar tags [PRIORITY]
			}

			for fname, fvalue := range db {
				if strings.Contains(fvalue.Tags, strings.ToLower(arg)) {
					result = append(result, fname)
				}
			}

			if len(result) == 0 {
				// TODO: Suggest similar tags [PRIORITY]
				fmt.Println("No such file found in taguh")
				os.Exit(1)
			} else {
				if len(result) == 1 {
					PrintOutput(result[0], db[result[0]].Tags, db[result[0]].CreatedOn)
				} else {
					if len(options) != 0 {
						// Do the sortings
						if !DataValidate(options, "option") {
							HandleError(errors.New(fmt.Sprintf("option %s does not exists", options)))
						} else {
							result = PerformOptions(options, result)
							
						}
					}
					fmt.Printf("Found %d matching files\n", len(result))
					for i := range result {
						PrintOutput(result[i], db[result[i]].Tags, db[result[i]].CreatedOn)
					}
				}
			}

		default:
			fmt.Fprintf(os.Stderr, "Invalid type %s", findType)
			_subCommandUsage("list")
		}

	}
}

func _subCommandUsage(cmd string) {
	fmt.Fprintf(os.Stderr, "Please provide arguments. For usage : taguh %s -h\n", cmd)
	os.Exit(1)
}
