package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Delims("{{", "}}")
	router.SetFuncMap(template.FuncMap{
		"NumTimesHelped": NumTimesHelped,
	})
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/status", handleStatus)
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"210ta": "210ta",
	}))
	authorized.GET("/ta", handleTAStatus)
	authorized.POST("/served", handleServed)
	router.POST("/join", handleJoinReq)
	router.Run(":" + port)
}

func handleJoinReq(c *gin.Context) {
	name := c.PostForm("name")
	CSid := c.PostForm("csid")
	if !IsValidCSid(CSid) || name == "" {
		hpv := HomePageValues{Error: "Invalid name or CS ID entered."}
		c.HTML(http.StatusOK, "index.tmpl.html", hpv)
		return
	}
	if HasJoinedQueue(CSid) {
		hpv := HomePageValues{Error: "You have already joined the queue!"}
		c.HTML(http.StatusOK, "index.tmpl.html", hpv)
		return
	}
	aheadOfMe, waitTime := JoinQueue(name, CSid)
	jpv := JoinedPageValues{
		AheadOfMe:         strconv.Itoa(aheadOfMe),
		HasEstimate:       waitTime != 0,
		EstimatedWaitTime: strconv.Itoa(waitTime) + " seconds",
		JoinedAt:          time.Now().String(),
	}
	c.HTML(http.StatusOK, "joined.tmpl.html", jpv)
}

func handleStatus(c *gin.Context) {
	spv := StatusPageValues{UnservedEntries()}
	c.HTML(http.StatusOK, "status.tmpl.html", spv)
}

func handleTAStatus(c *gin.Context) {
	spv := StatusPageValues{UnservedEntries()}
	c.HTML(http.StatusOK, "tastatus.tmpl.html", spv)
}

func handleServed(c *gin.Context) {
	csid := c.PostForm("csid")
	if csid == "" {
		handleTAStatus(c)
		return
	}
	ServeStudent(csid)
	handleTAStatus(c)
}
