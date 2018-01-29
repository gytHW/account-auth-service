package logic

import (
	"bytes"
	"encoding/base64"
	// "encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/emicklei/go-restful"
	"github.com/satori/go.uuid"
	"image/color"
	"image/png"
	"strings"
	"time"
)

type verticode struct {
	uuid string
	code string
}

func CreateVertificationCode(req *Req, resp *restful.Response) {
	//create uuid
	uuid, err := uuid.NewV4()
	if err != nil {
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	}
	fmt.Printf("UUIDv4: %s\n", uuid)

	cap := captcha.New()
	// Must set a font, Other settings have default values
	if err := cap.SetFont("comic.ttf"); err != nil {
		fmt.Println("ttf")
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	}
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	img, str := cap.Create(6, captcha.NUM)

	//convert img to base64
	var emptyBuff bytes.Buffer
	png.Encode(&emptyBuff, img)
	bscode := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	//connect redis to store code

	client := ConnectRedis()
	defer client.Close()
	err = client.Set(uuid.String(), str, time.Hour).Err()
	if err != nil {
		fmt.Println("=======")
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	} else {
		ret_data := map[string]interface{}{
			"Uuid":              uuid.String(),
			"VertificationCode": bscode,
		}
		// ret_data, _ := json.Marshal(ret_data_map)
		resp.WriteEntity(&Resp{0, "success", 1, ret_data})
		return
	}

}

func VerifyVertificationCode(uuid string, vertiCode string) (bool, string) {
	client := ConnectRedis()
	defer client.Close()
	got, err := client.Get(uuid).Result()
	if err != nil {
		return false, err.Error()
	}
	if !strings.EqualFold(got, vertiCode) {
		return false, "wrong vertification code"
	} else {
		return true, "sucess"
	}

}

// func VerifyVertificationCode(req *restful.Request, resp *restful.Response) {
// 	request_data := verticode{}
// 	err := mapstructure.WeakDecode(req.Data, &request_data)
// 	if req.Data != nil {
// 		err := mapstructure.WeakDecode(req.Data, &request_data)
// 		if err != nil {
// 			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
// 			return
// 		}
// 	}
// 	uuid = request_data.uuid
// 	vertiCode = request_data.code

// 	client := connectRedis()
// 	defer client.Close()
// 	got, err := client.Get(uuid).Result()
// 	if err != nil {
// 		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
// 		return
// 	}
// 	if !strings.EqualFold(got, code) {
// 		resp.WriteEntity(&Resp{1, "wrong vertification code", 0, nil})
// 		return
// 	} else {
// 		resp.WriteEntity(&Resp{0, "sucess", 0, nil})
// 		return
// 	}

// }
