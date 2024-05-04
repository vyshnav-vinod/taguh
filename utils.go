/*
	Contains all the utility/helper functions for taguh
*/

package main

import (
	"fmt"
	"os"
)

func HandleError(e error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
	os.Exit(-1)
}

