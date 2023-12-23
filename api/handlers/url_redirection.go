package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
)

func FetchUrls() (map[string]interface{}, error) {
	// Fetch data from the /urls endpoint
	endpoint := os.Getenv("MICROSERVICE_ENDPOINT")
	urlsEndpoint := endpoint + "/urls"
	response, err := http.Get(urlsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from the /urls endpoint: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal JSON data
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return data, nil
}

func CreateNginxConfig(data map[string]interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating NGINX configuration file: %v", err)
	}
	defer file.Close()

	// Define NGINX configuration template
	//Replace example.com with the actual domain or add it on env
	nginxConfigTemplate := `
	server {
		listen 80;
		server_name example.com;
	
	{{range .data}}
		location {{.currentUrl}} {
			return 301 https://example.com{{.redirectUrl}};
		}
	{{end}}
	
		location / {
			# Your default configuration for other requests
		}
	
		# Other server configurations...
	}
	`

	// Parse the template
	tmpl, err := template.New("nginxConfig").Parse(nginxConfigTemplate)
	if err != nil {
		return fmt.Errorf("error parsing NGINX configuration template: %v", err)
	}

	// Execute the template and write to the NGINX configuration file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template and writing to NGINX configuration file: %v", err)
	}

	return nil
}

func GenerateNginxConfig(c *gin.Context) {
	data, err := FetchUrls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//should replace the path with the actual nginx
	//file path on the server now this saves it inside this project
	err = CreateNginxConfig(data, "etc/nginx.conf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NGINX configuration file generated successfully"})
}
