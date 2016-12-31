package controllers

import (
	netHttp "net/http"

	"github.com/asaskevich/govalidator"
	httpClient "github.com/ddliu/go-httpclient"
	"github.com/gorilla/mux"
	"github.com/iftekhersunny/shortx/configs"
	"github.com/iftekhersunny/shortx/http"
	"github.com/iftekhersunny/shortx/repositories"
)

// Urls controller struct
type UrlsController struct {
	Controller
}

// Index method
func (urlsController *UrlsController) Index(writer netHttp.ResponseWriter, request *netHttp.Request) {

	// validating api token
	apiToken := request.URL.Query()["api_token"][0]

	if apiToken != configs.API_KEY {
		http.Response(writer, netHttp.StatusUnauthorized, map[string]string{"message": "Unauthorized access"})
		return
	}

	// get long urls from request.body
	UrlsRepository := repositories.UrlsRepository{}
	Urls := UrlsRepository.GetLongUrls(request.Body)

	// validate urls
	for longUrl, _ := range Urls {
		if !govalidator.IsURL(longUrl) {
			http.Response(writer, netHttp.StatusUnprocessableEntity, "The url [ "+longUrl+" ] is not a valid url")
			return
		}
	}

	// Generate and serve short urls
	ShortUrls, isDone := UrlsRepository.GenerateShortUrls(Urls, request)

	if isDone {
		http.Response(writer, netHttp.StatusOK, ShortUrls)
	} else {
		http.Response(writer, netHttp.StatusRequestTimeout, ShortUrls)
	}
}

// Redirect to long url method
func (urlsController *UrlsController) RedirectToLongUrl(writer netHttp.ResponseWriter, request *netHttp.Request) {
	query := mux.Vars(request)

	shortUrl := query["short-url"]
	broadcastId := query["broadcast-id"]
	subscriberEmail := query["subscriber-email"]

	// get long url from the requested short url
	UrlsRepository := repositories.UrlsRepository{}
	shortUrlExist, longUrl := UrlsRepository.GetLongUrl(shortUrl)

	if !shortUrlExist {
		http.Response(writer, netHttp.StatusNotFound, "The given short url is not found")
		return
	}

	// post a request to mailscout server for saving link clicked tracking information
	res, _ := httpClient.Post(configs.MAILSCOUT_WEBHOOK+broadcastId+"/"+subscriberEmail, map[string]string{
		"long_url": longUrl,
	})

	// if mailscout server return 201 status code then redirect to long url
	if res.StatusCode == 201 {
		netHttp.Redirect(writer, request, longUrl, netHttp.StatusMovedPermanently)
		return
	}

	// if mailscout server return without 201 status code
	http.Response(writer, netHttp.StatusNotFound, "The given short url is not found")
	return
}
