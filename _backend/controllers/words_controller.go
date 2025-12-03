package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zyaaco/wowdle_backend/models"
)

type Feedback struct {
	Correct bool `json:"correct"`
	// We return an array of states for each letter
	// 0 = Gray (Wrong), 1 = Yellow (Wrong spot), 2 = Green (Right spot)
	LetterStates []int `json:"letter_states"`
}

func CompareWord(c *gin.Context) {
	targetWord, err := models.GetWord()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	guess := strings.ToUpper(strings.TrimSpace(c.Query("guess")))
	targetWord = strings.ToUpper(targetWord)

	// Validation
	if len(guess) != len(targetWord) {
		c.JSON(400, gin.H{"error": "Word length mismatch"})
		return
	}

	// Calculate Wordle Logic (2 - Green/ 1- Yellow/ 0- Gray)
	states := calculateStates(guess, targetWord)

	// 5. Check strictly for win condition
	isCorrect := guess == targetWord

	c.JSON(200, Feedback{
		Correct:      isCorrect,
		LetterStates: states,
	})
}

func calculateStates(guess, target string) []int {
	states := make([]int, len(guess))
	targetBytes := []byte(target)

	// Find Greens
	for i := 0; i < len(guess); i++ {
		if guess[i] == targetBytes[i] {
			states[i] = 2      // Green
			targetBytes[i] = 0 // Mark this letter in target as "used"
		}
	}

	// Find Yellows
	for i := 0; i < len(guess); i++ {
		if states[i] == 2 {
			continue
		} // Skip already green

		// Look for the letter elsewhere in the target
		for j := 0; j < len(targetBytes); j++ {
			if guess[i] == targetBytes[j] && targetBytes[j] != 0 {
				states[i] = 1      // Yellow
				targetBytes[j] = 0 // Mark used so we don't count it twice
				break
			}
		}
	}

	return states
}

func GetValidWords(c *gin.Context) {
	// Serve valid_words.txt
	c.File("valid_words.txt")
}
