package image_merge

import (
	"log"
	"testing"
)

func TestNewMergeEngine(t *testing.T) {
	merge := NewMergeEngine("../assets/PNG/", "../assets/tmp")
	merge.Background("bg.png")
	merge.AddLayer("4_1_1.png")
	merge.AddLayer("4_2_1.png")
	merge.AddLayer("4_3_1.png")
	merge.AddLayer("4_4_1.png")
	err := merge.Merge("output_test.png")
	if err != nil {
		log.Fatal(err)
	}
}
