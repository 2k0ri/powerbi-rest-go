package powerbi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetToken(tenantId, clientId, clientSecret string) string {
	values, err := url.ParseQuery(`grant_type=client_credentials&client_id=` + clientId + `&client_secret=` + clientSecret + `&resource=https://analysis.windows.net/powerbi/api`)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.PostForm("https://login.microsoftonline.com/"+tenantId+"/oauth2/token", values)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Print(res.Status)
		r, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(r))
	}
	var f interface{}
	json.NewDecoder(res.Body).Decode(&f)
	return f.(map[string]interface{})["access_token"].(string)
}

func GetTokenUserPassword(tenantId, clientId, user, password string) string {
	values, err := url.ParseQuery(`grant_type=password&client_id=` + clientId + `&resource=https://analysis.windows.net/powerbi/api` + `&username=` + user + `&password=` + password)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.PostForm("https://login.microsoftonline.com/"+tenantId+"/oauth2/token", values)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Print(res.Status)
		r, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(r))
	}
	var f interface{}
	json.NewDecoder(res.Body).Decode(&f)
	return f.(map[string]interface{})["access_token"].(string)
}
