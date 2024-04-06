package elastic

import (
	"context"

	utils "malscan/core/utils/hash"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains all elasticsearch index related functions

*/

//createIndex - Responsible for creating the malscan index if it does not exist
func createIndex() {

	ctx := context.Background()

	log.Debug("checking if a malscan index exists")
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("malscan").Do(ctx)
	if err != nil {
		log.Error(errors.Wrap(err, "error while checking if malscan index exists"))
	}

	if !exists {
		log.Debug("no malscan index exists ... attempting to create malscan index")
		// Create a new index.
		createIndex, err := client.CreateIndex("malscan").Do(ctx)
		if err != nil {
			log.Error(errors.Wrap(err, "error while creating index"))
		}
		if !createIndex.Acknowledged {
			log.Warn("index creation not acknowledged")
		}
	} else {
		log.Debug("malscan index already exits")
	}
}

//Index - Responsible for indexing any structure passed into the function into elasticsearch
func Index(toIndex interface{}, filename *string) {

	fileID, err := utils.GenerateFileSha1(filename)
	if err != nil {
		log.Error(errors.Wrap(err, "error while generating file ID"))
	}

	put, err := client.Index().
		Index("malscan").
		Type("avscan").
		Id(fileID).
		BodyJson(toIndex).
		Do(context.Background())
	if err != nil {
		log.Warn(errors.Wrap(err, "error while attempting to index result"))
	}

	log.Debug("indexed scan of file: ", put.Id, " to index: ", put.Index, " of type: ", put.Type)
}
