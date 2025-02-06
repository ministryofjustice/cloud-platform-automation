package utils

import (
	"regexp"

	"github.com/google/go-github/v68/github"
)

var (
	skipfile bool
)

// Files checks if the file is a namespace file and returns a boolean value
func Files(sfname string, files []*github.CommitFile) bool {
	for _, file := range files {
		if sfname, _ := regexp.MatchString(sfname, *file.Filename); sfname {
			skipfile = true
		} else {
			skipfile = false
		}
	}

	if skipfile {
		if len(files) == 1 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
