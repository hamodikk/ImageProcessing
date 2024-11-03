package main

import (
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	err := checkFileExists("images/image1.jpeg")
	if err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}

	err = checkFileExists("images/nonexistent_image.jpg")
	if err == nil {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

func TestLoadImage(t *testing.T) {
	imagePaths := []string{"images/image1.jpeg"}

	jobChannel := loadImageConcurrent(imagePaths) // Testing concurrent version

	job, ok := <-jobChannel
	if !ok || job.Image == nil {
		t.Fatalf("expected image, got nil or closed channel")
	}
}

func BenchmarkImageProcessingConcurrent(b *testing.B) {
	imagePaths := []string{"images/image1.jpeg", "images/image2.jpeg", "images/image3.jpeg", "images/image4.jpeg"}

	for i := 0; i < b.N; i++ {
		channel1 := loadImageConcurrent(imagePaths)
		channel2 := resizeConcurrent(channel1)
		channel3 := convertToGrayscaleConcurrent(channel2)
		writeResults := saveImageConcurrent(channel3)

		for success := range writeResults {
			if !success {
				b.Fatalf("failed to save image in concurrent processing")
			}
		}
	}
}

func BenchmarkImageProcessingSequential(b *testing.B) {
	imagePaths := []string{"images/image1.jpeg", "images/image2.jpeg", "images/image3.jpeg", "images/image4.jpeg"}

	for i := 0; i < b.N; i++ {
		jobs := loadImageSequential(imagePaths)
		jobs = resizeSequential(jobs)
		jobs = convertToGrayscaleSequential(jobs)
		results := saveImageSequential(jobs)

		for _, success := range results {
			if !success {
				b.Fatalf("failed to save image in sequential processing")
			}
		}
	}
}
