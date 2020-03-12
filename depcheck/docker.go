package depcheck

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"encoding/json"
	"github.com/hmuendel/docker-registry-client/registry"
	"github.com/hmuendel/glog"
)

// The default registry url used for images like concourse/concourse
const DEFAULT_DOCKER_URL = "registry.hub.docker.com"

// Prefix of library images like e.g. just nginx
const DEFAULT_LIBRARY_PREFIX = "library"

// The docker tag used if no tag is provided
const DEFAULT_TAG = "latest"

type HubClient interface {
	Tags(image, url, username, password string) ([]string, error)
}

type Config struct {
	DefaultDockerUrl     string `desc:"the default registry url used for images like concourse/concourse"`
	DefaultLibraryPrefix string `desc:"prefix of library images like e.g. just nginx"`
	DefaultTag           string `desc:"the docker tag used if no tag is provided"`
}

var DefaultConfig = Config{
	DefaultDockerUrl:     DEFAULT_DOCKER_URL,
	DefaultLibraryPrefix: DEFAULT_LIBRARY_PREFIX,
	DefaultTag:           DEFAULT_TAG,
}

// Holding usrname and password for docker registries
type DockerCredential struct {
	Username string
	Password string
}

func splitDockerString(image string) (hubUrl, name, tag string, err error) {
	tagSplit := strings.Split(image, ":")
	// docker image should only contain one colon
	if len(tagSplit) > 2 {
		return "", "", "", errors.New("could not extract tag")
	}
	// if no tag is given, use default tag
	if len(tagSplit) == 1 {
		tag = DefaultConfig.DefaultTag
	}
	if len(tagSplit) == 2 {
		tag = tagSplit[1]
	}
	imgSplit := strings.Split(tagSplit[0], "/")
	// no dot, no registry domain, using default registry
	if !strings.Contains(tagSplit[0], ".") {
		hubUrl = DefaultConfig.DefaultDockerUrl
		// library image with no slashes, append prefix
		if len(imgSplit) == 1 {
			name = path.Join(DefaultConfig.DefaultLibraryPrefix, imgSplit[0])
			return
		} else {
			name = strings.Join(imgSplit, "/")
			return
		}
	}
	// registry url without a slash is invalid
	if len(imgSplit) == 1 {
		return "", "", "", errors.New("invalid docker image")
	}
	hubUrl = imgSplit[0]
	name = strings.Join(imgSplit[1:], "/")
	return
}

type HerokuClient struct {
}

func (hc *HerokuClient) Tags(image, url, username, password string) ([]string, error) {
	hub, err := registry.New(url, username, password, false)
	if err != nil {
		return nil, err
	}
	hub.Logf = glog.V(5).Infof
	tags, err := hub.Tags(image)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

type SimpleClient struct {
}

func (sc *SimpleClient) Tags(image, url, username, password string) ([]string, error) {
	var regResponse struct {
		Errors []struct {
			Code    string
			Message string
		}
		Name string
		Tags []string
	}
	if glog.V(9) {
		glog.Infof(
			"simpleClients tag function called with image %s, url: %s, username: %s, pw: *******",
			image, url, username)
	}
	tagsUrl := url + "/v2/" + image + "/tags" + "/list"
	if glog.V(4) {
		glog.Infof("making request to %s", tagsUrl)
	}
	resp, err := http.Get(tagsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&regResponse); err != nil {
		return nil, err
	}
	if glog.V(10) {
		glog.Infof("%+v", regResponse)
	}
	return regResponse.Tags, nil

}

//
type DockerDepChecker struct {
	credentials map[string]DockerCredential
	HubClient   HubClient
}

func NewDockerDepChecker() *DockerDepChecker {
	return &DockerDepChecker{
		HubClient: &HerokuClient{},
	}
}

func (c *DockerDepChecker) AddCredentials(hubUrl string, creds DockerCredential) {
	c.credentials[hubUrl] = creds
}

func (c *DockerDepChecker) Check(imageString string) (string, []string, error) {
	hubUrl, image, tag, err := splitDockerString(imageString)
	if err != nil {
		return "", nil, err
	}
	username := ""
	password := ""
	if creds, ok := c.credentials[hubUrl]; ok {
		username = creds.Username
		password = creds.Password
	}
	url := "https://" + hubUrl
	tags, err := c.HubClient.Tags(image, url, username, password)
	if err != nil {
		return "", nil, err
	}
	return tag, tags, nil
}
