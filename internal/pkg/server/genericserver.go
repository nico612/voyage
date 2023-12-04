package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/pkg/middleware"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/log"
	"github.com/nico612/voyage/pkg/version"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"
	"time"
)

// GenericAPIServer 包含一个 adminsrv API服务器的状态。包含基础的中间件和api、close 和 ping 功能
// 类型为 GenericAPIServer gin.Engine。
type GenericAPIServer struct {
	middlewares []string // 使用的中间件名

	// SecureServingInfo TLS 服务器的配置信息。
	SecureServingInfo *SecureServingInfo

	// InsecureServingInfo 保存不安全的 HTTP 服务器的配置信息。
	InsecureServingInfo *InsecureServingInfo

	// ShutdownTimeout 优雅关停服务的超时时间
	ShutdownTimeout time.Duration

	*gin.Engine
	healthz         bool
	enableMetrics   bool
	enableProfiling bool
	// wrapper for gin.Engine

	insecureServer, secureServer *http.Server
}

func initGenericAPIServer(s *GenericAPIServer) {
	// do some setup
	// s.GET(path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

func (s *GenericAPIServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares install generic middlewares
func (s *GenericAPIServer) InstallMiddlewares() {

	// 使用 request ID 中间件
	s.Use(gin.Recovery())
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	// install custom middlewares
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}
		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}
}

// InstallAPIs install generic apis
func (s *GenericAPIServer) InstallAPIs() {

	// install healthz handler
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			response.Success(c, map[string]string{"status": "ok"})
		})
	}

	// install metric handler
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// install pprof handler
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	s.GET("/version", func(c *gin.Context) {
		response.Success(c, version.Get())
	})
}

func (s *GenericAPIServer) Run() error {

	// 为了可扩展性，在这里使用自定义的HTTP配置模式。
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	//// https 服务
	//s.secureServer = &http.Server{
	//	Addr:    s.secureServer.Addr,
	//	Handler: s,
	//	//ReadTimeout:    10 * time.Second,
	//	//WriteTimeout:   10 * time.Second,
	//	//MaxHeaderBytes: 1 << 20,
	//}

	var eg errgroup.Group

	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)
		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	// 启动 https 服务
	//eg.Go(func() errors {
	//	key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
	//	if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
	//		return nil
	//	}
	//	log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())
	//
	//	if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		log.Fatal(err.Error())
	//
	//		return err
	//	}
	//	log.Infof("Server on %s stopped", s.SecureServingInfo.Address())
	//	return nil
	//})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

// Close 优雅关停服务
func (s *GenericAPIServer) Close() {

	// 关闭超时时间
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	// 关闭 HTTP 服务
	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}

	//// 关闭 HTTPS 服务
	//if err := s.secureServer.Shutdown(ctx); err != nil {
	//	log.Warnf("Shutdown secure server failed: %s", err.Error())
	//}

}

// ping pings the http server to make sure the router is working.
func (s *GenericAPIServer) ping(ctx context.Context) error {

	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)

	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		// Change NewRequest to NewRequestWithContext and pass context it
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// Ping the server by sending a GET request to `/healthz`.

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}
