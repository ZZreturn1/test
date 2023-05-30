package web

import (
	"context" // 上下文
	"crypto/tls" // 加密传输层安全性协议
	"embed" // 嵌入静态文件
	"html/template" // HTML模板
	"io" // 输入输出
	"io/fs" // 文件系统
	"net" // 网络操作
	"net/http" // HTTP协议
	"os" // 操作系统函数
	"strconv" // 字符串转换
	"strings" // 字符串操作
	"time" // 时间
	"x-ui/config" // 配置
	"x-ui/logger" // 日志
	"x-ui/util/common" // 通用工具
	"x-ui/web/controller" // 控制器
	"x-ui/web/job" // 任务
	"x-ui/web/network" // 网络
	"x-ui/web/service" // 服务

	"github.com/BurntSushi/toml" // 解析TOML配置文件
	"github.com/gin-contrib/sessions" // Gin框架的会话管理中间件
	"github.com/gin-contrib/sessions/cookie" // Gin框架的基于cookie的会话存储
	"github.com/gin-gonic/gin" // Gin框架
	"github.com/nicksnyder/go-i18n/v2/i18n" // 国际化和本地化
	"github.com/robfig/cron/v3" // 定时任务
	"golang.org/x/text/language" // 语言标签
)

// 嵌入静态文件系统
var assetsFS embed.FS

// 嵌入HTML文件系统
var htmlFS embed.FS

// 嵌入国际化翻译文件系统
var i18nFS embed.FS

// 服务器启动时间
var startTime = time.Now()

// 封装静态文件系统
type wrapAssetsFS struct {
	embed.FS
}

// 打开文件
func (f *wrapAssetsFS) Open(name string) (fs.File, error) {
	// 打开嵌入的静态文件
	file, err := f.FS.Open("assets/" + name)
	if err != nil {
		return nil, err
	}
	return &wrapAssetsFile{
		File: file,
	}, nil
}

// 封装文件
type wrapAssetsFile struct {
	fs.File
}

// 获取文件信息
func (f *wrapAssetsFile) Stat() (fs.FileInfo, error) {
	info, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return &wrapAssetsFileInfo{
		FileInfo: info,
	}, nil
}

// 封装文件信息
type wrapAssetsFileInfo struct {
	fs.FileInfo
}

/// 获取修改时间
func (f *wrapAssetsFileInfo) ModTime() time.Time {
	// 返回启动时间
	return startTime
}

// 服务器结构体
type Server struct {
	httpServer *http.Server     // HTTP服务器
	listener   net.Listener     // 网络监听器

	index  *controller.IndexController	// 首页控制器
	server     *controller.ServerController    // 服务器控制器
	xui        *controller.XUIController       // XUI控制器

	xrayService    service.XrayService        // Xray服务
	settingService service.SettingService     // 设置服务
	inboundService service.InboundService     // 入站服务

	cron           *cron.Cron                // 定时任务

	ctx            context.Context           // 上下文
	cancel         context.CancelFunc        // 取消上下文的函数
}

// 创建新的服务器实例
func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		ctx:    ctx,
		cancel: cancel,
	}
}

// 获取HTML文件列表
func (s *Server) getHtmlFiles() ([]string, error) {
	files := make([]string, 0)
	dir, _ := os.Getwd()
	err := fs.WalkDir(os.DirFS(dir), "web/html", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 返回错误
		}
		if d.IsDir() {
			return nil // 忽略目录
		}
		files = append(files, path) // 添加文件路径到列表
		return nil
	})
	if err != nil {
		return nil, err // 返回错误
	}
	return files, nil // 返回文件列表
}

// 获取HTML模板
func (s *Server) getHtmlTemplate(funcMap template.FuncMap) (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	err := fs.WalkDir(htmlFS, "html", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 返回错误
		}
		if d.IsDir() {
			newT, err := t.ParseFS(htmlFS, path+"/*.html") // 解析HTML模板文件
			if err != nil {
				// 忽略错误
				return nil
			}
			t = newT
		}
		return nil // 返回nil
	})
	if err != nil {
		return nil, err // 返回错误
	}
	return t, nil // 返回HTML模板
}

// 初始化路由
func (s *Server) initRouter() (*gin.Engine, error) {
	if config.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	secret, err := s.settingService.GetSecret()
	if err != nil {
		return nil, err // 返回错误
	}

	basePath, err := s.settingService.GetBasePath()
	if err != nil {
		return nil, err // 返回错误
	}

	assetsBasePath := basePath + "assets/"

	store := cookie.NewStore(secret) // 基于Cookie的会话存储
	engine.Use(sessions.Sessions("session", store)) // 使用会话中间件
	engine.Use(func(c *gin.Context) {
		c.Set("base_path", basePath) // 设置基础路径
	})
	//
	engine.Use(func(c *gin.Context) {
		uri := c.Request.RequestURI
		if strings.HasPrefix(uri, assetsBasePath) {
			c.Header("Cache-Control", "max-age=31536000") // 设置缓存控制头
		}
	})
	err = s.initI18n(engine)
	if err != nil {
		return nil, err // 返回错误
	}

	if config.IsDebug() {
		files, err := s.getHtmlFiles()
		if err != nil {
			return nil, err // 返回错误
		}
		engine.LoadHTMLFiles(files...) // 加载HTML文件
		engine.StaticFS(basePath+"assets", http.FS(os.DirFS("web/assets"))) // 设置静态文件系统
	} else {
		t, err := s.getHtmlTemplate(engine.FuncMap)
		if err != nil {
			return nil, err // 返回错误
		}
		engine.SetHTMLTemplate(t) // 设置HTML模板
		engine.StaticFS(basePath+"assets", http.FS(&wrapAssetsFS{FS: assetsFS})) // 设置静态文件系统
	}

	g := engine.Group(basePath) // 创建路由组

	s.index = controller.NewIndexController(g)   // 创建首页控制器
	s.server = controller.NewServerController(g) // 创建服务器控制器
	s.xui = controller.NewXUIController(g)       // 创建XUI控制器

	return engine, nil // 返回Gin引擎和nil错误
}

// 初始化国际化配置
func (s *Server) initI18n(engine *gin.Engine) error {
	// 创建语言包
	bundle := i18n.NewBundle(language.SimplifiedChinese) // 创建国际化Bundle实例
	// 注册解析函数
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 遍历翻译文件并解析
	err := fs.WalkDir(i18nFS, "translation", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := i18nFS.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = bundle.ParseMessageFileBytes(data, path)
		return err
	})
	if err != nil {
		return err
	}
	// 查找国际化参数名
	findI18nParamNames := func(key string) []string {
		names := make([]string, 0)	// 参数名列表
		keyLen := len(key)		// 键长度

		for i := 0; i < keyLen-1; i++ {
			if key[i:i+2] == "{{" { // 判断开头 "{{"
				j := i + 2
				isFind := false
				for ; j < keyLen-1; j++ {
					if key[j:j+2] == "}}" { // 结尾 "}}"
						isFind = true
						break
					}
				}
				if isFind {
					names = append(names, key[i+3:j])
				}
			}
		}
		return names
	}

	var localizer *i18n.Localizer		// 本地化器

	// 注册i18n函数到模板引擎
	engine.FuncMap["i18n"] = func(key string, params ...string) (string, error) {
		names := findI18nParamNames(key)		// 获取参数名列表
		if len(names) != len(params) {
			return "", common.NewError("find names:", names, "---------- params:", params, "---------- num not equal")
		}
		// 模板数据
		templateData := map[string]interface{}{}

		for i := range names {
			templateData[names[i]] = params[i]		// 绑定参数名和参数值
		}
		return localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    key,		// 键
			TemplateData: templateData,		// 模板数据
		})
	}

	// 添加中间件，设置语言环境
	engine.Use(func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")	// 获取语言标识
		localizer = i18n.NewLocalizer(bundle, accept)	// 创建本地化器
		c.Set("localizer", localizer)	// 设置本地化器到上下文
		c.Next()	// 执行下一个处理程序
	})

	return nil
}

// 启动任务
func (s *Server) startTask() {
	err := s.xrayService.RestartXray(true)
	if err != nil {
		logger.Warning("启动 xray 失败:", err)
	}
	// 每 30 秒检查一次 xray 是否在运行
	s.cron.AddJob("@every 30s", job.NewCheckXrayRunningJob())

	// 启动 goroutine，延迟 5 秒后每 10 秒统计一次流量
	go func() {
		time.Sleep(time.Second * 5)
		// 每 10 秒统计一次流量，首次启动延迟 5 秒，与重启 xray 的时间错开
		s.cron.AddJob("@every 10s", job.NewXrayTrafficJob())
	}()

	// 每 30 秒检查一次 inbound 流量超出和到期的情况
	s.cron.AddJob("@every 30s", job.NewCheckInboundJob())
	// 每一天提示一次流量情况,上海时间8点30
	var entry cron.EntryID

	isTgbotenabled, err := s.settingService.GetTgbotenabled()

	if (err == nil) && (isTgbotenabled) {
		runtime, err := s.settingService.GetTgbotRuntime()
		if err != nil || runtime == "" {
			// 运行时间无效时，使用默认值
			logger.Errorf("添加 NewStatsNotifyJob 错误[%s]，运行时间[%s]无效，将使用默认值", err, runtime)
			runtime = "@daily"
		}
		logger.Infof("启用 Tg 通知，运行时间：%s", runtime)
		entry, err = s.cron.AddJob(runtime, job.NewStatsNotifyJob())
		if err != nil {
			logger.Warning("添加 NewStatsNotifyJob 错误", err)
			return
		}
	} else {
		s.cron.Remove(entry)	// 移除任务
	}
}

// 启动服务器
func (s *Server) Start() (err error) {
	// 这是一个匿名函数，没没有函数名，延迟执行的函数，在返回错误时停止服务器
	defer func() {
		if err != nil {
			s.Stop()
		}
	}()

	// 获取时区
	loc, err := s.settingService.GetTimeLocation()
	if err != nil {
		return err
	}
	// 创建定时任务调度器
	s.cron = cron.New(cron.WithLocation(loc), cron.WithSeconds())
	s.cron.Start()

	// 初始化路由
	engine, err := s.initRouter()
	if err != nil {
		return err
	}

	// 获取证书文件路径和密钥文件路径
	certFile, err := s.settingService.GetCertFile()
	if err != nil {
		return err
	}
	keyFile, err := s.settingService.GetKeyFile()
	if err != nil {
		return err
	}

	// 获取监听地址和端口
	listen, err := s.settingService.GetListen()
	if err != nil {
		return err
	}
	port, err := s.settingService.GetPort()
	if err != nil {
		return err
	}
	listenAddr := net.JoinHostPort(listen, strconv.Itoa(port))
	
	// 监听网络连接
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	// 如果指定了证书文件或密钥文件，则启用 HTTPS
	if certFile != "" || keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			listener.Close()
			return err
		}
		// 创建 TLS 配置
		c := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		// 创建自动 HTTPS 监听器
		listener = network.NewAutoHttpsListener(listener)
		listener = tls.NewListener(listener, c)
	}

	if certFile != "" || keyFile != "" {
		logger.Info("Web 服务器运行在 HTTPS 模式，监听地址为", listener.Addr())
	} else {
		logger.Info("Web 服务器运行在 HTTP 模式，监听地址为", listener.Addr())
	}

	
	s.listener = listener

	s.startTask()
	
	// 设置 HTTP 服务器配置
	s.httpServer = &http.Server{
		Handler: engine,
	}

	// 在 goroutine 中启动 HTTP 服务器
	go func() {
		s.httpServer.Serve(listener)
	}()

	return nil
}

// 停止服务器
func (s *Server) Stop() error {
	s.cancel()		// 取消上下文
	s.xrayService.StopXray()	// 停止 Xray 服务

	// 停止定时任务调度器
	if s.cron != nil {
		s.cron.Stop()
	}
	var err1 error
	var err2 error

	// 关闭 HTTP 服务器
	if s.httpServer != nil {
		err1 = s.httpServer.Shutdown(s.ctx)
	}

	// 关闭监听器
	if s.listener != nil {
		err2 = s.listener.Close()
	}
	return common.Combine(err1, err2)
}

// 获取上下文
func (s *Server) GetCtx() context.Context {
	return s.ctx
}

// 获取定时任务调度器
func (s *Server) GetCron() *cron.Cron {
	return s.cron
}