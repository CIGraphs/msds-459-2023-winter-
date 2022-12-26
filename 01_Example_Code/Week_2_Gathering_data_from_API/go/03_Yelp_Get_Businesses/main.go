package main

import (
	"bufio" //File reading buffer
	"encoding/csv"
	"encoding/json" //This is used to interact with json responses and json files
	"fmt"           //The default module to print to console
	"io/ioutil"     // Reading input and saving output
	"net/http"      //used to call out to websites and APIs
	"os"            //Similar to python os
	"strconv"       //https://pkg.go.dev/strconv Somewhat painful to figure out which one to use
)

// -------------------------SETUP for Yelp_locations.json file--------------------------
// Define what the Yelp_locations.json file looks like as a structure
type locations struct {
	Locations []loc `json:"locations"`
}

// Each specifc location in the Yelp_locations.json is made up of these types
// FYI go / json package is case sensitive the first lette of the variable must be capital to work with json
type loc struct {
	Manualname string  `json:"manual_name"`
	Latitude   float64 `json:"Latitude"`
	Longitude  float64 `json:"Longitude"`
}

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

// ---------------Write summary results out to a CSV files seprated by | incase commas are in any of the fields-------------------
type csvSummary struct {
	webpullstatus    string
	fileName         string
	City             string
	Latitude         float64
	Longitude        float64
	CategoryAlias    string
	maxDistence      float64
	minDistence      float64
	maxRating        float64
	minRating        float64
	MaxQtyVotes      int64
	MinQtyVotes      int64
	PotentialMatches int64
}

// -------------Main Program------------------------------------
func main() {
	//Get the contents of Yelp_locations.json
	groupoflocations := import_json_locations("Yelp_Locations.json")
	//Uncomment for troubleshooting
	fmt.Println(groupoflocations)

	//Read the input textfile categories.txt into a "slice"
	sliceofcategories := read_text_file_line_by_line("categories.txt")
	//Uncomment for troubleshooting
	fmt.Println(sliceofcategories)

	//YelpSearchBusiness := extractYelpSearchBusiness("BusinessSearchResultsSan Francisco CA--wineries.json")

	//fmt.Println(YelpSearchBusiness)

	//Create a file object with json MarshalIndent, write the file to a file
	//file, _ := json.MarshalIndent(YelpSearchBusiness, "", " ")
	//_ = ioutil.WriteFile("test.json", file, 0644)

	//Setup a 2d slice for a CSV extract with a header row
	headerRow := []string{
		"fileName", "maxDistence", "minDistence", "maxRating", "minRating", "webpullstatus", "CategoryAlias", "City", "Latitude", "Longitude", "MaxQtyVotes", "MinQtyVotes", "PotentialMatches",
	}

	summaryResultsSlide := [][]string{
		headerRow,
	}

	fmt.Println("Loop Through Locations & Categories")
	for i := 0; i < len(groupoflocations.Locations); i++ {
		fmt.Println("Location Name: " + groupoflocations.Locations[i].Manualname)

		for j := 0; j < len(sliceofcategories); j++ {
			if sliceofcategories[j] != "" {
				var LocationCategoryRow csvSummary
				fmt.Println("    " + sliceofcategories[j])

				CityName := groupoflocations.Locations[i].Manualname
				CityLat := groupoflocations.Locations[i].Latitude
				CityLong := groupoflocations.Locations[i].Longitude
				Category := sliceofcategories[j]

				jsonstatus, maxDistence, minDistence, maxRating, minRating, PotentialMatches, MaxQtyVotes, MinQtyVotes := extractYelpSearchBusiness(CityName, CityLat, CityLong, Category)

				LocationCategoryRow.fileName = "BusinessSearchResults--" + CityName + "--" + Category + ".json"
				LocationCategoryRow.maxDistence = maxDistence
				LocationCategoryRow.minDistence = minDistence
				LocationCategoryRow.maxRating = maxRating
				LocationCategoryRow.minRating = minRating
				LocationCategoryRow.webpullstatus = jsonstatus
				LocationCategoryRow.CategoryAlias = Category
				LocationCategoryRow.City = CityName
				LocationCategoryRow.Latitude = CityLat
				LocationCategoryRow.Longitude = CityLong
				LocationCategoryRow.MaxQtyVotes = MaxQtyVotes
				LocationCategoryRow.MinQtyVotes = MinQtyVotes
				LocationCategoryRow.PotentialMatches = PotentialMatches

				//Reading/writing CSV files https://articles.wesionary.team/read-and-write-csv-file-in-go-b445e34968e9
				summaryResultsSlide = append(summaryResultsSlide, []string{
					LocationCategoryRow.fileName,
					strconv.FormatFloat(maxDistence, 'E', -1, 64),
					strconv.FormatFloat(LocationCategoryRow.minDistence, 'E', -1, 64),
					strconv.FormatFloat(LocationCategoryRow.maxRating, 'E', -1, 64),
					strconv.FormatFloat(LocationCategoryRow.minRating, 'E', -1, 64),
					LocationCategoryRow.webpullstatus,
					LocationCategoryRow.CategoryAlias,
					LocationCategoryRow.City,
					strconv.FormatFloat(LocationCategoryRow.Latitude, 'E', -1, 64),
					strconv.FormatFloat(LocationCategoryRow.Longitude, 'E', -1, 64),
					strconv.FormatInt(LocationCategoryRow.MaxQtyVotes, 10),
					strconv.FormatInt(LocationCategoryRow.MinQtyVotes, 10),
					strconv.FormatInt(LocationCategoryRow.PotentialMatches, 10),
				})
			}

		}
		//fmt.Println("   Latitude: " + strconv.FormatFloat(groupoflocations.Locations[i].Latitude, 'E', -1, 32))
		//fmt.Println("   Longitude: " + strconv.FormatFloat(groupoflocations.Locations[i].Longitude, 'E', -1, 32))
		//fmt.Println("User Name: " + strconv.Itoa(users.Users[i].Longitude)) # if converting INT to string
		//fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook) if pulling more nexsted info
	}

	csvFile, err := os.Create("APIpull_status.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	err = csvwriter.WriteAll(summaryResultsSlide) // calls Flush internally
	if err != nil {
		fmt.Println(err)
	}

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

func import_json_locations(FileLocationString string) (jsonlocations locations) {
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
	json.Unmarshal(byteValue, &jsonlocations)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	//fmt.Println("Made it to this line")
	for i := 0; i < len(jsonlocations.Locations); i++ {
		fmt.Println("Location Name: " + jsonlocations.Locations[i].Manualname)
		fmt.Println("   Latitude: " + strconv.FormatFloat(jsonlocations.Locations[i].Latitude, 'E', -1, 32))
		fmt.Println("   Longitude: " + strconv.FormatFloat(jsonlocations.Locations[i].Longitude, 'E', -1, 32))
		//fmt.Println("User Name: " + strconv.Itoa(users.Users[i].Longitude)) # if converting INT to string
		//fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook) if pulling more nexsted info
	}
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
