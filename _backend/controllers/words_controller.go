package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zyaaco/wowdle_backend/models"
)

func GetWord(c *gin.Context) {
	word, err := models.GetWord()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"word": word,
	})
}

func UpsertWord(c *gin.Context) {
	word := c.PostForm("word")
	err := models.ChangeWord(word)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "word updated",
	})
}
