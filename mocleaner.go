package main

import (
	"bufio"
	"flag"
	"fmt"
	"path"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


/* This function takes a string and returns
   a (potentially nil) error object */
func TTWS(filename string) error {
	/* Open the input file */
	inf, err := os.Open(filename);
	/* In case this function generates a "panic", be sure to close this file */
	defer inf.Close();
	/* Did we open it successfully?  If not, close and return. */
	if (err!=nil) { inf.Close(); return err; }

	/* Open the output file in system temp dir*/
	outf, err := ioutil.TempFile("","");
	/* In case this function generates a "panic", be sure to close this file */
	defer outf.Close();
	/* Did we open it succesfully?  If not, close all and return. */
	if (err!=nil) { inf.Close(); outf.Close(); return err; }
	/* Create a scanner object to break this in to lines */
	scanner := bufio.NewScanner(inf);
	/* Declare a variable for the line */
	var line string;
	/* Loop over lines */
	for scanner.Scan() {
		/* Trim right space and then add the \n back on the end before writing */
		line = strings.TrimRight(scanner.Text(), " \t")+"\n"
		outf.Write([]byte(line));
	}
	/* Close all open files */
	inf.Close();
	outf.Close();
	/* Replace the source file by the trimmed file */
	os.Rename(outf.Name(), filename);

	/* No errors, so we return nil */
	return nil;
}

var blacklist = []string{".bzr", ".cvs", ".git", ".hg", ".svn"}

func processFile(filename string, verbose bool) error {
	inf, err := os.Open(filename)
	defer inf.Close();
	if (err!=nil) { inf.Close(); return err; }

	readStart := io.LimitReader(inf, 512);

	data, err := ioutil.ReadAll(readStart);

	/* Close all open files */
	inf.Close();

	/* Determine file type */
	fileType := http.DetectContentType(data);

	/* only act on text files */
	if (strings.Contains(fileType, "text/plain")){
		if (verbose) { fmt.Printf("Trimming: %v\n", filename); }
		return TTWS(filename);
	} else {
		if (verbose) { fmt.Printf("Skipping file of type '%v': %v\n", fileType, filename); }
		return nil;
	}
}

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
		return processFile(node, verbose);
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
	err := processNode(root, true);
	fmt.Printf("processNode("+root+") returned %v\n", err);
}
