package depcheck_test

import (
	dc "github.com/hmuendel/deputyl/depcheck"
	"testing"
)

type mockHubClient struct {
	err             error
	response        []string
	requestImage    string
	requestUrl      string
	requestUsername string
	requestPassword string
}

func (mc *mockHubClient) Tags(image, url, username, password string) ([]string, error) {
	mc.requestImage = image
	mc.requestUrl = url
	mc.requestUsername = username
	mc.requestPassword = password
	return mc.response, mc.err
}

type shResult struct {
	err   bool
	url   string
	tag   string
	image string
}

func TestImageStingHandling(t *testing.T) {
	testsCases := [...]struct {
		name   string
		image  string
		result shResult
	}{
		{"dockerhub-library", "nginx",
			shResult{false, "https://" + dc.DEFAULT_DOCKER_URL, dc.DEFAULT_TAG, dc.DEFAULT_LIBRARY_PREFIX + "/nginx"}},
		{"dockerhub-library-tag", "nginx:tag",
			shResult{false, "https://" + dc.DEFAULT_DOCKER_URL, "tag", dc.DEFAULT_LIBRARY_PREFIX + "/nginx"}},
		{"dockerhub-library-2tags", "nginx:tag:2", shResult{true, "", "", ""}},
		{"dockerhub-non-library", "concourse/concourse",
			shResult{false, "https://" + dc.DEFAULT_DOCKER_URL, dc.DEFAULT_TAG, "concourse/concourse"}},
		{"dockerhub-non-library-tag", "concourse/concourse:tag",
			shResult{false, "https://" + dc.DEFAULT_DOCKER_URL, "tag", "concourse/concourse"}},
		{"dockerhub-non-library-2tags", "concourse/concourse:tag:2",
			shResult{true, "", "", ""}},
		{"gcr", "gcr.io/kubernetes-helm/tiller",
			shResult{false, "https://gcr.io", dc.DEFAULT_TAG, "kubernetes-helm/tiller"}},
		{"gcr-tag", "gcr.io/kubernetes-helm/tiller:tag",
			shResult{false, "https://gcr.io", "tag", "kubernetes-helm/tiller"}},
		{"gcr-2tags", "gcr.io/kubernetes-helm/tiller:tag:2",
			shResult{true, "", "", ""}},
		{"docker.io", "docker.io/kubernetes-helm/tiller:tag",
			shResult{false, "https://" + dc.DefaultConfig.DefaultDockerUrl, "tag", "kubernetes-helm/tiller"}},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			checker := dc.NewDockerDepChecker()
			mc := mockHubClient{}
			checker.HubClient = &mc
			tag, _, err := checker.Check(tc.image)
			if err == nil && tc.result.err {
				t.Errorf("%s, should have errored", tc.name)
			}
			if err != nil && !tc.result.err {
				t.Errorf("%s, should not errored", tc.name)
			}
			if mc.requestImage != tc.result.image {
				t.Errorf("expected %s, got %s", tc.result.image, mc.requestImage)
			}
			if mc.requestUrl != tc.result.url {
				t.Errorf("expected %s, got %s", tc.result.url, mc.requestUrl)
			}
			if tag != tc.result.tag {
				t.Errorf("expected %s, got %s", tc.result.tag, tag)
			}
		})
	}
}
