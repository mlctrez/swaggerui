package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/mlctrez/zipbackpack/httpfs"
)

func main() {

	var sf http.FileSystem
	var err error

	if sf, err = httpfs.NewStaticFileSystem(""); err != nil {
		sf, err = httpfs.NewStaticFileSystem("web")
		if err != nil {
			panic(err)
		}
	}

	handler := http.FileServer(sf)
	http.HandleFunc("/api.yaml", func(rw http.ResponseWriter, r *http.Request) {
		open, e := os.Open(os.Args[1])
		if e != nil {
			panic(e)
		}
		_, e = io.Copy(rw, open)
		if e != nil {
			panic(e)
		}
	})
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			t := template.Must(template.New("index").Parse(index))
			u := fmt.Sprintf("http://localhost:8100/api.yaml?ts=%d", time.Now().Unix())
			e := t.Execute(rw, map[string]string{"APIUrl": u})
			if e != nil {
				panic(e)
			}
			return
		}
		handler.ServeHTTP(rw, r)
	})

	log.Println("listening on :8100")

	exec.Command("open", "http://localhost:8100").CombinedOutput()

	err = http.ListenAndServe(":8100", nil)
	if err != nil {
		panic(err)
	}

}

var index = `
<!-- HTML for static distribution bundle build -->
<!-- 
https://github.com/swagger-api/swagger-ui/blob/master/dist/index.html
modified to template the url in the SwaggerUIBundle constructor 
-->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="./swagger-ui-bundle.js"> </script>
    <script src="./swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: "{{.APIUrl}}",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region

      window.ui = ui
    }
  </script>
  </body>
</html>
`
