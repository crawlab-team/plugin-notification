package main

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/models"
	plugin "github.com/crawlab-team/crawlab-plugin"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	*plugin.Internal
}

func (svc *Service) Init() (err error) {
	api := svc.GetApi()
	api.POST("/send/mail", svc.sendMail)
	api.POST("/send/mobile", svc.sendMobile)
	api.POST("/set", svc.set)
	api.POST("/set/:model/*oid", svc.set)
	api.POST("/delete", svc.delete)
	api.POST("/delete/:model/*oid", svc.delete)
	return nil
}

func (svc *Service) Start() (err error) {
	svc.StartApi()
	return nil
}

func (svc *Service) Stop() (err error) {
	svc.StopApi()
	return nil
}

func (svc *Service) sendMail(c *gin.Context) {
	controllers.HandleSuccess(c)
}

func (svc *Service) sendMobile(c *gin.Context) {
	var payload SendPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// models
	t, s, n, ts, err := svc._getModels(payload)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// title
	title := fmt.Sprintf("[Crawlab] \"%s\" 任务 %s", s.GetName(), t.GetStatus())
	content := GetMobileTaskMarkdownContent(t, s, n, ts)

	// TODO: test
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=7e08bf6f891b0ffc81fc91c4871f744251df5b143db242e271fb28c158a3176c"

	// send
	if err := SendMobileNotification(webhook, title, content); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccess(c)
}

func (svc *Service) set(c *gin.Context) {
	model, oid, query := svc._getParamsAndQuery(c)

	var value bson.M
	if err := c.ShouldBindJSON(&value); err != nil {
		controllers.HandleErrorBadRequest(c, err)
		return
	}

	// extra value model service
	evSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdExtraValue)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// extra value
	var ev interfaces.ExtraValue

	// attempt to find in database
	doc, err := evSvc.Get(query, nil)
	if err != nil {
		if err.Error() != mongo.ErrNoDocuments.Error() {
			// error
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
		// not exists, add a new one
		ev = &models.ExtraValue{
			ObjectId: oid,
			Model:    model,
			Type:     ExtraValueTypeNotification,
			Value:    nil,
		}
		if err := client.NewModelDelegate(ev).Add(); err != nil {
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		// exists, update
		ev = doc.(interfaces.ExtraValue)
		ev.SetValue(value)
		if err := client.NewModelDelegate(ev).Save(); err != nil {
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
	}

	controllers.HandleSuccess(c)
}

func (svc *Service) delete(c *gin.Context) {
	_, _, query := svc._getParamsAndQuery(c)

	// extra value model service
	evSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdExtraValue)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// attempt to find in database
	doc, err := evSvc.Get(query, nil)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// delete
	ev := doc.(interfaces.ExtraValue)
	if err := client.NewModelDelegate(ev).Delete(); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccess(c)
}

func (svc *Service) _getModels(payload SendPayload) (t interfaces.Task, s interfaces.Spider, n interfaces.Node, ts interfaces.TaskStat, err error) {
	var doc interfaces.Model

	// task
	taskSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdTask)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	doc, err = taskSvc.GetById(payload.TaskId)
	if err != nil {
		return
	}
	t = doc.(interfaces.Task)

	// spider
	spiderSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdSpider)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	doc, err = spiderSvc.GetById(t.GetSpiderId())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	s = doc.(interfaces.Spider)

	// node
	nodeSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdNode)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	doc, err = nodeSvc.GetById(t.GetNodeId())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	n = doc.(interfaces.Node)

	// task stat
	taskStatSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdTaskStat)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	doc, err = taskStatSvc.GetById(t.GetId())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	ts = doc.(interfaces.TaskStat)

	return t, s, n, ts, nil
}

func (svc *Service) _getParamsAndQuery(c *gin.Context) (model string, oid primitive.ObjectID, query bson.M) {
	model = c.Param("model")
	oidStr := c.Param("oid")

	// query
	query = bson.M{"t": ExtraValueTypeNotification}

	// if empty, set model as global
	if model == "" {
		model = ExtraValueModelGlobal
	}
	query["m"] = bson.M{"m": model}

	// object id (associated with model)
	if oidStr != "" {
		var err error
		oid, err = primitive.ObjectIDFromHex(oidStr)
		if err != nil {
			controllers.HandleErrorBadRequest(c, err)
			return
		}
		query["oid"] = oid
	}

	return model, oid, query
}

func NewService() *Service {
	// service
	svc := &Service{
		Internal: plugin.NewInternal(),
	}

	if err := svc.Init(); err != nil {
		panic(err)
	}

	return svc
}
