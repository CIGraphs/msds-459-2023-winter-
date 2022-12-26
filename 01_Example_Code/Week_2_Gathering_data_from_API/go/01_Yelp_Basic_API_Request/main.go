package main

import (
	"bufio"         //File reading buffer
	"encoding/json" //This is used to interact with json responses and json files
	"fmt"           //The default module to print to console
	"io/ioutil"     // Reading input and saving output
	"net/http"      //used to call out to websites and APIs
	"os"            //Similar to python os
	"strconv"       //https://pkg.go.dev/strconv Somewhat painful to figure out which one to use
)

// -----------------SETUP for the Yelp API Response----------------------------
type yelpBusSearch struct {
	Businesses []yelpbusinesses `json:"businesses"`
	Total      int64            `json:"total"`
	Region     yelpRegion       `json:"region"`
}

type yelpRegion struct {
	Center yelpLatLong `json:"center"`
}

type yelpLatLong struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"Latitude"`
}

type yelpbusinesses struct {
	Id           string      `json:"id"`
	Alias        string      `json:"alias"`
	Name         string      `json:"name"`
	ImageUrl     string      `json:"image_url"`
	IsClosed     bool        `json:"is_closed"`
	Url          string      `json:"url"`
	ReviewCount  int64       `json:"review_count"`
	Categories   []yelpCat   `json:"categories"`
	Rating       float64     `json:"rating"`
	Coordinates  yelpLatLong `json:"coordinates"`
	Transactions []string    `json:"transactions"`
	Price        string      `json:"price"`
	Locations    yelpAddress `json:"location"`
	Phone        string      `json:"phone"`
	DispPhone    string      `json:"display_phone"`
	Distance     float64     `json:"distance"`
}

type yelpCat struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

type yelpAddress struct {
	Address1    string   `json:"address1"`
	Address2    string   `json:"address2"`
	Address3    string   `json:"address3"`
	City        string   `json:"city"`
	Zipcode     string   `json:"zip_code"`
	Country     string   `json:"country"`
	State       string   `json:"state"`
	DispAddress []string `json:"display_address"`
	Phone       string   `json:"phone"`
	DispPhone   string   `json:"display_phone"`
	Distance    float64  `json:"distance"`
}

// -------------Main Program------------------------------------
func main() {

	CityName := "TEST-Seattle-WA-"
	CityLat := 47.6062
	CityLong := -122.3321
	Category := "hair"

	jsonstatus := extractYelpSearchBusiness(CityName, CityLat, CityLong, Category)

	fmt.Println(jsonstatus.Region.Center)

	fmt.Println(jsonstatus.Businesses[3].Alias)
	fmt.Println(jsonstatus.Businesses[3].Categories)

	fmt.Println(jsonstatus.Businesses[0].Alias)
	fmt.Println(jsonstatus.Businesses[0].ReviewCount)

}

func read_text_file_line_by_line(FileLocationString string) (SliceOfStrings []string) {

	f, err := os.Open(FileLocationString) //Opening a text file

	//error displaying
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	//"Scan the file line by line"
	for scanner.Scan() {

		fmt.Println(scanner.Text())
		SliceOfStrings = append(SliceOfStrings, scanner.Text())
	}
	//error displaying
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	//SliceOfStrings = append(SliceOfStrings, "This")
	return
}

func extractYelpSearchBusiness(CityName string, CityLat float64, CityLong float64, Category string) (jsonResponse yelpBusSearch) {
	//https://tutorialedge.net/golang/parsing-json-with-golang/
	sliceKeys := read_text_file_line_by_line("cred_user_pwd_keys.txt")

	APIKey := sliceKeys[0]

	filename := "BusinessSearchResults--" + CityName + "--" + Category + ".json"

	RadiusInMeters := 40000 //40000 is the max
	latitude := fmt.Sprintf("%f", CityLat)
	//longitude := strconv.FormatFloat(CityLong, 'E', -1, 32) #Note formatfloat with strconv uses scientific notation which will NOT WORK in the URL
	longitude := fmt.Sprintf("%f", CityLong)
	limit := strconv.FormatInt(40, 10)     //40 is the max results we can get from the API
	radius := strconv.Itoa(RadiusInMeters) //40 is the max results we can get from the API
	categories := Category
	sort_by := "review_count"

	url := "https://api.yelp.com/v3/businesses/search?latitude=" + latitude + "&longitude=" + longitude + "&radius=" + radius + "&categories=" + categories + "&sort_by=" + sort_by + "&limit=" + limit

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	json.Unmarshal(body, &jsonResponse)

	file, _ := json.MarshalIndent(jsonResponse, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644) // the 0644 thing I think is a unix/linux specific option??

	return
}
