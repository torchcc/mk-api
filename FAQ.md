## 记录本项目遇到的坑

### 1. IDE 和环境问题

- Goland 打开工程后，记得设置【File】【settings...】

  - 【Go】【GOROOT】设置SDK路径（go版本在 1.11 后），如果系统的环境变量里面已经设置，则可跳过

  - 【Go】【GOPATH】可不选

  - 【Go】【Go Modules (vgo)】

    > 参考：[Golang1.13.x 解决 go get 无法下载问题](https://www.cnblogs.com/zhangmingcheng/p/12294156.html)

    Proxy: https://goproxy.cn,direct

    或者在 Terminal 窗口里键入：

    ```bash
    go env -w GO111MODULE=on
    go env -w GOPROXY=https://goproxy.cn,direct
    ```


### 2. 开始

```bash
go mod tidy
```





