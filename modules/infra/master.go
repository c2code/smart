package infra

import (
	"fmt"
	"sync"
	"smart.com/weixin/smart/logp"
	"net/http"
	"time"
	"smart.com/weixin/smart/utils"
	"github.com/jinzhu/gorm"
	"smart.com/weixin/smart/modules/users"
	"github.com/gin-gonic/gin"
	"smart.com/weixin/smart/modules/courses"
)


type Manager struct {
	dispatchChan   chan WorkerRequest
	freeWorkerChan chan chan interface{}
	doneChan       chan bool
	config         Config
	logger         *logp.Logger
	waitgroup      sync.WaitGroup
	httpClient     *http.Client
}

// WorkerRequest request wrapper
type WorkerRequest struct {
	Type         string
	GUID         string
	GinContext   *gin.Context
	Body         interface{}
	RspChan      chan interface{}
	DoneChan     chan bool
}

var manager Manager

func getManagerObj() Manager {
	return manager
}

func Run() {
	//init master data
	config, _ := initConfig() //load config file tuna.json
	dispatchChan := make(chan WorkerRequest, 2000)
	freeWorkerChan := make(chan chan interface{}, config.MaxWorker)
	doneChan := make(chan bool)
	logger  := logp.NewLogger(ModuleName)
	fmt.Println("master begin to init config and data ")

	manager = Manager {
		config:           config,
		logger:           logger,
		dispatchChan:     dispatchChan,
		freeWorkerChan:   freeWorkerChan,
		doneChan:         doneChan,
	}

	manager.httpClient = &http.Client{Timeout: time.Second * 2}

	//to select a free worker  to handle task
	go manager.dispatch()

	for i := 0; i < config.MaxWorker; i++ {
		workerID := fmt.Sprintf("worker_%d", i)
		go manager.work(workerID)
	}

	// test 1 to 1
	db := utils.Init()
	manager.Migrate(db)
	defer db.Close()
	/*tx1 := db.Begin()
	userA := users.UserModel{
		Username: "AAAAAAAAAAAAAAAA",
		Email:    "aaaa@g.cn",
		Bio:      "hehddeda",
		Image:    nil,
	}
	tx1.Save(&userA)
	tx1.Commit()
	fmt.Println(userA)*/

	courses.InitCouses()

	//to liston to port 8088 by default
	manager.webListen()

	//set doneChan to close all tasks
	close(manager.doneChan)

	time.Sleep(time.Duration(5) * time.Second)
}

func (m Manager) dispatch() {
	logger := m.logger.Named("dispatch")

	var req WorkerRequest

	for {
		select {
		case req = <-m.dispatchChan: //receive a request
			logger.Debugf("recv req: %+v", req)
			select {
			case workerChan := <-m.freeWorkerChan: //find a free worker to handle the request
				workerChan <- req
			case <-m.doneChan:
				return
			}
		case <-m.doneChan:
			return
		}
	}
}

func (m Manager) Migrate(db *gorm.DB) {
	users.AutoMigrate()
	courses.AutoMigrate()
	/*db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{})*/
}
