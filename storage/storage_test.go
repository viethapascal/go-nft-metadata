package storage

import (
	"encoding/json"
	"fmt"
	"github.com/viethapascal/go-nft-metadata/image-merge"
	"log"
	"math/rand"
	url2 "net/url"
	"path/filepath"
	"testing"
	"time"
)

func TestGCSRepository_Write(t *testing.T) {
	gcs := NewGcsRepository("depoc-storage")
	merge := image_merge.NewMergeEngine("../assets/PNG/", "", image_merge.WithWriter(gcs.WriteImage))
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
	data := map[string]interface{}{
		"greeting": "someone",
	}
	bytes_, _ := json.Marshal(data)
	rand.Seed(time.Now().UnixNano())
	id := rand.Uint64()
	err := gcs.Write(bytes_, fmt.Sprintf("test/%d", id), JSONTYPE)
	if err != nil {
		log.Fatal(err)
	}

}

func TestName(t *testing.T) {
	u, _ := url2.Parse("https:/cdn.depoc.io/")
	u2 := url2.URL{
		Scheme: "https",
		Host:   u.Host,
		Path:   filepath.Join("static", "fiel asdf"),
	}
	log.Println(u2.String())
	log.Println(u.String())
	log.Println(u.Scheme)
	filePath := filepath.Join(u.Host, "iamge/asdf")
	log.Println(filePath)
}
