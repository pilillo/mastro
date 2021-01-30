package impala

import (
	_ "github.com/pilillo/mastro/sources/impala"
)

type impalaCrawler struct {
	connector *impala.Connector
}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &impalaCrawler{}
}

func (crawler *impalaCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	crawler.connector = impala.NewConnector()
	if err := crawler.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		return nil, err
	}
	crawler.connector.InitConnection(&cfg.DataSourceDefinition)
}

func (crawler *impalaCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	var assets []abstract.Asset

	levels := strings.SplitAndTrim(root, "/")

	var tables []string
	
	// check if a specific database and table was defined
	if levels != nil && len(levels) > 0 {	
		// a table is defined	
		if len(levels) > 1 {
			// list only provided table by appending to list of identified ones
			tables = []string{ levels(1) }
		}else{
			// list all tables in provided db
			tables, err = crawler.connector.ListTables(levels(0))
			if err != nil {
				// error while accessing the sole DB we desired to access
				return nil, err
			}
		}
	}else{
		// list all databases, skip those we can't access, as may be a right issue
		dbs, err := crawler.Connector.ListDatabases()
		if err != nil {
			return nil, err
		}
		// list all tables in available dbs
		for i, dbInfo := range dbs {
			tables, err = crawler.connector.ListTables(dbInfo))
			if err != nil {
				// skipping DB
				log.println(fmt.Sprintf("Error while accessing DB %s! Skipping..", dbInfo.Name))
			}
		}
	}
	
	// describe all tables
	for i, tableName := range tables {
		// map[string]abstract.ColumnInfo
		tableSchema, err := crawler.connector.DescribeTable(dbName, tableName)
		if err != nil {
			log.println(fmt.Sprintf("Error while accessing %s.%s! Skipping..", dbName, tableName))
		}else{
			log.Printf("Found table %s.%s", dbName, tableName)
			// convert table schema to actual Asset definition
			
		}
	}

	return assets, nil
}