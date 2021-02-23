package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

const (
	indexName = "tempusers2"
	typeName  = "store"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type User struct {
	Loc      Location `json:"location"`
	Username string   `json:"username"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
}

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	if err != nil {
		log.Fatal(err)
	}

	// Setup the index
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		// Drop the index
		_, err = client.DeleteIndex(indexName).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Create the index with a mapping
	_, err = client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Add a point
	_, err = client.Index().Index(indexName).Id("1").BodyJson(&User{
		Loc:      Location{Lat: 45.750389, Lon: 19.699748},
		Username: "Marienplatz",
		Country:  "Serbia",
		City:     "Novi Sad",
	}).Refresh("true").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	distanceQuery := elastic.NewGeoDistanceQuery("location")
	distanceQuery = distanceQuery.Lat(43.750389)
	distanceQuery = distanceQuery.Lon(19.699748)
	distanceQuery = distanceQuery.Distance("50km")

	geoDistanceSorter := elastic.NewGeoDistanceSort("location").
		Point(48.81433999999999, 2.3204906).
		Unit("km").
		GeoDistance("plane").
		Asc()
	_ = geoDistanceSorter
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(elastic.NewMatchAllQuery())
	boolQuery = boolQuery.Filter(distanceQuery)

	src, err := boolQuery.Source()
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(src)
	if err != nil {
		log.Fatal(err)
	}

	got := string(data)

	fmt.Println(got)

	searchResult, err := client.Search().Index(indexName).Query(boolQuery).Do(ctx) //SortBy(geoDistanceSorter)
	if err != nil {
		log.Fatal(err)
	}

	if searchResult.Hits.TotalHits.Value > 0 {
		fmt.Printf("Found a total of %d stores\n", searchResult.Hits.TotalHits.Value)

		for _, hit := range searchResult.Hits.Hits {
			var s User
			err := json.Unmarshal(hit.Source, &s)
			if err != nil {
				// deserialization failed
				log.Fatal(err)
			}

			fmt.Printf("Store from - Country : " + s.Country + " ; City : " + s.City)
		}
	} else {
		fmt.Print("Found no stores\n")
	}
}

const mapping = `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			
					"username":{
						"type":"text"
					},
					"city":{
						"type":"text"
					},
					"country":{
						"type":"text"
					},
					"location":{
						"type":"geo_point"
					}
		
		}
	}
}`
