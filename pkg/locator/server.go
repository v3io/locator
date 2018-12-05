package locator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/gin-gonic/gin.v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type server struct {
	config    *Config
	engine    *gin.Engine
	clientSet *kubernetes.Clientset
}

func (s *server) initDefaults() {
	s.config.Defaults()
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.Default()
}

func (s *server) registerMetrics() {
	prometheus.DefaultRegisterer.Register(prometheus.NewGoCollector()) // nolint: errcheck
	s.engine.Any("/metrics", gin.WrapH(
		promhttp.HandlerFor(
			prometheus.DefaultGatherer, promhttp.HandlerOpts{})))
}

func (s *server) registerHandlers() {
	s.engine.GET("/locate/:app/:key", func(c *gin.Context) {
		app := c.Param("app")
		key := c.Param("key")
		query := c.Request.URL.Query()
		selector := []string{fmt.Sprintf("app=%s", app)}
		for key, value := range query {
			selector = append(selector, fmt.Sprintf("%s=%s", key, value))
		}
		pods, err := s.clientSet.CoreV1().Pods(s.config.Namespace).List(v1.ListOptions{
			LabelSelector: strings.Join(selector, ","),
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err) // nolint: errcheck
		} else {
			for _, pod := range pods.Items {
				if pod.Status.HostIP == key {
					c.String(http.StatusOK, pod.Status.PodIP)
					return
				}
			}
		}
		c.String(http.StatusNoContent, "")
	})
}

func (s *server) run() error {
	s.initDefaults()
	s.registerMetrics()
	s.registerHandlers()
	fmt.Printf("Starting registry server listening on %d for namepace %s\n", s.config.Port, s.config.Namespace)
	return s.engine.Run(fmt.Sprintf(":%d", s.config.Port))
}

func RunServer(config *Config) error {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return err
	}
	srv := &server{
		config:    config,
		clientSet: clientSet,
	}
	return srv.run()
}
