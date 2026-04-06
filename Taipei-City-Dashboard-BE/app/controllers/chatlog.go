// // Package controllers stores all the controllers for the Gin router.
package controllers

import (
	"TaipeiCityDashboardBE/app/models"
	"html"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateChatLog godoc
// @Summary     Create a chat log entry
// @Tags        chatlog
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       session   formData  string  true  "Session ID"
// @Param       question  formData  string  true  "User question"
// @Param       answer    formData  string  true  "AI answer"
// @Success     200       {object}  map[string]interface{}
// @Failure     401       {object}  map[string]interface{}
// @Security    BearerAuth
// @Router      /chatlog [post]
func CreateChatLog(c *gin.Context) {
	var chatLog models.ChatLog
	
	accountID, exists  := c.Get("accountID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	// Sanitize input to prevent XSS
	session := c.PostForm("session")
	question := c.PostForm("question")
	answer := c.PostForm("answer")
	ipAddress := c.ClientIP()
	session = html.EscapeString(session)
	question = html.EscapeString(question)
	answer = html.EscapeString(answer)

	chatLog, _ = models.CreateChatLog(session, question, answer, ipAddress, accountID.(int))
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": chatLog})
}

// GetALLChatLog godoc
// @Summary     Get all chat log sessions
// @Tags        chatlog
// @Produce     json
// @Success     200  {object}  map[string]interface{}
// @Failure     401  {object}  map[string]interface{}
// @Security    BearerAuth
// @Router      /chatlog/session [get]
func GetALLChatLog(c *gin.Context) {
	var chatLogList []models.ChatLog
	
	accountID, exists  := c.Get("accountID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	chatLogList, _ = models.GetALLChatLogSession(accountID.(int))

	type ChatLogSummary struct {
		Session   string    `json:"session"`		
		CreatedAt time.Time `json:"created_at"`
	}

    var summaries []ChatLogSummary
    for _, log := range chatLogList {
        summaries = append(summaries, ChatLogSummary{
            Session:   log.Session,
            CreatedAt: log.CreatedAt,
        })
    }

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": summaries})
}

// GetChatLogDetailBySession godoc
// @Summary     Get chat log by session ID
// @Tags        chatlog
// @Produce     json
// @Param       session  path      string  true  "Session ID"
// @Success     200      {object}  map[string]interface{}
// @Failure     401      {object}  map[string]interface{}
// @Security    BearerAuth
// @Router      /chatlog/session/{session} [get]
func GetChatLogDetailBySession(c *gin.Context) {

	var chatLogList []models.ChatLog
	session := c.Param("session")
	session = html.EscapeString(session)
	accountID, exists  := c.Get("accountID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	chatLogList, _ = models.GetChatLogDetailBySession(session,accountID.(int))
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": chatLogList})
}