package types

// urls struct
type Urls struct {
	LongUrls []Url `json:"long_urls"`
}

// url struct
type Url struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}
