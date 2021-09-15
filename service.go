package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	mongo2 "github.com/crawlab-team/crawlab-db/mongo"
	grpc "github.com/crawlab-team/crawlab-grpc"
	plugin "github.com/crawlab-team/crawlab-plugin"
	"github.com/crawlab-team/go-trace"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"strings"
	"time"
)

type Service struct {
	*plugin.Internal
	col *mongo2.Col // notification settings
}

func (svc *Service) Init() (err error) {
	// handle events
	go svc.handleEvents()

	// api
	api := svc.GetApi()
	api.POST("/send", svc.send)
	api.GET("/settings", svc.getSettingList)
	api.GET("/settings/:id", svc.getSetting)
	api.PUT("/settings", svc.putSetting)
	api.POST("/settings/:id", svc.postSetting)
	api.DELETE("/settings/:id", svc.deleteSetting)
	api.POST("/settings/:id/enable", svc.enableSetting)
	api.POST("/settings/:id/disable", svc.disableSetting)

	return nil
}

func (svc *Service) Start() (err error) {
	// initialize data
	if err := svc.initData(); err != nil {
		return err
	}

	// start api
	svc.StartApi()

	return nil
}

func (svc *Service) Stop() (err error) {
	svc.StopApi()
	return nil
}

func (svc *Service) initData() (err error) {
	total, err := svc.col.Count(nil)
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	// data to initialize
	settings := []NotificationSetting{
		{
			Id:          primitive.NewObjectID(),
			Type:        NotificationTypeMail,
			Enabled:     false,
			Name:        "Mail Notification",
			Description: "This is the default mail notification. You can edit it with your own settings",
			Mail:        NotificationSettingMail{
				//Server:         "",
				//Port:           0,
				//User:           "",
				//Password:       "",
				//SenderEmail:    "",
				//SenderIdentity: "",
				//Title:          "",
				//Template:       "",
				//To:             "",
				//Cc:             "",
			},
		},
		{
			Id:          primitive.NewObjectID(),
			Type:        NotificationTypeMobile,
			Enabled:     false,
			Name:        "Mobile Notification",
			Description: "This is the default mobile notification. You can edit it with your own settings",
			Mobile:      NotificationSettingMobile{
				//Webhook:  "",
				//Title:    "",
				//Template: "",
			},
		},
	}
	var data []interface{}
	for _, s := range settings {
		data = append(data, s)
	}
	_, err = svc.col.InsertMany(data)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) send(c *gin.Context) {
	var payload SendPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// models
	m, err := svc._getModels(payload)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	// setting
	s, err := svc._getSetting(m)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	switch s.Type {
	case NotificationTypeMail:
		err = svc.sendMail(s, m)
	case NotificationTypeMobile:
		err = svc.sendMobile(s, m)
	default:
		controllers.HandleErrorInternalServerError(c, errors.New(fmt.Sprintf("%s is not supported", s.Type)))
		return
	}

	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccess(c)
}

func (svc *Service) sendMail(s *NotificationSetting, m *models.ModelMap) (err error) {
	// TODO: implement
	//SendMail(s.Mail.To)

	return nil
}

func (svc *Service) sendMobile(s *NotificationSetting, m *models.ModelMap) (err error) {
	// title
	title := fmt.Sprintf("[Crawlab] \"%s\" 任务 %s", m.Spider.GetName(), m.Task.GetStatus())
	content := GetMobileTaskMarkdownContent(&m.Task, &m.Spider, &m.Node, &m.TaskStat)

	// webhook
	//webhook := "https://oapi.dingtalk.com/robot/send?access_token=7e08bf6f891b0ffc81fc91c4871f744251df5b143db242e271fb28c158a3176c"
	webhook := s.Mobile.Webhook

	// send
	if err := SendMobileNotification(webhook, title, content); err != nil {
		return err
	}

	return nil
}

func (svc *Service) getSettingList(c *gin.Context) {
	// params
	pagination := controllers.MustGetPagination(c)
	query := controllers.MustGetFilterQuery(c)
	sort := controllers.MustGetSortOption(c)

	// get list
	var list []NotificationSetting
	if err := svc.col.Find(query, &mongo2.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	}).All(&list); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			controllers.HandleSuccessWithListData(c, nil, 0)
		} else {
			controllers.HandleErrorInternalServerError(c, err)
		}
		return
	}

	// total count
	total, err := svc.col.Count(query)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccessWithListData(c, list, total)
}

func (svc *Service) getSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		controllers.HandleErrorBadRequest(c, err)
		return
	}

	var s NotificationSetting
	if err := svc.col.FindId(id).One(&s); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccessWithData(c, s)
}

func (svc *Service) putSetting(c *gin.Context) {
	var s NotificationSetting
	if err := c.ShouldBindJSON(&s); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	s.Id = primitive.NewObjectID()
	if _, err := svc.col.Insert(s); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccessWithData(c, s)
}

func (svc *Service) postSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		controllers.HandleErrorBadRequest(c, err)
		return
	}

	var s NotificationSetting
	if err := svc.col.FindId(id).One(&s); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	if err := c.ShouldBindJSON(&s); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}
	s.Id = id

	if err := svc.col.ReplaceId(id, s); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccessWithData(c, s)
}

func (svc *Service) deleteSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		controllers.HandleErrorBadRequest(c, err)
		return
	}

	if err := svc.col.DeleteId(id); err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccess(c)
}

func (svc *Service) enableSetting(c *gin.Context) {
	svc._toggleSettingFunc(true)(c)
}

func (svc *Service) disableSetting(c *gin.Context) {
	svc._toggleSettingFunc(false)(c)
}

func (svc *Service) handleEvents() {
	log.Infof("start handling events")

	// get stream
	var stream grpc.PluginService_SubscribeClient
	for {
		stream = svc.Internal.GetEventService().GetStream()
		if stream == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	for {
		// receive stream message
		msg, err := stream.Recv()

		if err != nil {
			// TODO: re-connect

			// end
			if err == io.EOF {
				log.Infof("received EOF signal, disconnecting")
				return
			}

			trace.PrintError(err)
			time.Sleep(1 * time.Second)
			continue
		}

		var data entity.GrpcEventServiceMessage
		switch msg.Code {
		case grpc.StreamMessageCode_SEND_EVENT:
			// data
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				return
			}
			if len(data.Events) < 1 {
				continue
			}

			// event name
			eventName := data.Events[0]

			// settings
			var settings []NotificationSetting
			if err := svc.col.Find(bson.M{
				"enabled": true,
			}, nil).All(&settings); err != nil {
				continue
			}

			// triggers
			tSet := hashset.New()
			for _, s := range settings {
				for _, t := range s.Triggers {
					tSet.Add(t.Event)
				}
			}

			// filter
			// TODO: performance concern
			if !tSet.Contains(eventName) {
				continue
			}

			// handle events
			arr := strings.Split(eventName, ":")
			switch arr[0] {
			case "model":
				if len(arr) < 2 {
					continue
				}
				colName := arr[1]
				action := arr[2]
				if err := svc._handleEventModel(colName, action, data.Data); err != nil {
					trace.PrintError(err)
				}
			}
		default:
			continue
		}
	}
}

func (svc *Service) _handleEventModel(colName, action string, data []byte) (err error) {
	m := models.NewModelMap()
	switch colName {
	case interfaces.ModelColNameTask:
		_ = json.Unmarshal(data, &m.Task)
		if err := svc._sendByTaskId(m.Task.GetId()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) _sendByTaskId(taskId primitive.ObjectID) (err error) {
	// models
	m, err := svc._getModelsByTaskId(taskId)
	if err != nil {
		return err
	}

	// setting
	s, err := svc._getSetting(m)
	if err != nil {
		return err
	}

	// send
	switch s.Type {
	case NotificationTypeMail:
		err = svc.sendMail(s, m)
	case NotificationTypeMobile:
		err = svc.sendMobile(s, m)
	default:
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) _toggleSettingFunc(value bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			controllers.HandleErrorBadRequest(c, err)
			return
		}
		var s NotificationSetting
		if err := svc.col.FindId(id).One(&s); err != nil {
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
		s.Enabled = value
		if err := svc.col.ReplaceId(id, s); err != nil {
			controllers.HandleErrorInternalServerError(c, err)
			return
		}
		controllers.HandleSuccess(c)
	}
}

func (svc *Service) _getModels(payload SendPayload) (m *models.ModelMap, err error) {
	return svc._getModelsByTaskId(payload.TaskId)
}

func (svc *Service) _getModelsByTaskId(taskId primitive.ObjectID) (m *models.ModelMap, err error) {
	m = models.NewModelMap()
	var doc interfaces.Model

	// task
	taskSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdTask)
	if err != nil {
		return nil, err
	}
	doc, err = taskSvc.GetById(taskId)
	if err != nil {
		return
	}
	m.Task = *doc.(*models.Task)

	// spider
	spiderSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdSpider)
	if err != nil {
		return nil, err
	}
	doc, err = spiderSvc.GetById(m.Task.GetSpiderId())
	if err != nil {
		return nil, err
	}
	m.Spider = *doc.(*models.Spider)

	// node
	nodeSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdNode)
	if err != nil {
		return nil, err
	}
	doc, err = nodeSvc.GetById(m.Task.GetNodeId())
	if err != nil {
		return nil, err
	}
	m.Node = *doc.(*models.Node)

	// task stat
	taskStatSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdTaskStat)
	if err != nil {
		return nil, err
	}
	doc, err = taskStatSvc.GetById(m.Task.GetId())
	if err != nil {
		return nil, err
	}
	m.TaskStat = *doc.(*models.TaskStat)

	// project
	if !m.Spider.ProjectId.IsZero() {
		projectSvc, err := svc.GetModelService().NewBaseServiceDelegate(interfaces.ModelIdProject)
		if err != nil {
			return nil, err
		}
		doc, err = projectSvc.GetById(m.Spider.ProjectId)
		if err != nil {
			return nil, err
		}
		m.Project = *doc.(*models.Project)
	}

	// TODO: user

	return m, nil
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

func (svc *Service) _getSetting(m *models.ModelMap) (setting *NotificationSetting, err error) {
	targets := []bson.M{
		{
			"_id":   m.Spider.Id,
			"model": interfaces.ModelColNameSpider,
		},
		{
			"_id":   m.Project.Id,
			"model": interfaces.ModelColNameProject,
		},
		//{
		//	"_id":   m.User.Id,
		//	"model": interfaces.ModelColNameUser,
		//},
		{
			"global": true,
		},
	}

	for _, target := range targets {
		s, err := svc._getSettingByTarget(target)
		if err != nil {
			return nil, err
		}
		if s == nil {
			continue
		}

		// found setting
		return s, nil
	}

	// not found
	return nil, nil
}

func (svc *Service) _getSettingByTarget(target bson.M) (setting *NotificationSetting, err error) {
	isGlobal := false
	res, ok := target["global"]
	if ok {
		isGlobal, _ = res.(bool)
	}

	var query bson.M
	if isGlobal {
		query = bson.M{"global": true}
	} else {
		query = bson.M{"targets": target}
		res, ok := target["_id"]
		if !ok {
			return nil, nil
		}
		_id, ok := res.(primitive.ObjectID)
		if !ok {
			return nil, nil
		}
		if _id.IsZero() {
			return nil, nil
		}
	}

	var s NotificationSetting
	if err := svc.col.Find(query, nil).One(&s); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (svc *Service) _getSettingsFromExtraValue(ev *models.ExtraValue) (settings *NotificationSetting) {
	var s NotificationSetting
	data, _ := json.Marshal(ev.Value)
	_ = json.Unmarshal(data, &s)
	return &s
}

func NewService() *Service {
	// service
	svc := &Service{
		Internal: plugin.NewInternal(),
		col:      mongo2.GetMongoCol(NotificationSettingsColName),
	}

	if err := svc.Init(); err != nil {
		panic(err)
	}

	return svc
}
