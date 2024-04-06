package structs

import "encoding/json"

//FullFileReport - Used to fill in information for a file to be
//sent off for alerting and elasticsearch indexing
type FullFileReport struct {
	File fileinfo `structs:"file" json:"file"`
}

type fileinfo struct {
	Name    string   `structs:"filename" json:"filename"`
	Sha1    string   `structs:"sha1" json:"sha1"`
	Md5     string   `structs:"md5" json:"md5"`
	Date    string   `structs:"date" json:"date"`
	Tags    []string `structs:"tags" json:"tags"`
	Malware malware  `structs:"malware" json:"malware"`
}

type analyzers struct {
	Names       []string    `structs:"name" json:"name"`
	RawAnalysis rawAnalysis `structs:"raw-analysis" json:"raw-analysis"`
}

type malware struct {
	Infected  bool      `structs:"infected" json:"infected"`
	Results   []string  `structs:"variants" json:"variants"`
	Analyzers analyzers `structs:"analyzers" json:"analyzers"`
}

type rawAnalysis struct {
	AntiVirus map[string]json.RawMessage `structs:"detection" json:"detection"`
	Enricher  map[string]json.RawMessage `structs:"enricher" json:"enricher"`
}
