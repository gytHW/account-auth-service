package common

import (
	"encoding/json"
	"github.com/donnie4w/go-logger/logger"
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
)

func ApiPost(host string, data map[string]interface{}) ( interface{}, error)  {

	data["AuthKey"] = "cc5538a413147bce"
	logger.Debug("InterApiPost, params: ", data)
	fmt.Println(data)

	b, err := json.Marshal(&data)
	if err != nil {
		logger.Debug("json err:", err)
	}
	body := bytes.NewBuffer([]byte(b))


	res,err := http.Post(host, "application/json", body)
	if err != nil {
		logger.Debug(err)
		return nil, err
	}

	result, err1 := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var resp interface{}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		logger.Debug(" INNER_API_RESPONSE json error")
		return nil, err
	}


	logger.Debug(" INNER_API_RESPONSE | (%s.%s) \n    %+v", data["Backend"], data["Action"], resp)

	return resp, err1
}
