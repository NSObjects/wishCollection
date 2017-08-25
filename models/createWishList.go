package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type CWishList struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Permalink string `json:"permalink"`
		UserID    string `json:"user_id"`
		Name      string `json:"name"`
		Bid       string `json:"bid"`
		Private   bool   `json:"private"`
		ID        string `json:"id"`
		Size      int    `json:"size"`
	} `json:"data"`
	SweeperUUID string `json:"sweeper_uuid"`
	NotiCount   int    `json:"noti_count"`
}

func CreateWishList(user User) (created CWishList) {
	// 创建收藏列表 (POST http://www.wish.com/api/user/wishlist/create)

	params := url.Values{}
	params.Set("_capabilities[]", "11")
	params.Set("_capabilities[]", "12")
	params.Set("_capabilities[]", "13")
	params.Set("_capabilities[]", "15")
	params.Set("_capabilities[]", "2")
	params.Set("_capabilities[]", "21")
	params.Set("_capabilities[]", "24")
	params.Set("_capabilities[]", "25")
	params.Set("_capabilities[]", "28")
	params.Set("_capabilities[]", "3")
	params.Set("_capabilities[]", "32")
	params.Set("_capabilities[]", "35")
	params.Set("_capabilities[]", "39")
	params.Set("_capabilities[]", "4")
	params.Set("_capabilities[]", "40")
	params.Set("_capabilities[]", "43")
	params.Set("_capabilities[]", "47")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "8")
	params.Set("_capabilities[]", "9")
	params.Set("_xsrf", "1")
	params.Set("_app_type", "wish")
	params.Set("_version", "3.21.0")

	params.Set("advertiser_id", user.AdvertiserId)
	params.Set("app_device_id", user.AppDeviceID)
	params.Set("app_device_model", "iPhone9,2")

	params.Set("_client", "iosapp")

	names := []string{"Wishlist#1", "Birthday", "Miscellaneous", "Accessories"}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	params.Set("name", names[r.Intn(len(names)-1)])

	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://www.wish.com/api/user/wishlist/create", body)

	// Headers
	cookie := fmt.Sprintf("_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN; sweeper_session=\"%s\"; bsid=%s", user.SweeperSession, user.Baid)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(respBody, &created)

	return

}
