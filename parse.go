package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	links := make([]string, 10)
	// Открываем HTML страницу с вакансиями
	mainurl := "https://hh.ru/search/vacancy?text=Python&salary=&ored_clusters=true&area=113&hhtmFrom=vacancy_search_list&hhtmFromLabel=vacancy_search_line"
	findpage(mainurl, links)
	a := "https://hh.ru/vacancy/95473687?query=Python+developer&hhtmFrom=vacancy_search_list"
	all(a)
	fmt.Println(links)
}

func findpage(url string, links []string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	c := 0
	doc.Find("span.serp-item__title-link-wrapper>a").Each(func(i int, s *goquery.Selection) {
		hr, _ := s.Attr("href")
		// links = append(links, hr)
		links = append(links, hr)
		c++
		fmt.Println(hr)
	})
	fmt.Println(c)
}

func all(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Создаем новый документ goquery из файла HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	location := extractJobLocation(doc)
	title, salary := extractsalary(doc)

	// Вывод переменных
	fmt.Println(location)
	fmt.Println(title)
	fmt.Println(salary)
}
func extractsalary(doc *goquery.Document) (string, string) {

	// тестирую вывод заралпаты
	name := doc.Find("h1.bloko-header-section-1").Text()
	salary := doc.Find("div.vacancy-title>div>span").Text()
	return name, salary
}

// extractJobLocation извлекает местоположение вакансии из документа goquery
func extractJobLocation(doc *goquery.Document) string {
	var location string
	// Ищем все скрипты в документе
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		// Проверяем, является ли скрипт JSON-LD
		if val, _ := s.Attr("type"); val == "application/ld+json" {
			jsonData := strings.TrimSpace(s.Text())
			// Проверяем, содержит ли JSON данные о местоположении вакансии
			if strings.Contains(jsonData, "jobLocation") {
				location = extractLocationFromJSON(jsonData)
			}
		}
	})

	return location
}

func extractLocationFromJSON(jsonStr string) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	locationData, ok := data["jobLocation"].(map[string]interface{})
	if !ok {
		return ""
	}

	address, ok := locationData["address"].(map[string]interface{})
	if !ok {
		return ""
	}

	addressLocality, _ := address["addressLocality"].(string)
	addressRegion, _ := address["addressRegion"].(string)
	streetAddress, _ := address["streetAddress"].(string)

	location := fmt.Sprintf("%s, %s, %s", addressLocality, addressRegion, streetAddress)
	return location
}
