package powerbi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type PowerBIClient struct {
	Token   string
	GroupId string
}

func (c *PowerBIClient) CreateDataSet(d DataSet, disableRetention bool) (string, error) {
	j, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	q := "datasets"
	if disableRetention {
		q += "?dafeaultRetentionPolicy=None"
	}
	r := c.request("POST", "datasets?defaultRetentionPolicy=None", bytes.NewReader(j))
	defer r.Body.Close()
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *PowerBIClient) GetDataSets() []DataSet {
	r := c.request("GET", "datasets", nil)
	defer r.Body.Close()
	fmt.Println(ioutil.ReadAll(r.Body))
	return nil // @TODO
}

func (c *PowerBIClient) AddRows(dataSetId, tableName string, rows Rows) (string, error) {
	b, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	r := c.request("POST", "datasets/"+dataSetId+"/tables/"+tableName+"/rows", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *PowerBIClient) GetGroups() []DataSet {
	r := c.request("GET", "groups", nil)
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	return nil // @TODO
}

func (c *PowerBIClient) Authenticate(tenantId, clientId, clientSecret string) {
	c.Token = GetToken(tenantId, clientId, clientSecret)
}

func (c *PowerBIClient) AuthenticateUserPassword(tenantId, clientId, userName, password string) {
	c.Token = GetTokenUserPassword(tenantId, clientId, userName, password)
}

func (c *PowerBIClient) request(method, path string, body io.Reader) *http.Response {
	url := "https://api.powerbi.com/v1.0/myorg/"
	if c.GroupId != "" {
		url += "groups/" + c.GroupId + "/"
	}
	url += path
	h, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	err = h.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	if method == "POST" {
		h.Header.Set("Content-Type", "application/json")
	}
	h.Header.Set("Authorization", "Bearer "+c.Token)
	log.Print(h)
	r, err := http.DefaultClient.Do(h)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode >= 400 {
		log.Fatal(r.Status)
	}
	return r
}
