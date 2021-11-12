package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gonutz/dfm"
)

func main() {
	exitCode := 0

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, `usage: %s one.dfm two.dfm three.dfm [...]
Provide at least one .dfm file. All files will be stripped of blocks like this:

  ExplicitLeft = 0
  ExplicitTop = 0
  ExplicitWidth = 0
  ExplicitHeight = 0

If any of these is not 0 or does not exist, the object stays unchanged.`,
			filepath.Base(os.Args[0]))
		exitCode = 1
	}

	for _, path := range os.Args[1:] {
		if err := cleanse(path); err != nil {
			fmt.Fprintf(os.Stderr, "error for file '%s': %v\n", path, err)
			exitCode = 2
		}
	}

	os.Exit(exitCode)
}

func cleanse(path string) error {
	// Using this program with other tools, we might have a space at the end of
	// our path name. Trim it.
	path = strings.TrimRight(path, " \t\n")

	if !strings.HasSuffix(strings.ToLower(path), ".dfm") {
		return errors.New("only .dfm files are supported")
	}

	root, err := dfm.ParseFile(path)
	if err != nil {
		return err
	}

	if cleanseObject(root) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		return root.WriteTo(f)
	} else {
		// Nothing was changed.
		return nil
	}
}

func cleanseObject(obj *dfm.Object) (changedAny bool) {
	propIndex := func(name string) int {
		for i, prop := range obj.Properties {
			if prop.Name == name {
				return i
			}
		}
		return -1
	}

	isZero := func(name string) bool {
		i := propIndex(name)
		if i != -1 {
			if n, ok := obj.Properties[i].Value.(dfm.Int); ok {
				return n == 0
			}
		}
		// If it was not found or it was not 0, return false.
		return false
	}

	if isZero("ExplicitLeft") &&
		isZero("ExplicitTop") &&
		isZero("ExplicitWidth") &&
		isZero("ExplicitHeight") {

		remove := func(i int) {
			obj.Properties = append(obj.Properties[:i], obj.Properties[i+1:]...)
		}

		remove(propIndex("ExplicitLeft"))
		remove(propIndex("ExplicitTop"))
		remove(propIndex("ExplicitWidth"))
		remove(propIndex("ExplicitHeight"))

		changedAny = true
	}

	for i := range obj.Properties {
		if subObj, ok := obj.Properties[i].Value.(*dfm.Object); ok {
			// NOTE This cannot be
			//
			// 	changedAny || cleanseObject(subObj)
			//
			// because the || operator will shortcut if changedAny is already
			// true, thus not calling cleanseObject recursively as we want.
			changedAny = cleanseObject(subObj) || changedAny
		}
		if set, ok := obj.Properties[i].Value.(dfm.Set); ok {
			for i := range set {
				if subObj, ok := set[i].(*dfm.Object); ok {
					changedAny = cleanseObject(subObj) || changedAny
				}
			}
		}
		if tuple, ok := obj.Properties[i].Value.(dfm.Tuple); ok {
			for i := range tuple {
				if subObj, ok := tuple[i].(*dfm.Object); ok {
					changedAny = cleanseObject(subObj) || changedAny
				}
			}
		}
		if items, ok := obj.Properties[i].Value.(dfm.Items); ok {
			for i := range items {
				for j := range items[i] {
					if subObj, ok := items[i][j].Value.(*dfm.Object); ok {
						changedAny = cleanseObject(subObj) || changedAny
					}
				}
			}
		}
	}

	return changedAny
}
