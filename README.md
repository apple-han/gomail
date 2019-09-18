# goemail
相信这个项目，对于刚入微服务开发与web的你来说，帮助是巨大的。
##### 项目开始之前，你需要申请一个163的smpt的账号，密码。
### Install from docker 
    git clone https://github.com/apple-han/goemail.git
    cd pkg/configuration_center/server
    make build
    make docker
    cd email
    make build
    make docker
    cd user
    make build
    make docker
    docker-compose -f docker-compose.yml up -d
### 本源码视频教程地址（https://study.163.com/course/courseMain.htm?courseId=1209482821&share=2&shareId=400000000606033）
