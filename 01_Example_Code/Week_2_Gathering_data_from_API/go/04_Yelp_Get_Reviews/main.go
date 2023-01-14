package main

import (
	"bufio"         //File reading buffer
	"encoding/json" //This is used to interact with json responses and json files
	"fmt"           //The default module to print to console
	"io/ioutil"     // Reading input and saving output
	"net/http"      //used to call out to websites and APIs
	"os"            //Similar to python os
	"strconv"       //https://pkg.go.dev/strconv Somewhat painful to figure out which one to use
	"time"          //For date and times coming from json
)

// -----------------SETUP for the Yelp API json File Structure----------------------------
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

// ---------------Write summary results out to a json Struct files seprated by | incase commas are in any of the fields-------------------
type yelpReviewsResults struct {
	Reviews           []yelpReviewData `json:"reviews"`
	Total             int64            `json:"total"`
	PossibleLanguages []string         `json:"possible_languages"`
}

type yelpReviewData struct {
	ReviewId         string         `json:"id"`
	ReviewURL        string         `json:"url"`
	ReviewText       string         `json:"text"`
	ReviewRating     int32          `json:"rating"`
	ReviewCreationTS time.Time      `json:"time_created"`
	ReviewUser       yelpReviewUser `json:"user"`
}

type yelpReviewUser struct {
	UserId     string `json:"id"`
	UserUrl    string `json:"profile_url"`
	UserImgUrl string `json:"img_url"`
	UserName   string `json:"name"`
}

// -------------Main Program------------------------------------
func main() {

	groupoflocations := import_json_Businesses("BusinessSearchResults--Austin TX--dentists.json")

	TestIndexNumber := 0
	fmt.Println(groupoflocations.Businesses[TestIndexNumber].Name)
	fmt.Println(groupoflocations.Businesses[TestIndexNumber].Locations.Address1 + groupoflocations.Businesses[TestIndexNumber].Locations.Address2 + groupoflocations.Businesses[TestIndexNumber].Locations.City)
	fmt.Println(groupoflocations.Businesses[TestIndexNumber].ReviewCount)

	//BusID := groupoflocations.Businesses[TestIndexNumber].Id
	BusID := "NPCQRbpaO6w0cKswP2J7ig"
	var ResLimit int32 = 50
	var offset int32 = 0
	SortBy := "yelp_sort"
	Locale := "en_US"

	//BusID should be a valid Yelp Business ID,
	//ResLimit maxes out at 50
	//SortBY:  "yelp_sort" OR "newest"
	//locale for US people should probably always be: en_US
	MyJson := GetYelpAPIReviews(BusID, ResLimit, SortBy, Locale, offset)

	fmt.Println(MyJson)
}

func GetYelpAPIReviews(BusID string, ResLimit int32, SortBy string, Locale string, offset int32) (CurrentPageJsonResults yelpReviewsResults) {
	//https://tutorialedge.net/golang/parsing-json-with-golang/
	sliceKeys := read_text_file_line_by_line("cred_user_pwd_keys.txt")

	APIKey := sliceKeys[0]

	filename := "REVIEWS-" + BusID + "--.json"

	URLlimit := strconv.FormatInt(int64(ResLimit), 10) //40 is the max results we can get from the API
	URLoffset := strconv.FormatInt(int64(offset), 10)  //40 is the max results we can get from the API

	url := "https://api.yelp.com/v3/businesses/" + BusID + "/reviews?locale=" + Locale + "&offset=" + URLoffset + "&limit=" + URLlimit + "&sort_by=" + SortBy

	fmt.Println(url)

	//url := "https://api.yelp.com/v3/businesses/search?latitude=" + latitude + "&longitude=" + longitude + "&radius=" + radius + "&categories=" + categories + "&sort_by=" + sort_by + "&limit=" + limit

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println(string(body))

	json.Unmarshal(body, &CurrentPageJsonResults)

	file, _ := json.MarshalIndent(CurrentPageJsonResults, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

	return
}

func import_json_Businesses(FileLocationString string) (jsonStruct yelpBusSearch) {
	//https://tutorialedge.net/golang/parsing-json-with-golang/

	// Open our jsonFile
	jsonFile, err := os.Open(FileLocationString)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	//var yelplocs locations

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &jsonStruct)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	//fmt.Println("Made it to this line")
	for i := 0; i < len(jsonStruct.Businesses); i++ {
		fmt.Println("Business Yelp ID: " + jsonStruct.Businesses[i].Id)
		fmt.Println("   Alias: " + jsonStruct.Businesses[i].Alias)

	}
	return
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

func extractYelpSearchBusiness(CityName string, CityLat float64, CityLong float64, Category string) (jsonstatus string, maxDistence float64, minDistence float64, maxRating float64, minRating float64, PotentialMatches int64, maxqtyVotes int64, minqtyVotes int64) {
	//https://tutorialedge.net/golang/parsing-json-with-golang/
	sliceKeys := read_text_file_line_by_line("cred_user_pwd_keys.txt")

	APIKey := sliceKeys[0]

	// Open our jsonFile
	//jsonFile, err := os.Open(FileLocationString)
	// if we os.Open returns an error then handle it
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	//defer jsonFile.Close()
	var jsonResponse yelpBusSearch
	// read our opened jsonFile as a byte array.
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	filename := "BusinessSearchResults--" + CityName + "--" + Category + ".json"
	// we initialize our Users array
	//var yelplocs locations

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	//json.Unmarshal(byteValue, &results)

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

	//fmt.Println(res)
	//fmt.Println(string(body))

	json.Unmarshal(body, &jsonResponse)

	file, _ := json.MarshalIndent(jsonResponse, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

	jsonstatus = "200" //TO DO figure out how to pull this from Go
	PotentialMatches = jsonResponse.Total
	//Because we sort by vote count we can just get the first and last entry for the most and least amount of votes
	maxqtyVotes = jsonResponse.Businesses[0].ReviewCount
	minqtyVotes = jsonResponse.Businesses[len(jsonResponse.Businesses)-1].ReviewCount

	//Set these as out of bounds initial numbers
	maxDistence = 0
	minDistence = 100000
	maxRating = 0
	minRating = 100
	//Loop over all the Businesses returned ues comparison to get the highest and lowest values
	for i := 0; i < len(jsonResponse.Businesses); i++ {
		if jsonResponse.Businesses[i].Distance > maxDistence {
			maxDistence = jsonResponse.Businesses[i].Distance
		}

		if jsonResponse.Businesses[i].Distance < minDistence {
			minDistence = jsonResponse.Businesses[i].Distance
		}

		if jsonResponse.Businesses[i].Rating > maxRating {
			maxRating = jsonResponse.Businesses[i].Rating
		}

		if jsonResponse.Businesses[i].Rating < minRating {
			minRating = jsonResponse.Businesses[i].Rating
		}
	}
	return
}
