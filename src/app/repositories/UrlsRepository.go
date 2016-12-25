package repositories

import (
	"app/models"
	"app/types"
	"app/utils"
	"configs"
	"encoding/json"
	"io"

	"github.com/jinzhu/gorm"
)

// Urls repository struct
type UrlsRepository struct {
	//
}

// Get long urls from given data
func (urlsRepository *UrlsRepository) GetLongUrls(body io.ReadCloser) types.Urls {

	decoder := json.NewDecoder(body)

	urls := types.Urls{}

	err := decoder.Decode(&urls)

	if err != nil {
		panic(err)
	}

	return urls
}

// Generate short urls from the given long urls
func (urlsRepository *UrlsRepository) GenerateShortUrls(Urls types.Urls) types.Urls {
	db := DbRepository{}.init()

	for key, urls := range Urls.LongUrls {
		shortUrlExist, shortUrl := urlsRepository.shortUrlExist(db, urls.LongUrl)

		if shortUrlExist {
			Urls.LongUrls[key].ShortUrl = shortUrl
		} else {
			Urls.LongUrls[key].ShortUrl = urlsRepository.generateShortUrl(db, urls.LongUrl)
		}
	}

	defer db.Close()

	return Urls
}

// Generate short url from the given long url
func (urlsRepository *UrlsRepository) generateShortUrl(db *gorm.DB, longUrl string) string {

	for {
		randomString := utils.RandomString(configs.SHORT_URL_STRING_SIZE)

		if db.Where("short_url = ?", randomString).First(&models.Url{}).Error != nil {
			url := models.Url{LongUrl: longUrl, ShortUrl: randomString}
			db.Create(&url)

			return randomString
		}
	}

}

// Determine short url existence
func (urlsRepository *UrlsRepository) shortUrlExist(db *gorm.DB, longUrl string) (bool, string) {
	var url models.Url
	query := db.Where("long_url = ?", longUrl).First(&url)

	if query.Error == nil {
		return true, url.ShortUrl
	}

	return false, ""
}

// Get long url from the given short url
func (urlsRepository *UrlsRepository) GetLongUrl(shortUrl string) (bool, string) {
	db := DbRepository{}.init()
	var url models.Url

	query := db.Where("short_url = ?", shortUrl).First(&url)

	defer db.Close()

	if query.Error == nil {
		return true, url.LongUrl
	}

	return false, ""
}
