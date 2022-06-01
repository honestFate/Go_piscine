package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type DBReader interface {
	readFile(*os.File) error
}

func readData(dbReader DBReader, fileName string) error {
	dataFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	return dbReader.readFile(dataFile)
}

type Ingredient struct {
	IngredientCount string
	IngredientUnit  string
}

type Cake struct {
	Time          string
	IngredientMap map[string]Ingredient
}

type MapReciepes struct {
	Cake map[string]Cake
}

type RecipesJSON struct {
	Cake []struct {
		Name        string `json:"name"`
		Time        string `json:"time"`
		Ingredients []struct {
			IngredientCount string `json:"ingredient_count"`
			IngredientName  string `json:"ingredient_name"`
			IngredientUnit  string `json:"ingredient_unit,omitempty"`
		} `json:"ingredients"`
	} `json:"cake"`
}

type RecipesXML struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Text    string   `xml:",chardata" json:"-"`
	Cake    []struct {
		Text        string `xml:",chardata" json:"-"`
		Name        string `xml:"name"`
		Stovetime   string `xml:"stovetime"`
		Ingredients struct {
			Text string `xml:",chardata" json:"-"`
			Item []struct {
				Text      string `xml:",chardata" json:"-"`
				Itemname  string `xml:"itemname"`
				Itemcount string `xml:"itemcount"`
				Itemunit  string `xml:"itemunit,omitempty"`
			} `xml:"item"`
		} `xml:"ingredients"`
	} `xml:"cake"`
}

func (data *RecipesJSON) readFile(file *os.File) error {
	jsonParser := json.NewDecoder(file)
	err := jsonParser.Decode(data)
	return err
}

func (data *RecipesXML) readFile(file *os.File) error {
	xmlParser := xml.NewDecoder(file)
	err := xmlParser.Decode(data)
	return err
}

func main() {
	flagF := flag.Bool("f", false, "./readDB -f .json/.xml")
	flag.Parse()
	if !*flagF || flag.NArg() != 1 {
		flag.PrintDefaults()
		return
	}
	if strings.HasSuffix(os.Args[2], ".json") {
		jsonData := &RecipesJSON{}
		if err := readData(jsonData, os.Args[2]); err != nil {
			log.Fatal(err)
		}
		newXML, err := xml.Marshal(jsonData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(newXML))
	} else if strings.HasSuffix(os.Args[2], ".xml") {
		xmlData := &RecipesXML{}
		if err := readData(xmlData, os.Args[2]); err != nil {
			log.Fatal(err)
		}
		newJSON, err := json.Marshal(xmlData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(newJSON))
	} else {
		fmt.Println("Wrong file type")
		flag.PrintDefaults()
		return
	}
	return
}
