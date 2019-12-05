package main

import(
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"fmt"

)
func convertLogLevel(level string) int {

	switch (level) {
		case "debug":
			return logs.LevelDebug
		case "warn" : 
			return logs.LevelWarn
		case "info" : 
			return logs.LevelInfo
		case "trace":
			return logs.LevelTrace
		}
	return logs.LevelDebug
}

func initLogger() (err error){
	config :=make(map[string]interface{})
	config["filename"] = appConfig.logPath
	config["level"]= convertLogLevel(appConfig.logLevel)

	configStr, err :=json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	logs.Debug("initialize logger succ")
	return

}