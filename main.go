package main

import (
	"bytes"
	"compress/gzip"
	"doc/utility"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"wishCollection/models"
)

var stop bool

func main() {

	requestWishId()
	l := make(chan int)
	<-l

}

type CollectionJSON struct {
	Code int    `json:"code"`
	Id   string `json:"id"`
}

func requestWishId() {
	// 获取指定id (GET http://localhost:3384/api/collection)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "http://localhost:3384/api/collection", nil)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	var cJSON CollectionJSON
	if err := json.Unmarshal(respBody, &cJSON); err != nil {
		fmt.Println(err)
		requestWishId()
		return
	}

	if cJSON.Code != 0 {
		return
	}

ReRegister:
	if user, err := models.RegisterUser(); err == nil {
		getWishIdFromFeed("tabbed_feed_latest", user, cJSON.Id)
	} else {
		goto ReRegister
	}
}
//13672
func getWishIdFromFeed(categoryId string, user models.User, wishId string) {
	c := time.NewTicker(time.Minute * 2)
	go TimeOut(c)
	fmt.Println(wishId)
	page := 0
	for {
		if stop == true {
			stop = false
			if g := models.GetWisList(user); g.Code == 0 {
				if len(g.Data.Wishlists) > 0 {
					if a := models.AddProductToWishList(g.Data.Wishlists[0].ID, wishId, user); a.Code != 0 {
						utility.SendLog(a.Msg)
					}
				} else {
					utility.SendLog(fmt.Sprintln("创建收藏列表失败", wishId))
				}
			}
			requestWishId()
			return
		}

		if err := loadFeed(page, categoryId, user); err != nil {
			utility.SendLog(err.Error())
			continue
		}

		time.Sleep(time.Second * 10)
		page++
	}

}

func TimeOut(c *time.Ticker) {
	for now := range c.C {
		fmt.Println(now)
		stop = true
		c.Stop()
	}
}

func loadFeed(page int, categoryId string, user models.User) error {

	body := feedBodyWith(page, user, categoryId)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.wish.com/api/feed/get-filtered-feed", body)
	if err != nil {
		return err
	}

	req = headerWish(req, user)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			utility.Errorln(err)
		}

	default:
		reader = resp.Body
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)

	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d %s", resp.StatusCode, buf.Bytes())
	}

	var feeds Feeds

	if err = json.Unmarshal(buf.Bytes(), &feeds); err != nil {
		return err
	}

	if feeds.Code == 10 {
		return fmt.Errorf("not more product")
	} else {

		if len(feeds.Data.Products) <= 0 {
			return fmt.Errorf("not more product")
		}
	}

	return nil
}

func feedBodyWith(page int, user models.User, category string) *bytes.Buffer {

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
	params.Set("_capabilities[]", "32")
	params.Set("_capabilities[]", "35")
	params.Set("_capabilities[]", "39")
	params.Set("_capabilities[]", "4")
	params.Set("_capabilities[]", "40")
	params.Set("_capabilities[]", "43")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "8")
	params.Set("_capabilities[]", "9")

	params.Set("request_id", category)

	params.Set("_app_type", "wish")
	params.Set("_version", "3.20.6")
	params.Set("_client", "iosapp")
	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone9,2")

	params.Set("advertiser_id", user.AdvertiserId)
	params.Set("_riskified_session_token", user.RiskifiedSessionToken)
	params.Set("app_device_id", user.AppDeviceID)
	//params.Set("_threat_metrix_session_token", user.)

	params.Set("count", "30")
	params.Set("offset", fmt.Sprintf("%d", page*30))
	//params.Set("request_categories", "true")

	body := bytes.NewBufferString(params.Encode())
	return body
}

func headerWish(req *http.Request, user models.User) *http.Request {
	// Headers
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "zh-Hans-CN;q=1")
	cookie := fmt.Sprintf("_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN; sweeper_session=\"%s\"; bsid=%s", user.SweeperSession, user.Baid)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Wish/3.20.6 (iPhone; iOS 10.3.2; Scale/3.00)")

	return req
}

type Feeds struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Products []struct {
			ID string `json:"id"`
		} `json:"products"`
	} `json:"data"`
}
