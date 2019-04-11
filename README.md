events
[![Build Status](https://travis-ci.org/issue9/events.svg?branch=master)](https://travis-ci.org/issue9/events)
[![Go Report Card](https://goreportcard.com/badge/github.com/issue9/events)](https://goreportcard.com/report/github.com/issue9/events)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/issue9/events/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/events)
======


简单的事件订阅发布系统

```go
p, e := events.New()

e.Attach(sub1)
e.Attach(sub2)

p.Publish("触发事件1") // sub1 和 sub2 均会收事事件
```



### 安装

```shell
go get github.com/issue9/events
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/events)
[![GoDoc](https://godoc.org/github.com/issue9/events?status.svg)](https://godoc.org/github.com/issue9/unique)


### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
