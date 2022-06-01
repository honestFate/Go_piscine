package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	typeJSON = iota
	typeXML
	typeErr
)

type DBReader interface {
	readJSON(*os.File) error
	readXML(*os.File) error
	convertXMLToMap()
	convertJSONToMap()
}

func readData(dbReader DBReader, fileName string) error {
	fileType := getDataType(fileName)
	fmt.Println(fileType)
	dataFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	if fileType == typeJSON {
		if err := dbReader.readJSON(dataFile); err != nil {
			return err
		}
		dbReader.convertJSONToMap()
		return nil
	} else if fileType == typeXML {
		if err := dbReader.readXML(dataFile); err != nil {
			return err
		}
		dbReader.convertXMLToMap()
		return nil
	}
	err = errors.New("Wrong file type")
	return err
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
	Cake     map[string]Cake
	DataJSON RecipesJSON `xml:"-" json:"-"`
	DataXML  RecipesXML  `xml:"-" json:"-"`
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

func (data *MapReciepes) readJSON(file *os.File) error {
	fmt.Println("read JSON")
	jsonParser := json.NewDecoder(file)
	err := jsonParser.Decode(&data.DataJSON)
	return err
}

func (data *MapReciepes) readXML(file *os.File) error {
	fmt.Println("read XML")
	xmlParser := xml.NewDecoder(file)
	err := xmlParser.Decode(&data.DataXML)
	return err
}

func (dataMap *MapReciepes) convertXMLToMap() {
	dataMap.Cake = make(map[string]Cake)
	for _, cake := range dataMap.DataXML.Cake {
		if entry, ok := dataMap.Cake[cake.Name]; ok {
			entry.IngredientMap = make(map[string]Ingredient)
			for _, ingredient := range cake.Ingredients.Item {
				if ingredientCopy, ok := entry.IngredientMap[ingredient.Itemname]; ok {
					ingredientCopy.IngredientCount = ingredient.Itemcount
					ingredientCopy.IngredientCount = ingredient.Itemunit
					entry.IngredientMap[ingredient.Itemname] = ingredientCopy
				}
			}
			entry.Time = cake.Stovetime
			dataMap.Cake[cake.Name] = entry
		}
	}
}

func (dataMap *MapReciepes) convertJSONToMap() {
	dataMap.Cake = make(map[string]Cake)
	for _, cake := range dataMap.DataJSON.Cake {
		if entry, ok := dataMap.Cake[cake.Name]; ok {
			entry.IngredientMap = make(map[string]Ingredient)
			for _, ingredient := range cake.Ingredients {
				if ingredientCopy, ok := entry.IngredientMap[ingredient.IngredientName]; ok {
					ingredientCopy.IngredientCount = ingredient.IngredientCount
					ingredientCopy.IngredientCount = ingredient.IngredientUnit
					entry.IngredientMap[ingredient.IngredientName] = ingredientCopy
				}
			}
			entry.Time = cake.Time
			dataMap.Cake[cake.Name] = entry
		}
	}
}

func getDataType(fileName string) int {
	if strings.HasSuffix(fileName, ".json") {
		return typeJSON
	} else if strings.HasSuffix(fileName, ".xml") {
		return typeXML
	}
	return typeErr
}

func main() {
	flagF := flag.Bool("f", false, "./readDB -f .json/.xml")
	flag.Parse()
	if !*flagF || flag.NArg() != 1 {
		flag.PrintDefaults()
		return
	}
	data := &MapReciepes{}
	if err := readData(data, os.Args[2]); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
	switch getDataType(os.Args[2]) {
	case typeJSON:
		newXML, err := xml.Marshal(data.DataJSON)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(newXML))
	case typeXML:
		newJSON, err := json.Marshal(data.DataXML)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(newJSON))
	case typeErr:
		log.Fatal("Wrong file type")
	}
	return
}
