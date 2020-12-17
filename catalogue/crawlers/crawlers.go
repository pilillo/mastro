package crawlers

import (
	"fmt"

	"github.com/pilillo/mastro/catalogue/crawlers/hdfs"
	"github.com/pilillo/mastro/utils/conf"
	"github.com/pilillo/mastro/abstract"
)

var factories = map[string]func() abstract.Crawler{
	"hdfs": hdfs.NewCrawler,
}

func Start(cfg *conf.Config) (abstract.Crawler, error) {
	// start crawler defined in Config
	var crawler abstract.Crawler

	if crawlerFactory, ok := factories[cfg.CrawlerDefinition.Type]; ok {
		// call factory for selected crawler
		crawler = crawlerFactory()
		// todo: schedule crawler
		//crawler.Start()
		return crawler, nil
	}
	return nil, fmt.Errorf("Impossible to find specified Crawler %s", cfg.DataSourceDefinition.Type)
}
