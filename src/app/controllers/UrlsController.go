package controllers

import (
	"app/http"
	"app/repositories"
	"configs"
	netHttp "net/http"

	"github.com/asaskevich/govalidator"
	httpclient "github.com/ddliu/go-httpclient"
	"github.com/gorilla/mux"
)

// Urls controller struct
type UrlsController struct {
	Controller
}

// Index method
func (urlsController *UrlsController) Index(writer netHttp.ResponseWriter, request *netHttp.Request) {

	// get long urls from request.body
	UrlsRepository := repositories.UrlsRepository{}
	Urls := UrlsRepository.GetLongUrls(request.Body)

	// validate urls
	for _, urls := range Urls.LongUrls {
		if !govalidator.IsURL(urls.LongUrl) {
			http.Response(writer, netHttp.StatusUnprocessableEntity, "The url [ "+urls.LongUrl+" ] is not a valid url")
			return
		}
	}

	// Generate and serve short urls
	ShortUrls := UrlsRepository.GenerateShortUrls(Urls)

	http.Response(writer, netHttp.StatusOK, ShortUrls)
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
	res, _ := httpclient.Post(configs.MAILSCOUT_WEBHOOK_LINK+broadcastId+"/"+subscriberEmail, map[string]string{
		"long_url": longUrl,
	})

	if res.StatusCode == netHttp.StatusOK {
		netHttp.Redirect(writer, request, longUrl, netHttp.StatusMovedPermanently)
		return
	}

	http.Response(writer, netHttp.StatusNotFound, "The given short url is not found")
	return
}
