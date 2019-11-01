package api

import ("fmt"
	"github.com/kataras/iris"
	"github.com/lishimeng/jimeng/internal/etc"
	"github.com/lishimeng/jimeng/internal/monitor"
)

type WebOptions struct {
	Listen string
}

type WebEngine struct {
	app *iris.Application
	options *WebOptions
}

type WebServer interface {
	Start()
}

type welcome struct {
	Name string `json:"appName"`
	Author string `json:"author"`
	Version string `json:"version"`
	DataStatistics int64 `json:"dataStatistics"`
	StartTime string `json:"startTime"`
}

func Create(options *WebOptions) *WebEngine {
	app := iris.New()
	engine := WebEngine{}
	engine.app = app
	engine.options = options
	return &engine
}


func initialize(app *iris.Application) () {
	fmt.Printf("Web server start!")
	app.Get("/", func(cxt iris.Context) {
		appInfo := welcome{
			Name: etc.Config.Name,
			Version: etc.Config.Version,
			Author: "lism",
			DataStatistics: monitor.D,
			StartTime: fmt.Sprintf("%s", monitor.StartTime),
		}
		_, _ = cxt.JSON(appInfo)
	})
}

func (engine *WebEngine) Start() () {
	initialize(engine.app)
	_ = engine.app.Run(iris.Addr(engine.options.Listen))
}