package main

import (
	"encoding/json"
	"fmt"
	"github.com/viethapascal/go-nft-metadata/nft"
	"github.com/viethapascal/go-nft-metadata/storage"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func changeFileName() {
	files, err := ioutil.ReadDir("assets/PNG/")
	if err != nil {
		log.Fatal(err)
	}
	traitMap := map[string]string{
		"1": "body",
		"2": "eye",
		"3": "clothe",
		"4": "belt",
		"5": "eye_wear",
		"6": "misc",
		"7": "hat",
	}
	for _, file := range files {
		filename := file.Name()
		split := strings.Split(filename, "_")
		data, err := os.ReadFile("assets/PNG/" + filename)
		if err != nil {
			log.Fatal(err)
			continue
		}
		newFilename := filepath.Join("assets", "cat_"+traitMap[split[1]]+"_"+split[2])
		os.WriteFile(newFilename, data, 0644)
		//fmt.Println(file.Name(), file.IsDir())
	}
}
func main() {
	rand.Seed(time.Now().Unix())
	//id := rand.Int()
	//gcs := storage.NewGcsRepository("depoc-public")
	st := storage.LocalStorage{}
	//generator := nft.NewGenerator(st, 0, "nfts/collections/zodiac", "nfts/collections/zodiac/data")
	//err := generator.LoadConfig("nfts/collections/zodiac/config1.json")
	generator := nft.NewGenerator(st, 0, "assets/PNG", "assets/output/")
	err := generator.LoadConfig("nft/config1.json")
	if err != nil {
		log.Fatal(err)
	}
	res, _ := json.MarshalIndent(generator.Config, "", "\t")
	log.Println(string(res))
	metadata, err := generator.GenerateMetadata("zodiac", "test zodiac storage", fmt.Sprintf("%d", 1), "", ".png")
	if err != nil {
		log.Fatal(err)
	}
	res, _ = json.MarshalIndent(metadata, "", "\t")
	log.Println(string(res))

}
