package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Subscriber struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type Meta struct {
	Page       int `json:"page"`
	Count      int `json:"count"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}

type DripClient struct {
	AccountID string
	APIKey    string
	UserAgent string
}

func (dc DripClient) baseURL() string {
	return fmt.Sprintf("https://api.getdrip.com/v2/%s/", dc.AccountID)
}

func (dc DripClient) subscribersURL(page int) string {
	values := make(url.Values, 0)
	values.Add("per_page", "1000")
	values.Add("status", "all")
	values.Add("page", fmt.Sprintf("%d", page))
	return fmt.Sprintf("%s%s", dc.baseURL(), "subscribers?"+values.Encode())
}

func (dc DripClient) unsubscribeURL() string {
	return fmt.Sprintf("%s%s", dc.baseURL(), "unsubscribes/batches")
}

func (dc DripClient) deleteSubURL(id string) string {
	return fmt.Sprintf("%s%s%s", dc.baseURL(), "subscribers/", id)
}

func (dc DripClient) request(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if dc.UserAgent != "" {
		req.Header.Set("User-Agent", dc.UserAgent)
	}
	req.SetBasicAuth(dc.APIKey, "")
	return req, nil
}

func (dc DripClient) Subscribers(page int) ([]Subscriber, *Meta, error) {
	var resp struct {
		Subscribers []Subscriber `json:"subscribers"`
		Meta        Meta         `json:"meta"`
	}
	req, err := dc.request(http.MethodGet, dc.subscribersURL(page), nil)
	if err != nil {
		return nil, nil, err
	}
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&resp)
	if err != nil {
		return nil, nil, err
	}
	return resp.Subscribers, &resp.Meta, nil
}

func (dc DripClient) DeleteSub(s Subscriber) error {
	req, err := dc.request(http.MethodDelete, dc.deleteSubURL(s.ID), nil)
	if err != nil {
		return err
	}
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return fmt.Errorf("drip: bad status code in response - %d", res.StatusCode)
	}
	return nil
}

func (dc DripClient) Unsubscribe(subscribers []Subscriber) error {
	type batch struct {
		Subscribers []Subscriber `json:"subscribers"`
	}
	body := struct {
		Batches []batch `json:"batches"`
	}{
		Batches: []batch{
			{Subscribers: subscribers},
		},
	}
	// 204 is success response

	pr, pw := io.Pipe()
	enc := json.NewEncoder(pw)
	go func() {
		enc.Encode(body)
		pw.Close()
	}()
	req, err := dc.request(http.MethodPost, dc.unsubscribeURL(), pr)
	if err != nil {
		return err
	}
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return fmt.Errorf("drip: bad status code in response - %d", res.StatusCode)
	}
	return nil
}
