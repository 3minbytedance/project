# 第六届字节跳动青训营 后端 极简版抖音项目
项目荣获第六届青训营二等奖，感谢队友的支持
本项目是使用Go语言开发，基于Hertz + Kitex +  MySQL + MongoDB + Redis + Kafka + Gorm + Zap + Etcd +OSS等技术实现的极简版抖音APP后端项目，该项目部署在华为云服务器上，实现了基础功能以及互动和社交方向的全部功能。<br>
项目团队：起名起了3min<br>
项目文档：https://vish8y9znlg.feishu.cn/docx/XffIdI4sso6oGNx2yWEc4DV4nrh  <br>
在架构选型上，项目演进从gin -> hertz+kitex<br>
其中main分支为稳定大版本<br>
develop-rpc 为目前正在开发的rpc分支<br>
develop分支为最初的gin单体原型设计<br>
目前main分支和本地分支不同步，因为在找工作暂时搁置了，后面等不忙了再重新优化

**如何运行本项目？** <br>
使用 Linux 环境<br>
安装go 1.20 版本<br>
安装MySQL 8.0 及以上版本<br>
安装Redis 6.2 及以上版本<br>
安装Kafka 3.0 及以上版本<br>
安装MongoDB 4.4 及以上版本<br>

并在config/app.yaml 中修改配置<br>

**由于项目启用了https**<br>
需要在service-api-main.go里，将main函数换成别的名字，然后将下面的mainWithOutTls()函数改为main()函数，就是走的http请求<br>

最后 cd 到项目根目录<br>
sh ./build_all_services.sh  

sh ./start_all_services.sh  

如果执行sh ./build_all_services.sh  提示：
'\r': command not found <br>
这是由于goLand在windows环境下将换行符LF转为了CRLF换行 <br> 需要选中项目根目录，然后修改换行符
![image](https://github.com/3minbytedance/project/assets/42531412/88fd695e-422f-469e-9477-0ca0e35e2d38)






