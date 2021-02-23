package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bokimilinkovic/upp/model"
	"github.com/olivere/elastic/v7"
)

const mapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"testic":{
			"properties":{
				"isbn":{
					"type":"keyword"
				},
				"publishyear":{
					"type":"date"
				},
				"pagenumber":{
					"type":"integer"
				},
				"genre":{
					"type":"text"
				},
				"author":{
					"type":"text"
				},
				"title":{
					"type":"text"
				},
				"content":{
					"type":"text"
				}
			}
		}
	}
}
`

func main() {
	esClient, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)

	if err != nil {
		panic(err)
	}
	exists, err := esClient.IndexExists("testic").Pretty(true).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if exists {
		_, err := esClient.DeleteIndex("testic").Pretty(true).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	_, err = esClient.CreateIndex("testic").Body(mapping).Do(context.Background())
	if err != nil {
		panic(err)
	}

	// Add a book
	book := model.Book{
		ISBN:        "123123",
		PublishYear: time.Now(),
		PageNumber:  12,
		Genre:       "drama",
		Author:      "Boki Milinkovic",
		Title:       "Lepa knjiga",
		Content:     "lepa u 3 lepe",
	}

	_, err = esClient.Index().
		Index("testic").
		Id("1").
		BodyJson(&book).
		Refresh("true").
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	// Read the book
	doc, err := esClient.Get().
		Index("testic").
		Type("_doc").
		Id("1").
		Pretty(true).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	var book2 model.Book
	if err = json.Unmarshal(doc.Source, &book2); err != nil {
		panic(err)
	}
	fmt.Println(book)
}
