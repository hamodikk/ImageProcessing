# Image Processing Pipeline

This Go program loads, resizes, grayscales, and saves images in it's directory.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Usage](#usage)
- [Code Explanation](#code-explanation)
- [Analysis](#Analysis)
- [Observations](#Observations)

## Introduction

This program is created as part of the MSDS-431 class. It is modified from Amrit Singh's work, you can find the github repo [here](https://github.com/code-heim/go_21_goroutines_pipeline) and a video explaining their work [here](https://www.youtube.com/watch?v=8Rn8yOQH62k). It showcases the utility of goroutines in the form of an image processing pipeline

## Features

- Turns images in its directory to grayscale images
- Added test and benchmark.
- Can switch between sequential and concurrent approaches easily.
- Reports execution time.

## Installation

1. Make sure you have [Go installed](https://go.dev/doc/install).
2. Clone this repo to your local machine:
    ```bash
    git clone https://github.com/hamodikk/ImageProcessing.git
    ```
3. Navigate to the project directory
    ```bash
    cd <project-directory>
    ```

## Usage

There are two ways to run the program, sequentially, or concurrently using goroutines.
Use the following command in your terminal or Powershell to run the program sequentially:
```bash
go run .\main.go
```

Or use the following command to run the same program concurrently:
```bash
go run .\main.go -goroutines
```
This will override the default setting to use the goroutines.

The program will handle errors and return "Success!" when the pipeline is complete for each image, as well as return an execution time.

### Code Explanation

I will not explain the entire code as it is modified from an original pipeline. I will explain the parts that I have modified and their purposes.

- Added an init function to use all available CPU cores to better utilize goroutines.
```go
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
```

- Added error handling to make sure the files exist before the pipeline runs. This function is also added to each function (loadImage, resize, convertToGrayscale and saveImage) to handle potential errors.
```go
func checkFileExists(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filepath)
	}
	return nil
}
```

- Created a similar function for loadImage, resize, convertToGrayscale and saveImage to handle concurrent or sequential approaches. I will only give an example for one function as it is implemented in a similar fashion to all four of them.

  * Here is the sequential function for resize:
  ```go
  func resizeSequential(jobs []Job) []Job {
	for i := range jobs {
		jobs[i].Image = imageprocessing.Resize(jobs[i].Image)
	}
	return jobs
  }
  ```

  * Here is the same function running with goroutines:
  ```go
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
  ```


- Added a flag, which defaults the program to run sequentially, allowing us to override it to run the program concurrently using goroutines. I added conditionals later to check for useGoroutines to run the program accordingly.
```go
useGoroutines := flag.Bool("goroutines", false, "Run the pipeline with goroutines")
flag.Parse()
```

- Some other changes worth mentioning:
  * Added code to calculate the process time.
  * Added main_test.go with tests for loadImage and checkFileExists, as well as benchmark functions for concurrent and sequential approaches.

## Analysis

There are two different ways to look at the process times. One is using the benchmarks.

You can run the following in terminal to get the execution time as well as other details:
   Benchmark Concurrent:
   ```bash
   go test -benchmem -run=^$ -bench ^BenchmarkImageProcessingConcurrent$ goroutines_pipeline
   ```
   Benchmark Sequential:
   ```bash
   go test -benchmem -run=^$ -bench ^BenchmarkImageProcessingSequential$ goroutines_pipeline
   ```

You can also run the program with different approaches mentioned in [Usage](#usage). The program will report the execution time.

For . Following are the results:

| Code Language  | Execution time (milliseconds) | Benchmark (seconds) |
|----------------|-------------------------------|---------------------|
| Sequential     | 956.5199                      | 3.060               |
| Concurrent     | 817.3723                      | 2.606               |

The results for both execution time and benchmark show a relative improvement for the same task utilizing goroutines. While the difference might not seem significant, the scale of the program could contribute to the change and a bigger/more complex program could further prove the usefulness of goroutines.

## Observations

I did not get the chance to play around with the resizing of the images, so there is the issue of warped images after grayscale is applied. I might tackle this issue later. Otherwise, the program seems to run without issues.

It is important to note that running `go test -bench=.` returned only one process time. That is why I chose to use the VSCode UI to run the benchmark and found about the code to run the specific benchmark function.

I am not sure if there was a more efficient way to switch between concurrent and sequential approaches, as my thought process likely resulted in a longer code and duplicated lines. A more succint code could result in improved performance, and could be explored later.

Initially the benchmarks were returning results that showed that sequential approach was running faster. After doing some research, I decided to include the init function that allows the program to use all available CPU cores. I am not sure how effective this was in changing the result as I have made additional changes before benchmarking again. I am also unsure if this would be a good approach for more complex programs, but I imagine that it wouldn't cause any issues in my case.