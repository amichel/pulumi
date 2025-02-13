package test

import (
	"github.com/blang/semver"
	"github.com/pulumi/pulumi/pkg/resource/deploy/deploytest"
	"github.com/pulumi/pulumi/sdk/go/common/resource/plugin"
)

func NewHost(schemaDirectoryPath string) plugin.Host {
	return deploytest.NewPluginHost(nil, nil, nil,
		deploytest.NewProviderLoader("aws", semver.MustParse("1.0.0"), func() (plugin.Provider, error) {
			return AWS(schemaDirectoryPath)
		}))
}
