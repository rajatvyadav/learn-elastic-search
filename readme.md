 # Go - Gin - Elasticsearch

This is the sample client project for demostration of Elasticsearch using go-gin
You can also check my article - https://rajatvyadav.medium.com/elasticsearch-in-go-65388f8dbfad

## Running the application

**Prerequisite:**
 - Before installing Elasticsearch make sure Java is installed on your local machine.
 - Download the Elasticsearch and install link - https://www.elastic.co/downloads/elasticsearch.

_**Now it's time to run...**_
 - Extact the zip and cd bin and run elasticsearch.bat. it will take time setup elasticsearch.
 - Official website for Elasticsearch - https://github.com/elastic/go-elasticsearch

**Development**

    go mod init learn-elastic-search

    go get rajatvyadav/learn-elastic-search
 
    go run main.go

    check server is running or not
    GET http://localhost:4700/ping

    Insert multiple records in DB
    POST http://localhost:4700/bulkImportEmployeeDetails
    requestBody: Json
    body:
    {
        "fileName": "MOCK_DATA.json"
    }

    Search for record having first_name or last_name 
    POST http://localhost:4700/searchEmployee
    requestBody: Json
    body:
    {
        "key": "Francesca"
    }

_**Happy Coding...**_
