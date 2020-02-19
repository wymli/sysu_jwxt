# sysu_jwxt
written by go

---
依赖:  
- 非必须(用于自动验证码识别):
  1. tesseract(安装后加环境目录)
  2. python,以及库:pytesseract,opencv
- 必须:
  1. golang,以及库:goquery
---
- 运行方式:
> 在目录下:  
编译:`go build`  
运行:`./a`
---
- 写了一些函数,可以自己在main.go的main函数里面添加  
一些表单的参数要自己F12看然后加在代码里面   
---
已经完成的函数有:  
1. 登陆
1. 查询任课老师列表,下载任课老师照片,来源是评教系统  
1. 查询成绩(gpa / 课程成绩)
2. 查询课程(公选/专选)  
3. 选课退课(想抢课写个循环就行)  
4. 一键健康申报
---
- 待完成:
- [ ] 选课循环
- [ ] 查询单一课程,为选课做准备;因为服务器不支持查询单一课程,所以只能用page+pagesize=1来查询,要先不断获取课程列表直到找到

---
- 有两个配置文件:userConfig.txt 和  jksb_formdata.txt  
  - userConfig写用户名密码: 格式: username&password  
  - jksb_formdata直接复制健康申报网站(F12)提交的表单的其中一项:FormData
- 有两种登陆方式:校内直接访问jwxt和使用webvpn(从portal.sysu.edu.cn)登入,对应于mode: webvpn / normal  
设置方式:main.go里面有一个常量`mode`

---
- 注意的是尽量不要用windows记事本编辑配置文件,避免出现奇奇怪怪的bug!