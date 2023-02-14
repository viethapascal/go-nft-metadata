package nft

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func TestWeightedGenerate(t *testing.T) {
	choices := []Choice[string, uint]{
		{"Common", 1},
		{"Uncommon", 1},
		{"Rare", 1},
		{"Epic", 1},
		{"Legendary", 1},
		{"Mythic", 1},
	}
	//counter := map[string]int{"Common": 0, "Uncommon": 0, "Rare": 0, "Epic": 0, "Legendary": 0, "Mythic": 0}
	generator := RandomGenerator[string, uint]{}
	generator.FromChoices(choices...)
	traitChoice := []Choice[string, uint]{
		{"Cap", 75},
		{"Clothe", 50},
		{"Boot", 75},
		{"Wallet", 75},
		{"Brace", 5},
		{"Ring", 5},
	}
	generator.FromChoices(traitChoice...)

	for i := 0; i < 10000; i++ {
		generator.GenerateN(3, true, true)
		if len(generator.ResultArr) > 0 {
			log.Println(generator.ResultArr)
		}
		//if idx != -1 {
		//	counter[choices[idx].Item] += 1
		//}
	}
	//result, _ := json.MarshalIndent(counter, "", "\t")
	//log.Println(string(result))
}

func TestMetadataGenerator_GenCombination(t *testing.T) {
	config := &CollectionConfig{
		Decimal:    2,
		Collection: "test",
		Name:       "test name",
		RarityLevel: map[uint]uint{
			1: 100,
			2: 30,
			3: 20,
			4: 10,
			5: 1,
			6: 0,
		},
		TraitList: map[string]CollectionTrait{
			"head": {
				Mandatory: true,
				Rarity:    1,
				ValueList: []ValueRarity{
					{
						Value:  "1",
						Rarity: 2,
					},
					{
						Value:  "2",
						Rarity: 2,
					},
					{
						Value:  "3",
						Rarity: 2,
					},
				},
			},
			"leg": {
				Mandatory: true,
				Rarity:    1,
				ValueList: []ValueRarity{
					{
						Value:  "1",
						Rarity: 3,
					},
					{
						Value:  "2",
						Rarity: 2,
					},
					{
						Value:  "3",
						Rarity: 2,
					},
				},
			},
			"hand": {
				Mandatory: true,
				Rarity:    1,
				ValueList: []ValueRarity{
					{
						Value:  "1",
						Rarity: 3,
					},
					{
						Value:  "2",
						Rarity: 2,
					},
					{
						Value:  "3",
						Rarity: 2,
					},
				},
			},
			"clothes": {
				Mandatory: false,
				Rarity:    3,
				ValueList: []ValueRarity{
					{
						Value:  "1",
						Rarity: 3,
					},
					{
						Value:  "2",
						Rarity: 2,
					},
					{
						Value:  "3",
						Rarity: 2,
					},
				},
			},
			"shoes": {
				Mandatory: false,
				Rarity:    4,
				ValueList: []ValueRarity{
					{
						Value:  "1",
						Rarity: 3,
					},
					{
						Value:  "2",
						Rarity: 2,
					},
					{
						Value:  "3",
						Rarity: 2,
					},
				},
			},
		},
	}
	file, _ := json.MarshalIndent(config, "", " ")

	_ = ioutil.WriteFile("config.json", file, 0644)
	//conf, _ := json.MarshalIndent(config, "", "\t")
	//log.Println("config:", string(conf))
	stats := map[string]int{}
	for i := 0; i < 4; i++ {
		gen := NewGenerator(nil, 1, "", "").UseConfig(config)
		attrs := gen.GenCombination()
		metadata := Metadata{
			Name:        "test",
			Description: "test",
			Image:       "",
			MetadataUrl: "",
			Attributes:  attrs,
		}
		str, _ := json.MarshalIndent(metadata, "", "\t")
		log.Println(string(str))
		for _, a := range metadata.Attributes {
			if _, ok := stats[a.TraitType]; !ok {
				stats[a.TraitType] = 1
			} else {
				stats[a.TraitType] += 1
			}
		}
	}
	log.Println("stats:", stats)

	//metadataStr, _ := json.MarshalIndent(metadata, "", "\t")
	//log.Println(string(metadataStr))
}
