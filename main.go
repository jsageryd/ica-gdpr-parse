package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var tzSthlm *time.Location

func init() {
	var err error

	tzSthlm, err = time.LoadLocation("Europe/Stockholm")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ica-gdpr-parse <dir>, where <dir> is a directory containing exported xml files")
		return
	}

	data, err := readAll(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	t, err := totals(
		data,
		time.Date(2023, 1, 1, 0, 0, 0, 0, tzSthlm),
		time.Date(2024, 1, 1, 0, 0, 0, 0, tzSthlm),
	)
	if err != nil {
		log.Fatalf("calculate totals: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(t); err != nil {
		log.Fatal(err)
	}
}

func totals(data Data, from, to time.Time) (Totals, error) {
	include := make(map[string]bool)

	for _, txn := range data.ButikKvitto.Transactions {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", txn.Timestamp, tzSthlm)
		if err != nil {
			return Totals{}, fmt.Errorf("parse timestamp %q: %v", txn.Timestamp, err)
		}

		if timestamp.After(from.Add(-1*time.Nanosecond)) && timestamp.Before(to) {
			include[txn.ID] = true
		}
	}

	itemsMap := make(map[string]ItemTotal)

	for _, row := range data.ButikKvittorader.Rows {
		if !include[row.TransactionID] {
			continue
		}

		item := itemsMap[row.ItemDesc]

		item.ItemDesc = row.ItemDesc
		item.TotalQuantity += row.Quantity
		item.TotalPrice += row.Price
		item.TotalDiscountValue += row.DiscountValue
		item.TotalDiscountedPrice += row.Price + row.DiscountValue

		itemsMap[row.ItemDesc] = item
	}

	var items []ItemTotal
	for _, item := range itemsMap {
		items = append(items, item)
	}

	slices.SortFunc(items, func(a, b ItemTotal) int {
		return strings.Compare(a.ItemDesc, b.ItemDesc)
	})

	return Totals{
		From:  from,
		To:    to,
		Items: items,
	}, nil
}

func readAll(dir string) (Data, error) {
	var data Data

	for _, each := range []struct {
		filename string
		v        any
	}{
		{"Butik kvitto.xml", &data.ButikKvitto},
		{"Butik kvittorader.xml", &data.ButikKvittorader},
	} {
		if err := readFile(filepath.Join(dir, each.filename), each.v); err != nil {
			return Data{}, err
		}
	}

	return data, nil
}

func readFile(name string, v any) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("open %q: %v", filepath.Base(name), err)
	}
	defer f.Close()

	if err := xml.NewDecoder(f).Decode(&v); err != nil {
		return fmt.Errorf("decode %q: %v", filepath.Base(name), err)
	}

	return nil
}

type Totals struct {
	From  time.Time   `json:"from"`
	To    time.Time   `json:"to"`
	Items []ItemTotal `json:"items"`
}

type ItemTotal struct {
	ItemDesc             string  `json:"item"`
	TotalQuantity        float64 `json:"total_quantity"`
	TotalPrice           float64 `json:"total_price"`
	TotalDiscountValue   float64 `json:"total_discount_value"`
	TotalDiscountedPrice float64 `json:"total_discounted_price"`
}

type Data struct {
	ButikKvitto      ButikKvitto
	ButikKvittorader ButikKvittorader
}

type ButikKvitto struct {
	XMLName      xml.Name                 `xml:"businessObjectToFileArea"`
	Transactions []ButikKvittoTransaction `xml:"resObject>TransactionHeader>transactions"`
}

type ButikKvittoTransaction struct {
	ID            string  `xml:"transactionId"`
	Timestamp     string  `xml:"transactionTimestamp"`
	Value         float64 `xml:"transactionValue"`
	MarketingName string  `xml:"marketingName"`
}

type ButikKvittorader struct {
	XMLName xml.Name              `xml:"businessObjectToFileArea"`
	Rows    []ButikKvittoraderRow `xml:"resObject>LineItems>transactions"`
}

type ButikKvittoraderRow struct {
	Quantity      float64 `xml:"quantity"`
	Price         float64 `xml:"price"`
	ItemDesc      string  `xml:"itemDesc"`
	DiscountValue float64 `xml:"discountValue"`
	TransactionID string  `xml:"transactionId"`
}
