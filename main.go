package main

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

var config Config

func main() {

	config = ReadConfig()
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
	router.GET("/", handleIndex)
	router.GET("/status", handleStatus)
	router.POST("/status_for_id", handleStatusForID)
	router.GET("/leaveearly", handleLeave)
	router.POST("/isqueueopen", handleIsQueueOpen)
	authorized := router.Group("/", gin.BasicAuth(LoadPasswordsFromDisk()))
	authorized.GET("/ta", handleTAStatus)
	authorized.POST("/served", handleServed)
	authorized.GET("/jsondump", handleDump)
	authorized.POST("/openqueue", handleOpenQueue)
	authorized.POST("/closequeue", handleCloseQueue)
	router.POST("/join", handleJoinReq)
	err := router.Run(":" + config.ListenAt)
	if err != nil {
		log.Fatalln("Listening on port failed with error:", err)
	}
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", HomePageValues{TotalNumStudentsHelped(), ""})
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
	c.SetCookie("queue-csid", CSid, 0, "", "", true, false)
	c.SetCookie("queue-secret", GenerateSecretForCSid(CSid), 0, "", "", true, false)
	aheadOfMe, waitTime := JoinQueue(name, CSid, taskInfo)
	if waitTime != -1 {
		c.HTML(http.StatusOK, "status.tmpl.html", nil)
	} else {
		rpv := RejectedPageValues{
			NumTimesJoined: aheadOfMe,
			Name:           name,
		}
		c.HTML(http.StatusOK, "rejected.tmpl.html", rpv)
	}
}

func handleStatus(c *gin.Context) {
	c.HTML(http.StatusOK, "status.tmpl.html", nil)
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

func handleLeave(c *gin.Context) {
	CSid, err := c.Cookie("queue-csid")
	secret, err := c.Cookie("queue-secret")
	if err != nil || CSid == "" || !IsValidCSid(CSid) ||
		!CheckSecretForCSid(secret, CSid) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	ServeStudent(CSid)
	c.Redirect(http.StatusMovedPermanently, "/status")
}

func handleDump(c *gin.Context) {
	c.File("persistence.json")
}

func handleStatusForID(c *gin.Context) {
	CSid := getCSIDFromCookie(c)
	isWaiting, position := QueuePositionForCSID(CSid)
	c.JSON(http.StatusOK, map[string]interface{}{
		"success":  isWaiting,
		"csid":     CSid,
		"position": position,
		"waittime": uint(EstimatedWaitTime() / 60),
	})
}

func handleIsQueueOpen(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"open": IsQueueOpen(),
	})
}

func handleOpenQueue(c *gin.Context) {
	OpenQueue()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func handleCloseQueue(c *gin.Context) {
	CloseQueue()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func getCSIDFromCookie(c *gin.Context) string {
	CSid, err := c.Cookie("queue-csid")
	secret, err1 := c.Cookie("queue-secret")
	if err != nil || err1 != nil || CSid == "" || !IsValidCSid(CSid) || !CheckSecretForCSid(secret, CSid) {
		c.AbortWithStatus(http.StatusBadRequest)
		return ""
	}
	return CSid
}
