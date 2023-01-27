// Copyright 2019 Iguazio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package locator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	prometheus.DefaultRegisterer.Register(collectors.NewGoCollector()) // nolint: errcheck
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
		pods, err := s.clientSet.CoreV1().Pods(s.config.Namespace).List(c,
			v1.ListOptions{
				LabelSelector: strings.Join(selector, ","),
			})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err) // nolint: errcheck
			return
		}

		// find the pod with the matching host IP
		// and return its pod IP
		for _, pod := range pods.Items {
			if pod.Status.HostIP == key {
				c.String(http.StatusOK, pod.Status.PodIP)
				return
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
