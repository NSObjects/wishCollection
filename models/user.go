package models

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"errors"

	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wishCollection/utility"

	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

var (
	firstName []string
	lastName  []string
	emailType []string
	contrys   []string
)

func init() {
	initValue()
}

type LoginInfo struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		AlreadyHadApp        bool        `json:"already_had_app"`
		SessionToken         string      `json:"session_token"`
		NewUser              bool        `json:"new_user"`
		DeferredDeepLinkType interface{} `json:"deferred_deep_link_type"`
		SignupFlowType       string      `json:"signup_flow_type"`
		User                 string      `json:"user"`
		DeferredDeepLinkPid  interface{} `json:"deferred_deep_link_pid"`
	} `json:"data"`
	SweeperUUID string `json:"sweeper_uuid"`
	NotiCount   int    `json:"noti_count"`
}

type User struct {
	Id                    int64  `orm:"auto"`
	Baid                  string `json:"baid"`
	SweeperSession        string `json:"sweeper_session"`
	Email                 string
	Password              string
	RiskifiedSessionToken string `json:"riskified_session_token"`
	AdvertiserId          string `json:"advertiser_id"`
	AppDeviceID           string
	Country               string `json:"country"`
	FullName              string
	HasAddress            int
	Invalid               int //账号是否被封 , 0代表没有封。1代表账号被封
	UserId                string
	Updated               time.Time `orm:"auto_now;type(datetime)"`
}

func RegisterUser() (user User, err error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	f := firstName[r.Intn(len(firstName))]
	l := lastName[r.Intn(len(lastName))]
	e := fmt.Sprintf("%d%s", time.Now().Unix(), emailType[r.Intn(len(emailType))])
	if user, err = registIdWith(e, "1234567890", f, l); err == nil {
		if c := CreateWishList(user); c.Code != 0 {
			utility.SendLog(fmt.Sprint(c.Msg))
			time.Sleep(time.Minute * 5)
			return user, errors.New("创建失败")
		}

		return user, nil
	} else {
		return user, err
	}

	return user, errors.New("创建失败")
}

func registIdWith(email, password, firstName, lastName string) (User, error) {
	// 注册 (POST https://www.wish.com/api/email-signup)
	params := url.Values{}
	params.Set("_app_type", "wish")
	params.Set("_version", "3.20.0")
	params.Set("_client", "iosapp")

	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone9,2")
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
	params.Set("_capabilities[]", "47")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "8")
	params.Set("_capabilities[]", "9")

	AdvertiserID := strings.ToUpper(uuid.NewV4().String())
	params.Set("advertiser_id", AdvertiserID)

	riskifiedSessionToken := strings.ToUpper(uuid.NewV4().String())
	params.Set("_riskified_session_token", riskifiedSessionToken)

	key := []byte(strings.ToUpper(uuid.NewV4().String()))
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(time.Now().String()))
	appDeviceID := fmt.Sprintf("%x", mac.Sum(nil))
	params.Set("app_device_id", appDeviceID)

	params.Set("first_name", firstName)
	params.Set("last_name", lastName)
	params.Set("password", password)
	params.Set("email", email)
	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://www.wish.com/api/email-signup", body)

	// Headers
	req.Header.Add("Cookie", "_xsrf=1; _appLocale=zh-Hans-CN; _timezone=8")
	req.Header.Add("User-Agent", "Wish/3.20.0 (iPhone; iOS 10.3.1; Scale/3.00)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utility.Errorln(err)
		if e, ok := err.(*json.SyntaxError); ok {
			utility.Errorln(e)
		}
	}

	var loginInfo LoginInfo
	err = json.Unmarshal(respBody, &loginInfo)
	if err != nil {
		utility.Errorln(err)
		if e, ok := err.(*json.SyntaxError); ok {
			utility.Errorln(e)
		}
	}

	user := User{}

	if loginInfo.Code != 0 {
		fmt.Println(loginInfo, email)
		return user, errors.New("login error")
	}

	user.UserId = loginInfo.Data.User
	user.Email = email
	user.Password = password
	user.AppDeviceID = appDeviceID
	user.AdvertiserId = AdvertiserID
	user.RiskifiedSessionToken = riskifiedSessionToken
	user.FullName = firstName + " " + lastName
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "bsid":
			user.Baid = cookie.Value
		case "sweeper_session":
			user.SweeperSession = cookie.Value
		}

	}

	return user, nil

}

func initValue() {
	firstName = []string{
		"Aaron", "Abbott", "Abel", "Abner", "Abraham", "Adair", "Adam", "Addison",
		"Adolph", "Adonis", "Adrian", "Ahern", "Alan", "Albert", "Aldrich", "Alexander",
		"Alfred", "Alger", "Algernon", "Allen", "Alston", "Alva", "Alvin", "Alvis", "Amos", "Andre",
		"Andrew", "Andy", "Angelo", "Augus", "Ansel", "Antony", "Bevis", "Bill", "Bishop", "Blair", "Blake",
		"Bob", "Clarence", "Clark", "Claude", "Clyde", "Colin", "Dana", "Darnell", "Darcy", "Dempsey", "Dominic",
		"Edwiin", "Edward", "Elvis", "Fabian", "Frank", "Gale", "Gilbert", "Goddard", "Grover", "Hayden",
		"Hogan", "Hunter", "Isaac", "Ingram", "Isidore", "Jacob", "Jason", "Jay", "Jeff", "Jeremy", "Jesse",
		"Jerry", "Jim", "Jonathan", "Joseph", "Joshua", "Julian", "Julius", "Ken", "Kennedy", "Kent",
		"Kerr", "Kerwin", "Kevin", "Kirk", "King", "Lance", "Larry", "Leif", "Leonard", "Leopold", "Lewis",
		"Lionel", "Lucien", "Lyndon", "Magee", "Malcolm", "Mandel", "Marico", "Marsh", "Marvin", "Maximilian",
		"Meredith", "Merlin", "Mick", "Michell", "Monroe", "Montague", "Moore", "Mortimer", "Moses", "Nat",
		"Nathaniel", "Neil", "Nelson", "Newman", "Nicholas", "Nick", "Noah", "Noel", "Norton", "Ogden",
		"Oliver", "Omar", "Orville", "Osborn", "Oscar", "Osmond", "Oswald", "Otis", "Otto", "Owen", "Page", "Parker",
		"Paddy", "Patrick", "Paul", "Payne", "Perry", "Pete", "Peter", "Philip", "Phil",
		"Primo", "Quennel", "Quincy", "Quintion", "Rachel", "Ralap", "Randolph", "Robin", "Rodney", "Ron",
		"Roy", "Rupert", "Ryan", "Sampson", "Samuel", "Simon", "Stan", "Stanford", "Steward",
	}

	lastName = []string{
		"Baker", "Hunter", "Carter", "Smith", "Cook", "Turner", "Baker", "Miller", "Smith", "Turner", "Hall",
		"Hill", "Lake", "Field", "Green", "Wood", "Well", "Brown", "Longman", "Short", "White", "Sharp",
		"Hard", "Yonng", "Sterling", "Hand", "Bull", "Fox", "Hawk", "Bush", "Stock", "Cotton", "Reed",
		"George", "Henry", "David", "Clinton", "Macadam", "Abbot", "Abraham", "Acheson", "Ackerman", "Adam",
		"Addison", "Adela", "Adolph", "Agnes", "Albert", "Alcott", "Aldington", "Alerander", "Alick", "Amelia",
		"Adams",
	}

	emailType = []string{
		"@gmail.com", "@qq.com", "@126.com", "@163.com", "@vip.sina.com", "@sina.com", "@tom.com", "@263.com", "@189.com", "@outlook.com",
	}
	contrys = []string{
		"AU",
		"GB",
		"US",
		"FR",
		"DE",
		"CA",
		"HK",
		"VN",
		"SG",
		"MY",
		"JP",
		"KR",
		"IN",
		"ID",
		"GR",
		"BR",
		"FI",
		"AT",
		"ES",
		"RU",
		"NO",
		"SE",
		"NL",
		"CH",
		"DK",
		"IT",
	}

}
