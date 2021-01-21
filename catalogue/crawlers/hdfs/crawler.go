package hdfs

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/sources/hdfs"
	"github.com/pilillo/mastro/utils/conf"
	"github.com/pilillo/mastro/utils/strings"
)

type hadoopCrawler struct {
	connector *hdfs.Connector
}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &hadoopCrawler{}
}

func (crawler *hadoopCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	crawler.connector = hdfs.NewHDFSConnector()
	if err := crawler.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		log.Panicln(err)
	}
	// inits connection
	crawler.connector.InitConnection(&cfg.DataSourceDefinition)
	return crawler, nil
}

func (crawler *hadoopCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	var assets []abstract.Asset

	var walkFn filepath.WalkFunc = func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// check if it is a regular file (not dir) and the name is like the filter
		if info.Mode().IsRegular() && strings.MatchPattern(info.Name(), filter) {

			fileReader, err := crawler.connector.GetClient().Open(path)
			if err != nil {
				return err
			}
			defer fileReader.Close()

			buf := new(bytes.Buffer)
			if _, err := io.CopyN(buf, fileReader, info.Size()); err != nil {
				return err
			}

			a, err := abstract.ParseAsset(buf.Bytes())
			if err != nil {
				return err
			}
			assets = append(assets, *a)
		}
		return nil
	}

	crawler.connector.GetClient().Walk(root, walkFn)

	return assets, nil
}
