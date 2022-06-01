package vmshell

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
)

const HOST = "https://vmshell.com"
const LOGIN_API = HOST + "/index.php?rp=/login"
const SERVER_API = HOST + "/modules/servers/solusvmplus/get_client_data.php"

type vmShellClient struct {
	lock     *sync.Mutex
	username string
	password string
	client   *http.Client
	logged   bool
}

func (v *vmShellClient) GetServerInfo(serverId string, retry bool) (*ServerInfo, error) {
	if v.client == nil {
		panic("use FromConfig to create Client")
	}
	// 有几种情况是需要 getToken
	// 1. 第一次登录
	url, _ := url.Parse(SERVER_API)
	if len(v.client.Jar.Cookies(url)) == 0 {
		v.Login()
	}

	query := url.Query()
	query.Add("vserverid", serverId)
	query.Add("_", fmt.Sprintf("%d", time.Now().UnixMilli()))
	url.RawQuery = query.Encode()

	res, err := v.client.Get(url.String())

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		// cookies out of Date
		if !retry {
			log.Info("no retry on login, skip")
			return nil, fmt.Errorf("server unavailable")
		}
		log.Info("maybe need to login")
		// 2. 可能需要重新LogIn
		v.Login()
		v.GetServerInfo(serverId, false)
	}

	si := &ServerInfo{}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, si)
	if err != nil {
		// cookies out of date
		if !retry {
			log.Info("no retry on login, skip")
			return nil, fmt.Errorf("server unavailable")
		}
		log.Info("maybe need to login")
		// 2. 可能需要重新LogIn
		v.Login()
		v.GetServerInfo(serverId, false)
	}
	return si, nil
}

func (v *vmShellClient) Login() {
	if v.client == nil {
		panic("use FromConfig to create Client")
	}
	if !v.lock.TryLock() {
		// 没获取锁的可以等待其他线程完成Login
		times := 0
		for {
			time.Sleep(5 * time.Second)
			if times > 5 {
				log.Error("a lock is never been release!")
				return
			}
			if v.logged {
				return
			} else {
				times += 1
			}
		}
	}

	// clear cookies
	host, _ := url.Parse(HOST)
	v.client.Jar.SetCookies(host, v.client.Jar.Cookies(host)[:0])

	v.logged = false
	var err error
	var token string = ""
	var formData url.Values = url.Values{}
	var resp *http.Response

	token = v.GetCSRFToken()
	if token == "" {
		goto errHandler
	}

	formData.Add("username", v.username)
	formData.Add("password", v.password)
	formData.Add("token", token)

	// login
	resp, err = v.client.PostForm(LOGIN_API, formData)
	if err != nil {
		goto errHandler
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("login code is not 200, check you network")
		goto errHandler
	}

	v.client.Jar.SetCookies(resp.Request.URL, resp.Cookies())

	v.logged = true
	v.lock.Unlock()
	return

errHandler:
	if err != nil {
		log.Error(err)
	}
	log.Error("Get Token failed!")
	v.lock.Unlock()
	v.logged = true
}

// 获取 CSRF Token
func (v *vmShellClient) GetCSRFToken() string {
	if v.client == nil {
		panic("use FromConfig to create Client")
	}
	res, err := v.client.Get(LOGIN_API) // Store cookies
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	if res.StatusCode != 200 {
		log.Error("get csrf token failed!")
		return ""
	}

	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	doc = htmlquery.FindOne(doc, "//input[@name='token']")
	if doc == nil {
		log.Error("get csrf token from html failed!")
		return ""
	}
	for _, a := range doc.Attr {
		if a.Key == "value" {
			return a.Val
		}
	}
	log.Error("no csrf attr found!")
	return ""
}

func newClient(username, password string) *vmShellClient {
	jar, _ := cookiejar.New(nil)
	return &vmShellClient{username: username, password: password, lock: &sync.Mutex{}, client: &http.Client{Timeout: 30 * time.Second, Jar: jar}}
}
