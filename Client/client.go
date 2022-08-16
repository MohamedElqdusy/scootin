package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"scootin/models"
	"time"
)

// Client is meant to be a standalone package to be used by others
// to connect to the service
//without any dependencies on the service itself.

type (
	// Client connects to the service using its url
	Client struct {
		baseUrl string
	}

	CheckoutCreate struct {
		MerchantCode string
		Amount       int64
		Currency     string
	}

	Checkout struct {
		ID           string
		MerchantCode string
		Amount       int64
		Currency     string
		CreatedAt    time.Time
	}
)

// NewClient take the service base url, returns a new service's client
func NewClient(url string) *Client {
	return &Client{baseUrl: url}
}

// CreateUser creates a user, returns the user uuid.
func (c *Client) CreateUser(user *models.User) (*models.UUIDResponse, error) {
	var (
		j, body []byte
		resp    *http.Response
		err     error
		uuid    *models.UUIDResponse
	)

	url := fmt.Sprintf("%s%s", c.baseUrl, "/v0.1/user")
	if j, err = json.Marshal(user); err != nil {
		return nil, err
	}
	if resp, err = http.Post(url, "application/json", bytes.NewBuffer(j)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, nil
	}

	if err = json.Unmarshal(body, &uuid); err != nil {
		return nil, err
	}
	return uuid, nil
}

// CreateScooter creates a scooter, returns the scooter uuid.
func (c *Client) CreateScooter() (*models.UUIDResponse, error) {
	var (
		j, body []byte
		resp    *http.Response
		err     error
		uuid    *models.UUIDResponse
	)

	url := fmt.Sprintf("%s%s", c.baseUrl, "/v0.1/scooter")
	if j, err = json.Marshal(""); err != nil {
		return nil, err
	}
	if resp, err = http.Post(url, "application/json", bytes.NewBuffer(j)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, nil
	}

	if err = json.Unmarshal(body, &uuid); err != nil {
		return nil, err
	}
	return uuid, nil
}

// BookScooter books a scooter for a user.
func (c *Client) BookScooter(scooterID, userID string) error {
	var err error
	url := fmt.Sprintf("%s%s%s", c.baseUrl, "/v0.1/scooter/book/", scooterID)
	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	// set the request header
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("user-id", userID)

	// initialize http client
	h := &http.Client{}

	// execute the request
	if _, err = h.Do(req); err != nil {
		return err
	}
	return nil
}

// ReleaseScooter releases a scooter by a user.
func (c *Client) ReleaseScooter(userID string) error {
	var err error
	url := fmt.Sprintf("%s%s", c.baseUrl, "/v0.1/scooter/release/")
	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	// set the request header
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("user-id", userID)

	// initialize http client
	h := &http.Client{}

	// execute the request
	if _, err = h.Do(req); err != nil {
		return err
	}
	return nil
}

// ListAvailableScooter returns the available scooters to ride
func (c *Client) ListAvailableScooter() ([]models.ScooterInfo, error) {
	var (
		body []byte
		resp *http.Response
		err  error
		sco  []models.ScooterInfo
	)

	url := fmt.Sprintf("%s%s", c.baseUrl, "/v0.1/scooters")
	if resp, err = http.Get(url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, nil
	}

	if err = json.Unmarshal(body, &sco); err != nil {
		return sco, err
	}
	return nil, nil
}
