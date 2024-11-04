package main

import (
	"flag"
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Initialize to use all available CPU cores
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// Error handling to check if image files exist
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

// Concurrent version of loadImage
func loadImageConcurrent(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		for _, p := range paths {
			if err := checkFileExists(p); err != nil {
				log.Printf("Error loading image %s: %v\n", p, err)
				continue
			}
			job := Job{
				InputPath: p,
				Image:     imageprocessing.ReadImage(p),
				OutPath:   strings.Replace(p, "images/", "images/output/", 1),
			}
			out <- job
		}
	}()
	return out
}

// Sequential version of loadImage
func loadImageSequential(paths []string) []Job {
	var jobs []Job
	for _, p := range paths {
		if err := checkFileExists(p); err != nil {
			log.Printf("Error loading image %s: %v\n", p, err)
			continue
		}
		job := Job{
			InputPath: p,
			Image:     imageprocessing.ReadImage(p),
			OutPath:   strings.Replace(p, "images/", "images/output/", 1),
		}
		jobs = append(jobs, job)
	}
	return jobs
}

// Concurrent version of resize
func resizeConcurrent(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		for job := range input {
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
	}()
	return out
}

// Sequential version of resize
func resizeSequential(jobs []Job) []Job {
	for i := range jobs {
		jobs[i].Image = imageprocessing.Resize(jobs[i].Image)
	}
	return jobs
}

// Concurrent version of convertToGrayscale
func convertToGrayscaleConcurrent(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		for job := range input {
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
	}()
	return out
}

// Sequential version of convertToGrayscale
func convertToGrayscaleSequential(jobs []Job) []Job {
	for i := range jobs {
		jobs[i].Image = imageprocessing.Grayscale(jobs[i].Image)
	}
	return jobs
}

// Concurrent version of saveImage
func saveImageConcurrent(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		defer close(out)
		for job := range input {
			err := imageprocessing.WriteImage(job.OutPath, job.Image)
			out <- err == nil
		}
	}()
	return out
}

// Sequential version of saveImage
func saveImageSequential(jobs []Job) []bool {
	var results []bool
	for _, job := range jobs {
		err := imageprocessing.WriteImage(job.OutPath, job.Image)
		results = append(results, err == nil)
	}
	return results
}

func main() {
	useGoroutines := flag.Bool("goroutines", false, "Run the pipeline with goroutines")
	flag.Parse()

	imagePaths := []string{
		"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	start := time.Now()

	if *useGoroutines {
		fmt.Println("Running with goroutines...")

		// Run each stage with goroutines
		channel1 := loadImageConcurrent(imagePaths)
		channel2 := resizeConcurrent(channel1)
		channel3 := convertToGrayscaleConcurrent(channel2)
		writeResults := saveImageConcurrent(channel3)

		for success := range writeResults {
			if success {
				fmt.Println("Success!")
			} else {
				fmt.Println("Failed!")
			}
		}

	} else {
		fmt.Println("Running sequentially...")

		// Run each stage sequentially
		jobs := loadImageSequential(imagePaths)
		jobs = resizeSequential(jobs)
		jobs = convertToGrayscaleSequential(jobs)
		results := saveImageSequential(jobs)

		for _, success := range results {
			if success {
				fmt.Println("Success!")
			} else {
				fmt.Println("Failed!")
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Processing took %s\n", elapsed)
}
