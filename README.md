# 极简版抖音服务器
本项目是使用Go语言开发，基于Hertz + Kitex +  MySQL + MongoDB + Redis + Kafka + Gorm + Zap + Etcd +OSS等技术实现的极简版抖音APP后端项目，该项目部署在华为云服务器上，实现了基础功能以及互动和社交方向的全部功能。<br>
项目团队：起名起了3min<br>
项目文档：https://vish8y9znlg.feishu.cn/docx/XffIdI4sso6oGNx2yWEc4DV4nrh  <br>

**如何运行本项目？** <br>
使用 Linux 环境<br>
安装go 1.20 版本<br>
安装MySQL 8.0 以上版本<br>
安装Redis 6.2 以上版本<br>
安装Kafka 3.0 以上版本<br>
安装MongoDB 4.4 以上版本<br>

并在config/app.yaml 中修改配置<br>

**由于项目启用了https**<br>
需要在service-api-main.go里，将main函数换成别的名字，然后将下面的mainWithOutTls()函数改为main()函数，就是走的http请求<br>

最后 cd 到项目根目录<br>
sh ./build_all_services.sh  

sh ./start_all_services.sh  


