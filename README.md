# project
项目团队：起名起了3min

如何运行本项目？
使用 Linux 环境
安装go 1.20 版本
安装MySQL 8.0 以上版本
安装Redis 6.2 以上版本
安装Kafka 3.0 以上版本
安装MongoDB 4.4 以上版本

在config/app.yaml 中修改配置

由于项目启用了https
需要在service-api-main.go里，将main函数换成别的名字，然后将下面的mainWithOutTls()函数改为main()函数，就是走的http请求

cd 到项目根目录
sh ./build_all_services.sh
sh ./start_all_services.sh

