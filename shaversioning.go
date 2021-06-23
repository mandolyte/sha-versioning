package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	now := time.Now().UTC()
	repo := flag.String("r", "", "Repo")
	repoType := flag.String("rt", "", "Repo Type: tn, ta, tw, twl; required, no default")
	outputFile := flag.String("o", "", "Output path/filename, required")
	flag.Parse()
	baseUrl := "https://git.door43.org/api/v1/repos/unfoldingword"

	if *repo == "" {
		log.Fatalln("Repo argument missing")
	}

	if *repoType == "" {
		log.Fatalln("Repo Type argument missing. Must be one of [tn, ta, tw, twl]")
	}

	if *outputFile == "" {
		log.Fatalln("Output file name is missing")
	}

	//rversions := strings.Split(*repoVersions, ",")
	rversions := getTags(baseUrl, *repo)
	log.Printf("Start for repo %v\n", *repo)

	var results [][]string
	header := []string{"Repo", "Release", "Filename", "SHA"}
	results = append(results, header)

	if *repoType == "tn" || *repoType == "twl" {
		results = tsv_revisions(baseUrl, *repo, rversions, results)
	} else if *repoType == "tw" {
		results = tw_revisions(baseUrl, *repo, rversions, results)
	} else if *repoType == "ta" {
		results = ta_revisions(baseUrl, *repo, rversions, results)
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

// example url:
// https://qa.door43.org/api/v1/repos/unfoldingword/en_twl/releases
func getTags(baseUrl, repo string) []string {
	fullUrl := baseUrl + "/" + repo + "/releases"
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// convert json to a map
	var jsonArray []map[string]interface{}
	err = json.Unmarshal([]byte(body), &jsonArray)
	if err != nil {
		panic(err)
	}

	var tags []string
	for _, val := range jsonArray {
		tags = append(tags, fmt.Sprintf("%v", val["tag_name"]))
	}

	return tags
}

/* sample JSON from releases URL API

[
  {
    "id": 124428,
    "tag_name": "v3",
    "target_commitish": "master",
    "name": "Version 3",
    "body": "This is version 3 of the unfoldingWord®"
*/
