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
	jsonParser := json.NewDecoder(file)
	err := jsonParser.Decode(&data.DataJSON)
	return err
}

func (data *MapReciepes) readXML(file *os.File) error {
	xmlParser := xml.NewDecoder(file)
	err := xmlParser.Decode(&data.DataXML)
	return err
}

func (dataMap *MapReciepes) convertXMLToMap() {
	dataMap.Cake = make(map[string]Cake)
	for _, cake := range dataMap.DataXML.Cake {
		dataMap.Cake[cake.Name] = Cake{}
		if entry, ok := dataMap.Cake[cake.Name]; ok {
			entry.IngredientMap = make(map[string]Ingredient)
			for _, ingredient := range cake.Ingredients.Item {
				entry.IngredientMap[ingredient.Itemname] = Ingredient{}
				if ingredientCopy, ok := entry.IngredientMap[ingredient.Itemname]; ok {
					ingredientCopy.IngredientCount = ingredient.Itemcount
					ingredientCopy.IngredientUnit = ingredient.Itemunit
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
		dataMap.Cake[cake.Name] = Cake{}
		if entry, ok := dataMap.Cake[cake.Name]; ok {
			entry.IngredientMap = make(map[string]Ingredient)
			for _, ingredient := range cake.Ingredients {
				entry.IngredientMap[ingredient.IngredientName] = Ingredient{}
				if ingredientCopy, ok := entry.IngredientMap[ingredient.IngredientName]; ok {
					ingredientCopy.IngredientCount = ingredient.IngredientCount
					ingredientCopy.IngredientUnit = ingredient.IngredientUnit
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

func formatChange() {
	data := &MapReciepes{}
	if err := readData(data, os.Args[2]); err != nil {
		log.Fatal(err)
	}
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
}

func bdCompare(flagOld *string, flagNew *string) {
	oldData := MapReciepes{}
	newData := MapReciepes{}
	if err := readData(&oldData, *flagOld); err != nil {
		log.Fatal(err)
	}
	if err := readData(&newData, *flagNew); err != nil {
		log.Fatal(err)
	}
	for k := range newData.Cake {
		if _, ok := oldData.Cake[k]; !ok {
			fmt.Printf("ADDED cake \"%s\"\n", k)
		}
	}
	for cakeKey, cakeVal := range oldData.Cake {
		newCake, ok := newData.Cake[cakeKey]
		if !ok {
			fmt.Printf("REMOVED cake \"%s\"\n", cakeKey)
		} else {
			if cakeVal.Time != newCake.Time {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n",
					cakeKey, cakeVal.Time, newCake.Time)
			}
			for ingredientKey := range newCake.IngredientMap {
				_, ok = cakeVal.IngredientMap[ingredientKey]
				if !ok {
					fmt.Printf("ADDED ingredient \"%s\" for cake  \"%s\"\n",
						ingredientKey, cakeKey)
				}
			}
			for ingredientKey, ingredientVal := range cakeVal.IngredientMap {
				newIngredient, ok := newCake.IngredientMap[ingredientKey]
				if !ok {
					fmt.Printf("REMOVED ingredient \"%s\" for cake  \"%s\"\n",
						ingredientKey, cakeKey)
				} else {
					if ingredientVal.IngredientCount != newIngredient.IngredientCount {
						fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n",
							ingredientKey, cakeKey, newIngredient.IngredientCount, ingredientVal.IngredientCount)
					}
					if ingredientVal.IngredientUnit != "" {
						if newIngredient.IngredientUnit == "" {
							fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n",
								ingredientVal.IngredientUnit, ingredientKey, cakeKey)
						}
						if ingredientVal.IngredientUnit != newIngredient.IngredientUnit {
							fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
								ingredientKey, cakeKey, newIngredient.IngredientUnit, ingredientVal.IngredientUnit)
						}
					}
				}
			}
		}
	}
}

func flagAction() {
	flagF := flag.Bool("f", false, "./readDB -f .json/.xml")
	flagOld := flag.String("old", "", "./compareDB --old original_database.xml --new stolen_database.json")
	flagNew := flag.String("new", "", "./compareDB --old original_database.xml --new stolen_database.json")
	flag.Parse()
	if *flagF && flag.NArg() == 1 {
		formatChange()
	} else if flag.NArg() == 0 {
		bdCompare(flagOld, flagNew)
	} else {
		flag.PrintDefaults()
		log.Fatal("Wrong usage")
	}
}

func main() {
	flagAction()
}
