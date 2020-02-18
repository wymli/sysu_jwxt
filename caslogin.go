package main

import (
	// "encoding/json"
	// "encoding/json"

	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	// "image"
	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	// "os"

	// "os"
	// "os/exec"
	// "runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//
// var commands = map[string]string{
// 	"windows": "start",
// 	"darwin":  "open",
// 	"linux":   "xdg-open",
// }

//Login in and return *client
func casLogin() (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		// @ 禁止重定向 , 可以通过  len(via)  控制重定向次数
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 { //狗逼系统刚好重定向10次,go默认重定向10次
				return http.ErrUseLastResponse
			} else {
				fmt.Println("[Redirecting] via ", req.URL.Path, "method:", req.Method)
			}
			return nil
		},
		Jar: jar,
	}

	req, _ := http.NewRequest("GET", urlLists.loginURL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		log.Println(resp.Cookies()[0])
	}
	defer resp.Body.Close()

	//fill form value
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	//get execution
	execution, isexist := doc.Find("[name=execution]").Attr("value")
	if !isexist {
		return nil, errors.New("未找到页面参数\"execution\",可能登陆逻辑改变,前往 " + urlLists.loginURL + " 按F12查看表单内容")
	}
	//get captchacode
	captchaReq, err := http.NewRequest("GET", urlLists.captchaURL, nil)
	if err == nil {
		log.Println("captcha download successfully")
	} else {
		log.Println("fail to download captcha")
		return nil, err
	}
	captchaReq.AddCookie(resp.Cookies()[0])
	captchaResp, _ := client.Do(captchaReq)
	defer captchaResp.Body.Close()
	bytes, _ := ioutil.ReadAll(captchaResp.Body)
	// _, imgCategories, err := image.DecodeConfig(strings.NewReader(string(bytes)))
	// if err != nil {
	// 	return nil, nil
	// }
	_ = ioutil.WriteFile("captcha.jpg", bytes, 0777) //判断图片类型意义不大
	// // read captchaCode from stdin
	// var captchaCode string
	// var openEngine *exec.Cmd
	// switch runtime.GOOS {
	// case "windows":
	// 	openEngine = exec.Command("cmd", "/C", "start", "pbrush", "./captcha."+imgCategories)
	// case "darwin": //mac os and linux are not tested
	// 	openEngine = exec.Command("open", "./captcha."+imgCategories)
	// case "linux":
	// 	openEngine = exec.Command("xdg-open", "./captcha."+imgCategories)
	// }
	// err = openEngine.Start()
	// if err != nil {
	// 	log.Println("尝试调用系统指令打开验证码图片失败,请手动打开")
	// }
	// log.Println("输入验证码: (若打开图片失败,在当前文件夹下手动打开captcha图片)")
	// fmt.Scanln(&captchaCode)
	var captchaCode string
	log.Println("正在自动识别:")
	cmd := exec.Command("py", "auto_captcha.py")
	bytessssss, err := cmd.Output()
	if err != nil {
		log.Println(err)
		log.Println("调用tesserate失败,手动输入:")
		cmd2 := exec.Command("cmd", "/C", "start", "pbrush", "./captcha.jpg")
		cmd2.Start()
		fmt.Scanln(&captchaCode)
	}else{
		captchaCode = string(bytessssss)
		if len(captchaCode) < 4 {
			return nil, errors.New("Captcha ERROR!")
		}
		captchaCode = captchaCode[:4] //从命令行获得输入后面会有空格
		log.Println("captcha:", captchaCode, "  len:", len(captchaCode))
	}
	

	// construct request body
	data := map[string][]string{
		"username":    {myusername},
		"password":    {mypassword},
		"captcha":     {captchaCode},
		"_eventId":    {"submit"},
		"execution":   {execution},
		"geolocation": {},
	}
	d := url.Values(data).Encode()
	bodydata := strings.NewReader(d)
	req2, _ := http.NewRequest("POST", urlLists.loginURL, bodydata)

	// req2.AddCookie(resp.Cookies()[0])
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	// req2.Header.Add("Referer", "https%3A%2F%2Fcas.sysu.edu.cn%2Fcas%2Flogin%3Fservice%3Dhttps%3A%2F%2Fjwxt.sysu.edu.cn%2Fjwxt%2Fapi%2Fsso%2Fcas%2Flogin%3Fpattern%3Dstudent-login")
	// req2.Header.Add("Referer", "https%3A%2F%2Fcas-443.webvpn.sysu.edu.cn%2Fcas%2Flogin%3Fservice%3Dhttps%3A%2F%2Fportal.sysu.edu.cn%2Fmanagement%2Fshiro-cas")
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("Origin", "https://cas-443.webvpn.sysu.edu.cn")

	//第二次请求,登陆请求,获取token
	resp2, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()
	// log.Println("Login response statusCode:", resp2.StatusCode) //这个由于页面跳转的原因,意义不大,新写了isLogin()函数判断
	return client, nil
}

// //tested , can't filter
// func getOneCourseInfo(client *http.Client , courseId string){
// 	jsonBody := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"2019-2","selectedType":"1","selectedCate":"21","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","collectionStatus":"0" }}`
// 	getCourseList(client, jsonBody)
// }

func isLoginAndGetInfo(client *http.Client) (string, bool) {

	type respJson struct {
		Meta struct {
			Success    bool   `json:"success"`
			StatusCode int    `json:"statusCode"`
			Message    string `json:"message"`
		} `json:"meta"`
		Data struct {
			UserID         string `json:"userId"`
			UserName       string `json:"userName"`
			DepartmentName string `json:"departmentName"`
			TokenID        string `json:"tokenId"`
		} `json:"data"`
	}

	infoStruct := respJson{}

	infoUrl := "https://portal.sysu.edu.cn/tryLoginUserInfo"
	req, _ := http.NewRequest("POST", infoUrl, nil)

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse //禁止重定向,才能判定是否登陆成功
	}
	defer func(client *http.Client) {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return http.ErrUseLastResponse
			} else {
				return nil
			}
		}
	}(client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Login Status Code:", resp.StatusCode, "登陆失败,检查用户密码验证码")
		return "", false
	} else {
		log.Println("Login Status Code:", resp.StatusCode, "登陆成功")
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(bytes, &infoStruct)
	if err != nil {
		return "json2struct error parse!", true
	}
	info := fmt.Sprintf("[USERINFO] Id:%s,Name:%s,[%s]\n<TOKEN>_astraeus_session:%s\n-[o]-[0]-> (๑•̀ㅂ•́)و✧ | (u‿ฺu✿ฺ) --- <(￣︶￣)↗[GO!]",
		infoStruct.Data.UserID, infoStruct.Data.UserName, infoStruct.Data.DepartmentName, infoStruct.Data.TokenID)
	// 记录session
	cookie := "_astraeus_session=" + infoStruct.Data.TokenID + " time=" + fmt.Sprintf("%d", time.Now().Unix())
	ioutil.WriteFile("tmp", []byte(cookie), 0777)

	return info, true
}

func getJwxtCookieWithWebVpn(client *http.Client)  {
	url := "https://cas-443.webvpn.sysu.edu.cn/cas/login?service=https%3A%2F%2Fjwxt-443.webvpn.sysu.edu.cn%2Fjwxt%2Fapi%2FtoCasUrl%3FtoCasUrl%3D%252F"
	req , _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	// req.Header.Add("Referer","https://portal.sysu.edu.cn/")  //这个头服务端也没有检查
	client.Do(req)
	return
}

// 下面的函数暂时没有用到,因为意义不是很大
func isCookieValid() (string, bool) {
	file, err := os.Open("tmp")
	if err != nil {
		log.Println("SessionId is not VALID!")
		return "", false
	}
	bytes, _ := ioutil.ReadAll(file)
	str := string(bytes)
	var lastTimeStr string
	for i := 0; i < len(str); i++ {
		if str[i] == 't' {
			lastTimeStr = str[i+2:]
			break
		}
	}
	currTime := time.Now().Unix()
	lastTime, _ := strconv.Atoi(lastTimeStr)
	if currTime-int64(lastTime) > int64(time.Hour.Seconds()) {
		// 认为超过1小时就过期,不清楚服务器后台是怎么设置的
		log.Println("SessionId is not VALID!")
		//delete
		file.Close()
		os.Remove("tmp")
		return "", false
	}
	file.Close()
	return str, true
}

func setCookieAndHaveATry(client *http.Client, cookie string) *http.Client {
	pivot := strings.Index(cookie, "=")
	end := strings.LastIndex(cookie, " time=")
	url, _ := url.Parse(urlLists.loginURL)
	cookieHTTP := &http.Cookie{
		Name:   cookie[:pivot],
		Value:  cookie[pivot+1 : end],
		Domain: ".sysu.edu.cn",
		Path:   "/",
	}
	slice := []*http.Cookie{cookieHTTP}
	client.Jar.SetCookies(url, slice)
	return client
}
