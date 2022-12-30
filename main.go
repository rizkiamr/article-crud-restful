package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Article struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Content     string `form:"content" json:"content"`
	Deleted     bool   `json:"-"`
}

var articles []Article

func createArticle(c echo.Context) error {
	var article Article

	if err := c.Bind(&article); err != nil {
		return err
	}

	articles = append(articles, article)

	printArticles(articles)

	return c.NoContent(http.StatusCreated)
}

func printArticles(articles []Article) {
	for i, article := range articles {
		fmt.Printf("%d. %s\n", i+1, article.Title)
	}
	fmt.Printf("Total article: %d\n", len(articles))
}

func showArticle(c echo.Context) error {
	articleId, err := strconv.Atoi(c.Param("id"))

	if len(articles) < articleId {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	article := articles[articleId]

	if article.Deleted {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, article)
}

func updateArticle(c echo.Context) error {
	articleId, err := strconv.Atoi(c.Param("id"))

	if len(articles) < articleId {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	article := articles[articleId]

	if article.Deleted {
		return c.NoContent(http.StatusNotFound)
	}

	if err = c.Bind(&article); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	articles[articleId] = article

	return c.NoContent(http.StatusOK)
}

func deleteArticle(c echo.Context) error {
	articleId, err := strconv.Atoi(c.Param("id"))

	if len(articles) < articleId {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	article := articles[articleId]
	article.Deleted = true

	return c.NoContent(http.StatusGone)
}

func listArticles(c echo.Context) error {
	shownArticles := make([]Article, len(articles))
	for _, article := range articles {
		if !article.Deleted {
			shownArticles = append(shownArticles, article)
		}
	}
	return c.JSON(http.StatusOK, shownArticles)
}

func main() {
	articles = make([]Article, 0)
	e := echo.New()
	e.POST("/articles", createArticle)
	e.GET("/articles", listArticles)
	e.GET("/articles/:id", showArticle)
	e.PUT("/articles/:id", updateArticle)
	e.DELETE("/articles/:id", deleteArticle)
	e.Logger.Fatal(e.Start(":8080"))
}
