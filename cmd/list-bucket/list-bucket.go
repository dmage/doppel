package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/dmage/doppel/pkg/junit"
)

var (
	OutputDir = flag.String("output-dir", "output/", "")
	Prefix    = flag.String("prefix", "", "")
)

var reNonAlNum = regexp.MustCompile(`[^0-9A-Za-z]`)

func ForEachTestCase(suite *junit.TestSuite, f func(tc *junit.TestCase) error) error {
	for _, s := range suite.Children {
		err := ForEachTestCase(s, f)
		if err != nil {
			return err
		}
	}
	for _, tc := range suite.TestCases {
		err := f(tc)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *Prefix == "" {
		log.Fatal("empty prefix")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatal(err)
	}

	var files []string

	// https://prow.svc.ci.openshift.org/view/gcs/origin-ci-test/logs/release-openshift-ocp-installer-e2e-aws-fips-4.3/324
	objs := client.Bucket("origin-ci-test").Objects(ctx, &storage.Query{
		Prefix: *Prefix,
	})
	for {
		attrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(attrs.Name, ".xml") {
			fmt.Println(attrs.Name)
			files = append(files, attrs.Name)
		}
	}

	for _, filename := range files {
		h := client.Bucket("origin-ci-test").Object(filename)
		r, err := h.NewReader(ctx)
		if err != nil {
			log.Fatal(err)
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		outputNo := 0
		for _, attempt := range []func() error{
			func() error {
				var suites junit.TestSuites
				err = xml.Unmarshal(buf, &suites)
				if err != nil {
					return err
				}

				for _, suite := range suites.Suites {
					ForEachTestCase(suite, func(tc *junit.TestCase) error {
						if tc.FailureOutput == nil {
							return nil
						}

						fmt.Println(tc.Name)
						name := reNonAlNum.ReplaceAllString(tc.Name, "_")

						outputNo++
						f, err := os.Create(fmt.Sprintf("%s%04d_%s.txt", *OutputDir, outputNo, name))
						if err != nil {
							return err
						}
						defer f.Close()

						if tc.SystemOut == "" {
							f.Write([]byte(tc.FailureOutput.Output))
						} else {
							f.Write([]byte(tc.SystemOut))
						}

						return nil
					})
				}
				return nil
			},
			func() error {
				var suite junit.TestSuite
				err = xml.Unmarshal(buf, &suite)
				if err != nil {
					return err
				}

				ForEachTestCase(&suite, func(tc *junit.TestCase) error {
					if tc.FailureOutput == nil {
						return nil
					}

					fmt.Println(tc.Name)
					name := reNonAlNum.ReplaceAllString(tc.Name, "_")

					outputNo++
					f, err := os.Create(fmt.Sprintf("%s%04d_%s.txt", *OutputDir, outputNo, name))
					if err != nil {
						return err
					}
					defer f.Close()

					if tc.SystemOut == "" {
						f.Write([]byte(tc.FailureOutput.Output))
					} else {
						f.Write([]byte(tc.SystemOut))
					}

					return nil
				})
				return nil
			},
		} {
			if err := attempt(); err == nil {
				break
			}
			log.Println(err)
		}
	}
}
