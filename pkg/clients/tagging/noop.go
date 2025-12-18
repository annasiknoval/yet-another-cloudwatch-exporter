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
package tagging

import (
	"context"
	"log/slog"

	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/model"
)

// NoOpClient provides a tagging client that does nothing.
// This is useful when tagging functionality is not available or needed.
type NoOpClient struct {
	logger *slog.Logger
}

// NewNoOpClient creates a new no-op tagging client.
func NewNoOpClient(logger *slog.Logger) Client {
	return &NoOpClient{
		logger: logger,
	}
}

// GetResources returns an empty slice of resources.
// This effectively disables resource discovery via tagging.
func (c *NoOpClient) GetResources(ctx context.Context, job model.DiscoveryJob, region string) ([]*model.TaggedResource, error) {
	c.logger.Debug("Tagging disabled - skipping resource discovery", "region", region, "namespace", job.Namespace)
	return []*model.TaggedResource{}, nil
}
