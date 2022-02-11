package connectors

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// PostgresSSL is a type for setting connection ssl mode to PostgresDB.
type PostgresSSL uint8

// Enum.
const (
	_                     PostgresSSL = iota
	PostgresSSLDisable                // disable
	PostgresSSLAllow                  // allow
	PostgresSSLPrefer                 // prefer
	PostgresSSLRequire                // require
	PostgresSSLVerifyCa               // verify-ca
	PostgresSSLVerifyFull             // verify-full
)

type (
	// PostgresDB config for connecting to cockroachDB.
	PostgresDB struct {
		User     string `yaml:"user" json:"user" hcl:"user"`
		Password string `yaml:"password" json:"password" hcl:"password"`
		Host     string `yaml:"host" json:"host" hcl:"host"`
		Port     int    `yaml:"port" json:"port" hcl:"port"`
		Database string `yaml:"database" json:"database" hcl:"database"`

		Parameters *PostgresDBParameters `yaml:"parameters" json:"parameters" hcl:"parameters,block"`
	}

	// PostgresDBParameters contains url parameters for connecting to database.
	PostgresDBParameters struct {
		ApplicationName string        `yaml:"application_name" json:"application_name" hcl:"application_name"`
		Mode            PostgresSSL   `yaml:"mode" json:"mode" hcl:"mode"`
		SSLCert         string        `yaml:"ssl_cert" json:"ssl_cert" hcl:"ssl_cert"`
		SSLKey          string        `yaml:"ssl_key" json:"ssl_key" hcl:"ssl_key"`
		SSLRootCert     string        `yaml:"ssl_root_cert" json:"ssl_root_cert" hcl:"ssl_root_cert"`
		ConnectTimeout  time.Duration `yaml:"connect_timeout" json:"connect_timeout" hcl:"connect_timeout"`

		Options *CockroachDBOptions `yaml:"options" json:"options" hcl:"options,block"`
	}

	// PostgresDBOptions contains options for setting variables and cluster ID.
	PostgresDBOptions struct {
		Cluster  string             `yaml:"cluster" json:"cluster" hcl:"cluster"`
		Variable PostgresDBVariable `yaml:"variable" json:"variable" hcl:"variable,block"`
	}

	// PostgresDBVariable sets variable for connections.
	PostgresDBVariable struct {
		Name  string `yaml:"name" json:"name" hcl:"name"`
		Value string `yaml:"value" json:"value" hcl:"value"`
	}

	// TODO: Will we add?
	s struct {
		SearchPath                      string             // Specifies the order in which schemas are searched.
		DefaultTransactionIsolation     sql.IsolationLevel // One of: LevelDefault, LevelReadUncommitted, LevelReadCommitted, LevelRepeatableRead, LevelSerializable.
		StatementTimeout                time.Duration      // Round to milliseconds.
		LockTimeout                     time.Duration      // Round to milliseconds.
		IdleInTransactionSessionTimeout time.Duration      // Round to milliseconds.
		Other                           map[string]string  // Any other parameters from https://www.postgresql.org/docs/current/runtime-config-client.html.
	}
)

// DSN convert struct to DSN and returns connection string.
func (p PostgresDB) DSN() (string, error) {
	str := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
	)

	uri, err := url.Parse(str)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}

	if p.Parameters == nil {
		return uri.String(), nil
	}

	parameters := url.Values{}

	if p.Parameters.ApplicationName != "" {
		parameters.Add("application_name", p.Parameters.ApplicationName)
	}

	if p.Parameters.Mode != 0 {
		parameters.Add("sslmode", p.Parameters.Mode.String())
	}

	if p.Parameters.SSLRootCert != "" {
		parameters.Add("sslrootcert", p.Parameters.SSLRootCert)
	}

	if p.Parameters.SSLCert != "" {
		parameters.Add("sslcert", p.Parameters.SSLCert)
	}

	if p.Parameters.SSLKey != "" {
		parameters.Add("sslkey", p.Parameters.SSLKey)
	}

	uri.RawQuery = parameters.Encode()
	if p.Parameters.Options == nil {
		return uri.String(), nil
	}

	var options []string
	if p.Parameters.Options.Cluster != "" {
		options = append(options, fmt.Sprintf("--cluster=%s", p.Parameters.Options.Cluster))
	}

	if p.Parameters.Options.Variable.Name != "" {
		options = append(options, fmt.Sprintf("-c %s=%s", p.Parameters.Options.Variable.Name, p.Parameters.Options.Variable.Value))
	}

	parameters.Add("options", strings.Join(options, " "))
	uri.RawQuery = parameters.Encode()

	return uri.String(), nil
}
