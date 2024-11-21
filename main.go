package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const html = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>File Upload</title>
  </head>
  <body>
    <input type="file" name="file" id="file" /><br />
    <progress id="progress" value="0" max="100"></progress><br />
    <button>Click here to upload</button>
    <div id="msg"></div>
    <script>
      const input = document.querySelector("#file");
      const btn = document.querySelector("button");
      const progress = document.querySelector("#progress");
      const msg = document.querySelector("#msg");

      const sendFile = (e) => {
        e.preventDefault();
        const xhr = new XMLHttpRequest();
        xhr.upload.onprogress = (e) => {
          progress.value = parseInt((e.loaded / e.total) * 100);
          msg.innerHTML = parseInt((e.loaded / e.total) * 100) + "%";
        };
        xhr.onloadend = function () {
          if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
              msg.innerHTML = xhr.response;
            } else {
              msg.innerHTML ="error:"+ xhr.response;
            }
          }
        };

        xhr.open("POST", "/upload", true);
        const formData = new FormData();
        formData.append("file", input.files[0]);
        xhr.send(formData);
      };
      btn.addEventListener("click", sendFile);
    </script>
  </body>
</html>
`

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Browse: true,
	}))
	e.GET("/upload", func(c echo.Context) error {
		return c.HTML(200, html)
	})
	e.POST("/upload", upload)
	e.Logger.Fatal(e.Start(":1323"))
}

func upload(c echo.Context) error {

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully .</p>", file.Filename))
}
