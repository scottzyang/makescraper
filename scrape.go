package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

// EggItem represents the structure of each egg item
type EggItem struct {
	ItemPrice     string `json:"item_price"`
	PricePerCount string `json:"price_per_count"`
	ItemName      string `json:"item_name"`
	ItemImage     string `json:"item_image"`
}

func main() {
	// Create a new collector
	c := colly.NewCollector()

	// Slice to store scraped data
	var eggItems []EggItem

	// Find and extract item information
	c.OnHTML("[id=\"\\30 \"] > section > div", func(e *colly.HTMLElement) {
		// Increment the image selector ID with each iteration
		imageSelector := "#is-0-productImage-" + strconv.Itoa(len(eggItems))

		// Extracting data using the provided selectors
		itemImage := e.ChildAttr(imageSelector, "src")
		itemPrice := e.ChildText("div:nth-child(3) > div > div > div > div:nth-child(2) > div.flex.flex-wrap.justify-start.items-center.lh-title.mb1 > span")
		pricePerCount := e.ChildText("div:nth-child(3) > div > div > div > div:nth-child(2) > div.flex.flex-wrap.justify-start.items-center.lh-title.mb1 > div.gray.mr1.f6.f5-l.flex.items-end.mt1")
		itemName := e.ChildText("div:nth-child(3) > div > div > div > div:nth-child(2) > span > span")

		// Create a new EggItem with the extracted data
		eggItem := EggItem{
			ItemPrice:     itemPrice,
			PricePerCount: pricePerCount,
			ItemName:      itemName,
			ItemImage:     itemImage,
		}

		// Append the EggItem to the slice
		eggItems = append(eggItems, eggItem)
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
		fmt.Printf("Item Image: %s\n", item.ItemImage)
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
