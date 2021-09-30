// +build ignore

// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	css, err := os.Create("csscompacteds.go")
	if err != nil {
		panic(err)
	}
	defer css.Close()
	css.Write([]byte(fmt.Sprintf("// Code generated by go generate; DO NOT EDIT.\npackage %s\nvar css=`", os.Args[1])))
	for _, f := range files {
		fileExtension := filepath.Ext(f.Name())
		if fileExtension == ".css" {
			if cssfile, err := os.Open(f.Name()); err == nil {

				reader := bufio.NewReader(cssfile)

				for {
					var line string
					var err error
					if line, err = reader.ReadString('\n'); err == nil {

						line = strings.Trim(line, " ")
						line = strings.Trim(line, "\n")
						space := regexp.MustCompile(`\s+`)
						line = space.ReplaceAllString(line, " ")

						css.Write([]byte(line))
					}

					if err != nil {
						break
					}
				}

			}

		}
		if fileExtension == ".links" {
			if linkfile, err := os.Open(f.Name()); err == nil {

				reader := bufio.NewReader(linkfile)

				var url []byte
				var isPrefix bool
				for {
					url, isPrefix, err = reader.ReadLine()
					var rsp *http.Response
					if rsp, err = http.Get(string(url)); err != nil {
						panic(err)
					}

					defer rsp.Body.Close()

					if rsp.StatusCode == http.StatusOK {

						reader := bufio.NewReader(rsp.Body)
						for {
							var line string
							var err error
							if line, err = reader.ReadString('\n'); err == nil {

								line = strings.Trim(line, " ")
								line = strings.ReplaceAll(line, "\n", " ")
								//line = strings.Trim(line, "\n")
								space := regexp.MustCompile(`\s+`)
								line = space.ReplaceAllString(line, " ")

								css.Write([]byte(line))
							}

							if err != nil {
								break
							}
						}
					}
					// If we've reached the end of the line, stop reading.
					if !isPrefix {
						break
					}

					// If we're just at the EOF, break
					if err != nil {
						break
					}

				}

				/*

					for {
						var url string
						var err error
						if url, err = reader.ReadString('\n'); err == nil {
							var rsp *http.Response
							fmt.Printf("Get url %s\n", url)
							if rsp, err = http.Get(url); err != nil {
								panic(err)
							}

							defer rsp.Body.Close()

							if rsp.StatusCode == http.StatusOK {

								reader := bufio.NewReader(rsp.Body)
								for {
									var line string
									var err error
									if line, err = reader.ReadString('\n'); err == nil {

										line = strings.Trim(line, " ")
										line = strings.Trim(line, "\n")
										space := regexp.MustCompile(`\s+`)
										line = space.ReplaceAllString(line, " ")

										css.Write([]byte(line))
									}

									if err != nil {
										break
									}
								}
							}

						}
					}*/
			}
		}

	}
	css.Write([]byte(fmt.Sprint("`")))

}
