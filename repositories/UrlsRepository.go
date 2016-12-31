package repositories

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/iftekhersunny/shortx/configs"
	"github.com/iftekhersunny/shortx/models"
	"github.com/iftekhersunny/shortx/utils"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Urls repository struct
type UrlsRepository struct {

	// urls channel max length
	ChannelMaxLength int
}

// Get long urls from given data
func (urlsRepository *UrlsRepository) GetLongUrls(body io.ReadCloser) map[string]string {

	decoder := json.NewDecoder(body)

	var urls map[string]map[string]string

	err := decoder.Decode(&urls)

	if err != nil {
		panic(err)
	}

	return urls["long_urls"]
}

// Generate short urls from the given long urls
func (urlsRepository *UrlsRepository) GenerateShortUrls(Urls map[string]string, request *http.Request) (map[string]string, bool) {

	// database instance
	db, session := DbRepository{}.init()
	urlsRepository.ChannelMaxLength = len(Urls)

	// channels
	urls := make(chan map[string]string, urlsRepository.ChannelMaxLength)
	done := make(chan bool)

	// running go routines
	for longUrl, _ := range Urls {
		shortUrlExist, shortUrl := make(chan bool), make(chan string)

		go urlsRepository.shortUrlExist(db, longUrl, shortUrlExist, shortUrl)

		go urlsRepository.sendShortUrl(db, request, urls, longUrl, shortUrlExist, shortUrl, done)
	}

	select {
	// wait for completing all goroutines execution
	case <-done:
		session.Close()
		return urlsRepository.receiveData(urls), true

	// if all goroutines did not complete within 30 seconds
	// send a request timeout message to the client
	case <-time.After(time.Second * 30):
		session.Close()
		return map[string]string{"message": "Request timeout"}, false
	}
}

// Receive data from urls channel
func (urlsRepository *UrlsRepository) receiveData(urls chan map[string]string) map[string]string {
	data := make(map[string]string)

	close(urls)

	for url := range urls {
		for key, value := range url {
			data[key] = value
		}
	}

	return data
}

// Send short url into urls channel
func (urlsRepository *UrlsRepository) sendShortUrl(
	db *mgo.Database, request *http.Request, urls chan map[string]string, longUrl string,
	shortUrlExist chan bool, shortUrl chan string, done chan bool) {
	
	if <-shortUrlExist {
		urls <- map[string]string{longUrl: "http://" + request.Host + "/" + <-shortUrl}
	} else {
		urls <- map[string]string{longUrl: "http://" + request.Host + "/" + urlsRepository.generateShortUrl(db, longUrl)}
	}

	if len(urls) == urlsRepository.ChannelMaxLength {
		done <- true
	}
}

// Generate short url from the given long url
func (urlsRepository *UrlsRepository) generateShortUrl(db *mgo.Database, longUrl string) string {
	for {
		randomString := utils.RandomString(configs.SHORT_URL_STRING_SIZE)
		collection := db.C("urls")
		result := models.Url{}
		err := collection.Find(bson.M{"shorturl": randomString}).One(&result)
		if err != nil {
			url := models.Url{LongUrl: longUrl, ShortUrl: randomString}
			collection.Insert(&url)

			return randomString
		}
	}
}

// Determine short url existence
func (urlsRepository *UrlsRepository) shortUrlExist(db *mgo.Database, longUrl string, shortUrlExist chan bool, shortUrl chan string) {
	url := models.Url{}
	collection := db.C("urls")

	err := collection.Find(bson.M{"longurl": longUrl}).One(&url)

	if err == nil {
		shortUrlExist <- true
		shortUrl <- url.ShortUrl
	}

	shortUrlExist <- false
	shortUrl <- ""
}

// Get long url from the given short url
func (urlsRepository *UrlsRepository) GetLongUrl(shortUrl string) (bool, string) {
	db, session := DbRepository{}.init()
	var url models.Url

	collection := db.C("urls")
	err := collection.Find(bson.M{"shorturl": shortUrl}).One(&url)
	session.Close()

	if err == nil {
		return true, url.LongUrl
	}

	return false, ""
}
