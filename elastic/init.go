package elastic

import (
	"context"
	"malscan/config"
	"time"

	elastic "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Responsible for initializing elastic search

*/

var client *elastic.Client
var err error

//init 1 - intialize and test connection to elasticsearch
//log errors if the connection cannot be made
//this function only is only called if the elasticsearch package is called

func init() {

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	url1 := config.Values.Elasticsearch.URL1

	go func() {

		for {

			// Obtain a client and connect to the default Elasticsearch installation

			//TLS SET TO FALSE IS ONLY FOR NON PRODUCTION TESTING
			if config.Values.Elasticsearch.TLS != true {
				client, err = elastic.NewClient()
			} else {
				client, err = elastic.NewClient(
					elastic.SetURL(url1),
					elastic.SetSniff(false),
					elastic.SetBasicAuth(getUserName(), getUserPassword()),
					elastic.SetHttpClient(getHTTPSClient()),
				)
			}
			if err != nil {
				// Handle error
				log.Error(errors.Wrap(err, "error while creating elasticsearch client"))
				time.Sleep(time.Second * 20)
				continue
			} else {
				break
			}

		}
		// Ping the Elasticsearch server to get e.g. the version number
		_, _, err := client.Ping(url1).Do(ctx)
		if err != nil {
			// Handle error
			log.Error(errors.Wrap(err, "could not ping elasticsearch node: "), url1)
		}

		// Getting the ES version number is quite common, so there's a shortcut
		esversion, err := client.ElasticsearchVersion(url1)
		if err != nil {
			// Handle error
			log.Error(errors.Wrap(err, "could not get elasticsearch version ... node: "), url1)
		}
		log.Debug("elasticsearch version: ", esversion)

		go createIndex()

	}()

}
