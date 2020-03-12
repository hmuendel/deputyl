package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/hmuendel/deputyl/depcheck"
	"github.com/hmuendel/deputyl/discovery"
	"github.com/hmuendel/deputyl/health"
	"github.com/hmuendel/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stevenroose/gonfig"
)

var (
	// VERSION is passed during build time to be later displayed
	// during application start
	VERSION = "⍺"
	// COMMIT is the git commit hash of the source code used building
	// this and passed during build
	COMMIT string
)

// The part of the nested configuration describing the http server properties
type serverConfig = struct {
	Address string `id:"address" desc:"ip address to listen on defaults to 0.0.0.0"`
	Port    int    `id:"port" desc:"tcp port to listen on defaults to 8080"`
}

// Default values for http server config
var defaultServerConfig = serverConfig{
	Address: "0.0.0.0",
	Port:    8080,
}

// The global main config which aggregates all configs from sub packages
var config = struct {
	Config    string            `id:"config" desc:"path to config file"`
	Version   bool              `short:"v" desc:"print the version and exits"`
	Log       *glog.LogConfig   `id:"log" desc:"configuration for verbose leveled logging"`
	Server    *serverConfig     `id:"server" desc:"http server config"`
	Discovery *discovery.Config `id:"discovery" desc:"discovery options"`
	Depckeck  *depcheck.Config  `id:"depcheck" desc:"options for upstream dependency checker"`
}{
	Log:       &glog.DefaultConfig,
	Server:    &defaultServerConfig,
	Discovery: &discovery.DefaultConfig,
	Depckeck:  &depcheck.DefaultConfig,
}

func main() {
	//
	err := gonfig.Load(&config, gonfig.Conf{
		ConfigFileVariable:  "config",
		FileDefaultFilename: "/config.yaml",
		FileDecoder:         gonfig.DecoderYAML,
		EnvPrefix:           "DPU_",
	})
	// cheapest exit with the version flag provided
	if config.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if err != nil {
		glog.Fatal(err)
	}
	// always logging startup without even considering log config
	glog.Infof("starting deputyl in version: %s, commit: %s", VERSION, COMMIT)
	// log config gets applied and is respected from here on
	//todo
	glog.Init(config.Log)
	if glog.V(8) {
		glog.Infof("config: %#v \n", config)
		glog.Infof("log config: %#v \n", config.Log)
		glog.Infof("server config: %#v \n", config.Server)
		glog.Infof("discovery config: %#v \n", config.Discovery)
		glog.Infof("depcheck config: %#v \n", config.Depckeck)
	}
	if err != nil {
		fmt.Printf("error getting proc: %s", err)
		os.Exit(1)
	}

	// configuring http server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", health.HandleHealth)
	address := config.Server.Address + ":" + strconv.Itoa(config.Server.Port)
	server := &http.Server{Addr: address, Handler: mux}

	// staring disvovery and web server
	if glog.V(3) {
		glog.Infof("starting  discovery")
	}
	d := discovery.NewDiscoverEmitter()

	d.StartDiscovery()
	if glog.V(3) {
		glog.Infof("starting http server on %s", address)
	}
	err = server.ListenAndServe()
	glog.Error(err)
}
