package conf

// CrawlerDefinition ... Config for a Crawler service
type CrawlerDefinition struct {
	RootFolder        string `yaml:"root-folder"`
	FilterFilename    string `yaml:"filter-filename"`
	ScheduleEvery     Period `yaml:"schedule-period"`
	ScheduleValue     uint64 `yaml:"schedule-value"`
	CatalogueEndpoint string `yaml:"catalogue-endpoint"`
}

// Period ... time period to schedule the crawler for
type Period string

const (
	// Seconds ... time interval
	Seconds Period = "seconds"
	// Minutes ... schedule time interval
	Minutes = "minutes"
	// Hours ... schedule time interval
	Hours = "hours"
	// Days ... schedule time interval
	Days = "days"
	// Weeks ... schedule time interval
	Weeks = "weeks"

	// Monday ...
	Monday = "monday"
	// Tuesday ...
	Tuesday = "tuesday"
	// Wednesday ...
	Wednesday = "wednesday"
	// Thursday ...
	Thursday = "thursday"
	// Friday ...
	Friday = "friday"
	// Saturday ...
	Saturday = "saturday"
	// Sunday = "sunday"
	Sunday = "sunday"
)
