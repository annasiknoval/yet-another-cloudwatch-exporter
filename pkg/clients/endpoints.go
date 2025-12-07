// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package clients

import (
	"os"
	"strings"
)

var (
	serviceEndpointKeyAliases = map[string]string{
		"cloudwatch":               "monitoring",
		"monitoring":               "monitoring",
		"sts":                      "sts",
		"iam":                      "iam",
		"tagging":                  "tagging",
		"resourcegroupstaggingapi": "tagging",
		"autoscaling":              "autoscaling",
		"apigateway":               "apigateway",
		"apigatewayv2":             "apigateway",
		"ec2":                      "ec2",
		"dms":                      "dms",
		"databasemigrationservice": "dms",
		"aps":                      "aps",
		"amp":                      "aps",
		"prometheus":               "aps",
		"prometheusservice":        "aps",
		"storagegateway":           "storagegateway",
		"shield":                   "shield",
	}
	canonicalServiceEndpointKeys = []string{
		"apigateway",
		"aps",
		"autoscaling",
		"monitoring",
		"dms",
		"ec2",
		"iam",
		"shield",
		"storagegateway",
		"sts",
		"tagging",
	}
)

// CanonicalServiceEndpointKey returns the canonical service key for an alias, or false if unsupported.
func CanonicalServiceEndpointKey(key string) (string, bool) {
	canonicalKey, ok := serviceEndpointKeyAliases[strings.ToLower(key)]
	return canonicalKey, ok
}

// CanonicalServiceEndpointKeys returns the list of supported canonical service keys.
func CanonicalServiceEndpointKeys() []string {
	keys := make([]string, 0, len(canonicalServiceEndpointKeys))
	keys = append(keys, canonicalServiceEndpointKeys...)
	return keys
}

// ServiceEndpointKeyAliases returns a copy of the alias map.
func ServiceEndpointKeyAliases() map[string]string {
	aliases := make(map[string]string, len(serviceEndpointKeyAliases))
	for k, v := range serviceEndpointKeyAliases {
		aliases[k] = v
	}
	return aliases
}

// LoadServiceEndpointsFromEnv reads environment variables in the form AWS_ENDPOINT_URL_<SERVICE>
// (e.g. AWS_ENDPOINT_URL_CLOUDWATCH) and returns a canonical service->endpoint map.
// Canonical env vars override alias env vars, and later overrides replace earlier values.
func LoadServiceEndpointsFromEnv() map[string]string {
	endpoints := map[string]string{}

	for _, key := range canonicalServiceEndpointKeys {
		if val := os.Getenv(endpointEnvVarName(key)); val != "" {
			endpoints[key] = val
		}
	}

	for alias := range serviceEndpointKeyAliases {
		envName := endpointEnvVarName(alias)
		if val := os.Getenv(envName); val != "" {
			if canonical, ok := CanonicalServiceEndpointKey(alias); ok {
				if _, exists := endpoints[canonical]; !exists {
					endpoints[canonical] = val
				}
			}
		}
	}

	return endpoints
}

func endpointEnvVarName(key string) string {
	return "AWS_ENDPOINT_URL_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}
