package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"log"
	"os"
	"strings"
)

func checkFileExists(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filepath)
	}
	return nil
}

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)

		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			if err := checkFileExists(p); err != nil {
				log.Printf("Error loading image %s: %v\n", p, err)
				continue
			}

			job := Job{InputPath: p,
				Image:   imageprocessing.ReadImage(p),
				OutPath: strings.Replace(p, "images/", "images/output/", 1),
			}

			out <- job
		}
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		defer close(out)
		for job := range input { // Read from the channel
			err := imageprocessing.WriteImage(job.OutPath, job.Image)
			if err != nil {
				log.Printf("Error saving image to %s: %v\n", job.OutPath, err)
				out <- false
			} else {
				out <- true
			}
		}
	}()
	return out
}

func main() {

	imagePaths := []string{"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	for success := range writeResults {
		if success {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!")
		}
	}
}
