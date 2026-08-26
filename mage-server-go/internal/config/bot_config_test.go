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
		// The Phase 5 cases below are all about the Anthropic path, so the
		// base fixture pins that provider explicitly. Gemini -- the DEFAULT --
		// gets its own cases.
		Provider:       BotProviderAnthropic,
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

		// Phase 5b: the provider dimension.
		{"unknown provider", func(c *Config) { c.Bot.Provider = "hal9000" }, "bot.provider"},
		{
			"gemini is valid", func(c *Config) {
				c.Bot.Provider = BotProviderGemini
				c.Bot.Model = "gemini-3.7-flash"
			}, "",
		},
		{
			// A model id from the WRONG provider is a 404 on the first request
			// of an unattended run. It is catchable here for free.
			"claude model under gemini", func(c *Config) {
				c.Bot.Provider = BotProviderGemini
			}, "bot.model",
		},
		{
			"gemini model under anthropic", func(c *Config) {
				c.Bot.Model = "gemini-3.7-flash"
			}, "bot.model",
		},
		{
			// reasoning_effort on the OpenAI-compatible endpoint takes
			// none|low|medium|high. xhigh and max are Anthropic-only.
			"xhigh effort under gemini", func(c *Config) {
				c.Bot.Provider = BotProviderGemini
				c.Bot.Model = "gemini-3.7-flash"
				c.Bot.Effort = "xhigh"
			}, "bot.effort",
		},
		{
			"high effort under gemini is fine", func(c *Config) {
				c.Bot.Provider = BotProviderGemini
				c.Bot.Model = "gemini-3.7-flash"
				c.Bot.Effort = "high"
			}, "",
		},
		{
			"gemini enabled without a key names GEMINI_API_KEY", func(c *Config) {
				c.Bot.Provider = BotProviderGemini
				c.Bot.Model = "gemini-3.7-flash"
				c.Bot.Enabled = true
			}, "GEMINI_API_KEY",
		},
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

func TestBotAPIKeyEnvIsProviderSpecific(t *testing.T) {
	// The key rule of Phase 5, extended rather than loosened: each provider
	// reads exactly one environment variable, and neither is reachable from
	// YAML (see TestBotAPIKeyIsNotUnmarshalable).
	if got := BotAPIKeyEnv(BotProviderGemini); got != "GEMINI_API_KEY" {
		t.Errorf("gemini key env = %q", got)
	}
	if got := BotAPIKeyEnv(BotProviderAnthropic); got != "ANTHROPIC_API_KEY" {
		t.Errorf("anthropic key env = %q", got)
	}
	// An unset/unknown provider must not fall back to the anthropic variable:
	// gemini is the default, so the default key must follow it.
	if got := BotAPIKeyEnv(""); got != "GEMINI_API_KEY" {
		t.Errorf("default key env = %q, want GEMINI_API_KEY", got)
	}
}

func TestBotDefaultsAreGemini(t *testing.T) {
	// Load() fills the provider and the provider-dependent model default. A
	// config with neither set must come out as a valid gemini config, not as a
	// validation error.
	c := baseConfig()
	c.Bot.Provider = ""
	c.Bot.Model = ""
	if c.Bot.Provider == "" {
		c.Bot.Provider = BotProviderGemini
	}
	if c.Bot.Model == "" {
		c.Bot.Model = botDefaultModel[c.Bot.Provider]
	}
	if c.Bot.Model != "gemini-3.7-flash" {
		t.Fatalf("default model = %q, want gemini-3.7-flash", c.Bot.Model)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("the default bot config does not validate: %v", err)
	}
}

func TestEveryProviderHasAValidDefaultModel(t *testing.T) {
	// The tables in config.go and the constants in internal/bot/llm are
	// duplicated on purpose (config imports nothing). This is what keeps them
	// from drifting apart silently.
	for _, p := range botProviders {
		def, ok := botDefaultModel[p]
		if !ok {
			t.Errorf("provider %q has no default model", p)
			continue
		}
		if !contains(botModels[p], def) {
			t.Errorf("provider %q default model %q is not in its own model list", p, def)
		}
		if !contains(botEfforts[p], "") {
			t.Errorf("provider %q does not allow an empty effort", p)
		}
	}
}
