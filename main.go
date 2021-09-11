package main

import (
	"fmt"
	dao "learnElasticSearch/dao"
	datasource "learnElasticSearch/dataSource"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World!!!")
	// initialize gin
	r := gin.Default()

	esClient, err := datasource.GetElasticSearch()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	res, err := esClient.Info()
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}

	defer res.Body.Close()
	if res.IsError() {
		fmt.Println(res.String())
	}
	// routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/bulkImportEmployeeDetails", dao.InsertEmployeeDetailsHandler(esClient))
	r.POST("/searchEmployee", dao.SearchEmployeeHandler(esClient))
	r.POST("/searchEmployeeWithPrefix", dao.SearchEmployeeHavingPrefixHandler(esClient))

	r.Run(":4700")
}
