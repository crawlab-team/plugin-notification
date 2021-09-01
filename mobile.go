package main

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/imroc/req"
	"runtime/debug"
)

type ResBody struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func GetMobileTaskMarkdownContent(t interfaces.Task, s interfaces.Spider, n interfaces.Node, ts interfaces.TaskStat) string {
	errMsg := ""
	errLog := "-"
	statusMsg := fmt.Sprintf(`<font color="#00FF00">%s</font>`, t.GetStatus())
	if t.GetStatus() == constants.TaskStatusError {
		errMsg = `（有错误）`
		errLog = fmt.Sprintf(`<font color="#FF0000">%s</font>`, t.GetError())
		statusMsg = fmt.Sprintf(`<font color="#FF0000">%s</font>`, t.GetStatus())
	}
	return fmt.Sprintf(`
您的任务已完成%s，请查看任务信息如下。

- **任务ID:** %s
- **任务状态:** %s
- **任务参数:** %s
- **爬虫ID:** %s
- **爬虫名称:** %s
- **节点:** %s
- **创建时间:** %s
- **开始时间:** %s
- **完成时间:** %s
- **等待时间:** %d秒
- **运行时间:** %d秒
- **总时间:** %d秒
- **结果数:** %d
- **错误:** %s

请登录Crawlab查看详情。
`,
		errMsg,
		t.GetId().Hex(),
		statusMsg,
		t.GetParam(),
		s.GetId().Hex(),
		s.GetName(),
		n.GetName(),
		utils.GetLocalTimeString(ts.GetCreateTs()),
		utils.GetLocalTimeString(ts.GetStartTs()),
		utils.GetLocalTimeString(ts.GetEndTs()),
		ts.GetWaitDuration()/1000,
		ts.GetRuntimeDuration()/1000,
		ts.GetTotalDuration()/1000,
		ts.GetResultCount(),
		errLog,
	)
}

func SendMobileNotification(webhook string, title string, content string) error {
	// request header
	header := req.Header{
		"Content-Type": "application/json; charset=utf-8",
	}

	// request data
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"title":   title,
			"text":    content,
			"content": content,
		},
		"at": req.Param{
			"atMobiles": []string{},
			"isAtAll":   false,
		},
	}

	// perform request
	res, err := req.Post(webhook, header, req.BodyJSON(&data))
	if err != nil {
		log.Errorf("dingtalk notification error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// parse response
	var resBody ResBody
	if err := res.ToJSON(&resBody); err != nil {
		log.Errorf("dingtalk notification error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// validate response code
	if resBody.ErrCode != 0 {
		log.Errorf("dingtalk notification error: " + resBody.ErrMsg)
		debug.PrintStack()
		return errors.New(resBody.ErrMsg)
	}

	return nil
}
