package main

import (
	"flag"
	"fmt"
	"path"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/* This function takes a string and returns
   a (potentially nil) error object */
func TTWS(filename string, verbose bool) error {
	/* Open the input file */
	inf, err := os.Open(filename);
	/* In case this function generates a "panic", be sure to close this file */
	defer inf.Close();
	/* Did we open it successfully?  If not, close and return. */
	if (err!=nil) { return err; }

	data, err := ioutil.ReadAll(inf);
	inf.Close();

	fileType := http.DetectContentType(data);
	if (!strings.Contains(fileType, "text/plain")) {
		if (verbose) { fmt.Printf("Skipping file of type '%v': %v\n", fileType, filename); }
		return nil;
	}

	/* Open the output file in system temp dir*/
	outf, err := ioutil.TempFile("","");
	/* In case this function generates a "panic", be sure to close this file */
	defer outf.Close();
	/* Did we open it succesfully?  If not, close all and return. */
	if (err!=nil) { return err; }

	for _, line := range(strings.Split(string(data), "\n")) {
		line = strings.TrimRight(line, " \t")+"\n"
		outf.Write([]byte(line));
	}

	outf.Close();

	/* Replace the source file by the trimmed file */
	err = os.Rename(outf.Name(), filename);
	if (err!=nil) { return err; }

	if (verbose) { fmt.Printf("Trimmed %s\n", filename); }
	/* No errors, so we return nil */
	return nil;
}

var blacklist = []string{".bzr", ".cvs", ".git", ".hg", ".svn"}

func processNode(node string, verbose bool) error {
	fi, err := os.Lstat(node)
	if (err!=nil) { return err; }

	if (fi.IsDir()) {
		if contains(fi.Name(), blacklist) { return nil; }
		contents, err := ioutil.ReadDir(node);
		if (err!=nil) { return err; }
		for _, n := range(contents) {
			serr := processNode(path.Join(node, n.Name()), verbose);
			if (serr!=nil) { return serr; }
		}
		return nil;
	} else {
		return TTWS(node, verbose);
	}
}

func contains(x string, a []string) bool {
	for _, e := range(a) {
		if (x==e) { return true; }
	}
	return false;
}


func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := processNode(root, false);
	fmt.Printf("processNode("+root+") returned %v\n", err);
}
