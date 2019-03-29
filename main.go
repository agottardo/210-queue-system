package main

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	port := "8888"

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
	router.POST("/status_for_id", handleStatusForID)
	authorized := router.Group("/", gin.BasicAuth(LoadPasswordsFromDisk()))
	authorized.GET("/ta", handleTAStatus)
	authorized.POST("/served", handleServed)
	authorized.POST("/nuke", handleNuke)
	authorized.GET("/jsondump", handleDump)
	router.POST("/join", handleJoinReq)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalln("Listening on port failed with error:", err)
	}
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
		c.SetCookie("queue-csid", CSid, 0, "",
			"", false, false)
		c.HTML(http.StatusOK, "index.tmpl.html", hpv)
		return
	}
	c.SetCookie("queue-csid", CSid, 0, "",
		"cpsc210queue.ugrad.cs.ubc.ca", true, false)
	aheadOfMe, waitTime := JoinQueue(name, CSid, taskInfo)
	if waitTime != -1 {
		jpv := JoinedPageValues{
			AheadOfMe:         strconv.Itoa(aheadOfMe),
			HasEstimate:       waitTime != 0,
			EstimatedWaitTime: strconv.Itoa(waitTime/60) + " minutes",
			JoinedAt:          time.Now().String(),
			Name:              name,
		}
		c.HTML(http.StatusOK, "joined.tmpl.html", jpv)
	} else {
		rpv := RejectedPageValues{
			NumTimesJoined: aheadOfMe,
			Name:           name,
		}
		c.HTML(http.StatusOK, "rejected.tmpl.html", rpv)
	}
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
			"success":       true,
			"nukedDatabase": true,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"success":       false,
		"nukedDatabase": false,
	})
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
	c.Redirect(http.StatusMovedPermanently, "/ta")
}

func handleDump(c *gin.Context) {
	c.File("persistence.json")
}

func handleStatusForID(c *gin.Context) {
	CSid, err := c.Cookie("queue-csid")
	if err != nil || CSid == "" || !IsValidCSid(CSid) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isWaiting, position := QueuePositionForCSID(CSid)
	c.JSON(http.StatusOK, map[string]interface{}{"success": isWaiting, "csid": CSid, "position": position, "waittime": uint(EstimatedWaitTime() / 60)})
}
