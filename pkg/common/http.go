package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//
func SendGetRequest(req *http.Request) (string, error) {

	//
	client := &http.Client{Timeout: timeoutSecond}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := string(respBytes)

	return result, nil
}

//
func PrintResponse(rs *http.Response) {
	fmt.Println("rs.Status:", rs.Status)
	fmt.Println("rs.StatusCode:", rs.StatusCode)
	fmt.Println("rs.TransferEncoding:", rs.TransferEncoding)
	fmt.Println("rs.Proto:", rs.Proto)
	fmt.Println("rs.ProtoMajor:", rs.ProtoMajor)
	fmt.Println("rs.ProtoMinor:", rs.ProtoMinor)
	fmt.Println("rs.Request:", rs.Request)
	fmt.Println("rs.ContentLength:", rs.ContentLength)
	fmt.Println("rs.Header:", rs.Header)
	fmt.Println("rs.Cookie:", rs.Cookies())
	fmt.Println("rs.Close:", rs.Close)
	PrintResponseBody(rs)
}

//
func PrintResponseBody(rs *http.Response) {
	respBytes, err := ioutil.ReadAll(rs.Body)
	if err == nil {
		result := string(respBytes)
		fmt.Println(result)
	}
}
