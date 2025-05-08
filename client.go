package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Client struct {
	BaseUrl     string
	Client      *http.Client
	AccessToken string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewClient(consumerKey, consumerSecret string) *Client {
	credentials := consumerKey + ":" + consumerSecret
	credentials = base64.StdEncoding.EncodeToString([]byte(credentials))
	baseUrl := os.Getenv("KAONAVI_BASE_URL")
	url := baseUrl + "token"

	data := "grant_type=client_credentials"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Unmarshal the JSON response
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON response: %v", err)
	}

	return &Client{
		BaseUrl:     baseUrl,
		Client:      client,
		AccessToken: tokenResponse.AccessToken,
	}
}

func (c *Client) GetMembers() string {
	url := c.BaseUrl + "members"

	return c.get(url)
}

func (c *Client) GetDepartments() string {
	url := c.BaseUrl + "departments"

	return c.get(url)
}

func (c *Client) GetCustom(endpoint string) string {
	url := c.BaseUrl + endpoint

	return c.get(url)
}

func (c *Client) get(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Kaonavi-Token", c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return string(body)
}
