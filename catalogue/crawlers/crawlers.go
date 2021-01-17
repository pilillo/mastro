package crawlers

import (
	"fmt"

	"log"

	"github.com/jasonlvhit/gocron"
	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/catalogue/crawlers/hdfs"
	"github.com/pilillo/mastro/catalogue/crawlers/s3"
	"github.com/pilillo/mastro/utils/conf"
)

var factories = map[string]func() abstract.Crawler{
	"hdfs": hdfs.NewCrawler,
	"s3":   s3.NewCrawler,
}

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

func Reconcile(crawler abstract.Crawler, cfg *conf.Config) {
	assets, err := crawler.WalkWithFilter(cfg.DataSourceDefinition.CrawlerDefinition.RootFolder, cfg.DataSourceDefinition.CrawlerDefinition.FilterFilename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("Found %d assets to merge in catalogue", len(assets))
	// todo: call a remote catalogue endpoint to add those assets that were just found
}
