package plugins

import (
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains functions and structs related to loading malscans plugins

*/

func parseResult(buf []byte, plugName string) (parsed json.RawMessage) {

	log.Debugf("parsing raw results for:%s", plugName)

	if err := json.Unmarshal(buf, &parsed); err != nil {
		log.Error(errors.Wrap(err, "error unmarshaling while parsing results"))
	}

	return parsed

}

func testInfected(buf []byte, plugName *string) bool {

	log.Debugf("testing infected for:%s", *plugName)

	var avresult map[string]map[string]interface{}

	if err := json.Unmarshal(buf, &avresult); err != nil {
		log.Error(errors.Wrap(err, "error unmarshaling while testing infected"))
	}

	infected, ok := avresult["analysis"]["infected"].(bool)

	if !ok {
		log.Error("error while trying to test infected for: ", *plugName)
	}

	return infected

}

func parseAnalysisResult(buf []byte, plugName *string) (result string) {

	log.Debugf("parsing result for:%s", *plugName)

	var avresult map[string]map[string]interface{}

	if err := json.Unmarshal(buf, &avresult); err != nil {
		log.Error(errors.Wrap(err, "error unmarshaling while testing infected"))
	}

	result, ok := avresult["analysis"]["result"].(string)

	if !ok {
		log.Error("error while trying to test infected for: ", *plugName)
	}

	return result

}
