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
package job

import (
	"context"
	"log/slog"
	"sync"

	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/clients/cloudwatch"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/config"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/model"
)

func runStaticJob(
	ctx context.Context,
	logger *slog.Logger,
	resource model.StaticJob,
	clientCloudwatch cloudwatch.Client,
) []*model.CloudwatchData {
	cw := []*model.CloudwatchData{}
	mux := &sync.Mutex{}
	var wg sync.WaitGroup

	for j := range resource.Metrics {
		metric := resource.Metrics[j]
		wg.Add(1)
		go func() {
			defer wg.Done()

			data := model.CloudwatchData{
				MetricName:   metric.Name,
				ResourceName: resource.Name,
				Namespace:    resource.Namespace,
				Dimensions:   createStaticDimensions(resource.Dimensions),
				MetricMigrationParams: model.MetricMigrationParams{
					NilToZero:              metric.NilToZero,
					AddCloudwatchTimestamp: metric.AddCloudwatchTimestamp,
				},
				Tags:                          nil,
				GetMetricDataProcessingParams: nil,
				GetMetricDataResult:           nil,
				GetMetricStatisticsResult:     nil,
			}

			// Use GetMetricData if the feature flag is enabled, otherwise use GetMetricStatistics
			flags := config.FlagsFromCtx(ctx)
			if flags.IsFeatureEnabled("use-getmetricdata-for-static") {
				// Use GetMetricData API for static jobs
				// Create a separate CloudwatchData for each statistic
				for _, statistic := range metric.Statistics {
					statisticData := data
					statisticData.GetMetricDataProcessingParams = &model.GetMetricDataProcessingParams{
						Period:    metric.Period,
						Length:    metric.Length,
						Delay:     metric.Delay,
						Statistic: statistic,
					}
					statisticData.GetMetricDataResult = &model.GetMetricDataResult{
						Statistic: statistic,
						// GetMetricData will be called later in the processor
					}
					mux.Lock()
					cw = append(cw, &statisticData)
					mux.Unlock()
				}
			} else {
				// Use traditional GetMetricStatistics API
				data.GetMetricStatisticsResult = &model.GetMetricStatisticsResult{
					Results:    clientCloudwatch.GetMetricStatistics(ctx, logger, data.Dimensions, resource.Namespace, metric),
					Statistics: metric.Statistics,
				}

				if data.GetMetricStatisticsResult.Results != nil {
					mux.Lock()
					cw = append(cw, &data)
					mux.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	return cw
}

func createStaticDimensions(dimensions []model.Dimension) []model.Dimension {
	out := make([]model.Dimension, 0, len(dimensions))
	for _, d := range dimensions {
		out = append(out, model.Dimension{
			Name:  d.Name,
			Value: d.Value,
		})
	}

	return out
}
