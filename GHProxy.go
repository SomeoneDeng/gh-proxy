package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	en := gin.Default()
	en.Use(cors.Default())

	en.StaticFS("./gpx/static", http.Dir("./gpx/static"))
	en.LoadHTMLGlob("./gpx/html/*")
	en.Handle("GET", "/gpx", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	g := en.Group("/gpx/api")
	{
		g.GET("/download", func(c *gin.Context) {
			c.Request.ParseForm()
			urlRaw := c.Request.Form.Get("url")
			// log.Println("base --> ", urlRaw)
			raw, _ := base64.StdEncoding.DecodeString(urlRaw)
			// log.Println("url --> ", string(raw))
			target, _ := url.PathUnescape(string(raw))
			log.Println("target --> ", target)
			resp, err := http.Get(target)
			if err != nil {
				log.Println("proxy error --> ", err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
			for k, v := range resp.Header {
				c.Header(k, v[0])
			}
			for {
				_, err = io.CopyN(c.Writer, resp.Body, 128)
				// log.Println("proxy n --> ", n)
				if err != nil {
					log.Println("read error --> ", err.Error())
					break
				}
			}
			resp.Body.Close()

			// c.Writer.Write(data)
		})
	}

	en.Run(":9999")
}
