package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	now := time.Now().UTC()
	repoVersions := flag.String("rv", "", "Comma delimited list of Repo Versions: required, no default")
	repo := flag.String("r", "", "Repo")
	repoType := flag.String("rt", "", "Repo Type: tn, ta, tw, twl; required, no default")
	outputFile := flag.String("o", "", "Output path/filename, required")
	flag.Parse()

	if *repo == "" {
		log.Fatalln("Repo argument missing")
	}

	if *repoVersions == "" {
		log.Fatalln("Repo Versions argument missing")
	}

	if *repoType == "" {
		log.Fatalln("Repo Type argument missing. Must be one of [tn, ta, tw, twl]")
	}

	if *outputFile == "" {
		log.Fatalln("Output file name is missing")
	}

	rversions := strings.Split(*repoVersions, ",")

	baseUrl := "https://git.door43.org/api/v1/repos/unfoldingword"
	log.Println("Start.")

	var results [][]string
	header := []string{"Repo", "Release", "Filename", "SHA"}
	results = append(results, header)

	if *repoType == "tn" || *repoType == "twl" {
		results = tsv_revisions(baseUrl, *repo, rversions, results)
	} else if *repoType == "tw" {
		results = tw_revisions(baseUrl, *repo, rversions, results)
	} else {
		log.Fatalf("Resource not supported yet:%v", *repoType)
	}

	f, err := os.Create(*outputFile)
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(results) // calls Flush internally

	if err != nil {
		log.Fatal(err)
	}

	stop := time.Since(now)
	log.Printf("Done. %v", fmt.Sprintf("Elapsed Time: %v\n", stop))
}

func remove(slice [][]string, s int) [][]string {
	return append(slice[:s], slice[s+1:]...)
}
