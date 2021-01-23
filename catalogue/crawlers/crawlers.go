package crawlers

import (
	"fmt"

	"log"

	"github.com/jasonlvhit/gocron"
	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/catalogue/crawlers/hdfs"
	"github.com/pilillo/mastro/catalogue/crawlers/local"
	"github.com/pilillo/mastro/catalogue/crawlers/s3"
	"github.com/pilillo/mastro/utils/conf"

	"github.com/go-resty/resty/v2"
)

var factories = map[string]func() abstract.Crawler{
	"local": local.NewCrawler,
	"hdfs":  hdfs.NewCrawler,
	"s3":    s3.NewCrawler,
}

var client = resty.New()

// Start ... Starts the crawler defined in the provided config
func Start(cfg *conf.Config) (abstract.Crawler, error) {
	// start crawler defined in Config
	var crawler abstract.Crawler

	if crawlerFactory, ok := factories[cfg.DataSourceDefinition.Type]; ok {
		// call factory for selected crawler
		crawler = crawlerFactory()
		// init connection on the selected crawler
		crawler.InitConnection(cfg)
		// schedule crawler
		//every := gocron.Every(cfg.CrawlerDefinition.ScheduleValue)
		every := gocron.Every(cfg.DataSourceDefinition.CrawlerDefinition.ScheduleValue)
		switch cfg.DataSourceDefinition.CrawlerDefinition.ScheduleEvery {
		case conf.Seconds:
			every = every.Seconds()
		case conf.Minutes:
			every = every.Minutes()
		case conf.Hours:
			every = every.Hours()
		case conf.Days:
			every = every.Days()
		case conf.Weeks:
			every = every.Weeks()
		case conf.Monday:
			every = every.Monday()
		case conf.Tuesday:
			every = every.Tuesday()
		case conf.Wednesday:
			every = every.Wednesday()
		case conf.Thursday:
			every = every.Thursday()
		case conf.Friday:
			every = every.Friday()
		case conf.Saturday:
			every = every.Saturday()
		case conf.Sunday:
			every = every.Sunday()
		default:
			return nil, fmt.Errorf("crawler: schedule period %s not found", cfg.DataSourceDefinition.CrawlerDefinition.ScheduleEvery)
		}
		// spawn crawler for the selected schedule period
		every.Do(Reconcile, cfg)
		// start gocron - move outside if we decide to start multiple crawlers within the same agent
		<-gocron.Start()
		return crawler, nil
	}
	return nil, fmt.Errorf("Impossible to find specified Crawler %s", cfg.DataSourceDefinition.Type)
}

// Reconcile ... call to walkWithFilter to traverse the FS tree and post all found assets to the catalogue endpoint
func Reconcile(crawler abstract.Crawler, cfg *conf.Config) {
	assets, err := crawler.WalkWithFilter(cfg.DataSourceDefinition.CrawlerDefinition.RootFolder, cfg.DataSourceDefinition.CrawlerDefinition.FilterFilename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("Found %d assets to merge in catalogue", len(assets))
	// call a remote catalogue endpoint to add those assets that were just found
	// https://github.com/go-resty/resty/blob/master/example_test.go
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(assets).
		Post(cfg.DataSourceDefinition.CrawlerDefinition.CatalogueEndpoint)
		//Put(cfg.DataSourceDefinition.CrawlerDefinition.CatalogueEndpoint)
	// print response
	log.Println("Post to catalogue:", resp)

}
