package impala

import(
	"fmt"
	"log"
	"os"
	_ "github.com/koblas/impalathing"
)

var requiredFields = map[string]string{
	"host":     "host",
	"port":     "port",
}

func NewImpalaConnector() *Connector {
	return &ImpalaConnector{}
}

type Connector struct {
	connection *_.Connection
}

func (c *Connector) ValidateDataSourceDefinition(def *conf.DataSourceDefinition) error {
	// check all required fields are available
	var missingFields []string
	for _, reqvalue := range requiredFields {
		if _, exist := def.Settings[reqvalue]; !exist {
			missingFields = append(missingFields, reqvalue)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("The following fields are missing from the data source configuration: %s", strings.Join(missingFields, ","))
	}

	log.Println("Successfully validated data source definition")
	return nil
}


func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error

	host := def.Settings[requiredFields["host"]]
	port := def.Settings[requiredFields["port"]]

	if def.KerberosDetails != nil {
		options := impalathing.WithGSSAPISaslTransport()
		c.connection, err = impalathing.Connect(host, port, options)
	}else{
		c.connection, err = impalathing.Connect(host, port)
	}

	if err != nil {
		panic(err)
	}	
}

func (c *Connector) CloseConnection() {
	c.connection.Close()
}

// Impala specific methods and structs


func (c *Connector) ListDatabases() ([]DBInfo, error) {
	var result = make([]DBInfo)
	
	query, err := c.connection.Query("show databases")
	if err != nil {
		return nil, err
	}
	for query.Next() {
		db := DBInfo{}
		query.Scan(&db.Name, &db.Comment)

		result = append(result, db)
	}
	return result, nil
}

for (c *Connector) ListTables(dbName string) ([]string, error) {
	var result = make([]string)

	query, err := c.connection.Query(fmt.Sprintf("show tables in %s", dbName))
	for query.Next() {
		var tableName string
		query.Scan(&tableName)
		result = append(result, tableName)
	}

	return result, nil
}

func (c *Connector) DescribeTable(dbName string, tableName string) (map[string]abstract.ColumnInfo, error){
	var result = make(map[string]ColumnInfo)

	query, err := c.connection.Query(fmt.Sprintf("describe %s.%s", dbName, tableName))
	if err != nil {
		return nil, err
	}
	
	for query.Next() {
		var cName string
		cInfo := ColumnInfo{}

		query.Scan(&cName, &(cInfo.Type), &(cInfo.Comment))
		result[cName] = cInfo
	}

	return result, nil
}