package nft

import (
	"bytes"
	"encoding/json"
	"github.com/viethapascal/go-nft-metadata/image-merge"
	"github.com/viethapascal/go-nft-metadata/storage"
	"github.com/viethapascal/go-nft-metadata/utils"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type MetadataType int

const (
	Combination MetadataType = iota
	Singleton
)

type RarityString string

type Attribute struct {
	TraitType string      `json:"trait_type,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	//Rarity int
}

func (a Attribute) StringValue() string {
	rv := reflect.ValueOf(a.Value)
	v := reflect.Indirect(rv)
	return v.String()
}

type AttrImage struct {
	AttrKey   string
	ImagePath string
}
type ResourceMap struct {
	Resources []AttrImage `json:"resources"`
}

// Metadata example
/*
{
  "name": "Template metadata",
  "description": "something here",
  "image_url": "https://apt.nft/<collection>/<token_id>.jpg",
  "metadata_url": "https://apt.nft/<collection>/<token_id>.json",
  "attributes": [
    {
      "trait_type": "name",
      "value": "cat"
    },
    {
     "trait_type": "background",
      "value": "default"
    },
    {
     "trait_type": "body",
      "value": 1
    },
    {
     "trait_type": "jacket",
      "value": 1
    },
    {
     "trait_type": "eye",
      "value": 2
    },
    {
     "trait_type": "eye wear",
      "value": 3
    },
    {
     "trait_type": "head wear",
      "value": 4
    }
  ]
}
*/
type Metadata struct {
	Name        string      `json:"name"`
	TokenId     string      `json:"token_id"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	MetadataUrl string      `json:"metadata_url"`
	ExternalUrl string      `json:"external_url"`
	Attributes  []Attribute `json:"attributes"`
}

type ValueRarity struct {
	Value  string `json:"value"`
	Rarity uint   `json:"rarity"`
}
type CollectionTrait struct {
	Mandatory bool          `json:"mandatory"`
	Rarity    uint          `json:"rarity"`
	ValueList []ValueRarity `json:"value_list"`
}
type RarityLevel map[uint]uint

// CollectionConfig
type CollectionConfig struct {
	Collection  string                     `json:"collection"`
	Name        string                     `json:"name"`
	RarityLevel RarityLevel                `json:"rarity_level"`
	Decimal     uint8                      `json:"decimal"`
	TraitList   map[string]CollectionTrait `json:"trait_list"`
}

type MetadataGenerator struct {
	Storage    storage.NFTStorage
	ImageMerge *image_merge.MergeEngine
	Type       MetadataType
	Config     *CollectionConfig
	//Generator    *RandomGenerator[interface{}, uint]
	AssetPath  string
	TargetPath string
	CDNPrefix  string
}

func (m *MetadataGenerator) UseConfig(conf *CollectionConfig) *MetadataGenerator {
	m.Config = conf
	return m
}
func NewGenerator(storage storage.NFTStorage, metadata MetadataType, assetPath string, targetPath string) *MetadataGenerator {
	merge := image_merge.NewMergeEngine(assetPath, targetPath,
		image_merge.WithReader(storage.ReadImage),
		image_merge.WithWriter(storage.WriteImage))
	return &MetadataGenerator{Storage: storage, Type: metadata, ImageMerge: merge, AssetPath: assetPath, TargetPath: targetPath, CDNPrefix: os.Getenv("APP_CDN_PREFIX")}
}

func (m *MetadataGenerator) GenCombination() []Attribute {
	// Gen trait Collection
	attrs := make([]Attribute, 0)
	//selectedTraits := make([]string, 0)
	g1 := &RandomGenerator[string, uint]{Decimal: 2}

	for name, trait := range m.Config.TraitList {
		if !trait.Mandatory {
			rarity := m.Config.RarityLevel[trait.Rarity]
			g1.AddChoice(NewChoice(name, rarity))
		} else {
			g1.ResultArr = append(g1.ResultArr, name)
		}
	}
	g1.GenerateN(len(g1.Choices), true, true)
	selectedTraits := g1.ResultArr
	// Generate value
	for _, name := range selectedTraits {
		g2 := &RandomGenerator[string, uint]{Decimal: 2}
		for _, val := range m.Config.TraitList[name].ValueList {
			rarity := m.Config.RarityLevel[val.Rarity]
			g2.AddChoice(NewChoice(val.Value, rarity))
		}
		value := g2.PickOne(false)
		attrs = append(attrs, Attribute{name, value})
	}
	return attrs
}
func (m *MetadataGenerator) LoadConfig(path string) error {
	data, err := m.Storage.Read(path)
	if err != nil {
		return err
	}
	conf := CollectionConfig{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&conf)
	m.UseConfig(&conf)
	return nil
}

func (m *MetadataGenerator) MergeItems(metadata *Metadata) error {
	layers := make([]string, 0)
	mType := ""
	name := ""
	for _, attr := range metadata.Attributes {
		if attr.TraitType == "type" {
			mType = attr.StringValue()
		}
		if attr.TraitType == "name" {
			name = attr.StringValue()
		}
		if attr.TraitType == "background" {
			m.ImageMerge.Background(attr.StringValue() + ".png")

		}
	}
	m.AssetPath = filepath.Join(m.AssetPath, name)
	m.ImageMerge.AssetPath = m.AssetPath
	log.Println(m.AssetPath)
	for _, attr := range metadata.Attributes {
		rv := reflect.ValueOf(attr.Value)
		v := reflect.Indirect(rv)
		switch attr.TraitType {
		case "name", "type", "background":
			continue
		default:
			filename := strings.Join([]string{mType, attr.TraitType, v.String()}, "_") + ".png"
			layers = append(layers, filename)
		}
	}
	sort.Strings(layers)
	m.ImageMerge.FromLayers(layers...)
	err := m.ImageMerge.Merge(metadata.TokenId + ".png")
	if err != nil {
		return err
	}
	return nil

}
func (m *MetadataGenerator) GenerateMetadata(name, description, tokenId string, metadataExt, imageExt string) (*Metadata, error) {
	attrs := m.GenCombination()
	dat := &Metadata{
		Name:        name,
		Description: description,
		TokenId:     tokenId,
		Image:       "",
		MetadataUrl: "",
		Attributes:  attrs,
	}
	err := m.MergeItems(dat)
	if err != nil {
		return nil, err
	}
	metadataPath := filepath.Join(m.TargetPath, dat.TokenId+metadataExt)
	dat.Image = utils.BuildUrl(m.CDNPrefix, m.TargetPath, dat.TokenId+imageExt)
	dat.MetadataUrl = utils.BuildUrl(m.CDNPrefix, m.TargetPath, dat.TokenId)
	b, _ := json.Marshal(dat)
	err = m.Storage.Write(b, metadataPath, storage.JSONTYPE)
	if err != nil {
		return nil, err
	}
	return dat, nil
}
