package config

// Environment variable names. Loading and validation both use these,
// so the  two can never drift apart.
const (
	// application environment
	envAppName = "APP_NAME"
	envAppEnv  = "APP_ENV"

	// http
	envHTTPAddr              = "HTTP_ADDR"
	envHTTPReadTimeout       = "HTTP_READ_TIMEOUT"
	envHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
	//nolint:gosec
	envHTTPWriteTimeout    = "HTTP_WRITE_TIMEOUT"
	envHTTPIdleTimeout     = "HTTP_IDLE_TIMEOUT"
	envHTTPShutdownTimeout = "HTTP_SHUTDOWN_TIMEOUT"
	envHTTPDrainDelay      = "HTTP_DRAIN_DELAY"
	envHTTPMaxBodyBytes    = "HTTP_MAX_BODY_BYTES"
	envHTTPMaxHeaderBytes  = "HTTP_MAX_HEADER_BYTES"

	// logger
	envLogLevel     = "LOG_LEVEL"
	envLogFormat    = "LOG_FORMAT"
	envLogAddSource = "LOG_ADD_SOURCE"

	// database
	envDatabaseURL      = "DATABASE_URL"
	envDBMaxConns       = "DB_MAX_CONNS"
	envDBMinConns       = "DB_MIN_CONNS"
	envDBMaxConnLife    = "DB_MAX_CONN_LIFETIME"
	envDBMaxConnIdle    = "DB_MAX_CONN_IDLE_TIME"
	envDBConnectTimeout = "DB_CONNECT_TIMEOUT"
	envDBStartupTimeout = "DB_STARTUP_TIMEOUT"
	envDBHealthTimeout  = "DB_HEALTH_TIMEOUT"
)
