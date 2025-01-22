package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

// AI model cost
// $/M token
type ModelCost struct {
	Input, Output float64
}

// Configファイル情報
type Config struct {
	ConfigPath    string    `toml:"-"`
	Key           string    `toml:"key"`
	Endpoint      string    `toml:"endpoint"`
	Version       string    `toml:"version"`
	Model         string    `toml:"model"`
	ModelCost     ModelCost `toml:"modelCost"`
	Target        string    `toml:"target"`
	Output        string    `toml:"output"`
	State         string    `toml:"state"`
	Ignores       []string  `toml:"ignores"`
	Prompt        string    `toml:"prompt"`
	Type          string    `toml:"type"`
	Collector     string    `toml:"collector"`
	Previewer     string    `toml:"previewer"`
	Glamour       string    `toml:"glamour"`
	MaxTokens     int       `toml:"max_tokens"`
	TmpReviewPath string    `toml:"-"`
	Opener        string    `toml:"opener"`
}

var c Config

func loadConfig(filePath string, config *Config) (Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return *config, err
	}
	if _, err := toml.Decode(string(data), &config); err != nil {
		return *config, err
	}
	return *config, nil
}

func NewConfig() Config {
	tFlagDefault := "."
	oFlagDefault := "reviews.json"
	mFlagDefault := "gpt-3.5-0125"
	configPath := flag.String("config", "", "Path to the config file")
	targetFlag := flag.String("target", tFlagDefault, "Target path")
	outputFlag := flag.String("output", oFlagDefault, "Output path")
	modelFlag := flag.String("model", mFlagDefault, "Model to use")

	flag.Parse()

	// コンフィグファイルの読み込み
	if *configPath != "" {
		var err error
		c, err = loadConfig(*configPath, &c)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	// フラグ値でコンフィグ値を上書き（優先順位：フラグ > コンフィグファイル）
	if *targetFlag != tFlagDefault || c.Target == "" {
		c.Target = *targetFlag
	}
	if *outputFlag != oFlagDefault || c.Output == "" {
		c.Output = *outputFlag
	}
	if *modelFlag != mFlagDefault || c.Model == "" {
		c.Model = *modelFlag
	}

	// 必須フィールドの確認
	if c.Target == "" || c.Output == "" || c.Model == "" {
		log.Fatalf("Missing required configuration fields. Ensure `token`, `target`, `output`, and `model` are provided.")
	}

	if c.State == "" {
		c.State = path.Join(xdg.StateHome, "lazyreview", "state.json")
	}

	c.ConfigPath = *configPath
	c.TmpReviewPath = path.Join(xdg.CacheHome, "lazyreview", "tmp_review.md")

	return c
}

func (c Config) ToStringArray() []string {
	var result []string

	// Replace Key value with asterisks
	result = append(result, "key=******")

	// Append other fields as key=value pairs
	result = append(result, fmt.Sprintf("endpoint=%s", c.Endpoint))
	result = append(result, fmt.Sprintf("version=%s", c.Version))
	result = append(result, fmt.Sprintf("model=%s", c.Model))
	result = append(result, fmt.Sprintf("target=%s", c.Target))
	result = append(result, fmt.Sprintf("output=%s", c.Output))
	result = append(result, fmt.Sprintf("state=%s", c.State))
	result = append(result, fmt.Sprintf("ignores=%s", strings.Join(c.Ignores, ",")))
	result = append(result, fmt.Sprintf("prompt=%s", c.Prompt))
	result = append(result, fmt.Sprintf("type=%s", c.Type))
	result = append(result, fmt.Sprintf("collector=%s", c.Collector))
	result = append(result, fmt.Sprintf("glamour=%s", c.Glamour))
	result = append(result, fmt.Sprintf("max_tokens=%d", c.MaxTokens))
	result = append(result, fmt.Sprintf("tmp_review_path=%s", c.TmpReviewPath))
	result = append(result, fmt.Sprintf("opener=%s", c.Opener))

	return result
}
