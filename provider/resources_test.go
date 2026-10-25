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

package volcenginecc

import (
	"context"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func TestStripCredentialProperties(t *testing.T) {
	ctx := context.Background()

	t.Run("top-level computed secret", func(t *testing.T) {
		in := resource.PropertyMap{
			"id":         resource.NewStringProperty("rk-123"),
			"secretKey":  resource.NewStringProperty("SUPERSECRET"),
			"access_key": resource.NewStringProperty("AK"),
		}
		out, err := stripCredentialProperties(ctx, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := out["secretKey"]; ok {
			t.Errorf("secretKey was not stripped")
		}
		if v, ok := out["access_key"]; !ok || !v.IsString() {
			t.Errorf("non-credential property access_key was lost")
		}
	})

	t.Run("nested single object (iam user loginProfile.password)", func(t *testing.T) {
		in := resource.PropertyMap{
			"userName": resource.NewStringProperty("alice"),
			"loginProfile": resource.NewObjectProperty(resource.PropertyMap{
				"password":      resource.NewStringProperty("hunter2"),
				"loginAllowed":  resource.NewBoolProperty(true),
				"lastLoginDate": resource.NewStringProperty("2026-10-01"),
			}),
		}
		out, err := stripCredentialProperties(ctx, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		lp := out["loginProfile"].ObjectValue()
		if _, ok := lp["password"]; ok {
			t.Errorf("nested password was not stripped")
		}
		if _, ok := lp["loginAllowed"]; !ok {
			t.Errorf("nested non-credential property was lost")
		}
	})

	t.Run("nested array of objects (rocketmq topic accessPolicies)", func(t *testing.T) {
		in := resource.PropertyMap{
			"topicName": resource.NewStringProperty("t1"),
			"accessPolicies": resource.NewArrayProperty([]resource.PropertyValue{
				resource.NewObjectProperty(resource.PropertyMap{
					"accessKey": resource.NewStringProperty("AKID"),
					"authority": resource.NewStringProperty("ALL"),
				}),
			}),
		}
		out, err := stripCredentialProperties(ctx, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		pol := out["accessPolicies"].ArrayValue()[0].ObjectValue()
		if _, ok := pol["accessKey"]; ok {
			t.Errorf("nested accessKey was not stripped")
		}
		if _, ok := pol["authority"]; !ok {
			t.Errorf("nested authority was lost")
		}
	})

	t.Run("deeply nested (cdn domain https cert privateKey)", func(t *testing.T) {
		in := resource.PropertyMap{
			"domain": resource.NewStringProperty("d.example.com"),
			"https": resource.NewObjectProperty(resource.PropertyMap{
				"certificates": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewObjectProperty(resource.PropertyMap{
						"privateKey":    resource.NewStringProperty("-----BEGIN"),
						"certification": resource.NewStringProperty("certdata"),
					}),
				}),
			}),
		}
		out, err := stripCredentialProperties(ctx, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cert := out["https"].ObjectValue()["certificates"].ArrayValue()[0].ObjectValue()
		if _, ok := cert["privateKey"]; ok {
			t.Errorf("deeply nested privateKey was not stripped")
		}
		if _, ok := cert["certification"]; !ok {
			t.Errorf("deeply nested certification was lost")
		}
	})
}
