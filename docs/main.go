// Copyright 2022 Chainguard, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"chainguard.dev/melange/pkg/cli"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
draft: false
images: []
type: "article"
toc: true
---
`

func main() {
	melange := cli.New()

	var pathout string
	var baseURL string
	flag.StringVar(&pathout, "out", "./md", "Path to the output directory.")
	flag.StringVar(&baseURL, "baseurl", "/chainguard/chainguard-enforce/melange-docs/", "Base URL for melange-docs on Academy site.")
	flag.Parse()

	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.Split(strings.TrimSuffix(name, path.Ext(name)), "/")[:1][0]
		url := baseURL + strings.ToLower(base) + "/"
		return fmt.Sprintf(fmTemplate, now, strings.ReplaceAll(base, "_", " "), base, url)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return baseURL + strings.ToLower(base) + "/"
	}

	if err := os.MkdirAll(pathout, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatalf("error creating directory %#v: %#v", pathout, err)
	}

	fmt.Printf("Generating Markdown documentation into directory %#v\n", pathout)
	err := doc.GenMarkdownTreeCustom(melange, pathout, filePrepender, linkHandler)
	if err != nil {
		log.Fatalf("error creating documentation: %#v", err)
	}
}
