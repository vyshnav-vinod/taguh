package main

import (
	"fmt"
	"os"
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
			tagDesc := addArgs[1]

			Tags[tagName] = TagDbSchema{
				Description: tagDesc,
				CreatedOn:   time.Now().Format("2006-01-02 15:04:05"),
			}
			WriteJsonToFile(TagsFileName, Tags)

		}
	}
}

func _subCommandUsage(cmd string) {
	fmt.Fprintf(os.Stderr, "Please provide arguments. For usage : taguh %s -h\n", cmd)
	os.Exit(1)
}
