package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

func FetchUrls() {
	// Fetch data from the /urls endpoint
	urlsEndpoint := "http://localhost:8080/urls"
	response, err := http.Get(urlsEndpoint)
	if err != nil {
		fmt.Println("Error fetching data from the /urls endpoint:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal JSON data
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Create NGINX configuration file
	err = CreateNginxConfig(data)
	if err != nil {
		fmt.Println("Error creating NGINX configuration file:", err)
		return
	}

	fmt.Println("NGINX configuration file created successfully.")
}

// createNginxConfig generates the NGINX configuration file based on the provided data
func CreateNginxConfig(data map[string]interface{}) error {
	// Open or create NGINX configuration file
	file, err := os.Create("/etc/dynamic_redirects")
	if err != nil {
		return err
	}
	defer file.Close()

	// Define NGINX configuration template
	nginxConfigTemplate := `
server {
    listen 80;
    server_name example.com;

{{range .data}}
    location {{.CurrentURL}} {
        return 301 http://example.com{{.RedirectURL}};
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
		return err
	}

	// Execute the template and write to the NGINX configuration file
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
