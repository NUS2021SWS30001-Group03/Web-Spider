package main

import (
	"fmt"
	g "github.com/serpapi/google-search-results-golang"
	"sws3001spider/controller"
	"sws3001spider/utils"
)

func main() {
	StopFlag := 1
	StartFlag := 1
	NextPageToken := ""
	for {
		var parameter map[string]string
		if StartFlag == 1 {
			parameter = map[string]string{
				"engine":   "google_scholar_profiles",
				"mauthors": "Huazhong University of Science and Technology",
				"api_key":  utils.URLKey,
			}
			StartFlag = 0
		} else {
			parameter = map[string]string{
				"engine":       "google_scholar_profiles",
				"mauthors":     "Huazhong University of Science and Technology",
				"api_key":      utils.URLKey,
				"after_author": NextPageToken,
			}
		}
		search := g.NewGoogleSearch(parameter, utils.URLKey)
		results, _ := search.GetJSON()
		organicResults := results["profiles"].([]interface{})
		for _, cur := range organicResults {
			curResult := cur.(map[string]interface{})
			name := curResult["name"].(string)
			affliations := curResult["affiliations"].(string)
			author_id := curResult["author_id"].(string)
			counts := curResult["cited_by"].(float64)
			controller.AddToQueue(name, affliations, author_id, counts)
			if counts < 8000.0 {
				StopFlag = 0
			}
		}
		if StopFlag == 0 {
			break
		}
		nextPageResult := results["pagination"].(interface{})
		NextPageToken = nextPageResult.(map[string]interface{})["next_page_token"].(string)
		fmt.Println(NextPageToken)
	}
	controller.Search()
}
