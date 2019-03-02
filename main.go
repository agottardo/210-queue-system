package main

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	LoadDataFromDisk()

	router := gin.New()
	router.Use(gin.Logger())
	router.Delims("{{", "}}")
	router.SetFuncMap(template.FuncMap{
		"NumTimesHelped": NumTimesHelped,
		"RelativeTime":   humanize.Time,
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
	authorized.POST("/nuke", handleNuke)
	authorized.GET("/jsondump", handleDump)
	router.POST("/join", handleJoinReq)
	router.Run(":" + port)
}

func handleJoinReq(c *gin.Context) {
	name := c.PostForm("name")
	CSid := c.PostForm("csid")
	taskInfo := c.PostForm("task")
	if !IsValidCSid(CSid) || name == "" {
		hpv := HomePageValues{Error: "Invalid name or CS ID entered."}
		c.HTML(http.StatusOK, "index.tmpl.html", hpv)
		return
	}
	if HasJoinedQueue(CSid) {
		hpv := HomePageValues{Error: "You have already joined the queue! Click above to see your status."}
		c.HTML(http.StatusOK, "index.tmpl.html", hpv)
		return
	}
	aheadOfMe, waitTime := JoinQueue(name, CSid, taskInfo)
	jpv := JoinedPageValues{
		AheadOfMe:         strconv.Itoa(aheadOfMe),
		HasEstimate:       waitTime != 0,
		EstimatedWaitTime: strconv.Itoa(waitTime/60) + " minutes",
		JoinedAt:          time.Now().String(),
		Name:              name,
	}
	c.HTML(http.StatusOK, "joined.tmpl.html", jpv)
}

func handleStatus(c *gin.Context) {
	spv := StatusPageValues{UnservedEntries()}
	c.HTML(http.StatusOK, "status.tmpl.html", spv)
}

func handleNuke(c *gin.Context) {
	confirm, _ := c.GetPostForm("confirm")
	if confirm == "true" {
		NukeAllTheThings(c.Request.RemoteAddr)
		c.JSON(http.StatusAccepted, gin.H{
			"success": true,
			"nukedDatabase": true,
		})
		return
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"success": false,
			"nukedDatabase": false,
		})
	}
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

func handleDump(c *gin.Context) {
	c.File("persistence.json")
}