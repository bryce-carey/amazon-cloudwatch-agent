// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package ec2taggerprocessor

import (
	"time"

	"github.com/aws/private-amazon-cloudwatch-agent-staging/plugins/processors/ec2tagger"
	"github.com/aws/private-amazon-cloudwatch-agent-staging/translator/translate/otel/common"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/confmap"
)

const (
	AppendDimensionsKey = "append_dimensions"
)

type translator struct {
	factory component.ProcessorFactory
}

var _ common.Translator[config.Processor] = (*translator)(nil)

func NewTranslator() common.Translator[config.Processor] {
	return &translator{ec2tagger.NewFactory()}
}

func (t *translator) Type() config.Type {
	return t.factory.Type()
}

// Translate creates an processor config based on the fields in the
// Metrics section of the JSON config.
func (t *translator) Translate(conf *confmap.Conf) (config.Processor, error) {
	taggerKey := common.ConfigKey(common.MetricsKey, AppendDimensionsKey)
	if conf == nil || !conf.IsSet(taggerKey) {
		return nil, &common.MissingKeyError{Type: t.Type(), JsonKey: taggerKey}
	}

	cfg := t.factory.CreateDefaultConfig().(*ec2tagger.Config)
	for k, v := range ec2tagger.SupportedAppendDimensions {
		value, ok := common.GetString(conf, common.ConfigKey(common.MetricsKey, AppendDimensionsKey, k))
		if ok && v == value {
			if k == "AutoScalingGroupName" {
				cfg.EC2InstanceTagKeys = append(cfg.EC2InstanceTagKeys, k)
			} else {
				cfg.EC2MetadataTags = append(cfg.EC2MetadataTags, k)
			}
		}
	}
	cfg.RefreshIntervalSeconds = 0 * time.Second

	return cfg, nil
}
