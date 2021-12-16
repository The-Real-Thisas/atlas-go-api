package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/thinkerou/favicon"
)

type Data struct {
	Questions []struct {
		Subject      string   `json:"subject"`
		Topic        string   `json:"topic"`
		ID           string   `json:"id"`
		PaperNumber  string   `json:"paper_number"`
		Topics       string   `json:"topics"`
		QuestionUrls []string `json:"question_urls"`
		AnswerUrls   []string `json:"answer_urls"`
	} `json:"questions"`
}

type Topics struct {
	Physics struct {
		Topics     []string `json:"topics"`
		TopicNames []string `json:"topicNames"`
	} `json:"physics"`
	Biology struct {
		Topics     []string `json:"topics"`
		TopicNames []string `json:"topicNames"`
	} `json:"biology"`
	Chemistry struct {
		Topics     []string `json:"topics"`
		TopicNames []string `json:"topicNames"`
	} `json:"chemistry"`
	ComputerScience struct {
		Topics     []string `json:"topics"`
		TopicNames []string `json:"topicNames"`
	} `json:"computer-science"`
	Psychology struct {
		Topics     []string `json:"topics"`
		TopicNames []string `json:"topicNames"`
	} `json:"psychology"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Print Banner
	fmt.Println(`
    ___   __  __           ___    ____  ____
   /   | / /_/ /___  _____/   |  / __ \/  _/
  / /| |/ __/ / __ \/ ___/ /| | / /_/ // /  
 / ___ / /_/ / /_/ (__  ) ___ |/ ____// /   
/_/  |_\__/_/\__,_/____/_/  |_/_/   /___/   
                                     
	`)
	fmt.Println("[*] Listening on http://localhost:" + port)

	//-------------------------------------------------------
	// Load Questions, Answers and Options
	//-------------------------------------------------------
	content, err := ioutil.ReadFile("./data/questions.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	topicContent, err := ioutil.ReadFile("./data/topics.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var topics Topics
	err = json.Unmarshal(topicContent, &topics)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	//-------------------------------------------------------

	router := gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}), gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.Use(favicon.New("./favicon.ico"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/:subject", func(c *gin.Context) {
		subject := c.Param("subject")
		// Send the topics for the subject
		if subject == "physics" {
			c.JSON(http.StatusOK, topics.Physics)
		} else if subject == "biology" {
			c.JSON(http.StatusOK, topics.Biology)
		} else if subject == "chemistry" {
			c.JSON(http.StatusOK, topics.Chemistry)
		} else if subject == "computer-science" {
			c.JSON(http.StatusOK, topics.ComputerScience)
		} else if subject == "psychology" {
			c.JSON(http.StatusOK, topics.Psychology)
		} else {
			c.JSON(http.StatusOK, nil)
		}

	})

	router.GET("/all/:subject", func(c *gin.Context) {
		subject := c.Param("subject")
		var subjects = []string{"physics", "chemistry", "biology", "computer-science", "psychology"}
		if !contains(subjects, subject) {
			c.HTML(http.StatusNotFound, "lost.tmpl.html", nil)
			return
		}

		var questions []struct {
			Subject      string   `json:"subject"`
			Topic        string   `json:"topic"`
			ID           string   `json:"id"`
			PaperNumber  string   `json:"paper_number"`
			Topics       string   `json:"topics"`
			QuestionUrls []string `json:"question_urls"`
			AnswerUrls   []string `json:"answer_urls"`
		}

		// Get only questions for the subject
		for _, question := range payload.Questions {
			if question.Subject == subject {
				questions = append(questions, question)
			}
		}

		c.JSON(200, questions)
	})

	router.GET("/:subject/*topic", func(c *gin.Context) {
		subject := c.Param("subject")
		topic := c.Param("topic")
		var subjects = []string{"physics", "chemistry", "biology", "computer-science", "psychology"}
		if !contains(subjects, subject) {
			c.HTML(http.StatusNotFound, "lost.tmpl.html", nil)
			return
		}

		var questions []struct {
			Subject      string   `json:"subject"`
			Topic        string   `json:"topic"`
			ID           string   `json:"id"`
			PaperNumber  string   `json:"paper_number"`
			Topics       string   `json:"topics"`
			QuestionUrls []string `json:"question_urls"`
			AnswerUrls   []string `json:"answer_urls"`
		}

		// Get only questions for the subject and topic combination
		for _, question := range payload.Questions {
			if question.Subject == subject && ("/"+question.Topic) == topic {
				questions = append(questions, question)
			}
		}

		c.JSON(200, questions)
	})

	router.GET("/id/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Get question with the given id
		for _, question := range payload.Questions {
			if question.ID == id {
				c.JSON(200, question)
				return
			}
		}
		c.JSON(404, gin.H{"error": "Question not found"})
	})

	router.GET("/filter/:customquery", func(c *gin.Context) {
		query := c.Param("customquery")
		customquery := strings.Split(query, "&")
		subject := customquery[0]
		var subjects = []string{"physics", "chemistry", "biology", "computer-science", "psychology"}
		if !contains(subjects, subject) {
			c.HTML(http.StatusNotFound, "lost.tmpl.html", nil)
		}
		topics := strings.Split(customquery[1], ",")
		years := strings.Split(customquery[2], ",")

		var questions []struct {
			Subject      string   `json:"subject"`
			Topic        string   `json:"topic"`
			ID           string   `json:"id"`
			PaperNumber  string   `json:"paper_number"`
			Topics       string   `json:"topics"`
			QuestionUrls []string `json:"question_urls"`
			AnswerUrls   []string `json:"answer_urls"`
		}

		if topics[0] == "ANY" && years[0] == "ANY" {
			// Get only questions for the subject
			for _, question := range payload.Questions {
				if question.Subject == subject {
					questions = append(questions, question)
				}
			}
			c.JSON(200, questions)
		}

		if topics[0] != "ANY" && years[0] == "ANY" {
			// Get only questions for the subject and topic combination
			for _, question := range payload.Questions {
				if question.Subject == subject && contains(topics, question.Topic) == true {
					questions = append(questions, question)
				}
			}
			c.JSON(200, questions)
		}

		if topics[0] == "ANY" && years[0] != "ANY" {
			// Get only questions for the subject and year combination
			for _, question := range payload.Questions {
				Questionyear := question.PaperNumber[8:12]
				if question.Subject == subject && contains(years, Questionyear) == true {
					questions = append(questions, question)
				}
			}
			c.JSON(200, questions)
		}

		if topics[0] != "ANY" && years[0] != "ANY" {
			// Get only questions for the subject, topic and year combination
			for _, question := range payload.Questions {
				Questionyear := question.PaperNumber[8:12]
				if question.Subject == subject && contains(topics, question.Topic) == true && contains(years, Questionyear) == true {
					questions = append(questions, question)
				}
			}
			c.JSON(200, questions)
		}
	})

	router.Run(":" + port)
}
