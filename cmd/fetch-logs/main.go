package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/petergtz/alexa-wikipedia/persistence"
	"golang.org/x/sync/semaphore"
)

var s3Client *s3.S3

func main() {
	s3Client = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), ""),
	})))

	var filenames []string
	e := s3Client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String("alexa-wikipedia"),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, object := range p.Contents {
			filename := *object.Key
			if _, e := os.Stat(filename); os.IsNotExist(e) {
				filenames = append(filenames, filename)
			}
		}
		return true
	})
	if e != nil {
		fmt.Printf("%#v\n", e)
		os.Exit(1)
	}
	DownloadInParallel(filenames, 10)

	fileInfos, e := ioutil.ReadDir("cache")
	PanicOnError(e)
	var allLogEntries []persistence.LogEntry
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		file, e := os.Open(filepath.Join("cache", fileInfo.Name()))
		PanicOnError(e)
		content, e := ioutil.ReadAll(file)
		PanicOnError(e)
		var logEntries []persistence.LogEntry
		e = json.Unmarshal(content, &logEntries)
		PanicOnError(e)
		allLogEntries = append(allLogEntries, logEntries...)
	}
	output, e := json.Marshal(allLogEntries)
	PanicOnError(e)
	fmt.Println(string(output))
}

func fail(e error) {
	fmt.Printf("%#v\n", e)
	os.Exit(1)
}

func DownloadInParallel(names []string, numWorkers int64 /*, deletetionFunc func(name string) error*/) []error {
	// var errMutex sync.Mutex
	e := os.MkdirAll("cache", 0755)
	if e != nil {
		fmt.Printf("%#v\n", e)
		os.Exit(1)
	}

	deletionErrs := []error{}

	ctx := context.TODO()
	sem := semaphore.NewWeighted(numWorkers)
	for _, name := range names {
		Must(sem.Acquire(ctx, 1))

		go func(name string) {
			defer sem.Release(1)

			object, e := s3Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String("alexa-wikipedia"),
				Key:    aws.String(name),
			})
			if e != nil {
				fmt.Printf("Could not download object %#v", e)
				return
			}
			defer object.Body.Close()
			file, e := os.Create(filepath.Join("cache", name))
			if e != nil {
				fmt.Printf("Could not create local file %#v", e)
				return
			}
			io.Copy(file, object.Body)

		}(name)
	}
	Must(sem.Acquire(ctx, numWorkers))

	return deletionErrs
}

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

var Must = PanicOnError
