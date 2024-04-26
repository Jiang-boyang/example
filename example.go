package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CODE_LANGUAGE int32

const (
	PHP CODE_LANGUAGE = iota + 1
	Go
	Java
	NodeJS
	Python
	Other
)

type Application struct {
	Id              uint           `json:"id" gorm:"primarykey"`
	ApplicationName string         `json:"application_name" gorm:"column:application_name"`
	Description     string         `json:"description" gorm:"column:description"`
	GitAddress      string         `json:"git_address" gorm:"column:git_address"`
	CodeLanguage    CODE_LANGUAGE  `json:"code_language" gorm:"column:code_language"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (Application) TableName() string {
	return "application"
}

func main() {
	dsn := "devops:xxxxxx@tcp(rm-xxxxxx.mysql.rds.aliyuncs.com:3306)/devops?charset=utf8&parseTime=True&loc=Local&timeout=1000ms"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	apps := []Application{}
	db.Find(&apps)

	for _, app := range apps {
		if app.CodeLanguage == NodeJS {
			continue
		}
		fmt.Println(app.ApplicationName)

		var response interface{}
		b, _ := json.Marshal(map[string]interface{}{
			"envId": 7,
			"appId": app.Id,
		})
		resp := PublicHttpRequest(
			"PUT",
			"http://ops.xxxxx.com/api/v1/cicd/config/map",
			b,
			map[string]string{"Content-Type": "application/json",
				"Authorization": "Bearer gmMe_gdjn1pcS3eycgH_3eKtOFDqEmcPEdzEpPEOkKE",
			},
		)

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response)
	}
}

func PublicHttpRequest(method, url string, values []byte, header map[string]string) *http.Response {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(values))) // URL-encoded payload
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for headKey, headValue := range header {
		req.Header.Add(headKey, headValue)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	return resp
}
