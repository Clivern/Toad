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

	"github.com/clivern/toad/internal/app/module"

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

	state := module.State{}

	if get == "health" {
		if state.IsStateless() {
			fmt.Println("i am ok")
			return
		}

		err := state.Init()

		if err != nil {
			panic(err)
		}

		if state.IsDown() {
			panic("I am not ok")
		}

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

		if state.IsStateless() {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
			return
		}

		err := state.Init()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
			return
		}

		if state.IsDown() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "down",
			})
			return
		}

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

	r.GET("/do/:name", func(c *gin.Context) {
		u := uuid.Must(uuid.NewV4(), nil)
		host, _ := os.Hostname()

		log.WithFields(log.Fields{
			"time":          time.Now().Format("Mon Jan 2 15:04:05 2006"),
			"host":          host,
			"uri":           c.Request.URL.Path,
			"method":        c.Request.Method,
			"correlationId": u.String(),
		}).Info("Incoming Request")

		if state.IsStateless() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error! Application is stateless",
			})
			return
		}

		name := c.Param("name")

		err := state.Init()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
			return
		}

		if name == "change" {
			state.Change()
		}

		if name == "reset" {
			state.Reset()
		}

		if name == "host_up" {
			state.HostUp()
		}

		if name == "host_down" {
			state.HostDown()
		}

		if name == "all_up" {
			state.AllUp()
		}

		if name == "all_down" {
			state.AllDown()
		}

		c.JSON(http.StatusOK, gin.H{
			"current": state.Get(),
		})
	})

	runerr := r.Run(
		fmt.Sprintf(":%s", port),
	)

	if runerr != nil {
		panic(runerr.Error())
	}
}
