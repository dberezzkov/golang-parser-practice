package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	// Открываем HTML страницу с вакансиями
	mainurl := "https://hh.ru/search/vacancy?from=suggest_post&ored_clusters=true&area=113&search_field=name&enable_snippets=false&text=Python+developer&customDomain=1"

	// test
	resp, err := http.Get(mainurl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	extracturls(doc)
}

func extracturls(doc *goquery.Document) {

	jsonstr := doc.Find("noindex>template").Text()

	extractlinks(jsonstr)
}

func extractlinks(jsonstr string) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonstr), &data)
	if err != nil {
		log.Fatal(err)
	}

	vacancySearchResult := data["vacancySearchResult"].(map[string]interface{})
	vacancies := vacancySearchResult["vacancies"].([]interface{})
	for i := range vacancies {
		vacancy := vacancies[i].(map[string]interface{})
		link := vacancy["links"].(map[string]interface{})["desktop"]
		fmt.Println(link)

	}
}

var a string = ` {
	"logos": {
		"logo": [
			{
				"@type": "vacancyPage",
				"@url": "/employer-logo/6656699.jpeg"
			},
			{
				"@type": "medium",
				"@url": "/employer-logo/6656699.jpeg"
			},
			{
				"@type": "ORIGINAL",
				"@url": "/employer-logo-original/1259093.jpg"
			}
			
		]}
	}`

func test(a string) {

	var data map[string]interface{}
	err := json.Unmarshal([]byte(a), &data)
	if err != nil {
		log.Fatal(err)
	}

	logos := data["logos"].(map[string]interface{})
	logo := logos["logo"].([]interface{})

	for i := range logo {
		t := logo[i].(map[string]interface{})
		fmt.Println(t["@type"])
	}
}

/*
func PrettyString(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return ""
	}
	return prettyJSON.String()

}
*/
// "desktop":"https://hh.ru/vacancy/" что что нужно искать
