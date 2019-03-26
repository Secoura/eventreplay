package config

import "time"

type Replacement struct {
	Token       string `yaml:"token"`
	Type        string `yaml:"type"`
	Replacement interface{} `yaml:"replacement"`
}

type SampleConfig struct {
	InputFile    string        `yaml:"input_file"`
	Identifier   string        `yaml:"identifier"`
	OutputMode   string        `yaml:"output_mode"`
	OutputCodec  string        `yaml:"output_codec"`
	SpoolDir     string        `yaml:"spool_dir"`
	Count        int           `yaml:"count"`
	Delimiter    string        `yaml:"delimiter"`
	Interval     time.Duration `yaml:"interval"`
	EarliestTime string        `yaml:"earliest_time"`
	LatestTime   string        `yaml:"latest_time"`
	Replacements []Replacement `yaml:"replacements"`
}

func (c *SampleConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawSampleConfig SampleConfig
	raw := rawSampleConfig(DefaultSampleConfig)
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*c = SampleConfig(raw)
	return nil
}

type Config struct {
	Samples []SampleConfig
}

var DefaultSampleConfig = SampleConfig{
	Identifier:   "eventreplay",
	OutputMode:   "console",
	OutputCodec:  "raw",
	Delimiter:    "[^\r\n\\s]+",
	Count:        -1,
	EarliestTime: "now",
	LatestTime:   "now",
	Interval:     60 * time.Second,
}
