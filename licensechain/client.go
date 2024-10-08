
package licensechain

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type Client struct {
    BaseURL    string
    HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL:    baseURL,
        HTTPClient: &http.Client{},
    }
}

func (c *Client) ValidateLicense(licenseKey string) (bool, error) {
    url := c.BaseURL + "/api/license/validate"
    requestBody, err := json.Marshal(map[string]string{"licenseKey": licenseKey})
    if err != nil {
        return false, err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return false, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    var result struct {
        Valid bool `json:"valid"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false, err
    }

    return result.Valid, nil
}
