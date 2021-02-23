package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bokimilinkovic/upp/geolocation"
	"github.com/bokimilinkovic/upp/handler/dto"
	"github.com/bokimilinkovic/upp/model"
	"github.com/bokimilinkovic/upp/pdf"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic/v7"
)

type Elastic struct {
	esClient *elastic.Client
	esPort   string
	geoLoc   *geolocation.Geo
}

func NewElastic(esPort string, geo *geolocation.Geo) *Elastic {
	esClient, err := elastic.NewClient(
		elastic.SetURL(esPort),
		elastic.SetSniff(false),
	)

	if err != nil {
		return nil
	}

	return &Elastic{esClient: esClient, esPort: esPort, geoLoc: geo}
}

func (e *Elastic) CheckConnection(c echo.Context) error {
	info, code, err := e.esClient.Ping(e.esPort).Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	text := fmt.Sprintf("elasticsearch reruned wit hcode %d and version %v\n", code, info)
	return c.String(http.StatusOK, text)
}

func (e *Elastic) CreateIndex(c echo.Context) error {
	exists, err := e.esClient.IndexExists("books").Do(c.Request().Context())
	if err != nil {
		return err
	}
	if !exists {
		_, err = e.esClient.CreateIndex("books").Do(c.Request().Context())
		if err != nil {
			return err
		}

		//Add a book to the index
		book := model.Book{
			ISBN:        "1234231",
			PublishYear: time.Date(1980, time.April, 20, 1, 1, 1, 1, time.Local),
			PageNumber:  190,
			Genre:       "comedy",
			Author:      "John Doe",
			Title:       "Funny Summer",
			Content:     "Some cool content",
		}

		_, err = e.esClient.Index().
			Index("books").
			Id("1").
			BodyJson(book).
			Refresh("wait_for").
			Do(c.Request().Context())

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating book index: %v", err))
		}
	}

	return c.String(http.StatusOK, "All good")
}

func (e *Elastic) Search(c echo.Context) error {
	term := c.QueryParam("term")
	//q2 := elastic.NewMatchQuery("author", term)
	matchMultiQuery := elastic.NewMultiMatchQuery(term, "author", "genre", "title", "content", "isbn")
	sr, err := e.esClient.Search().
		Index("books").
		Query(matchMultiQuery).
		From(0).
		Size(5).
		Pretty(true).
		Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("can not search index %v", err))
	}

	fmt.Printf("Query took %d milliseconds\n", sr.TookInMillis)

	books := make([]model.Book, 0)
	for _, hit := range sr.Hits.Hits {
		var book model.Book
		err := json.Unmarshal(hit.Source, &book)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		fmt.Println("BOOK ", book.Author)
		books = append(books, book)
	}

	fmt.Printf("found a total of %d books \n", len(sr.Hits.Hits))

	return c.JSON(http.StatusOK, books)
}

func (e *Elastic) AddBook(c echo.Context) error {
	var bookDto dto.BookDto

	if err := c.Bind(&bookDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad payload provided: "+err.Error())
	}
	fmt.Println(bookDto.Genre)
	publishYear, err := time.Parse(time.RFC3339, bookDto.PublishYear)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad date provided : "+err.Error())
	}

	book := &model.Book{
		ISBN:        bookDto.ISBN,
		PublishYear: publishYear,
		PageNumber:  bookDto.PageNumber,
		Genre:       bookDto.Genre,
		Author:      bookDto.Author,
		Title:       bookDto.Title,
	}

	pdfReader := &pdf.Reader{}
	content, err := pdfReader.ReadPdf("books/" + book.Title + ".pdf")
	if err != nil {
		fmt.Errorf("error reading content " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error reading pdf content "+err.Error())
	}

	book.Content = content
	fmt.Println(book)
	// add book to index
	dataJson, err := json.Marshal(book)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error marshaling book : ", err.Error())
	}

	js := string(dataJson)
	_, err = e.esClient.Index().
		Index("books").
		BodyJson(js).Do(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error inserting book to index :"+err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (e *Elastic) GetAllBooks(c echo.Context) error {
	res, err := e.esClient.Search().Index("books").Highlight(elastic.NewHighlight()).Size(10).Pretty(true).Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error connecting to elastic "+err.Error())
	}

	books := []model.Book{}
	for _, hit := range res.Hits.Hits {
		var book model.Book
		err := json.Unmarshal(hit.Source, &book)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error unmarhslaing book "+err.Error())
		}
		books = append(books, book)
	}

	return c.JSON(http.StatusOK, books)
}

func (e *Elastic) HighlightSearch(c echo.Context) error {
	query := elastic.NewBoolQuery()
	query = query.Should(elastic.NewTermQuery("author", "doe"))
	query = query.Should(elastic.NewTermQuery("content", "cool"))

	highlight := elastic.NewHighlight()
	highlight = highlight.Fields(elastic.NewHighlighterField("author"), elastic.NewHighlighterField("content"))
	highlight = highlight.PreTags("<b>").PostTags("</b>")

	searchResult, err := e.esClient.Search().
		Index("books").
		Highlight(highlight).
		Query(query).
		Pretty(true).
		Do(c.Request().Context())

	if err != nil {
		return err
	}

	hit := searchResult.Hits.Hits[0]
	var book model.Book
	if err = json.Unmarshal(hit.Source, &book); err != nil {
		return err
	}
	var highlitghed string
	for key, _ := range hit.Highlight {
		for _, v := range hit.Highlight[key] {
			highlitghed = highlitghed + v
		}
		highlitghed += "<br></br>"
	}

	return c.HTML(http.StatusOK, highlitghed)
}

// CreateUserIndex checks if index exists, dropping it and recreating new one
func (e *Elastic) CreateUserIndex(c echo.Context) error {
	// Setup the index
	exists, err := e.esClient.IndexExists("users").Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "error checking index :%v", err.Error())
	}

	if exists {
		// Drop the index
		_, err = e.esClient.DeleteIndex("users").Do(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "error deleting index :%v", err.Error())
		}
	}
	// Create the index with a mapping
	_, err = e.esClient.CreateIndex("users").BodyString(userMapping).Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "error creting index :%v", err.Error())
	}
	// Add a point
	_, err = e.esClient.Index().Index("users").Id("1").BodyJson(&model.User{
		Loc:      model.Location{Lat: 45.750389, Lon: 19.699748},
		Username: "Novosadjanin",
		Country:  "Serbia",
		City:     "Novi Sad",
	}).Refresh("true").Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "error creating point in index :%v", err.Error())
	}

	return c.String(http.StatusOK, "index created!")
}

func (e *Elastic) GetAllUsers(c echo.Context) error {
	res, err := e.esClient.Search().Index("users").Highlight(elastic.NewHighlight()).Size(10).Pretty(true).Do(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error connecting to elastic "+err.Error())
	}

	users := []model.User{}
	for _, hit := range res.Hits.Hits {
		var user model.User
		err := json.Unmarshal(hit.Source, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error unmarhslaing user "+err.Error())
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

func (e *Elastic) AddUser(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad payload provided: "+err.Error())
	}

	lat, lon, err := e.geoLoc.GetLatAndLon(user.City)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can not find that city, try another one! %s", err.Error())
	}
	fmt.Println(lat, lon)
	user.Loc = model.Location{Lat: lat, Lon: lon}

	// add book to index
	dataJson, err := json.Marshal(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error marshaling user : ", err.Error())
	}

	js := string(dataJson)
	_, err = e.esClient.Index().
		Index("users").
		BodyJson(js).Do(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error inserting user to index :"+err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (e *Elastic) SearchDistance(c echo.Context) error {
	var searchDistanceDto dto.SearchByDistanceDto
	if err := c.Bind(&searchDistanceDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad payload provided %v", err.Error())
	}

	// Get geo location of wanted city
	lat, lon, err := e.geoLoc.GetLatAndLon(searchDistanceDto.City)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request %v", err.Error())
	}
	// This lat and lon are starting point, we are sending range query from these points.

	distanceQuery := elastic.NewGeoDistanceQuery("location")
	distanceQuery = distanceQuery.Lat(lat)
	distanceQuery = distanceQuery.Lon(lon)
	distanceQuery = distanceQuery.Distance(fmt.Sprintf("%dkm", searchDistanceDto.Range)) // 50km, 100km...

	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(elastic.NewMatchAllQuery(), distanceQuery)
	boolQuery = boolQuery.Filter(distanceQuery)

	src, err := boolQuery.Source()
	if err != nil {
		return err
	}
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	searchResult, err := e.esClient.Search().Index("users").Query(boolQuery).Do(c.Request().Context())
	if err != nil {
		return fmt.Errorf("Error searching for users %s", err.Error())
	}

	var users []model.User
	if searchResult.Hits.TotalHits.Value > 0 {
		fmt.Printf("Found a total of %d users\n", searchResult.Hits.TotalHits.Value)

		for _, hit := range searchResult.Hits.Hits {
			var u model.User
			err := json.Unmarshal(hit.Source, &u)
			if err != nil {
				// deserialization failed
				return err
			}

			fmt.Printf("User from - City : " + u.City + " ; Country : " + u.Country)
			users = append(users, u)
		}
	} else {
		fmt.Print("Found no users\n")
	}

	return c.JSON(http.StatusOK, users)
}

func bookToString(book *model.Book) string {
	return fmt.Sprintf("Title: %s \n Content: %s \n Author:%s\n", book.Title, book.Content, book.Author)
}

const userMapping = `{
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
