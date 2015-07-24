package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"mime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"io"
)

// Global to hold the bucket name parsed from command-line args
var bucket string

func walkpath(path string, f os.FileInfo, err error) error {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	ext := filepath.Ext(path)
	// This tests if something is a directory; basically, look for a file extension. If that doesn't exist, assume
	// it's a directory and don't upload it to S3
	if len(ext) > 0 {
		key := path[2:len(path)]

		mtype := mime.TypeByExtension(ext)
		fmt.Println(key)

		err = s3Upload(&key, &bucket, fi, &mtype)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Directory, skipping")
	}

	return err
}

func s3Upload(k *string, bucket *string, rd io.ReadSeeker, ct *string) error {
	var err error

	// Will eventually rewrite the program so I can write this as a struct and use s3upload() to receive it
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{
				Filename: "/Users/cyarie/.aws/credentials",
				Profile: "suchgop",
			},
		},
	)

	cl := s3.New(&aws.Config{
		Credentials: creds,
		Region: "us-east-1",
	})

	params := &s3.PutObjectInput{
		Bucket: bucket,
		Key: k,
		Body: rd,
		ContentType: ct,
	}

	resp, err := cl.PutObject(params)
	// This is basically straight from the AWS Go SDK Documentation.
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	fmt.Println(awsutil.StringValue(resp))

	return err

}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	bucket = flag.Arg(1)
	filepath.Walk(root, walkpath)
}
