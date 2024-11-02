Hello everyone, the project will undergo a refactoring in November. The main goals of the refactoring include the introduction of Redis Sentinel, traceability, observability, Docker, a recommendation system implemented with Spark and Flink, Feed streams, and im system, as well as restructuring the directory and interfaces.If you have any questions or good ideas, please contact me via QQ 1914163770.

大家好，该项目将会在11月进行重构，主要重构的目标包括引入Redis Sentinel、链路追踪、可观测性、Docker、Spark及Flink实现的推荐系统、Feed流、im系统等，并对目录及接口进行重构。如果你有任何的疑问或者好的想法，请联系QQ 1914163770 

# 第六届字节跳动青训营 后端 极简版抖音项目
项目荣获第六届青训营二等奖，感谢队友的支持<br>
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






