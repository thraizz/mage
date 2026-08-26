package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete server configuration
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Validation ValidationConfig `mapstructure:"validation"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Mail       MailConfig       `mapstructure:"mail"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Cache      CacheConfig      `mapstructure:"cache"`
	Plugins    PluginConfig     `mapstructure:"plugins"`
	Health     HealthConfig     `mapstructure:"health"`
	Metrics    MetricsConfig    `mapstructure:"metrics"`
	Bot        BotConfig        `mapstructure:"bot"`
}

// ServerConfig contains server-related settings
type ServerConfig struct {
	Name           string          `mapstructure:"name"`
	GRPC           GRPCConfig      `mapstructure:"grpc"`
	WebSocket      WebSocketConfig `mapstructure:"websocket"`
	MaxSessions    int             `mapstructure:"max_sessions"`
	LeasePeriod    time.Duration   `mapstructure:"lease_period"`
	MaxIdleSeconds int             `mapstructure:"max_idle_seconds"`
	MaxGameThreads int             `mapstructure:"max_game_threads"`
}

// GRPCConfig contains gRPC server settings
type GRPCConfig struct {
	Address              string        `mapstructure:"address"`
	MaxConcurrentStreams int           `mapstructure:"max_concurrent_streams"`
	MaxConnectionAge     time.Duration `mapstructure:"max_connection_age"`
}

// WebSocketConfig contains WebSocket server settings
type WebSocketConfig struct {
	Address      string        `mapstructure:"address"`
	PingInterval time.Duration `mapstructure:"ping_interval"`
	PongTimeout  time.Duration `mapstructure:"pong_timeout"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// ValidationConfig contains validation rules
type ValidationConfig struct {
	Username UsernameValidation `mapstructure:"username"`
	Password PasswordValidation `mapstructure:"password"`
}

// UsernameValidation contains username validation rules
type UsernameValidation struct {
	MinLength int    `mapstructure:"min_length"`
	MaxLength int    `mapstructure:"max_length"`
	Pattern   string `mapstructure:"pattern"`
}

// PasswordValidation contains password validation rules
type PasswordValidation struct {
	MinLength int `mapstructure:"min_length"`
	MaxLength int `mapstructure:"max_length"`
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	Mode                  string        `mapstructure:"mode"`
	RequireEmail          bool          `mapstructure:"require_email"`
	PasswordResetTokenTTL time.Duration `mapstructure:"password_reset_token_ttl"`
	AdminPassword         string        `mapstructure:"admin_password"`
}

// MailConfig contains email service settings
type MailConfig struct {
	Provider string        `mapstructure:"provider"`
	SMTP     SMTPConfig    `mapstructure:"smtp"`
	Mailgun  MailgunConfig `mapstructure:"mailgun"`
}

// SMTPConfig contains SMTP server settings
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// MailgunConfig contains Mailgun API settings
type MailgunConfig struct {
	Domain string `mapstructure:"domain"`
	APIKey string `mapstructure:"api_key"`
	From   string `mapstructure:"from"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// CacheConfig contains caching settings
type CacheConfig struct {
	Cards CardCacheConfig `mapstructure:"cards"`
}

// CardCacheConfig contains card cache settings
type CardCacheConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	TTL     time.Duration `mapstructure:"ttl"`
	MaxSize int           `mapstructure:"max_size"`
}

// PluginConfig contains plugin settings
type PluginConfig struct {
	GameTypes       []string `mapstructure:"game_types"`
	TournamentTypes []string `mapstructure:"tournament_types"`
	PlayerTypes     []string `mapstructure:"player_types"`
}

// HealthConfig contains health check settings
type HealthConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address"`
}

// MetricsConfig contains metrics settings
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

// BotConfig configures the LLM-backed bot players (internal/bot/llm).
//
// THE API KEY IS NOT IN THIS FILE'S WORLD AT ALL. APIKey is tagged
// mapstructure:"-", so viper cannot populate it from any YAML under any key
// name, and Load fills it from the ANTHROPIC_API_KEY environment variable and
// nowhere else. Two reasons, both concrete:
//
//  1. config/config.yaml is a HARD LINK to config.dev.yaml -- one inode, both
//     names, and it is the file cmd/server/main.go:35 loads by default. A
//     secret written to either is a secret committed to the repository.
//  2. Load calls v.AutomaticEnv() with no prefix and NO SetEnvKeyReplacer, so a
//     nested key like bot.api_key is NOT reachable as BOT_API_KEY. Wiring the
//     key "through config" would need v.SetEnvKeyReplacer(strings.NewReplacer(
//     ".", "_")) added first, and would still put the value one careless
//     `v.Set` away from a YAML dump. An explicit os.Getenv has neither problem.
type BotConfig struct {
	// Enabled turns LLM bot seats on. When true, ANTHROPIC_API_KEY must be set.
	Enabled bool `mapstructure:"enabled"`
	// Model is the model ID. See internal/bot/llm for the cost table.
	Model string `mapstructure:"model"`
	// Effort is the OutputConfig effort level: low|medium|high|xhigh|max, or
	// empty. IT MUST BE EMPTY FOR claude-haiku-4-5, which does not support the
	// parameter; Validate rejects the combination rather than letting the API
	// reject the first request of an overnight run.
	Effort string `mapstructure:"effort"`
	// MaxTokens is the output cap per request. mage-bench measured 20000.
	MaxTokens int `mapstructure:"max_tokens"`
	// RequestTimeout bounds ONE LLM request attempt. mage-bench raised this to
	// 120s at harness epoch #15; 45s was empirically too short.
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	// StallTimeout is the table-stall guard: the total wall clock a bot may
	// hold priority before force-passing. It exists for the other players, so
	// it is deliberately much shorter than RequestTimeout.
	StallTimeout time.Duration `mapstructure:"stall_timeout"`
	// MaxRetries is passed to the SDK. Worst-case request time is
	// RequestTimeout x (MaxRetries+1), which is why StallTimeout is separate.
	MaxRetries int `mapstructure:"max_retries"`

	// APIKey comes from ANTHROPIC_API_KEY. Never from YAML -- see the type doc.
	APIKey string `mapstructure:"-"`
}

// Load loads configuration from a file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
		v.AddConfigPath(".")
	}

	// Enable environment variable override
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// The bot API key comes from the environment only (see BotConfig).
	cfg.Bot.APIKey = os.Getenv("ANTHROPIC_API_KEY")

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.name", "mage-server")
	v.SetDefault("server.grpc.address", "0.0.0.0:17171")
	v.SetDefault("server.grpc.max_concurrent_streams", 1000)
	v.SetDefault("server.grpc.max_connection_age", "1h")
	v.SetDefault("server.websocket.address", "0.0.0.0:17179")
	v.SetDefault("server.websocket.ping_interval", "30s")
	v.SetDefault("server.websocket.pong_timeout", "10s")
	v.SetDefault("server.max_sessions", 10000)
	v.SetDefault("server.lease_period", "5s")
	v.SetDefault("server.max_idle_seconds", 300)
	v.SetDefault("server.max_game_threads", 10)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.database", "mage")
	v.SetDefault("database.user", "mage")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	// Validation defaults
	v.SetDefault("validation.username.min_length", 3)
	v.SetDefault("validation.username.max_length", 14)
	v.SetDefault("validation.username.pattern", "^[a-z0-9_]+$")
	v.SetDefault("validation.password.min_length", 8)
	v.SetDefault("validation.password.max_length", 100)

	// Auth defaults
	v.SetDefault("auth.mode", "optional")
	v.SetDefault("auth.require_email", false)
	v.SetDefault("auth.password_reset_token_ttl", "1h")
	v.SetDefault("auth.admin_password", "")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")

	// Cache defaults
	v.SetDefault("cache.cards.enabled", true)
	v.SetDefault("cache.cards.ttl", "24h")
	v.SetDefault("cache.cards.max_size", 100000)

	// Health defaults
	v.SetDefault("health.enabled", true)
	v.SetDefault("health.address", "0.0.0.0:8080")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.port", 9090)
	v.SetDefault("metrics.path", "/metrics")

	// Bot defaults. No api_key default of any kind, not even an empty one:
	// registering the key as a config path at all invites someone to fill it in
	// (see BotConfig).
	v.SetDefault("bot.enabled", false)
	v.SetDefault("bot.model", "claude-sonnet-5")
	v.SetDefault("bot.effort", "")
	v.SetDefault("bot.max_tokens", 20000)
	v.SetDefault("bot.request_timeout", "120s")
	v.SetDefault("bot.stall_timeout", "50s")
	v.SetDefault("bot.max_retries", 2)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.MaxSessions <= 0 {
		return fmt.Errorf("server.max_sessions must be positive")
	}
	if c.Server.LeasePeriod <= 0 {
		return fmt.Errorf("server.lease_period must be positive")
	}

	// Validate database config
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database.database is required")
	}

	// Validate auth mode
	if c.Auth.Mode != "optional" && c.Auth.Mode != "required" {
		return fmt.Errorf("auth.mode must be 'optional' or 'required'")
	}

	// Validate mail provider
	if c.Mail.Provider != "" && c.Mail.Provider != "smtp" && c.Mail.Provider != "mailgun" {
		return fmt.Errorf("mail.provider must be 'smtp' or 'mailgun'")
	}

	// Validate logging level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("logging.level must be one of: debug, info, warn, error")
	}

	// Validate bot config
	validModels := map[string]bool{
		"claude-sonnet-5": true, "claude-opus-5": true, "claude-haiku-4-5": true,
	}
	if !validModels[c.Bot.Model] {
		return fmt.Errorf("bot.model must be one of: claude-sonnet-5, claude-opus-5, claude-haiku-4-5")
	}
	validEfforts := map[string]bool{
		"": true, "low": true, "medium": true, "high": true, "xhigh": true, "max": true,
	}
	if !validEfforts[c.Bot.Effort] {
		return fmt.Errorf("bot.effort must be one of: low, medium, high, xhigh, max (or empty)")
	}
	if c.Bot.Model == "claude-haiku-4-5" && c.Bot.Effort != "" {
		return fmt.Errorf("bot.effort is not supported on claude-haiku-4-5; leave it empty")
	}
	if c.Bot.MaxTokens <= 0 {
		return fmt.Errorf("bot.max_tokens must be positive")
	}
	if c.Bot.RequestTimeout <= 0 {
		return fmt.Errorf("bot.request_timeout must be positive")
	}
	if c.Bot.StallTimeout <= 0 {
		return fmt.Errorf("bot.stall_timeout must be positive")
	}
	if c.Bot.Enabled && c.Bot.APIKey == "" {
		return fmt.Errorf("bot.enabled requires the ANTHROPIC_API_KEY environment variable")
	}

	return nil
}
