package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AddProduct struct {
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

func AddProductToWishList(wishlistId, productIds string, user User) (addProduct AddProduct) {
	// 加入收藏 (POST http://www.wish.com/api/user/wishlist/add-product)

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

	params.Set("wishlist_id", wishlistId)
	params.Set("advertiser_id", user.AdvertiserId)
	params.Set("app_device_id", user.AppDeviceID)
	params.Set("product_ids[]", productIds)

	params.Set("_app_type", "wish")
	params.Set("_version", "3.21.0")
	params.Set("_client", "iosapp")
	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone9,2")

	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://www.wish.com/api/user/wishlist/add-product", body)

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

	json.Unmarshal(respBody, &addProduct)

	return

}
