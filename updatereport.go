package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Account struct {
	Id         string
	Name       string
	Slug       string
	Active     bool
	Created_at string
	Owner_id   string
	Owner      string
}

type Deployment struct {
	Id                       string
	Name                     string
	Provider                 string
	Region                   string
	Type                     string
	Plan                     string
	Current_primary          string
	Members                  []string
	Ignored_members          []string
	Allow_multiple_databases bool
	Status                   string
	Databases                []Database
}

type Database struct {
	Id               string
	Name             string
	Deprovision_date string
	Status           string
	Deployment_id    string
	Plan             string
}

type Version struct {
	Type                     string
	Version                  string
	Messages                 []string
	Eligible_upgrade_version UpgradeVersion
	Upgrade_path             []UpgradeVersion
}

type UpgradeVersion struct {
	Upgrade_type string
	Version      string
	Messages     []string
}

var baseurl = "https://beta-api.mongohq.com"

var token = "<YOUR PERSONAL ACCESS TOKEN HERE>"

func main() {
	fmt.Println("Update report")

	var accounts []Account

	err := getAndUnmarshal(fmt.Sprintf("%s/accounts", baseurl), &accounts)
	if err != nil {
		panic(err.Error())
	}

	for _, account := range accounts {
		fmt.Printf("Account name:%s id:%s slug:%s\n\n", account.Name, account.Id, account.Slug)

		var deployments []Deployment

		err := getAndUnmarshal(fmt.Sprintf("%s/accounts/%s/deployments", baseurl, account.Slug), &deployments)
		if err != nil {
			panic(err.Error())
		}

		for _, deployment := range deployments {
			fmt.Printf("Deployment name/id: %s/%s\n", deployment.Name, deployment.Id)
			for _, database := range deployment.Databases {
				fmt.Printf("Database name/id: %s/%s is %s\n", database.Name, database.Id, database.Status)

				var version Version

				err := getAndUnmarshal(fmt.Sprintf("%s/deployments/%s/%s/version", baseurl, account.Slug, database.Deployment_id), &version)

				if err != nil {
					panic(err.Error())
				}

				fmt.Printf("The instance is %s version %s\n", version.Type, version.Version)

				for _, msg := range version.Messages {
					fmt.Println(msg)
				}
				fmt.Println()
			}
		}
	}
}

type ErrorMessage struct {
	Error string
}

func getAndUnmarshal(url string, result interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application-json")
	req.Header.Set("Accept-Version", "2014-06")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		var errorMessage ErrorMessage
		json.Unmarshal(body, &errorMessage)
		if errorMessage.Error != "" {
			return errors.New(errorMessage.Error)
		}
		return errors.New(res.Status)
	}

	err = json.Unmarshal(body, result)
	return err
}
