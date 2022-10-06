# go-swjtu-yunpan
交大云盘(Anyshare)的go实现
## 说明
本项目是交大云盘(Anyshare)的go实现，目前仅支持文件上传和下载，其他功能正在开发中。目前登陆尚未完善，需要抓取加密后的密码进行登陆。
<br>(第一次写go，用来练手)
## 使用
### 安装
在release页面下载对应平台的可执行文件，或者自行编译。
### 使用
首先创建.env文件，内容如下：
```
STUID=学号
PASSWORD=密码(加密后的再base64一层)
```
然后执行
```
./go-swjtu-yunpan
```
即可使用。
### 命令
```
cd [path]   切换目录
ls [path]  列出文件
upload [path]   上传文件到当前目录
download [filename] 下载当前目录下的文件
pwd 显示当前目录
exit    退出
```

