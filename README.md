# sysu_jwxt
written by go

---
依赖:  
1. tesseract(安装后记得加环境目录)
2. python,以及库:pytesseract
3. golang,以及库:goquery

> 在目录下:  
编译:`go build`  
运行:`./a`
---
函数自己在main.go里面加  
一些表单的参数要自己F12看然后加在代码里面   

---
已经完成的函数有: 
1. 查询任课老师列表,下载任课老师照片,来源是评教系统  
2. 查询课程(公选/专选)  
3. 选课退课(想抢课写个循环就行)  
4. 一键健康申报

---
- 有两个配置文件:userConfig.txt 和  jksb_formdata.txt  
  - userConfig写用户名密码: 格式: username&password  
  - jksb_formdata直接复制健康申报网站按F12(没记错的话是异步请求的一个叫nextuserlist的,doaction也行,那个网站请求了两次)提交的formdata(注意不是全部表单,就一项命名为formdata的,详情看代码)
