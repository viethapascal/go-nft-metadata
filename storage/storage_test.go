package storage

import (
	"fmt"
	"github.com/go-nft-metadata/image-merge"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestGCSRepository_Write(t *testing.T) {
	gcs := NewGcsRepository("depoc-storage")
	merge := image_merge.NewMergeEngine("../assets/PNG/", image_merge.WithWriter(gcs.Write))
	merge.Background("bg.png")
	merge.AddLayer("4_1_1.png")
	merge.AddLayer("4_2_1.png")
	merge.AddLayer("4_3_1.png")
	merge.AddLayer("4_4_1.png")
	rand.Seed(time.Now().UnixMilli())
	outputPath := fmt.Sprintf("nft/zodiac/%v/mew_new.png", rand.Int())
	err := merge.Merge(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("write success to:", outputPath)
}

func TestGCSRepository_Read(t *testing.T) {
	gcs := NewGcsRepository("depoc-public")
	gcs.Read("nfts/collections/zodiac/config.json")

}
