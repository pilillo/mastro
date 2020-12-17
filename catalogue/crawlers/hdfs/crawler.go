package hdfs

import (
	"fmt"
	"os"

	gohdfs "github.com/colinmarc/hdfs/v2"
	"github.com/colinmarc/hdfs/v2/hadoopconf"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/utils/conf"

	krbClient "github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
)

func GetKerberosClient(details *conf.KerberosDetails) *krbClient.Client {
	// https://github.com/jcmturner/gokrb5/blob/master/v8/USAGE.md
	// Replace with a valid credentialed client.
	cfg, err := config.Load(details.KrbConfigPath)
	if err != nil {
		panic(err)
	}

	var krb5Client *krbClient.Client

	if len(details.KeytabPath) > 0 {
		kt, err := keytab.Load(details.KeytabPath)
		if err != nil {
			panic(err)
		}
		krb5Client = krbClient.NewWithKeytab(
			details.Username,
			details.Realm, kt, cfg,
			krbClient.DisablePAFXFAST(details.DisablePAFXFAST), krbClient.AssumePreAuthentication(false),
		)
	} else {
		krb5Client = krbClient.NewWithPassword(
			details.Username,
			details.Realm, details.Password, cfg,
			krbClient.DisablePAFXFAST(true),
		)
	}

	return krb5Client
}

type hadoopCrawler struct {
	client *gohdfs.Client
}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &hadoopCrawler{}
}

func (crawler *hadoopCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	krbDetails := cfg.DataSourceDefinition.KerberosDetails

	// "HADOOP_CONF_DIR" should be set for this to work
	_, present := os.LookupEnv("HADOOP_CONF_DIR")
	if !present {
		panic("HADOOP_CONF_DIR not set!")
	}

	/*
		LoadFromEnvironment tries to locate the Hadoop configuration files based on the environment,
		and returns a HadoopConf object representing the parsed configuration.
		If the HADOOP_CONF_DIR environment variable is specified, it uses that, or if HADOOP_HOME is specified, it uses $HADOOP_HOME/conf.
	*/
	hadoopConf, err := hadoopconf.LoadFromEnvironment()
	if err != nil {
		panic(err)
	}

	// https://godoc.org/github.com/colinmarc/hdfs#ClientOptionsFromConf
	clientOptions := gohdfs.ClientOptionsFromConf(hadoopConf)

	if clientOptions.KerberosClient != nil {
		clientOptions.KerberosClient = GetKerberosClient(krbDetails)
	}

	client, err := gohdfs.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	crawler.client = client
	return crawler, nil
}

func (crawler *hadoopCrawler) Open(path string) error {
	file, err := crawler.client.Open(path)

	if err != nil {
		return err
	}

	buf := make([]byte, 59)
	file.ReadAt(buf, 48847)

	fmt.Println(string(buf))
	return nil
}

func (crawler *hadoopCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	return nil, nil
}
