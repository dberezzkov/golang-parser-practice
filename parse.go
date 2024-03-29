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
	// Открываем файл HTML с данными о вакансии
	url := "https://hh.ru/vacancy/95473687?query=Python+developer&hhtmFrom=vacancy_search_list"
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
	if location == "" {
		fmt.Println("Местоположение вакансии не найдено.")
	} else {
		fmt.Println("Местоположение вакансии:", location)
	}
	// тестирую вывод заралпаты
	name := doc.Find("h1.bloko-header-section-1").Text()
	fmt.Println(name)
	salary := doc.Find("div.vacancy-title>div>span").Text()
	fmt.Println(salary)
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
