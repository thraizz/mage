package config

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// bot_config_test.go pins the one rule that has a real blast radius: the LLM
// API key must not be reachable from any YAML file.
//
// config/config.yaml is a HARD LINK to config.dev.yaml -- one inode, two names,
// and the one cmd/server/main.go loads by default. A key that could be set from
// config is a key that ends up committed.

// structTag returns the mapstructure tag of a named field.
func structTag(v any, field string) (string, bool) {
	f, ok := reflect.TypeOf(v).FieldByName(field)
	if !ok {
		return "", false
	}
	return f.Tag.Get("mapstructure"), true
}

func baseConfig() *Config {
	c := &Config{}
	c.Server.MaxSessions = 10
	c.Server.LeasePeriod = time.Second
	c.Database.Host = "localhost"
	c.Database.Database = "mage"
	c.Auth.Mode = "optional"
	c.Logging.Level = "info"
	c.Bot = BotConfig{
		Model:          "claude-sonnet-5",
		MaxTokens:      20000,
		RequestTimeout: 120 * time.Second,
		StallTimeout:   50 * time.Second,
	}
	return c
}

func TestBotAPIKeyIsNotUnmarshalable(t *testing.T) {
	// mapstructure:"-" is what makes this true; assert it structurally rather
	// than trusting a comment.
	f, ok := structTag(BotConfig{}, "APIKey")
	if !ok {
		t.Fatal("BotConfig has no APIKey field")
	}
	if f != "-" {
		t.Fatalf(`APIKey mapstructure tag is %q, want "-": a key settable from YAML is a key that gets committed`, f)
	}
}

func TestBotValidation(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{"good", func(*Config) {}, ""},
		{"bad model", func(c *Config) { c.Bot.Model = "gpt-4" }, "bot.model"},
		{"bad effort", func(c *Config) { c.Bot.Effort = "turbo" }, "bot.effort"},
		{
			// Haiku 4.5 does not support the effort parameter at all. Catching
			// it here beats catching it as a 400 on the first request of an
			// overnight run.
			"effort on haiku",
			func(c *Config) { c.Bot.Model = "claude-haiku-4-5"; c.Bot.Effort = "high" },
			"not supported on claude-haiku-4-5",
		},
		{"effort on sonnet is fine", func(c *Config) { c.Bot.Effort = "medium" }, ""},
		{"zero max tokens", func(c *Config) { c.Bot.MaxTokens = 0 }, "bot.max_tokens"},
		{"zero stall timeout", func(c *Config) { c.Bot.StallTimeout = 0 }, "bot.stall_timeout"},
		{"enabled without a key", func(c *Config) { c.Bot.Enabled = true }, "ANTHROPIC_API_KEY"},
		{"enabled with a key", func(c *Config) { c.Bot.Enabled = true; c.Bot.APIKey = "x" }, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := baseConfig()
			tc.mutate(c)
			err := c.Validate()
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("expected an error mentioning %q", tc.wantErr)
			case tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr):
				t.Fatalf("error = %v, want it to mention %q", err, tc.wantErr)
			}
		})
	}
}
