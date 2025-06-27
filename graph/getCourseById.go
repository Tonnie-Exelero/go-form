package graph

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"inquire/now-microservice/models"
)

func GetCourseByDocumentID(documentId string) models.Course {
	query := `
    query($documentId: ID!) {
        course(documentId: $documentId) {
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

	variables := map[string]string{"documentId": documentId}
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		log.Println("Failed to marshal request body:", err)
		return models.Course{}
	}

	log.Println("GraphQL request body:", string(requestBody))

	req, err := http.NewRequest("POST", os.Getenv("STRAPI_GRAPHQL"), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Failed to create request:", err)
		return models.Course{}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
		return models.Course{}
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
		return models.Course{}
	}

	log.Println("GraphQL response body:", string(body))

	var response struct {
		Data struct {
			Course models.Course `json:"course"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Failed to unmarshal response body:", err)
		return models.Course{}
	}

	return response.Data.Course
}
