// Copyright 2016-2023, Pulumi Corporation.
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
// SPDX-License-Identifier: MPL-2.0
// This file has been modified by Beijing Volcano Engine Technology Co., Ltd. on 2025

package volcenginecc

import (
	// embed is used to store bridge-metadata.json in the compiled binary
	"context"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/info"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/volcengine/pulumi-volcenginecc/provider/cmd/pulumi-resource-volcenginecc/token"
	"github.com/volcengine/pulumi-volcenginecc/provider/pkg/version"
	"github.com/volcengine/terraform-provider-volcenginecc/shim"

	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

// all of the random token components used below.
const (
	volcengeccpkg = "volcenginecc"
)

//go:embed cmd/pulumi-resource-volcenginecc/bridge-metadata.json
var metadata []byte

// Provider returns additional overlaid schema and metadata associated with the random package.
func Provider() tfbridge.ProviderInfo {
	prov := tfbridge.ProviderInfo{
		P:           pf.ShimProvider(shim.NewProvider()),
		Name:        "volcenginecc",
		Description: "A Pulumi package to safely use volcengine resource in Pulumi programs.",
		Keywords:    []string{"volcengine", "volcenginecc", "category/cloud"},
		License:     "MPL-2.0",
		Homepage:    "https://github.com/volcengine/pulumi-volcenginecc",
		Publisher:   "Volcengine",
		LogoURL:     "https://avatars.githubusercontent.com/u/67365215",
		Repository:  "https://github.com/volcengine/pulumi-volcenginecc",
		Version:     version.Version,
		// PluginDownloadURL is an optional URL used to download the Provider
		// for use in Pulumi programs
		PluginDownloadURL: "github://api.github.com/volcengine",
		MetadataInfo:      tfbridge.NewProviderMetadata(metadata),
		Resources:         map[string]*tfbridge.ResourceInfo{},
		GitHubHost:        "github.com",
		GitHubOrg:         "volcengine",
		Config: map[string]*tfbridge.SchemaInfo{
			"region": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_REGION"},
				},
			},
			"access_key": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_ACCESS_KEY"},
				},
			},
			"secret_key": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_SECRET_KEY"},
				},
			},
			"session_token": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_SESSION_TOKEN"},
				},
			},
			"file_path": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_FILE_PATH"},
				},
			},
			"profile": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_PROFILE"},
				},
			},
			"disable_ssl": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_DISABLE_SSL"},
				},
			},
			"customer_headers": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_CUSTOMER_HEADERS"},
				},
			},
			"proxy_url": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"VOLCENGINE_PROXY_URL"},
				},
			},
			"endpoints": {
				Fields: map[string]*info.Schema{
					"cloudcontrolapi": {
						Default: &tfbridge.DefaultInfo{
							EnvVars: []string{"VOLCENGINE_CC_ENDPOINT"},
						},
					},
					"sts": {
						Default: &tfbridge.DefaultInfo{
							EnvVars: []string{"VOLCENGINE_STS_ENDPOINT"},
						},
					},
				},
			},
			"assume_role": {
				Fields: map[string]*info.Schema{
					"assume_role_trn": {
						Default: &tfbridge.DefaultInfo{
							EnvVars: []string{"VOLCENGINE_ASSUME_ROLE_TRN"},
						},
					},
					"policy": {
						Default: &tfbridge.DefaultInfo{
							EnvVars: []string{"VOLCENGINE_ASSUME_ROLE_POLICY"},
						},
					},
					"duration_seconds": {
						Default: &tfbridge.DefaultInfo{
							EnvVars: []string{"VOLCENGINE_ASSUME_ROLE_DURATION_SECONDS"},
						},
					},
				},
			},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/volcengine/pulumi-%[1]s/sdk/", volcengeccpkg),
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				volcengeccpkg,
			),

			GenerateResourceContainerTypes: true,
			RespectSchemaVersion:           true,
		},
		CSharp: &tfbridge.CSharpInfo{
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
			RootNamespace: "Volcengine.Pulumi",
		},
		Python: &tfbridge.PythonInfo{
			PackageName: "pulumi_volcenginecc",
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			// List any npm dependencies and their versions

			DevDependencies: map[string]string{
				"@types/node": "^10.0.0", // so we can access strongly typed node definitions.
				"@types/mime": "^2.0.0",
			},
			// See the documentation for tfbridge.OverlayInfo for how to lay out this
			// section, or refer to the AWS provider. Delete this section if there are
			// no overlay files.
			// Overlay: &tfbridge.OverlayInfo{},
			PackageName: "@volcengine/pulumi-volcenginecc",
		},
		Java: &tfbridge.JavaInfo{
			BasePackage: "com.volcengine",
		},
		EnableAccurateBridgePreview: true,
		EnableRawStateDelta:         true,
	}
	makeToken := func(module, name string) (string, error) {
		return tokens.MakeStandard(volcengeccpkg)(module, name)
	}
	prov.MustComputeTokens(token.VolcengineToken("volcenginecc_", makeToken))

	// Avoid generating a TypeScript class named Object, which conflicts with
	// JavaScript's built-in Object in NodeNext/CommonJS modules.
	prov.Resources["volcenginecc_tos_object"].Tok = tfbridge.MakeResource(volcengeccpkg, "tos", "TosObject")

	// GN fork: credential projection (Phase 2 spike). Strip credential-bearing
	// output properties before they are recorded in Pulumi state, so secrets
	// never enter the statefile (state-hygiene gate). Applied via bridge-level
	// TransformOutputs, which runs on Create, Update and Read output paths.
	credentialStripResources := []string{
		"volcenginecc_cdn_domain",                             // nested cert config: privateKey, encryptionKey
		"volcenginecc_certificateservice_child_cert_instance", // privateKey
		"volcenginecc_iam_user",                               // loginProfile.password
		"volcenginecc_rocketmq_access_key",                    // secretKey
		"volcenginecc_rocketmq_topic",                         // accessPolicies[].accessKey
	}
	for _, tfName := range credentialStripResources {
		if r, ok := prov.Resources[tfName]; ok {
			r.TransformOutputs = stripCredentialProperties
		}
	}

	for k := range prov.Resources {
		// 获取第二跟斜杠后面的值
		lastPart := extractAndConvertToCamelCase(k)
		//如果当前类的字段存在同名，比如domain的class下有domain字段，则domain字段 + Value
		if prov.P.ResourcesMap().Get(k).Schema().Get(lastPart) != nil {
			prov.Resources[k].Fields = map[string]*info.Schema{
				lastPart: {
					CSharpName: toCamelCase(lastPart) + "Value",
				},
			}
		}
		//如果当前类的字段存在同名，比如domain的class下有domain_state字段，则domain_state字段 + Value
		if prov.P.ResourcesMap().Get(k).Schema().Get(lastPart+"_state") != nil {
			prov.Resources[k].Fields = map[string]*info.Schema{
				lastPart + "_state": {
					CSharpName: toCamelCase(lastPart+"_state") + "Value",
				},
			}
		}
		//如果当前类的字段存在同名，比如domain的class下有domain_args字段，则domain_args字段 + Value
		if prov.P.ResourcesMap().Get(k).Schema().Get(lastPart+"_args") != nil {
			prov.Resources[k].Fields = map[string]*info.Schema{
				lastPart + "_args": {
					CSharpName: toCamelCase(lastPart+"_args") + "Value",
				},
			}
		}
	}
	prov.MustApplyAutoAliases()

	return prov
}

// stripCredentialProperties removes credential-bearing properties (by Pulumi
// property name, at any nesting depth) from a resource output map before it is
// stored in Pulumi state.
func stripCredentialProperties(_ context.Context, props resource.PropertyMap) (resource.PropertyMap, error) {
	var stripValue func(v resource.PropertyValue) resource.PropertyValue
	stripValue = func(v resource.PropertyValue) resource.PropertyValue {
		if v.IsObject() {
			out := resource.PropertyMap{}
			for k, elem := range v.ObjectValue() {
				if credentialPropertyNames[k] {
					continue
				}
				out[k] = stripValue(elem)
			}
			return resource.NewObjectProperty(out)
		}
		if v.IsArray() {
			items := v.ArrayValue()
			out := make([]resource.PropertyValue, len(items))
			for i, item := range items {
				out[i] = stripValue(item)
			}
			return resource.NewArrayProperty(out)
		}
		return v
	}
	for k, elem := range props {
		if credentialPropertyNames[k] {
			delete(props, k)
			continue
		}
		props[k] = stripValue(elem)
	}
	return props, nil
}

var credentialPropertyNames = map[resource.PropertyKey]bool{
	"secretKey":     true,
	"privateKey":    true,
	"encryptionKey": true,
	"password":      true,
	"accessKey":     true,
}

func extractAndConvertToCamelCase(input string) string {
	// 按下划线分割字符串
	parts := strings.Split(input, "_")
	// 如果少于3个部分，返回空字符串或原字符串
	if len(parts) < 3 {
		return input
	}
	// 获取第二根下划线后的所有部分
	remainingParts := parts[2:]
	// 直接返回第二根下划线后的所有内容，保持下划线分隔符
	return strings.Join(remainingParts, "_")
}

// toCamelCase 将字符串转换为驼峰命名，首字母大写
func toCamelCase(s string) string {
	if len(s) == 0 {
		return ""
	}

	// 按下划线分割字符串
	parts := strings.Split(s, "_")
	var result strings.Builder

	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		if i == 0 {
			// 首字母大写
			result.WriteString(capitalizeFirst(part))
		} else {
			// 其余部分首字母大写，其余小写
			result.WriteString(capitalizeFirst(strings.ToLower(part)))
		}
	}

	return result.String()
}

// capitalizeFirst 将字符串首字母大写
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
