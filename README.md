# beautyPictureCollection

批量采集美女写真等图片到typecho上，图片可选择本地化，内含多种数据源

![1.png](https://github.com/Anderyly/beautyPictureCollection/blob/master/1.png?raw=true)

👉 http://www.52rm.cc

# 编译
```shell
    git clone https://github.com/Anderyly/beautyPictureCollection.git warm
    cd warm
    go mod init && go mod tidy
    go build
```

# 数据库

### 在typecho_contents表增加link字段 vachar 255即可 用于存储采集来源页

# 运行
### 请配置set.ini内容
```shell
    chomd +x warm
    ./warm
```

# Issues

- [https://github.com/Anderyly/beautyPictureCollection/issues](https://github.com/Anderyly/beautyPictureCollection/issues)
