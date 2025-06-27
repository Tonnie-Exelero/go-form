package handlers

import (
	"inquire/now-microservice/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler renders the main index page using Templ
func HomeHandler(c *gin.Context) {
	courseid := c.Query("courseid")
    restriction := c.Query("restriction")
    
    // Determine if restriction should apply
    showRestriction := restriction == "warning" || restriction == "restriction"
    
    // Prepare template data
    data := templates.GeoData{
        Restriction: showRestriction,
        CourseID:    courseid,
    }
    
    c.Writer.Header().Set("Content-Type", "text/html")
    templates.Home(data).Render(c.Request.Context(), c.Writer)
}

func CloseModal(c *gin.Context) {
	c.Status(http.StatusOK)
}