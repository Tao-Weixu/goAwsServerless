package tailf

import(
	"github.com/hpcloud/tail"
	"github.com/astaxie/beego/logs"
	//"fmt"
	"time"
	"sync"
)
const(
	StatusNormal = 1
	StatusDelete = 2
)
type CollectConf struct {
	LogPath string
	Topic string
}
type TextMsg struct {
	Msg string
	Topic string
}
type TailObj struct{
	tail *tail.Tail
	confOne CollectConf
	status int
	exitChan chan int
}
type TailObjMgr struct{ //gérer tous les objects du système
	tailObjs []*TailObj
	msgChan chan *TextMsg //channel envoie le message à la module "main", c'est module "main" 驱动的程序
	lock sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)
func UpdateConfigEtcd(confArry []CollectConf)(err error){
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()

	for _, obtainConf := range confArry{
		var isRunning = false
		for _, obj :=range tailObjMgr.tailObjs{
			if obtainConf.LogPath == obj.confOne.LogPath {
				isRunning = true
				obj.status = StatusNormal
				break
			}
		}
		if isRunning {
			continue
		}
		createNewTask(obtainConf)
	}
	var tailObjsSelect  []*TailObj
	for _, obj :=range tailObjMgr.tailObjs{
		obj.status  = StatusDelete
		for _, obtainConf := range confArry{
			if obtainConf.LogPath == obj.confOne.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue	
		}
		tailObjsSelect = append(tailObjsSelect, obj)
	}
	tailObjMgr.tailObjs = tailObjsSelect
	return
}

func createNewTask(obtainConf CollectConf)(err error){
	obj := &TailObj{
		confOne: obtainConf,	
		exitChan: make(chan int, 1),
	}
	tailInit, err:= tail.TailFile(obtainConf.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,	
		//Location: &tail.SeekInfo{Offset: 0, Whence:2},
		MustExist: false,
		Poll: true,
	})
	if err !=nil{
		logs.Error("collect filename[%s] failed, err:%v", obtainConf.LogPath, err)		
	}
	obj.tail = tailInit
	tailObjMgr.tailObjs =append(tailObjMgr.tailObjs, obj)

	go readFromTail(obj)
	return
}

func InitTail(confArray []CollectConf, chanSize int)(err error){
	tailObjMgr = &TailObjMgr{//il faut mettre cette partie en premier
		msgChan: make(chan *TextMsg, chanSize),
	}
	if len(confArray)==0{
		logs.Error("invalid config for log collect, confArray:%v", confArray)
		return
	}
	for _, unConf :=range confArray{
		createNewTask(unConf)
	}

	return
}
func readFromTail(tailObj *TailObj){

	for true{
		select {
		case line, ok :=<- tailObj.tail.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename: %s\n", tailObj.tail.Filename)
				time.Sleep(100*time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Msg: line.Text,
				Topic: tailObj.confOne.Topic,
			}
			tailObjMgr.msgChan<-textMsg
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited, conf: %v", tailObj.confOne)
			return
		}
	}
}
func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}
