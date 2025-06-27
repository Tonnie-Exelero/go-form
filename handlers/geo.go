package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"inquire/now-microservice/templates"

	"github.com/gin-gonic/gin"
)

// Define a local interface to avoid importing templ
type renderable interface {
	Render(ctx context.Context, w io.Writer) error
}

// GeoResponse defines the structure we expect from the external service
type GeoResponse struct {
	IsAllowed       bool   `json:"isAllowed"`
	Message         string `json:"message"`
	RestrictionType string `json:"restrictionType"`
	GeoTargeting    string `json:"geoTargeting"`
	Postcode        string `json:"postcode"`
	Locality        string `json:"locality"`
	State           string `json:"state"`
	Distance        string `json:"distance"`
	CenterPoint     string `json:"centerPoint"`
	RadiusKm        string `json:"radiusKm"`
	LocationText    string `json:"locationText"`
}

// GeoCheckHandler checks if a postcode is allowed for a course
func GeoCheckHandler(c *gin.Context) {
	// 1. Get and validate parameters
	courseID := c.Query("courseid")
	postcode := c.Query("location")

	if courseID == "" || postcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseid and postcode are required"})
		return
	}

	// 2. Call external geo service
	geoData, err := fetchGeoData(c.Request.Context(), courseID, postcode)
	if err != nil {
		log.Printf("Geo check failed: %v", err)
		c.HTML(http.StatusBadGateway, "error.html", gin.H{
			"error": "Geolocation service unavailable",
		})
		return
	}

	// 3. Prepare template data
	templateData := templates.GeoData{
		CourseID:   courseID,
		Postcode:   geoData.Postcode,
		Locality:   geoData.Locality,
		State:      geoData.State,
		IsAllowed:  geoData.IsAllowed,
	}

	// 4. Conditionally render templates
	c.Header("Content-Type", "text/html")
	
	var component renderable

	if geoData.IsAllowed {
		component = templates.GeoForm(templateData)
	} else {
		component = templates.GeoChecker(templateData)
	}

	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("Failed to render template: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering page")
	}
}

// fetchGeoData retrieves geo information from external service
func fetchGeoData(ctx context.Context, courseID, postcode string) (*GeoResponse, error) {
	// Build URL with timeout protection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Construct URL
	apiPrefix := os.Getenv("DEV_API_PREFIX")
	if apiPrefix == "" {
		log.Fatal("DEV_API_PREFIX environment variable not set")
	}
	endpoint := "/geo-check"
	baseURL := apiPrefix + endpoint
	params := url.Values{
		"id":       {courseID},
		"postcode": {postcode},
	}
	finalURL := baseURL + "?" + params.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute request
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Parse response
	return parseGeoResponse(resp.Body)
}

// parseGeoResponse handles JSON decoding
func parseGeoResponse(r io.Reader) (*GeoResponse, error) {
	var result GeoResponse
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		// Read body for error message
		body, _ := io.ReadAll(io.LimitReader(r, 1024))
		return nil, fmt.Errorf("parse JSON: %w\nResponse: %s", err, string(body))
	}
	return &result, nil
}