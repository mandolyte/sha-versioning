package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
					nodeMap := tval.(map[string]interface{})
					fmt.Printf("%v,%v,%v,%v\n", *repo, rversions[rv], nodeMap["path"], nodeMap["sha"])
				}
			}
		}
	}
	log.Println("Done.")
}
