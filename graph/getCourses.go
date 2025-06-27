package graph

import (
	"bytes"
	"encoding/json"
	"inquire/now-microservice/models" // Adjust the import path as necessary
	"io"
	"log"
	"net/http"
	"os"
)

func GetCourses(tagFilter string) ([]models.Course, error) {
	query := `
    query( $tagFilter: String ) {
        courses(filters: { tags: { name: { containsi: $tagFilter } } }) {
            documentId
            title
            slug
            description
            instructor
            duration
			learningOutcomes
            startDate
            price {
                course_fee
            }
            mode {
                online
                handsOn
                selfPaced
                lecture
                classroom
            }
            categories {
                name
            }
            provider {
                companyName
                website
            }
            features {
                lifetimeAccess
                oneOnOneMentorship
                referralProgram
                liveSupport
            }
        }
    }`

	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"tagFilter": tagFilter,
		},
	})
	if err != nil {
		log.Println("Failed to marshal request body:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", os.Getenv("STRAPI_GRAPHQL"), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Failed to create request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
		return nil, err
	}
	//defer resp.Body.Close()
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Println("Failed to close response body:", cerr)
			if err == nil {
				err = cerr
			}
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil, err
	}

	var response struct {
		Data struct {
			Courses []models.Course `json:"courses"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Failed to unmarshal response body:", err)
		return nil, err
	}

	return response.Data.Courses, nil
}
