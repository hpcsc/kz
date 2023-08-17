package gateway

import "github.com/stretchr/testify/mock"

var _ Github = (*fakeGithub)(nil)

type fakeGithub struct {
	mock.Mock
}

func NewFakeGithubGateway() *fakeGithub {
	return &fakeGithub{}
}

func (a *fakeGithub) LatestRelease() (*LatestReleaseResponse, error) {
	args := a.Called()
	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}

	return r.(*LatestReleaseResponse), args.Error(1)
}

func (a *fakeGithub) Download(url string) ([]byte, error) {
	args := a.Called(url)
	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}

	return r.([]byte), args.Error(1)
}

func (a *fakeGithub) OnLatestRelease() *mock.Call {
	return a.On("LatestRelease")
}

func (a *fakeGithub) OnDownload(url interface{}) *mock.Call {
	return a.On("Download", url)
}
