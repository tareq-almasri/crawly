package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	dbInitOnce sync.Once
)

type scannedPage struct {
	ID                     int    `gorm:"primaryKey;autoIncrement" json:"id"`
	HtmlVersion            string `json:"htmlVersion"`
	PageTitle              string `json:"pageTitle"`
	CountOfH1              uint8  `json:"countOfH1"`
	CountOfH2              uint8  `json:"countOfH2"`
	CountOfH3              uint8  `json:"countOfH3"`
	CountOfH4              uint8  `json:"countOfH4"`
	CountOfH5              uint8  `json:"countOfH5"`
	CountOfH6              uint8  `json:"countOfH6"`
	InternalLinksCount     uint8  `json:"internalLinksCount"`
	ExternalLinksCount     uint8  `json:"externalLinksCount"`
	InaccessibleLinksCount uint8  `json:"inaccessibleLinksCount"`
	HasLoginForm           bool   `json:"hasLoginForm"`
	URL                    string `json:"url"`
	Status                 string `json:"status"`
}

func initDB() {
	var err error
	dsn := "admin:password@tcp(127.0.0.1:3306)/crawly_schema?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	if err := db.AutoMigrate(&scannedPage{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

func main() {
	dbInitOnce.Do(initDB)
	fmt.Println("hello world")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("../dist", true)))

	api := router.Group("/api")
	{
		api.GET("/results", func(ctx *gin.Context) {
			var allResults []scannedPage
			if err := db.Find(&allResults).Error; err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to fetch results"})
				return
			}
			ctx.JSON(200, allResults)
		})

		api.POST("/add", func(ctx *gin.Context) {
			var req struct {
				UserInput string `json:"userInput"`
			}
			if err := ctx.BindJSON(&req); err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid input"})
				return
			}
			page := scannedPage{URL: req.UserInput, Status: "pending"}
			if err := db.Create(&page).Error; err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to save"})
				return
			}
			// page.ID is now set by GORM after Create
			go crawl(page.ID, page.URL, make(chan scannedPage))
			ctx.JSON(200, gin.H{"status": "added"})
		})

		api.DELETE("/delete/:id", func(ctx *gin.Context) {
			idStr := ctx.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid ID"})
				return
			}
			if err := db.Delete(&scannedPage{}, id).Error; err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to delete"})
				return
			}
			ctx.JSON(200, gin.H{"status": "deleted"})
		})
	}
	router.Run()
}

func crawl(pageID int, url string, resultChannel chan<- scannedPage) {
	c := colly.NewCollector(colly.MaxDepth(1))
	page := scannedPage{
		ID:  pageID,
		URL: url,
	}
	var internalLinks, externalLinks, inaccessibleLinks uint8
	var hasLoginForm bool
	var htmlVersion, pageTitle string
	var countOfH1, countOfH2, countOfH3, countOfH4, countOfH5, countOfH6 uint8

	visitedUrls := make(map[string]bool)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		htmlVersion = e.DOM.Nodes[0].Data
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		pageTitle = e.Text
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) { countOfH1++ })
	c.OnHTML("h2", func(e *colly.HTMLElement) { countOfH2++ })
	c.OnHTML("h3", func(e *colly.HTMLElement) { countOfH3++ })
	c.OnHTML("h4", func(e *colly.HTMLElement) { countOfH4++ })
	c.OnHTML("h5", func(e *colly.HTMLElement) { countOfH5++ })
	c.OnHTML("h6", func(e *colly.HTMLElement) { countOfH6++ })

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}
		if visitedUrls[link] {
			return
		}
		visitedUrls[link] = true
		if isInternal(url, link) {
			internalLinks++
		} else {
			externalLinks++
		}
	})

	c.OnHTML("form", func(e *colly.HTMLElement) {
		if e.ChildAttr("input[type='password']", "name") != "" {
			hasLoginForm = true
		}
	})

	c.OnError(func(e *colly.Response, err error) {
		inaccessibleLinks++
	})

	c.OnScraped(func(r *colly.Response) {
		page.HtmlVersion = htmlVersion
		page.PageTitle = pageTitle
		page.InternalLinksCount = internalLinks
		page.ExternalLinksCount = externalLinks
		page.InaccessibleLinksCount = inaccessibleLinks
		page.HasLoginForm = hasLoginForm
		page.CountOfH1 = countOfH1
		page.CountOfH2 = countOfH2
		page.CountOfH3 = countOfH3
		page.CountOfH4 = countOfH4
		page.CountOfH5 = countOfH5
		page.CountOfH6 = countOfH6
		page.Status = "Done"
		// Update the existing scannedPage by ID
		db.Model(&scannedPage{}).Where("id = ?", pageID).Updates(page)
		resultChannel <- page
		resultChannel <- page
	})

	c.Visit(url)
}

// Helper to check if a link is internal
func isInternal(base, link string) bool {
	return len(link) >= len(base) && link[:len(base)] == base
}
