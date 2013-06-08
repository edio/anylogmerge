package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"logmerge"
	"os"
)

// flags
var sortingKeyRegex string
var inFileNames []string
var outFileName string
var forceOverwrite bool

// files
var outFile io.Writer
var inFiles []io.Reader

// misc
var logger = log.New(os.Stderr, "", 0)

func init() {
	const (
		defaultSortingKeyRegex = `^(\S+).*`
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [OPTION]... [FILE]...\n", os.Args[0])
		flag.VisitAll(func(f *flag.Flag) {
			if len(f.DefValue) == 0 {
				fmt.Fprintf(os.Stdout, "  -%-8s %s\n", f.Name, f.Usage)
			} else {
				fmt.Fprintf(os.Stdout, "  -%-8s %s.\n%12sDefault: %s\n", f.Name, f.Usage, "", f.DefValue)
			}
		})
	}

	flag.StringVar(&sortingKeyRegex, "s", defaultSortingKeyRegex, "regular expression with capturing groups indicating sorting key")
	flag.StringVar(&outFileName, "o", "", "name of the output file. Write to stdout if not set")
	flag.BoolVar(&forceOverwrite, "f", false, "force overwrite if output file exists (when -o is used)")
}

func postInit() {
	flag.Parse()
	inFileNames = flag.Args()

	if len(inFileNames) < 2 {
		logger.Fatal("Nothing to merge. At least 2 files are required")
	}

	inFiles = make([]io.Reader, len(inFileNames))
	for index, name := range inFileNames {
		file := openFile(name)
		inFiles[index] = file
	}

	if len(outFileName) == 0 {
		outFile = os.Stdout
	} else {
		fileFlags := os.O_RDWR | os.O_CREATE
		if !forceOverwrite {
			fileFlags |= os.O_EXCL
		}

		var err error
		outFile, err = os.OpenFile(outFileName, fileFlags, 0666)
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func openFile(name string) io.Reader {
	inFile, err := os.Open(name)
	if err != nil {
		logger.Fatal(err)
	}
	return inFile
}

func main() {
	postInit()

	merger := logmerge.NewMerger(outFile)
	merger.Merge(logmerge.LexicographicMinimum, inFiles)
}
