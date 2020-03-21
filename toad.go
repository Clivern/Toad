// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var port string
	var get string

	flag.StringVar(&port, "port", "8080", "port")
	flag.StringVar(&get, "get", "", "get")
	flag.Parse()

	if get == "health" {
		fmt.Println("i am ok")
		return
	}

	if get == "release" {
		fmt.Println(
			fmt.Sprintf(
				`Toad Version %v Commit %v, Built @%v`,
				version,
				commit,
				date,
			),
		)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	gin.DisableConsoleColor()

	r := gin.Default()

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	r.GET("/_health", func(c *gin.Context) {
		u := uuid.Must(uuid.NewV4(), nil)
		host, _ := os.Hostname()

		log.WithFields(log.Fields{
			"time":          time.Now().Format("Mon Jan 2 15:04:05 2006"),
			"host":          host,
			"uri":           c.Request.URL.Path,
			"method":        c.Request.Method,
			"correlationId": u.String(),
		}).Info("Incoming Request")

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/", func(c *gin.Context) {
		u := uuid.Must(uuid.NewV4(), nil)
		host, _ := os.Hostname()

		log.WithFields(log.Fields{
			"time":          time.Now().Format("Mon Jan 2 15:04:05 2006"),
			"host":          host,
			"uri":           c.Request.URL.Path,
			"method":        c.Request.Method,
			"correlationId": u.String(),
		}).Info("Incoming Request")

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().Format("Mon Jan 2 15:04:05 2006"),
			"host":    host,
			"release": fmt.Sprintf(`Toad Version %v Commit %v, Built @%v`, version, commit, date),
		})
	})

	runerr := r.Run(
		fmt.Sprintf(":%s", port),
	)

	if runerr != nil {
		panic(runerr.Error())
	}
}
