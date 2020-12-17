package conf

// CrawlerDefinition ... Config for a Crawler service
type CrawlerDefinition struct {
	Type           string `json:"type"`
	RootFolder     string `json:"root_folder"`
	FilterFilename string `json:"filter_filename"`
	//UpdateInterval
}
