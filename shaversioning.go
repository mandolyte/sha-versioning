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
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

func main() {
	now := time.Now().UTC()
	repoVersions := flag.String("rv", "v1,v2", "Comma delimited list of Repo Versions: default 'v1,v2'")
	repo := flag.String("r", "", "Repo")
	flag.Parse()

	if *repo == "" {
		log.Fatalln("Repo argument missing")
	}

	rversions := strings.Split(*repoVersions, ",")

	baseUrl := "https://git.door43.org/api/v1/repos/unfoldingword"
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
	// now remove them; go backwards so indexes are stable
	for i := len(positionsToRemove) - 1; i > -1; i-- {
		results = remove(results, positionsToRemove[i])
	}

	// sort it by file(1) and release(2)
	sort.Slice(results, func(i, j int) bool {
		if results[i][2] < results[j][2] {
			return true
		} else if results[i][2] > results[j][2] {
			return false
		}
		// equal filenames; test on release version
		ival, jval := results[i][1], results[j][1]
		comparevals := semver.Compare(ival, jval)
		if comparevals == 1 {
			return false
		}
		if comparevals == -1 {
			return true
		}
		// note: this cannot happen!
		log.Fatal("data has equal filename and release semver - not possible!")
		return true
	})

	// Filter out the dups, taking only the first occurrence of each SHA
	// start comparing on row 2 (zero based)
	var _results [][]string
	for i := 0; i < len(results); i++ {
		if i == 0 {
			// this is the header
			_results = append(_results, results[i])
			continue
		}
		if i == 1 {
			// first row, always take it
			_results = append(_results, results[i])
			continue
		}
		// if the SHA is the same as last row, discard it
		if results[i-1][3] == results[i][3] {
			//_results = append(_results, results[i])
			continue
		}
		// take the row
		_results = append(_results, results[i])
	}

	// finally add the revision column to each row
	revision := 0
	for i := 0; i < len(_results); i++ {
		if i == 0 {
			// add a column
			_results[i] = append(_results[i], "Revision")
			continue
		}
		if i == 1 {
			// this is first row, always rev 1
			revision++
			_results[i] = append(_results[i], strconv.Itoa(revision))
			continue
		}
		// if filename is unchanged from previous row
		if _results[i][2] == _results[i-1][2] {
			// still on same file, increment revision
			revision++
			_results[i] = append(_results[i], strconv.Itoa(revision))
			continue
		}
		// filename changed, reset revision to zero
		revision = 1
		_results[i] = append(_results[i], strconv.Itoa(revision))
	}
	output := *repo + "_revs.csv"
	f, err := os.Create(output)
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(_results) // calls Flush internally

	if err != nil {
		log.Fatal(err)
	}

	stop := time.Since(now)
	log.Printf("Done. %v", fmt.Sprintf("Elapsed Time: %v\n", stop))
}

func remove(slice [][]string, s int) [][]string {
	return append(slice[:s], slice[s+1:]...)
}
