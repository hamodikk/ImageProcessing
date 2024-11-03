package main

import (
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	err := checkFileExists("images/image1.jpg")
	if err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}

	err = checkFileExists("images/image5.jpg")
	if err == nil {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

func TestLoadImage(t *testing.T) {
	imagePaths := []string{"images/image1.jpg"}

	jobChannel := loadImage(imagePaths)

	job := <-jobChannel
	if job.Image == nil {
		t.Fatalf("expected image, got nil")
	}
}

func BenchmarkImageProcessing(b *testing.B) {
	imagePaths := []string{"images/image1.jpg"}

	for i := 0; i < b.N; i++ {
		channel1 := loadImage(imagePaths)
		channel2 := resize(channel1)
		channel3 := convertToGrayscale(channel2)
		writeResults := saveImage(channel3)

		for success := range writeResults {
			if !success {
				b.Fatalf("failed to save image")
			}
		}
	}
}
