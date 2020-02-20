package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "regexp"
	"strconv"
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

type course struct {
	classId    string
	courseName string
}

type CourseListStruct struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
	Param    struct {
		SemesterYear         string `json:"semesterYear"`
		SelectedType         string `json:"selectedType"`
		SelectedCate         string `json:"selectedCate"`
		HiddenConflictStatus string `json:"hiddenConflictStatus"`
		HiddenSelectedStatus string `json:"hiddenSelectedStatus"`
		CollectionStatus     string `json:"collectionStatus"`
	} `json:"param"`
}

func DefaulyCourseListPayload() func(pageNo, pageSize int, semesterYear, selectedType, selectedCate string) CourseListStruct {
	return func(no, size int, semesterYear, selectedType, selectedCate string) CourseListStruct {
		return CourseListStruct{
			no, size, struct {
				SemesterYear         string `json:"semesterYear"`
				SelectedType         string `json:"selectedType"`
				SelectedCate         string `json:"selectedCate"`
				HiddenConflictStatus string `json:"hiddenConflictStatus"`
				HiddenSelectedStatus string `json:"hiddenSelectedStatus"`
				CollectionStatus     string `json:"collectionStatus"`
			}{
				semesterYear, selectedType, selectedCate, "0", "0", "0",
			},
		}
	}
}

func getCourseList(client *http.Client, payloadStruct CourseListStruct) ([]course, string, error) {
	payload, _ := json.Marshal(payloadStruct)
	// jsonBody := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"2019-2","selectedType":"4","selectedCate":"11","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","collectionStatus":"0"}}`
	// jsonBody := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"2019-2","selectedType":"1","selectedCate":"21","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","collectionStatus":"0"}}`
	var retCourseList = make([]course, 0)
	//时间戳,不加也行
	// timestamp := time.Now().Unix()
	// var courseListUrl_t =  courseListUrl + "?_t=" + fmt.Sprintf("%d", timestamp)
	log.Println("courseListUrl : ", urlLists.courseListUrl)
	log.Println("Query params :", string(payload))
	courseListReq, _ := http.NewRequest("POST", urlLists.courseListUrl, strings.NewReader(string(payload)))
	courseListReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	courseListReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	courseListReq.Header.Add("Referer", urlLists.baseUrl+"jwxt/mk/courseSelection/")
	courseListResp, _ := client.Do(courseListReq)
	defer courseListResp.Body.Close()
	b, _ := ioutil.ReadAll(courseListResp.Body)

	var courseList courseSelectList
	err := json.Unmarshal(b, &courseList)
	if err != nil {
		return nil, "", err
	}
	//添加进返回列表
	for _, it := range courseList.Data.Rows {
		id := it.TeachingClassID
		name := it.CourseName
		retCourseList = append(retCourseList, course{
			id, name,
		})
	}
	//添加结束

	//遍历所有课程
	totalCourses := courseList.Data.Total
	var times = totalCourses / 10 // 10 is one page size
	for i := 0; i < times; i++ {
		payloadS := strings.ReplaceAll(string(payload), `"pageNo":`+fmt.Sprintf("%d", i+1), `"pageNo":`+fmt.Sprintf("%d", i+2))
		courseListReq2, _ := http.NewRequest("POST", urlLists.courseListUrl, strings.NewReader(payloadS))
		courseListReq2.Header.Add("Content-Type", "application/json;charset=UTF-8")
		courseListReq2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
		courseListReq2.Header.Add("Referer", urlLists.baseUrl+"jwxt/mk/courseSelection/")

		resp2, _ := client.Do(courseListReq2)
		defer resp2.Body.Close()
		b2, _ := ioutil.ReadAll(resp2.Body)
		b = append(b, '\n')
		b = append(b, b2...)

		var courseList courseSelectList
		err = json.Unmarshal(b2, &courseList)
		if err != nil {
			return nil, "", err
		}
		//添加进返回列表
		for _, it := range courseList.Data.Rows {
			id := it.TeachingClassID
			name := it.CourseName
			retCourseList = append(retCourseList, course{
				id, name,
			})
		}
		//添加结束
	}
	filename := "CourseList_t=" + fmt.Sprintf("%d", time.Now().Unix())
	ioutil.WriteFile(filename, b, 0777)
	log.Println("写入文件成功 |" + filename)
	return retCourseList, string(payload), nil
}

type chooseCourseStruct struct {
	ClazzID      string `json:"clazzId"`
	SelectedType string `json:"selectedType"`
	SelectedCate string `json:"selectedCate"`
	Check        bool   `json:"check"`
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

// jsonBody := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"2019-2","selectedType":"4","selectedCate":"11","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","collectionStatus":"0"}}`
func grabCourse(client *http.Client, payloadStruct CourseListStruct, classId string, timeSeperate int) { //查询课程的payload
	log.Println("开始抢课---|->")
	// // 选课payload
	// reg, _ := regexp.Compile(`(?<="selectedType":").*?(?=")`)
	// SelectedType := reg.FindString(payload)
	// reg, _ = regexp.Compile(`(?<="selectedCate":").*?(?=")`)
	// SelectedCate := reg.FindString(payload)
	payloadStructC := chooseCourseStruct{
		classId, payloadStruct.Param.SelectedType, payloadStruct.Param.SelectedCate, true,
	}
	payload2, _ := json.Marshal(payloadStructC)

	page, payload := getOneCoursePage(client, payloadStruct, classId)
	for {
		if ok := queryWithOneCourse(client, payload, page); ok {
			ans := courseChoose(client, string(payload2))
			if ans[8] == 2 { //200
				log.Println("抢课成功")
				break
			} else { // 52021104  52021107
				log.Println("失败")
			}
		}
		log.Println("睡眠" + fmt.Sprintf("%d", timeSeperate) + "s")
		time.Sleep(time.Duration(timeSeperate) * time.Second)
	}

	for {
		ans := courseChoose(client, string(payload2))
		if ans[8] == 2 { //200
			log.Println("抢课成功")
			break
		} else { // 52021104  52021107
			log.Println("失败,睡眠" + fmt.Sprintf("%d", timeSeperate) + "s")
			time.Sleep(time.Duration(timeSeperate) * time.Second)
		}
	}
}

func getOneCoursePage(client *http.Client, payloadStruct CourseListStruct, classId string) (int, string) { //查询课程的payload
	courseList, payload, err := getCourseList(client, payloadStruct)
	if err != nil {
		log.Println(err)
	}
	cnt := 0
	for _, it := range courseList {
		cnt++
		if it.classId == classId {
			log.Println("找到课程:", it.courseName, "classId:", classId)
			return cnt, payload
		}
	}
	log.Println("未找到课程:", "classId:", classId)
	return 0, ""
}

func queryWithOneCourse(client *http.Client, payload string, page int) bool {
	payload = strings.ReplaceAll(payload, `"pageNo":1,"pageSize":10,`, `"pageNo":`+strconv.Itoa(page)+`,"pageSize":1,`)
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
		return false
	}
	vacantNum, _ := strconv.Atoi(courseList.Data.Rows[0].CourseSelectedNum)
	if courseList.Data.Rows[0].BaseReceiveNum-vacantNum == 0 {
		//满
		log.Println("课程:", courseList.Data.Rows[0].CourseName, "classId:", courseList.Data.Rows[0].TeachingClassID, "课程容量:", courseList.Data.Rows[0].BaseReceiveNum, "已选人数:", courseList.Data.Rows[0].CourseSelectedNum, "无空位")
		return false
	} else {
		log.Println("课程:", courseList.Data.Rows[0].CourseName, "classId:", courseList.Data.Rows[0].TeachingClassID, "课程容量:", courseList.Data.Rows[0].BaseReceiveNum, "已选人数:", courseList.Data.Rows[0].CourseSelectedNum, "有空位")
		return true
	}
}
