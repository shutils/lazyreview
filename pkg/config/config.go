package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

// Tmp file name
const (
	tmpReviewFileName = "tmp_review.md"
	tmpPromptFileName = "tmp_prompt.md"
)

// ModelCost represents the cost associated with AI model usage.
type ModelCost struct {
	Input, Output float64
}

// StringOrSlice is a custom type that can hold either a string or a slice of strings.
type StringOrSlice []string

// UnmarshalText decodes a single string into a StringOrSlice.
func (s *StringOrSlice) UnmarshalText(text []byte) error {
	*s = []string{string(text)}
	return nil
}

// UnmarshalTOML decodes either a string or a slice of strings into a StringOrSlice.
func (s *StringOrSlice) UnmarshalTOML(data any) error {
	switch v := data.(type) {
	case string:
		*s = strings.Fields(v) // Split the string by spaces into a slice
	case []any:
		var result []string
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				return fmt.Errorf("invalid type in array: %T", item)
			}
		}
		*s = result
	default:
		return fmt.Errorf("unexpected type: %T", v)
	}
	return nil
}

type Source struct {
	Name      string        `toml:"name"`
	Collector StringOrSlice `toml:"collector"`
	Previewer StringOrSlice `toml:"previewer"`
	Prompt    string        `toml:"prompt"`
	Enabled   bool          `toml:"enabled"`
}

func (i Source) Title() string {
	if i.Enabled {
		return "☑ " + i.Name
	}
	return "☐ " + i.Name
}
func (i Source) Description() string {
	if i.Enabled {
		return "☑ collector: " + strings.Join(i.Collector, ", ") + " previewer: " + strings.Join(i.Previewer, ", ")
	}
	return "☐ collector: " + strings.Join(i.Collector, ", ") + " previewer: " + strings.Join(i.Previewer, ", ")
}
func (i Source) FilterValue() string { return i.Name }

func (i Source) String() string {
	return fmt.Sprintf(
		"Name: %s\n"+
			"Collector: %s\n"+
			"Previewer: %s\n"+
			"Prompt: %s\n"+
			"Enabled: %v",
		i.Name,
		strings.Join(i.Collector, " "),
		strings.Join(i.Previewer, " "),
		i.Prompt,
		i.Enabled,
	)
}

const projectName = "lazyreview"

// Config holds the configuration details for the application.
type Config struct {
	ConfigPath    string        `toml:"-"`
	Key           string        `toml:"key"`
	Endpoint      string        `toml:"endpoint"`
	Version       string        `toml:"version"`
	Model         string        `toml:"model"`
	ModelCost     ModelCost     `toml:"modelCost"`
	Target        string        `toml:"target"`
	Output        string        `toml:"output"`
	State         string        `toml:"state"`
	Ignores       []string      `toml:"ignores"`
	Prompt        string        `toml:"prompt"`
	Type          string        `toml:"type"`
	Collector     StringOrSlice `toml:"collector"`
	Previewer     StringOrSlice `toml:"previewer"`
	Glamour       string        `toml:"glamour"`
	MaxTokens     int           `toml:"max_tokens"`
	Opener        StringOrSlice `toml:"opener"`
	Sources       []Source      `toml:"sources"`
	TmpReviewPath string        `toml:"-"`
	TmpPromptPath string        `toml:"-"`
}

// loadConfig reads the configuration from the specified file.
func loadConfig(filePath string, config *Config) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(data), config)
	return err
}

// NewConfig initializes a new Config instance and populates it with values.
func NewConfig() Config {
	c := Config{}
	defaults := map[string]string{
		"target": ".",
		"output": path.Join(xdg.DataHome, projectName, "reviews.json"),
		"model":  "gpt-4o-mini",
	}

	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()

	c.ConfigPath = determineConfigPath(*configPath)
	setConfigFlags(&c, defaults)

	if err := loadConfigFile(&c); err != nil {
		if os.IsNotExist(err) {
			saveConfig(c.ConfigPath, Config{})
			log.Println("Config file not found. A new one has been created at", c.ConfigPath)
		} else {
			log.Fatalln("Failed to load config file:", err)
		}
	}

	validateConfig(&c)

	c.State = setDefaultState(c.State)
	c.TmpReviewPath = path.Join(xdg.CacheHome, projectName, tmpReviewFileName)
	c.TmpPromptPath = path.Join(xdg.CacheHome, projectName, tmpPromptFileName)

	return c
}

func determineConfigPath(configPath string) string {
	if configPath == "" {
		return path.Join(xdg.ConfigHome, projectName, "config.toml")
	}
	return configPath
}

func setConfigFlags(c *Config, defaults map[string]string) {
	for flagKey, defaultValue := range defaults {
		switch flagKey {
		case "target":
			if val := flag.String("target", defaultValue, "Target path"); *val != defaultValue || c.Target == "" {
				c.Target = *val
			}
		case "output":
			if val := flag.String("output", defaultValue, "Output path"); *val != defaultValue || c.Output == "" {
				c.Output = *val
			}
		case "model":
			if val := flag.String("model", defaultValue, "Model to use"); *val != defaultValue || c.Model == "" {
				c.Model = *val
			}
		}
	}
}

func loadConfigFile(c *Config) error {
	if c.ConfigPath == "" {
		return nil
	}
	return loadConfig(c.ConfigPath, c)
}

func validateConfig(c *Config) {
	if c.Target == "" || c.Output == "" || c.Model == "" {
		log.Fatalf("Missing required configuration fields. Ensure `token`, `target`, `output`, and `model` are provided.")
	}
}

func setDefaultState(state string) string {
	if state == "" {
		return path.Join(xdg.StateHome, projectName, "state.json")
	}
	return state
}

// ToStringArray converts the Config struct to a string array.
func (c Config) ToStringArray() []string {
	result := []string{"key=******"}
	// Append other fields as key=value pairs
	result = append(result,
		fmt.Sprintf("endpoint=%s", c.Endpoint),
		fmt.Sprintf("version=%s", c.Version),
		fmt.Sprintf("model=%s", c.Model),
		fmt.Sprintf("target=%s", c.Target),
		fmt.Sprintf("output=%s", c.Output),
		fmt.Sprintf("state=%s", c.State),
		fmt.Sprintf("ignores=%s", strings.Join(c.Ignores, ",")),
		fmt.Sprintf("prompt=%s", c.Prompt),
		fmt.Sprintf("type=%s", c.Type),
		fmt.Sprintf("collector=%s", c.Collector),
		fmt.Sprintf("glamour=%s", c.Glamour),
		fmt.Sprintf("max_tokens=%d", c.MaxTokens),
		fmt.Sprintf("tmp_review_path=%s", c.TmpReviewPath),
		fmt.Sprintf("opener=%s", c.Opener),
		"\n",
	)

	// Append Sources
	for _, source := range c.Sources {
		result = append(result, source.String(), "\n")
	}

	return result
}

// saveConfig writes the Config data to a specified file.
func saveConfig(filePath string, config Config) {
	if filePath == "" {
		log.Fatalf("File path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create or open config file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		log.Fatalf("Failed to encode config to TOML: %v", err)
	}

	log.Printf("Config saved to %s", filePath)
}

func (c *Config) GetSources() []Source {
	return c.Sources
}

func (c *Config) ToggleSourceEnabled(sourceName string) {
	for i, source := range c.Sources {
		if source.Name == sourceName {
			c.Sources[i].Enabled = !source.Enabled
			return
		}
	}
}

func (c *Config) GetSourceFromName(sourceName string) Source {
	for _, source := range c.Sources {
		if source.Name == sourceName {
			return source
		}
	}
	return Source{}

}

// ParseCommand splits a command into the command name and arguments.
func ParseCommand(command []string) (string, []string) {
	if len(command) == 0 {
		return "", nil
	}
	if len(command) == 1 {
		return command[0], nil
	}
	return command[0], command[1:]
}
