package datasource

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func GetElasticSearch() (*elasticsearch.Client, error) {
	// cfg := elasticsearch.Config{
	// 	Addresses: []string{
	// 		"http://localhost:9200",
	// 	},
	// }
	// es, err := elasticsearch.NewClient(cfg)
	es, err := elasticsearch.NewDefaultClient()
	fmt.Println("ES initialized...")
	return es, err
}
