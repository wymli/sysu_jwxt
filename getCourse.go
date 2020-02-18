package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	// "os"
	"strings"
)

type courseSelectList struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    struct {
		Total int `json:"total"`
		Rows  []struct {
			MainClassesID        string      `json:"mainClassesID"`
			TeachingClassID      string      `json:"teachingClassId"`
			TeachingClassNum     string      `json:"teachingClassNum"`
			TeachingClassName    interface{} `json:"teachingClassName"`
			CourseNum            string      `json:"courseNum"`
			CourseName           string      `json:"courseName"`
			Credit               float64     `json:"credit"`
			ExamFormName         string      `json:"examFormName"`
			CourseUnitNum        string      `json:"courseUnitNum"`
			CourseUnitName       string      `json:"courseUnitName"`
			TeachingTeacherNum   string      `json:"teachingTeacherNum"`
			TeachingTeacherName  string      `json:"teachingTeacherName"`
			BaseReceiveNum       int         `json:"baseReceiveNum"`
			AddReceiveNum        int         `json:"addReceiveNum"`
			TeachingTimePlace    string      `json:"teachingTimePlace"`
			StudyCampusID        string      `json:"studyCampusId"`
			Week                 string      `json:"week"`
			ClassTimes           string      `json:"classTimes"`
			CourseSelectedNum    string      `json:"courseSelectedNum"`
			FilterSelectedNum    string      `json:"filterSelectedNum"`
			SelectedStatus       string      `json:"selectedStatus"`
			CollectionStatus     string      `json:"collectionStatus"`
			TeachingLanguageCode string      `json:"teachingLanguageCode"`
			PubCourseTypeCode    interface{} `json:"pubCourseTypeCode"`
			CourseCateCode       string      `json:"courseCateCode"`
			SpecialClassCode     interface{} `json:"specialClassCode"`
			SportItemID          interface{} `json:"sportItemId"`
			RecordMode           string      `json:"recordMode"`
			ClazzNum             string      `json:"clazzNum"`
			ExamFormCode         string      `json:"examFormCode"`
			CourseID             string      `json:"courseId"`
			ScheduleExamTime     interface{} `json:"scheduleExamTime"`
		} `json:"rows"`
	} `json:"data"`
}

func getCourseList(client *http.Client, payload string) error {
	// jsonBody := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"2019-2","selectedType":"4","selectedCate":"11","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","collectionStatus":"0"}}`
	// jsonBody := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"2019-2","selectedType":"1","selectedCate":"21","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","collectionStatus":"0"}}`

	//时间戳,不加也行
	// timestamp := time.Now().Unix()
	// var courseListUrl_t =  courseListUrl + "?_t=" + fmt.Sprintf("%d", timestamp)
	log.Println("courseListUrl : ", urlLists.courseListUrl)
	log.Println("Query params :", payload)
	courseListReq, _ := http.NewRequest("POST", urlLists.courseListUrl, strings.NewReader(payload))
	courseListReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	courseListReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	courseListReq.Header.Add("Referer", urlLists.baseUrl+"jwxt/mk/courseSelection/")
	courseListResp, _ := client.Do(courseListReq)
	defer courseListResp.Body.Close()
	b, _ := ioutil.ReadAll(courseListResp.Body)

	var courseList courseSelectList
	err := json.Unmarshal(b, &courseList)
	if err != nil {
		return err
	}
	//寻找有空位的课
	totalCourses := courseList.Data.Total
	var times = totalCourses / 10 // 10 is one page size
	for i := 0; i < times; i++ {
		payload = strings.ReplaceAll(payload, `"pageNo":`+fmt.Sprintf("%d", i+1), `"pageNo":`+fmt.Sprintf("%d", i+2))
		courseListReq2, _ := http.NewRequest("POST", urlLists.courseListUrl, strings.NewReader(payload))
		courseListReq2.Header.Add("Content-Type", "application/json;charset=UTF-8")
		courseListReq2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
		courseListReq2.Header.Add("Referer", urlLists.baseUrl+"jwxt/mk/courseSelection/")

		resp2, _ := client.Do(courseListReq2)
		defer resp2.Body.Close()
		b2, _ := ioutil.ReadAll(resp2.Body)
		b = append(b, '\n')
		b = append(b, b2...)
	}
	filename := "CourseList_t=" + fmt.Sprintf("%d", time.Now().Unix())
	ioutil.WriteFile(filename, b, 0777)
	log.Println("写入文件成功 |" + filename)
	return nil
}

func courseChoose(client *http.Client, courseSelectionChooseBody string) string {
	// courseSelectionChooseBody := `{"clazzId":"1201412705275330561","selectedType":"1","selectedCate":"21","check":true}` //专选 ,classid是teachingclassid
	// {"clazzId":"1208910925716574209","selectedType":"4","selectedCate":"21","check":true} //公选
	// {"code":52021104,"message":"你已选择过该课程，不能再选！","data":null}
	// {"code":200,"message":null,"data":"选课成功!"}
	// {"code":52021107,"message":"不能超过公选课限选的最大门数！","data":null}
	log.Println("Query params :", courseSelectionChooseBody)
	courseSelectionChooseReq, _ := http.NewRequest("POST", urlLists.courseSelectionChooseUrl, strings.NewReader(courseSelectionChooseBody))
	courseSelectionChooseReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	courseSelectionChooseReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	courseSelectionChooseReq.Header.Add("Referer", urlLists.baseUrl+"jwxt/mk/courseSelection/")
	courseSelectionChooseResp, _ := client.Do(courseSelectionChooseReq)
	defer courseSelectionChooseResp.Body.Close()
	courseSelectionChooseRespJsonBytes, _ := ioutil.ReadAll(courseSelectionChooseResp.Body)
	log.Println(string(courseSelectionChooseRespJsonBytes))
	return string(courseSelectionChooseRespJsonBytes)
}

func cancelCourse(client *http.Client, courseSelectionCancelBody string) string {
	// courseSelectionCancelBody := `{"courseId":"206169488","clazzId":"1201412705275330561","selectedType":"1"}`
	// {"code":200,"message":null,"data":"退课成功！"}
	// 多次退课都是退课成功
	log.Println("Query params :", courseSelectionCancelBody)
	courseSelectionCancelReq, _ := http.NewRequest("POST", urlLists.courseSelectionCancelUrl, strings.NewReader(courseSelectionCancelBody))
	courseSelectionCancelReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	courseSelectionCancelReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	courseSelectionCancelReq.Header.Add("Referer", urlLists.baseUrl+"jwxt/mk/courseSelection/")
	courseSelectionCancelResp, _ := client.Do(courseSelectionCancelReq)
	defer courseSelectionCancelResp.Body.Close()
	courseSelectionCancelRespJsonBytes, _ := ioutil.ReadAll(courseSelectionCancelResp.Body)
	log.Println(string(courseSelectionCancelRespJsonBytes))
	return string(courseSelectionCancelRespJsonBytes)
}

func grabCourse(client *http.Client, payload string, timeSeperate int) {
	log.Println("开始抢课---|->")
	for {
		ans := courseChoose(client, payload)
		if ans[8] == 2 { //200
			log.Println("抢课成功")
			break
		} else { // 52021104  52021107
			log.Println("失败,睡眠" + fmt.Sprintf("%d", timeSeperate) + "s")
			time.Sleep(time.Duration(timeSeperate) * time.Second)
		}
	}
}
