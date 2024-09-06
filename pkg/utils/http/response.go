package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/errors"
)

type Response interface {
	ParseErrorFromHTTPResponse(body []byte) error
}

type BaseResponse struct {
}

type ErrorResponse struct {
	Response struct {
		Error struct {
			Code   string `json:"Code"`
			TaskId string `json:"TaskId"`
		} `json:"Error" omitempty`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

type DeprecatedAPIErrorResponse struct {
	Code    string `json:"Code"`
	Message string `json:"message"`
	TaskId  string `json:"TaskId"`
}

func (r *BaseResponse) ParseErrorFromHTTPResponse(body []byte) (err error) {
	resp := &ErrorResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors.NewCdsSDKError("ClientError.ParseJsonError", msg, "")
	}
	//if resp.Response.Error.Code != "" {
	//	return errors.NewCdsSDKError(resp.Response.Error.Code, resp.Response.Error.Message, resp.Response.RequestId)
	//}

	deprecated := &DeprecatedAPIErrorResponse{}
	err = json.Unmarshal(body, deprecated)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors.NewCdsSDKError("ClientError.ParseJsonError", msg, "")
	}
	if deprecated.Code != "Success" {
		return errors.NewCdsSDKError(deprecated.Code, deprecated.Message, "")
	}
	return nil
}

func ParseFromHttpResponse(hr *http.Response, response Response, r Request) (err error) {
	defer hr.Body.Close()
	body, err := ioutil.ReadAll(hr.Body)
	reqParams, _ := json.Marshal(r)
	if err != nil {
		msg := fmt.Sprintf("Fail to read response body because %s", err)
		return errors.NewCdsSDKError(fmt.Sprintf("ClientError.IOError Action:%s  msg:", r.GetAction()), msg, "")
	}
	if hr.StatusCode != 200 {
		if len(body) != 0 {
			err = json.Unmarshal(body, &response)
			if err != nil {
				msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
				return errors.NewCdsSDKError(fmt.Sprintf("ClientError.ParseJsonError Action: %s  msg:", r.GetAction()), msg, "")
			}
		}
		fmt.Println(r.GetAction())
		msg := fmt.Sprintf("Request fail with http status code: %s,request params:%s, with body: %s", hr.Status, string(reqParams), body)
		return errors.NewCdsSDKError(fmt.Sprintf("ClientError.HttpStatusCodeError Action:%s   msg:", r.GetAction()), msg, "")
	}
	//log.Printf("[DEBUG] Response Body=%s", body)
	err = response.ParseErrorFromHTTPResponse(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors.NewCdsSDKError(fmt.Sprintf("ClientError.ParseJsonError Action:%s msg:", r.GetAction()), msg, "")
	}

	log.Printf("action：%s  request body: %s response body: %s", r.GetAction(), string(reqParams), string(body))
	return
}
