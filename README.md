# htmlrender

[![Build Status](https://travis-ci.org/felipeweb/htmlrender.svg?branch=master)](https://travis-ci.org/felipeweb/htmlrender) [![GoDoc](https://godoc.org/github.com/felipeweb/htmlrender?status.svg)](https://godoc.org/github.com/felipeweb/htmlrender)

Package htmlrender is a package that provides functionality for easily rendering HTML templates.

  ```go
  package main

  import (
      "net/http"
      "github.com/felipeweb/htmlrender"
  )

  func main() {
      r := htmlrender.New()
      mux := http.NewServeMux()
      mux.HandleFunc("/html", func(w http.ResponseWriter, req *http.Request) {
          // Assumes you have a template in ./views called "example.html".
          // $ mkdir -p templates && echo "<h1>Hello HTML world.</h1>" > views/example.html
          r.HTML(w, http.StatusOK, "example", nil)
      })
      http.ListenAndServe("0.0.0.0:3000", mux)
  }
  ```
  