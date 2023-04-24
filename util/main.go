package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := "e:\\fp2ideas\\nouns"

	files, err := os.ReadDir(dir)
	logErr(err)

	outDir := "C:\\dev\\projects\\fp2\\fp2web\\src\\assets\\icons"
	outTs := "C:\\dev\\projects\\fp2\\fp2web\\src\\app\\icons.ts"
	items := make(map[string]string)
	for _, f := range files {
		// Convert the name
		// noun-<NAME>-<id>.svg
		//noun-abnormal-wisdom-decreased-2359993.svg
		base := filepath.Base(f.Name())
		if strings.HasPrefix(base, "noun-") {
			parts := strings.Split(base, "-")
			keep := parts[1 : len(parts)-1]
			name := strings.Join(keep, "-")
			fname := name + ".svg"
			path := filepath.Join(outDir, fname)

			_, err := os.Stat(path)
			if err != nil {
				// Not there
				src := filepath.Join(dir, f.Name())
				err := copy(src, path)
				logErr(err)
			}

			items[name] = fname
		}

	}

	// Write the TS
	os.Remove(outTs)
	w, err := os.OpenFile(outTs, os.O_CREATE, 0600)
	logErr(err)

	w.WriteString("export const SvgIcons : any = {		\n")
	for k, v := range items {
		w.WriteString(fmt.Sprintf("'%v' : '/assets/icons/%v',\n", k, v))
	}
	w.WriteString("}")
	w.Close()
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(sourceFile string, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		return err
	}
	return nil
}
