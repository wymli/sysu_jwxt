package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// "os"
	"io/ioutil"
	// "net/http"
	// "github.com/kataras/iris/cache/client"
	// "net/http"
	// "net/http"
	// "net/http/cookiejar"
	// "os"
)

var (
	myusername string
	mypassword string
)

// const (
// 	myusername = "liwm29"
// 	mypassword = "`S0a0l0vare"
// )

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	log.SetPrefix("[@STrelitziA@]")

	bytes, _ := ioutil.ReadFile("userConfig.txt")
	config := strings.Split(string(bytes), "&")
	myusername = config[0]
	mypassword = config[1]
	log.Println("Username:"+myusername, "len:", len(myusername))
	log.Println("Password:"+mypassword, "len:", len(mypassword))

}

const mode = "webvpn"

// const mode = "normal"  // 校内

func main() {
	// cookie,mayValid := isCookieValid() //这段代码是检查cookie的,意义不大,直接请求新cookie就行


	urlLists.init(mode)
	var client *http.Client
	var cnt int = 1
	for {
		var err error
		client, err = casLogin() // 这里不能用 :=  ,否则会创建一个局部的client
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if info, ok := isLoginAndGetInfo(client); ok {
			fmt.Println(info)
			if mode == "webvpn" { //获取webvpn模式下 jwxt的cookie
				getJwxtCookieWithWebVpn(client)
			}
			break
		} else {
			log.Println("重启，尝试次数：", cnt)
		}
		cnt++ //避免重复太多次
		if cnt > 5 {
			panic("失败次数过多!")
		}
	}

	// getCourseList(client, urlLists.getPublicCoursejsonBody)

	// str, _ := getMyTeachersInfo(client)
	// fmt.Println(str)

	jksb(client)

	// payload := `{"clazzId":"1201412705275330561","selectedType":"1","selectedCate":"21","check":true}`
	// grabCourse(client, payload, 10)

	// if gpaInfo,err := getGPA(client);err==nil{
	// 	fmt.Println(gpaInfo)
	// }else{
	// 	log.Println(gpaInfo , err)
	// }

	// // courseSelectionChooseBody := `{"clazzId":"1201412705275330561","selectedType":"1","selectedCate":"21","check":true}` //专选
	// if teacherInfo,err := getMyTeachersInfo(client);err==nil{
	// 	fmt.Println(teacherInfo)
	// }else{
	// 	log.Println(teacherInfo , err)
	// }

}
