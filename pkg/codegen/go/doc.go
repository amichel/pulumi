// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Pulling out some of the repeated strings tokens into constants would harm readability, so we just ignore the
// goconst linter's warning.
//
// nolint: lll, goconst
package gen

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi/pkg/codegen"
	"github.com/pulumi/pulumi/pkg/codegen/schema"
)

// DocLanguageHelper is the Go-specific implementation of the DocLanguageHelper.
type DocLanguageHelper struct {
	packages map[string]*pkgContext
}

var _ codegen.DocLanguageHelper = DocLanguageHelper{}

// GetDocLinkForResourceType returns the godoc URL for a type belonging to a resource provider.
func (d DocLanguageHelper) GetDocLinkForResourceType(packageName string, moduleName string, typeName string) string {
	var path string
	if packageName != "" {
		path = fmt.Sprintf("%s/%s", packageName, moduleName)
	} else {
		path = moduleName
	}
	typeNameParts := strings.Split(typeName, ".")
	typeName = typeNameParts[len(typeNameParts)-1]
	typeName = strings.TrimPrefix(typeName, "*")

	if packageName != "" {
		return fmt.Sprintf("https://pkg.go.dev/github.com/pulumi/pulumi-%s/sdk/go/%s?tab=doc#%s", packageName, path, typeName)
	}
	return fmt.Sprintf("https://pkg.go.dev/github.com/pulumi/pulumi/sdk/go/%s?tab=doc#%s", path, typeName)
}

// GetDocLinkForResourceInputOrOutputType returns the godoc URL for an input or output type.
func (d DocLanguageHelper) GetDocLinkForResourceInputOrOutputType(packageName, moduleName, typeName string, input bool) string {
	link := d.GetDocLinkForResourceType(packageName, moduleName, typeName)
	if !input {
		return link + "Output"
	}
	return link + "Args"
}

// GetDocLinkForFunctionInputOrOutputType returns the doc link for an input or output type of a Function.
func (d DocLanguageHelper) GetDocLinkForFunctionInputOrOutputType(packageName, moduleName, typeName string, input bool) string {
	link := d.GetDocLinkForResourceType(packageName, moduleName, typeName)
	if !input {
		return link
	}
	return link + "Args"
}

// GetDocLinkForBuiltInType returns the godoc URL for a built-in type.
func GetDocLinkForBuiltInType(typeName string) string {
	return fmt.Sprintf("https://golang.org/pkg/builtin/#%s", typeName)
}

// GetLanguageTypeString returns the Go-specific type given a Pulumi schema type.
func (d DocLanguageHelper) GetLanguageTypeString(pkg *schema.Package, moduleName string, t schema.Type, input, optional bool) string {
	if moduleName == "" && pkg.Name == "kubernetes" {
		moduleName = "providers"
	}
	modPkg, ok := d.packages[moduleName]
	if !ok {
		glog.Errorf("cannot calculate type string for type %q. could not find a package for module %q", t.String(), moduleName)
		os.Exit(1)
	}
	return modPkg.plainType(t, optional)
}

// GeneratePackagesMap generates a map of Go packages for resources, functions and types.
func (d *DocLanguageHelper) GeneratePackagesMap(pkg *schema.Package, tool string, goInfo GoInfo) {
	d.packages = generatePackageContextMap(tool, pkg, goInfo)
}

// GetPropertyName returns the property name specific to Go.
func (d DocLanguageHelper) GetPropertyName(p *schema.Property) (string, error) {
	return strings.Title(p.Name), nil
}

// GetResourceFunctionResultName returns the name of the result type when a function is used to lookup
// an existing resource.
func (d DocLanguageHelper) GetResourceFunctionResultName(resourceName string) string {
	return "Lookup" + resourceName + "Result"
}
