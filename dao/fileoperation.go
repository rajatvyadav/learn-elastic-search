package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"learnElasticSearch/model"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

//read json file from path given in parameters
// return the array of type EmployeeDetails and error
func ReadJsonFile(filePath string) ([]*model.EmployeeDetails, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var employeeDetails []*model.EmployeeDetails
	err = json.Unmarshal([]byte(byteValue), &employeeDetails)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return employeeDetails, nil
}

func InsertEmployeeDetailsHandler(es *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// bind
		var fileName model.FileDetails
		c.Bind(&fileName)
		fmt.Println("Filename", fileName.FileName)
		empDetailsArray, err := ReadJsonFile(fileName.FileName)
		if err != nil {
			fmt.Println("Error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
		for i, employee := range empDetailsArray {
			dataJSON, err := json.Marshal(employee)
			if err != nil {
				fmt.Println("Error", err)
				continue
			}
			js := string(dataJSON)
			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "employees",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(js),
				Refresh:    "true",
			}
			// Perform the request with the client.
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			timeout := 5 * time.Millisecond
			if err == nil {
				// The request has a timeout, so create a context that is
				// canceled automatically when the timeout expires.
				ctx, cancel = context.WithTimeout(context.Background(), timeout)
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}
			defer cancel() // Cancel ctx as soon as handleSearch returns.
			res, err := req.Do(ctx, es)
			if err != nil {
				fmt.Printf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				fmt.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"Info": "Inserted!!!",
		})
	}
}

// Search for document containing firstName or lastname matches to key in parameters
// return the entire document
func SearchEmployeeHandler(es *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var field model.Field
		err := c.Bind(&field)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"error": err,
			})
		}
		// Build the request body.
		query := `{
			"query": {
			  "bool": {
				"should": [
				  {
					"multi_match": {
					  "query": "` + field.Key + `",
					  "type": "cross_fields",
					  "fields": [
						"first_name",
						"last_name"
					  ],
					  "minimum_should_match": "50%"
					}
				  },
				  {
					"multi_match": {
					  "query": "` + field.Key + `",
					  "type": "cross_fields",
					  "fields": [
						"*.edge"
					  ]
					}
				  }
				]
			  }
			}
		  }`
		// Perform the search request.
		res, err := es.Search(
			es.Search.WithIndex("employees"),
			es.Search.WithBody(strings.NewReader(query)),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
		if err != nil {
			fmt.Printf("Error getting response: %s", err)
			if err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"error": err,
				})
			}
		}
		defer res.Body.Close()
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			fmt.Printf("Error parsing the response body: %s", err)
			if err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"error": err,
				})
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"details": r,
		})
	}
}

// Search for document containing firstName or lastname matches or having prefix to key in parameters
// return the entire document
func SearchEmployeeHavingPrefixHandler(es *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var field model.Field
		err := c.Bind(&field)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"error": err,
			})
		}
		// Build the request body.
		query := `{
			"query": {
			  "match_phrase_prefix": {
				"first_name": {
				  "query": "` + field.Key + `"
				}
			  }
			}
		  }`
		// Perform the search request.
		res, err := es.Search(
			es.Search.WithIndex("employees"),
			es.Search.WithBody(strings.NewReader(query)),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
		if err != nil {
			fmt.Printf("Error getting response: %s", err)
			if err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"error": err,
				})
			}
		}
		defer res.Body.Close()
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			fmt.Printf("Error parsing the response body: %s", err)
			if err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"error": err,
				})
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"details": r,
		})
	}
}
