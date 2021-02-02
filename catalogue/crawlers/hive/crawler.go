package hive

import (
	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/sources/hive"
	"github.com/pilillo/mastro/utils/conf"

	"fmt"
	"log"

	"github.com/pilillo/mastro/utils/strings"
)

type hiveCrawler struct {
	connector *hive.Connector
}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &hiveCrawler{}
}

func (crawler *hiveCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	crawler.connector = hive.NewHiveConnector()
	if err := crawler.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		return nil, err
	}
	crawler.connector.InitConnection(&cfg.DataSourceDefinition)
	return crawler, nil
}

func (crawler *hiveCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	var assets []abstract.Asset

	levels := strings.SplitAndTrim(root, "/")

	// create empty map of kind, db -> []tablenames
	dbTables := map[string][]string{}

	// check if a specific database and table was defined
	// N.B. golang split returns a slice with one element, the empty string so len is 1 and we gotta check it
	// https://stackoverflow.com/questions/28330908/how-to-string-split-an-empty-string-in-go
	if levels != nil && len(levels) > 0 && levels[0] != "" {
		log.Printf("Provided specific db levels to locate: '%s'", root)

		// a table is defined
		if len(levels) > 1 {
			// list only provided table by appending to list of identified ones
			dbTables[levels[0]] = []string{levels[1]}
			//tables = []string{ levels[1] }
		} else {
			// list all tables in provided db
			tables, err := crawler.connector.ListTables(levels[0])
			if err != nil {
				// error while accessing the sole DB we desired to access
				return nil, err
			}
			dbTables[levels[0]] = tables
			log.Printf("Found %d tables in requested database %s: %v", len(tables), levels[0], tables)
		}
	} else {
		// list all databases, skip those we can't access, as may be a right issue
		dbs, err := crawler.connector.ListDatabases()

		if err != nil {
			return nil, err
		}

		// list all tables in available dbs
		for _, dbInfo := range dbs {
			tables, err := crawler.connector.ListTables(dbInfo.Name)
			if err != nil {
				// skipping DB
				log.Println(fmt.Sprintf("Error while accessing DB %s! Skipping..", dbInfo.Name))
			} else {
				// add all found tables to map for given db name
				dbTables[dbInfo.Name] = tables
				log.Printf("Found %d tables in database %s: %v", len(tables), dbInfo.Name, tables)
			}
		}
	}

	// visit each found db
	for dbName, tableNames := range dbTables {
		// describe each table in the db
		for _, tableName := range tableNames {
			// map[string]abstract.ColumnInfo
			_, err := crawler.connector.DescribeTable(dbName, tableName)
			if err != nil {
				log.Print(fmt.Sprintf("Error while accessing %s.%s! Skipping..", dbName, tableName))
			} else {
				log.Printf("Retrieved schema for table %s.%s", dbName, tableName)
				// convert table schema to actual Asset definition
				//tableSchema
			}
		}
	}

	return assets, nil
}
