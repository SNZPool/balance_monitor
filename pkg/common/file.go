package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitly/go-simplejson"
)

//
func ReadFile(configPath string) ([]byte, error) {
	//
	flag, err := PathExists(configPath)
	if err != nil {
		fmt.Printf("PathExists(%s),err(%v)\n", configPath, err)
	}
	if !flag {
		tip := fmt.Sprintf("%s is not exist\n", configPath)
		err = errors.New(tip)
		return nil, err
	}

	//
	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	fileData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	return fileData, nil
}

//
func ReadJsonToSimpleJson(configPath string) (*simplejson.Json, error) {

	fileData, err := ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
	}
	res, err := simplejson.NewJson([]byte(fileData))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return res, err
}

//
func PrintSimpleJson(json *simplejson.Json) (string, error) {
	if json == nil {
		return "", fmt.Errorf("json is empty")
	}

	output := ""
	for k, v := range json.MustMap() {
		output = output + string(k) + ": "
		switch v.(type) {
		case string:
			output = output + v.(string)
			break
		case int:
			output = output + fmt.Sprintf("%d", v)
			break
		case float64:
			output = output + fmt.Sprintf("%f", v)
			break
		default:
			output = output + fmt.Sprint(v)
			break
		}
		output = output + "\n"
	}

	return output, nil
}

//
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
