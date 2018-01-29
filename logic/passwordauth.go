package logic

import (
	"account-auth-service/model"
	"crypto/md5"
	"encoding/hex"
	// "encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/mitchellh/mapstructure"
	"math/rand"
	"strconv"
	"time"
)

type PwdReq struct {
	Password    string `json:"Password"`
	Uuid        string `json:"Uuid"`
	VertifyCode string `json:"VertificationCode"`
}

type ChangePwdReq struct {
	OldPassword string `json:"OldPassword"`
	NewPassword string `json:"NewPassword"`
}

func CreatePassword(req *Req, resp *restful.Response) {
	request_data := PwdReq{}
	err := mapstructure.WeakDecode(req.Data, &request_data)
	if req.Data != nil {
		err := mapstructure.WeakDecode(req.Data, &request_data)
		if err != nil {
			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
			return
		}
	}
	account_id := req.Account_id
	password := request_data.Password
	save_password := encrypt(password, account_id)

	db, _ := ConnectMysql()
	defer db.Close()

	//check if key exists
	fmt.Print("accountid:")
	fmt.Print(account_id)
	var user_auth model.User_auth
	db.First(&user_auth, account_id)
	if user_auth.Id == 0 {
		user_auth := model.User_auth{Id: account_id, Password: save_password, Is_initial: 0}
		err = db.Create(&user_auth).Error
		if err != nil {
			fmt.Println("db close")
			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
			return
		}
		resp.WriteEntity(&Resp{0, "success", 0, nil})
		return
	} else {
		resp.WriteEntity(&Resp{1, "user exists", 0, nil})
		return
	}
}

func ValidatePassword(req *Req, resp *restful.Response) {
	request_data := PwdReq{}
	err := mapstructure.WeakDecode(req.Data, &request_data)
	if req.Data != nil {
		err := mapstructure.WeakDecode(req.Data, &request_data)
		if err != nil {
			fmt.Println("decode error")
			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
			return
		}
	}

	account_id := req.Account_id
	uuid := request_data.Uuid
	input_password := request_data.Password
	en_input_password := encrypt(input_password, account_id)

	//verify vertification code
	if uuid != "" {
		verify_code := request_data.VertifyCode
		result, msg := VerifyVertificationCode(uuid, verify_code)
		if !result {
			resp.WriteEntity(&Resp{1, msg, 0, nil})
			return
		}
	}

	//select password from db
	db, err := ConnectMysql()
	if err != nil {
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	}
	defer db.Close()

	user_auth := model.User_auth{Id: account_id}
	db.First(&user_auth)

	saved_password := user_auth.Password
	if en_input_password != saved_password {
		resp.WriteEntity(&Resp{1, "wrong password", 0, nil})
		return
	}
	resp.WriteEntity(&Resp{0, "success", 0, nil})
	return

}

func CreateInitialPassword(req *Req, resp *restful.Response) {
	account_id := req.Account_id
	initial_password := createRandomString()
	fmt.Print("initial_password:")
	fmt.Println(initial_password)
	save_password := encrypt(initial_password, account_id)
	db, _ := ConnectMysql()
	defer db.Close()

	user_auth := model.User_auth{Id: account_id, Password: save_password, Is_initial: 1}
	if db.Where("Id = ?", account_id).First(&user_auth).RecordNotFound() {
		err := db.Create(&user_auth).Error
		if err != nil {
			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
			return
		}
		ret_data := map[string]string{
			"initial_password": initial_password,
		}
		resp.WriteEntity(&Resp{0, "success", 0, ret_data})
		return
	} else {
		resp.WriteEntity(&Resp{1, "user exists", 0, nil})
		return
	}
}

func ChangePassword(req *Req, resp *restful.Response) {
	account_id := req.Account_id

	request_data := ChangePwdReq{}
	err := mapstructure.WeakDecode(req.Data, &request_data)
	if req.Data != nil {
		err := mapstructure.WeakDecode(req.Data, &request_data)
		if err != nil {
			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
			return
		}
	}
	old_password := request_data.OldPassword
	new_password := request_data.NewPassword
	en_old_password := encrypt(old_password, account_id)
	en_new_password := encrypt(new_password, account_id)
	//select password from db
	db, err := ConnectMysql()
	if err != nil {
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	}
	defer db.Close()

	user_auth := model.User_auth{Id: account_id}
	db.First(&user_auth)
	if en_old_password != user_auth.Password {
		resp.WriteEntity(&Resp{1, "wrong old password", 0, nil})
		return
	}
	err = db.Model(&user_auth).Updates(map[string]interface{}{"password": en_new_password, "is_initial": 0}).Error
	if err != nil {
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	} else {
		resp.WriteEntity(&Resp{0, "success", 0, nil})
		return
	}

}

func CheckInitialPassword(req *Req, resp *restful.Response) {
	account_id := req.Account_id
	//select password from db
	db, err := ConnectMysql()
	if err != nil {
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	}
	defer db.Close()

	user_auth := model.User_auth{Id: account_id}
	db.First(&user_auth)

	var status bool
	if user_auth.Is_initial == 1 {
		status = true
	}
	ret_data := map[string]bool{
		"is_initial_password": status,
	}
	resp.WriteEntity(&Resp{0, "success", 0, ret_data})
	return
}

func ResetPassword(req *Req, resp *restful.Response) {
	account_id := req.Account_id
	//select password from db
	db, err := ConnectMysql()
	if err != nil {
		resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
		return
	}
	defer db.Close()

	var user_auth model.User_auth
	if db.Where("Id = ?", account_id).First(&user_auth).RecordNotFound() {
		resp.WriteEntity(&Resp{0, "user not exist", 0, nil})
		return
	} else {
		initial_password := createRandomString()
		fmt.Print("initial_password:")
		fmt.Println(initial_password)
		save_password := encrypt(initial_password, account_id)
		err = db.Model(&user_auth).Updates(map[string]interface{}{"password": save_password, "is_initial": 1}).Error
		if err != nil {
			resp.WriteEntity(&Resp{1, err.Error(), 0, nil})
			return
		}
		ret_data := map[string]string{
			"action":           "ResetPassword",
			"initial_password": initial_password,
		}
		resp.WriteEntity(&Resp{0, "success", 0, ret_data})
		return
	}
}

//encrypt password
func encrypt(password string, account_id int) string {
	salt := strconv.Itoa(account_id)
	salted_pwd := password + salt
	data := []byte(salted_pwd)
	h := md5.New()
	h.Write(data)
	str := hex.EncodeToString(h.Sum(nil))

	return str
}

func createRandomString() string {
	rand.Seed(time.Now().Unix())
	var str string
	for i := 0; i < 6; i++ {
		num := rand.Intn(10)
		str += strconv.Itoa(num)
	}
	return str
}
