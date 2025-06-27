package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"inquire/now-microservice/templates"

	"github.com/gin-gonic/gin"
)

// EnquiryRequest defines the structure for the external API request
type EnquiryRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CourseID    string `json:"course_id"`
	Message     string `json:"message"`
	GASource    string `json:"ga_source"`
	GAMedium    string `json:"ga_medium"`
	OriginalURL string `json:"original_url"`
	Location    string `json:"suburb"`
	Education   string `json:"education_level"`
	StartDate   string `json:"start_preference"`
	Reason      string `json:"reason"`
}

// EnquiryResponse defines the structure for the external API response
type EnquiryResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	EnquiryID  int    `json:"enquiry_id,omitempty"`
	CourseID   int    `json:"course_id,omitempty"`
	CourseUUID string `json:"course_uuid,omitempty"`
	Error      string `json:"error,omitempty"`
	Code       string `json:"code,omitempty"`
}

func SubmitFormHandler(c *gin.Context) {
	// 1. Collect form data
	queryParams := c.Request.URL.Query()
	
	formData := map[string]string{
		"firstname":  queryParams.Get("firstname"),
		"lastname":   queryParams.Get("lastname"),
		"phone":      queryParams.Get("phone"),
		"email":      queryParams.Get("email"),
		"location":   queryParams.Get("location"),
		"education":  queryParams.Get("education"),
		"startdate":  queryParams.Get("startdate"),
		"reason":     queryParams.Get("reason"),
		"courseid":   queryParams.Get("courseid"),
	}

	// 2. Prepare external API request
	reqBody := EnquiryRequest{
		FirstName:   formData["firstname"],
		LastName:    formData["lastname"],
		Email:       formData["email"],
		Phone:       formData["phone"],
		CourseID:    formData["courseid"],
		Reason:      formData["reason"],
		Location:    formData["location"],
		Education:   formData["education"],
		StartDate:   formData["startdate"],
		GASource:    "website",
		GAMedium:    "form",
		OriginalURL: c.Request.Referer(),
	}

	// 3. Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		renderErrorTemplate(c, "Failed to process form data")
		return
	}

	// 4. Create HTTP request
	apiPrefix := os.Getenv("DEV_API_PREFIX")
	if apiPrefix == "" {
		log.Fatal("DEV_API_PREFIX environment variable not set")
	}
	endpoint := "/enquiry"
	fullURL := apiPrefix + endpoint

	req, err := http.NewRequest(
		"POST",
		fullURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		renderErrorTemplate(c, "Failed to create submission request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 5. Send request with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("API request error: %v", err)
		renderErrorTemplate(c, "Failed to connect to submission service")
		return
	}
	defer resp.Body.Close()

	// 6. Handle response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Response read error: %v", err)
		renderErrorTemplate(c, "Failed to read submission response")
		return
	}

	// 7. Parse JSON response
	var apiResponse EnquiryResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("JSON parse error: %v\nResponse: %s", err, string(body))
		renderErrorTemplate(c, "Invalid response from submission service")
		return
	}

	// 8. Prepare template data
	templateData := templates.CallbackData{
		Success:    apiResponse.Success,
		Title:      "Submission Result",
		Message:    apiResponse.Message,
		EnquiryID:  apiResponse.EnquiryID,
		CourseID:   apiResponse.CourseID,
		CourseUUID: apiResponse.CourseUUID,
		Error:      apiResponse.Error,
		ErrorCode:  apiResponse.Code,
	}

	// 9. Render template
	c.Writer.Header().Set("Content-Type", "text/html")
	if err := templates.Callback(templateData).Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("Template render error: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering response page")
	}
}

// Helper function for error rendering
func renderErrorTemplate(c *gin.Context, message string) {
	templateData := templates.CallbackData{
		Success: false,
		Title:   "Submission Error",
		Error:   message,
	}
	c.Writer.Header().Set("Content-Type", "text/html")
	if err := templates.Callback(templateData).Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Critical error rendering page")
	}
}