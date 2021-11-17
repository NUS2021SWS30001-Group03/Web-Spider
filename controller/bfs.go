package controller

import (
	"fmt"
	g "github.com/serpapi/google-search-results-golang"
	"sws3001spider/model"
	"sws3001spider/utils"
)

var ResearchQueue []ResearcherQueData

type ResearcherQueData struct {
	Name        string
	Affliations string
	AuthorID    string
	Counts      int
	Depth       int
}

func AddToQueue(name string, affliations string, authorID string, counts float64) {
	data := ResearcherQueData{
		Name:        name,
		Affliations: affliations,
		AuthorID:    authorID,
		Counts:      int(counts),
		Depth:       0,
	}
	ResearchQueue = append(ResearchQueue, data)
}

func TestCoAuthor(mp map[string]interface{}) bool {
	if _, ok := mp["co_authors"]; ok {
		return ok
	}
	return false
}

func TestAffiliations(mp map[string]interface{}) bool {
	if _, ok := mp["affiliations"]; ok {
		return ok
	}
	return false
}

func TestPage(mp map[string]interface{}) bool {
	if _, ok := mp["pagination"]; ok {
		return ok
	}
	return false
}

func TestMaxPage(mp map[string]interface{}) int {
	cur := 500
	for {
		if cur == 0 {
			break
		}
		s := fmt.Sprintf("%d", cur)
		if _, ok := mp[s]; ok {
			return cur
		}
		cur = cur - 1
	}
	return -1
}

func Search() {
	//fmt.Println(ResearchQueue)
	for {
		if len(ResearchQueue) == 0 {
			break
		}
		cur := ResearchQueue[0]

		parameter := map[string]string{
			"engine":    "google_scholar_author",
			"author_id": cur.AuthorID,
			"api_key":   utils.URLKey,
		}
		search := g.NewGoogleSearch(parameter, utils.URLKey)
		results, _ := search.GetJSON()
		fmt.Println(cur.Name)
		//找跟他相连的边和边权，好痛苦
		if TestCoAuthor(results) {

			var current_co_author_list []model.CoAuthors
			co_authors := results["co_authors"].([]interface{})
			for _, cur_person := range co_authors {
				coauthor_name := cur_person.(map[string]interface{})["name"].(string)
				coauthor_id := cur_person.(map[string]interface{})["author_id"].(string)
				affiliations:=""
				if TestAffiliations(cur_person.(map[string]interface{})) {
					affiliations=cur_person.(map[string]interface{})["affiliations"].(string)
				}
				CheckExist:= model.CheckExist(coauthor_id)
				if CheckExist {
					continue
				}
				text := cur.Name + " " + coauthor_name
				//fmt.Println(text)
				parameter := map[string]string{
					"engine":  "google_scholar",
					"q":       text,
					"api_key": utils.URLKey,
				}
				search := g.NewGoogleSearch(parameter, utils.URLKey)
				results, _ := search.GetJSON()
				if TestPage(results)==false {
					current_co_author_list = append(current_co_author_list, model.CoAuthors{
						CoAuthorName: coauthor_name,
						Weight:       1,
					})
					queuedata := ResearcherQueData{
						Name:        coauthor_name,
						Affliations: affiliations,
						AuthorID:    coauthor_id,
						Counts:      0,
						Depth:       cur.Depth + 1,
					}
					if cur.Depth+1 <=2 {
						ResearchQueue = append(ResearchQueue, queuedata)
					}
					continue
				}
				organic_results := results["pagination"].(interface{})
				max_page := TestMaxPage(organic_results.(map[string]interface{})["other_pages"].(map[string]interface{}))
				sum := max_page
				for {
					sum += max_page
					parameter := map[string]string{
						"engine":  "google_scholar",
						"q":       text,
						"api_key": utils.URLKey,
						"start":   fmt.Sprintf("%d", sum*10),
					}
					search := g.NewGoogleSearch(parameter, utils.URLKey)
					results, _ := search.GetJSON()
					if TestPage(results)==false {
						break
					}
					organic_results := results["pagination"].(interface{})
					max_page = TestMaxPage(organic_results.(map[string]interface{})["other_pages"].(map[string]interface{}))
				}
				current_co_author_list = append(current_co_author_list, model.CoAuthors{
					CoAuthorName: coauthor_name,
					Weight:       sum * 10,
				})
				queuedata := ResearcherQueData{
					Name:        coauthor_name,
					Affliations: affiliations,
					AuthorID:    coauthor_id,
					Counts:      0,
					Depth:       cur.Depth + 1,
				}
				if cur.Depth+1 <=2 {
					ResearchQueue = append(ResearchQueue, queuedata)
				}
			}
			//写完了，希望人没事

			DBData := model.Researcher{
				Name:         cur.Name,
				Affiliation:  cur.Affliations,
				AuthorID:     cur.AuthorID,
				CoAuthorList: current_co_author_list,
			}
			err := model.InsertNode(DBData)
			if err != nil {
				panic(err)
			}
			fmt.Println(DBData)
		}
		ResearchQueue = ResearchQueue[1:]
	}
}
