package geolocation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Geo struct {
	Apikey string
}

type GeoResponse struct {
	Data []Info `json:"data"`
}

type Info struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (g *Geo) GetLatAndLon(city string) (float64, float64, error) {
	url := fmt.Sprintf("http://api.positionstack.com/v1/forward?access_key=%s&query=%s", g.Apikey, city)
	url = strings.Replace(url, " ", "%20", -1)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("error createing request %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("error sending request %v", err)
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading body %v", err)
	}
	fmt.Println(string(bodyBytes))

	var data GeoResponse
	_ = json.Unmarshal(bodyBytes, &data)
	if data.Data != nil {
		return data.Data[0].Latitude, data.Data[0].Longitude, nil
	} else {
		return g.GetLatAndLon(city)
	}

}
