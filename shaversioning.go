package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	repoVersions := flag.String("rv", "v1,v2", "Comma delimited list of Repo Versions: default 'v1,v2'")
	repo := flag.String("r", "", "Repo")
	flag.Parse()

	if *repo == "" {
		log.Fatalln("Repo argument missing")
	}

	rversions := strings.Split(*repoVersions, ",")

	baseUrl := "https://qa.door43.org/api/v1/repos/unfoldingword"
	log.Println("Start.")

	var results [][]string
	header := []string{"Repo", "Release", "Filename", "SHA"}
	results = append(results, header)
	for rv := range rversions {
		//log.Printf("Working on verison:%v", rversions[rv])
		fullUrl := baseUrl + "/" + *repo + "/git/trees/" + rversions[rv]
		resp, err := http.Get(fullUrl)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// convert json to a map
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(body), &jsonMap)
		if err != nil {
			panic(err)
		}

		for key, val := range jsonMap {
			if key == "tree" {
				v, ok := val.([]interface{})
				if !ok {
					log.Fatalln("failed to cast")
				}
				for _, tval := range v {
					var row []string
					nodeMap := tval.(map[string]interface{})
					//fmt.Printf("%v,%v,%v,%v\n", *repo, rversions[rv], nodeMap["path"], nodeMap["sha"])
					row = append(row, *repo, rversions[rv], fmt.Sprintf("%v", nodeMap["path"]), fmt.Sprintf("%v", nodeMap["sha"]))
					results = append(results, row)
				}
			}
		}
	}

	// remove any entries that are not TSV files (need to adjust this for markdown content)
	// identify rows to remove
	var positionsToRemove []int
	for i := 1; i < len(results); i++ {
		if !strings.HasSuffix(results[i][2], ".tsv") {
			positionsToRemove = append(positionsToRemove, i)
		}
	}
	// now remove them
	for i := len(positionsToRemove) - 1; i > -1; i-- {
		results = remove(results, positionsToRemove[i])
	}

	output := *repo + "_revs.csv"
	f, err := os.Create(output)
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(results) // calls Flush internally

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(results)
	log.Println("Done.")
}

func remove(slice [][]string, s int) [][]string {
	return append(slice[:s], slice[s+1:]...)
}
