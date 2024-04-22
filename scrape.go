package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// EggItem represents the structure of each egg item
type EggItem struct {
	ItemPrice     string `json:"item_price"`
	PricePerCount string `json:"price_per_count"`
	ItemName      string `json:"item_name"`
	Url           string `json:"url"`
}

func main() {
	// Create a new collector
	c := colly.NewCollector()

	// Slice to store scraped data
	var eggItems []EggItem

	// Find and extract item information
	c.OnHTML("[id=\"\\30 \"] > section > div > div", func(e *colly.HTMLElement) {
		// Check if the div has the class "w-100" or "w_3jM4"
		adclass := e.Attr("class")
		if strings.Contains(adclass, "w-100") || strings.Contains(adclass, "w_3jM4") {
			// Skip processing this div
			fmt.Println("Skipping div with class 'w-100' or 'w_3jM4'")
			return
		}

		if e.ChildAttr("div", "style") == "contain-intrinsic-size:198px 340px" {
			// Extracting data using the provided selectors)
			itemPrice := e.ChildText("div > div > div > div > div > div.flex.flex-wrap.justify-start.items-center.lh-title.mb1 > span")
			pricePerCount := e.ChildText("div > div > div > div > div > div.flex.flex-wrap.justify-start.items-center.lh-title.mb1 > div.gray.mr1.f6.f5-l.flex.items-end.mt1")
			itemName := e.ChildText("div > div > div > div > div > span > span")
			url := e.ChildAttr(`div > div > a`, "href")

			// Create a new EggItem with the extracted data
			eggItem := EggItem{
				ItemPrice:     itemPrice,
				PricePerCount: pricePerCount,
				ItemName:      itemName,
				Url:           url,
			}

			// Append the EggItem to the slice
			eggItems = append(eggItems, eggItem)

		}
	})

	// Before making a request, log "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on Walmart's website
	c.Visit("https://www.walmart.com/search?q=eggs")

	// Print the scraped data to console
	fmt.Println("Scraped Data:")
	for _, item := range eggItems {
		fmt.Printf("Item Name: %s\n", item.ItemName)
		fmt.Printf("Item Price: %s\n", item.ItemPrice)
		fmt.Printf("Price Per Count: %s\n", item.PricePerCount)
		fmt.Printf("Url: %s\n", item.Url)
		fmt.Println("-----------------------------------")
	}

	// Save the scraped data to a JSON file
	saveJSON(eggItems)
}

// saveJSON saves the scraped data to a JSON file
func saveJSON(data interface{}) {
	file, err := os.Create("output.json")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		log.Fatal("Cannot encode to JSON", err)
	}

	fmt.Println("Data saved to output.json")
}
