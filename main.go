package main

import (
	"crypto/md5"
	"fmt"
	"time"
	"os"
	"strconv"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)


func main() {
	apikey := os.Getenv("APIKEY")
	privatekey := os.Getenv("PRIVATEKEY")
	t := int(time.Now().Unix())
	data := []byte(strconv.Itoa(t)+privatekey+apikey)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	serv := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body := getCharacters(t, hash, apikey)
			go getCharactersbyID(body, t, hash, apikey)
		} else {
			w.WriteHeader(405)
			fmt.Fprintf(w, "Method not Allowed\n")
		}
	}
	http.HandleFunc("/", serv)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCharacters(t int, hash string, apikey string) []byte {
	url := fmt.Sprintf("https://gateway.marvel.com/v1/public/characters?ts=%v&apikey=%v&hash=%v", t, apikey, hash)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	return body
}

func getCharactersbyID(body []byte, t int, hash string, apikey string) {
	var read result
	err := json.Unmarshal(body, &read)
	if err != nil {
		fmt.Println("error:", err)
	}
	for i := 0; i < len(read.Data.Results); i++ {
		id := read.Data.Results[i].ID
		urlChar := fmt.Sprintf("https://gateway.marvel.com/v1/public/characters/%v?ts=%v&apikey=%v&hash=%v", id, t, apikey, hash)
		resp, err := http.Get(urlChar)
		if err != nil {
			log.Fatal(err)
		}
		bodyChar, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
		}
		req, err := http.Post("https://c44b7111-0323-4bc6-aaa4-90719a420d9c.mock.pstmn.io", "application/json", bytes.NewBuffer(bodyChar))
		if err != nil {
			fmt.Println("error:", err)
		}
		response := req.StatusCode
		req.Body.Close()
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("response Status:", response)
	}
}


type result struct {
	Code            int64  `json:"code"`           
	Status          string `json:"status"`         
	Copyright       string `json:"copyright"`      
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Etag            string `json:"etag"`           
	Data            Data   `json:"data"`           
}

type Data struct {
	Offset  int64    `json:"offset"` 
	Limit   int64    `json:"limit"`  
	Total   int64    `json:"total"`  
	Count   int64    `json:"count"`  
	Results []Result `json:"results"`
}

type Result struct {
	ID          int64     `json:"id"`         
	Name        string    `json:"name"`       
	Description string    `json:"description"`
	Modified    string    `json:"modified"`   
	Thumbnail   Thumbnail `json:"thumbnail"`  
	ResourceURI string    `json:"resourceURI"`
	Comics      Comics    `json:"comics"`     
	Series      Comics    `json:"series"`     
	Stories     Stories   `json:"stories"`    
	Events      Comics    `json:"events"`     
	Urls        []URL     `json:"urls"`       
}

type Comics struct {
	Available     int64        `json:"available"`    
	CollectionURI string       `json:"collectionURI"`
	Items         []ComicsItem `json:"items"`        
	Returned      int64        `json:"returned"`     
}

type ComicsItem struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`       
}

type Stories struct {
	Available     int64         `json:"available"`    
	CollectionURI string        `json:"collectionURI"`
	Items         []StoriesItem `json:"items"`        
	Returned      int64         `json:"returned"`     
}

type StoriesItem struct {
	ResourceURI string   `json:"resourceURI"`
	Name        string   `json:"name"`       
	Type        ItemType `json:"type"`       
}

type Thumbnail struct {
	Path      string    `json:"path"`     
	Extension Extension `json:"extension"`
}

type URL struct {
	Type URLType `json:"type"`
	URL  string  `json:"url"` 
}

type ItemType string
const (
	Cover ItemType = "cover"
	Empty ItemType = ""
	InteriorStory ItemType = "interiorStory"
)

type Extension string
const (
	GIF Extension = "gif"
	Jpg Extension = "jpg"
)

type URLType string
const (
	Comiclink URLType = "comiclink"
	Detail URLType = "detail"
	Wiki URLType = "wiki"
)