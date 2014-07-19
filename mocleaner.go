package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"runtime/pprof"
)

/* This function takes a string and returns
   a (potentially nil) error object */
func TTWS(filename string, verbose bool) error {
	data, err := ioutil.ReadFile(filename);
	/* Did we manage to read the contents successfully? */

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

	lines := strings.Split(string(data), "\n")
	nlines := len(lines)
	for i, line := range(lines) {
		//fmt.Println("Original line: '"+line+"'");

		/* Trim whitespace */
		line = strings.TrimRight(line, " \t")
		/* Don't add a \n to the last line if it is empty */
		if (i<nlines-1 || len(line)>0) { line = line+"\n" }
		outf.Write([]byte(line));

		//fmt.Println(" Trimmed line: '"+line+"'");
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

func run(root string, verbose bool) {
	err := processNode(root, verbose);
	fmt.Printf("processNode("+root+") returned %v\n", err);
}

func main() {
	var verbose = flag.Bool("verbose", false, "request verbose output")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	
	flag.Parse()

	root := flag.Arg(0)

	if (*cpuprofile!="") {
		f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }

		fmt.Println("Starting profiling");
		pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()

		run(root, *verbose);

		pprof.StopCPUProfile();
		fmt.Println("...done profiling");
		f.Close();
	} else {
		run(root, *verbose);
	}
}
