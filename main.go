package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))

	fsys, err := fs.Sub(f, "assets")
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.SetHTMLTemplate(templ)
	router.StaticFS("/assets", http.FS(fsys))

	router.GET("/", IndexHandler)
	router.GET("/recipes/:id", RecipeHandler)

	router.Run()
}

var recipes []Recipe

//go:embed assets/* templates/* 404.html recipes.json
var f embed.FS

func init() {
	data, _ := f.ReadFile("recipes.json")
	json.Unmarshal(data, &recipes)
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"recipes": recipes,
	})
}

func RecipeHandler(c *gin.Context) {
	for _, recipe := range recipes {
		if recipe.ID == c.Param("id") {
			c.HTML(http.StatusOK, "recipe.tmpl", gin.H{
				"recipe": recipe,
			})
			return
		}
	}
	c.File("404.html")
}

type Recipe struct {
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
	Steps       []string     `json:"steps,omitempty"`
	Picture     string       `json:"imageURL,omitempty"`
}

type Ingredient struct {
	Quantity string `json:"quantity,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
}
