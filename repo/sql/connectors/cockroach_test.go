package connectors_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/Meat-Hook/framework/repo/sql/connectors"
)

func TestCockroachDB_Unmarshal(t *testing.T) {
	t.Parallel()

	var (
		all = connectors.CockroachDB{
			User:     "user",
			Password: "password",
			Host:     "127.0.0.1",
			Port:     26257,
			Database: "defaultdb",
			Parameters: &connectors.CockroachDBParameters{
				ApplicationName: "application_name",
				Mode:            connectors.CockroachSSLDisable,
				SSLRootCert:     "path/to/ssl/root",
				SSLCert:         "path/to/ssl/cert",
				SSLKey:          "path/to/ssl/key",
				Options: &connectors.CockroachDBOptions{
					Cluster: "cluster_id",
					Variable: connectors.CockroachDBVariable{
						Name:  "name",
						Value: "value",
					},
				},
			},
		}
	)

	testCases := map[string]struct {
		path    string
		decoder func([]byte, interface{}) error
	}{
		"json": {"testdata/cockroach_db.json", func(b []byte, i interface{}) error { return json.Unmarshal(b, i) }},
		"yaml": {"testdata/cockroach_db.yaml", func(b []byte, i interface{}) error { return yaml.Unmarshal(b, i) }},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			b, err := os.ReadFile(tc.path)
			r.NoError(err)
			value := connectors.CockroachDB{}
			err = tc.decoder(b, &value)
			r.NoError(err)
			r.Equal(all, value)
		})
	}
}

func TestCockroachDB_DSN(t *testing.T) {
	t.Parallel()

	type T = connectors.CockroachDB
	change := func(t T, fn func(*T)) T {
		var parameters *connectors.CockroachDBParameters
		if t.Parameters != nil {
			p := *t.Parameters
			parameters = &p
		}

		if parameters != nil && parameters.Options != nil {
			op := *t.Parameters.Options
			parameters.Options = &op
		}
		t.Parameters = parameters

		fn(&t)
		return t
	}

	var (
		all = connectors.CockroachDB{
			User:     "user",
			Password: "password",
			Host:     "127.0.0.1",
			Port:     26257,
			Database: "defaultdb",
			Parameters: &connectors.CockroachDBParameters{
				ApplicationName: "application_name",
				Mode:            connectors.CockroachSSLDisable,
				SSLRootCert:     "path/to/ssl/root",
				SSLCert:         "path/to/ssl/cert",
				SSLKey:          "path/to/ssl/key",
				Options: &connectors.CockroachDBOptions{
					Cluster: "cluster_id",
					Variable: connectors.CockroachDBVariable{
						Name:  "name",
						Value: "value",
					},
				},
			},
		}
		allDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&options=--cluster%3Dcluster_id+-c+name%3Dvalue&sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersOptionsVariable       = change(all, func(t *T) { t.Parameters.Options.Variable = connectors.CockroachDBVariable{} })
		withoutParametersOptionsVariableDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&options=--cluster%3Dcluster_id&sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersOptionsCluster       = change(all, func(t *T) { t.Parameters.Options.Cluster = "" })
		withoutParametersOptionsClusterDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&options=-c+name%3Dvalue&sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersOptions       = change(all, func(t *T) { t.Parameters.Options = nil })
		withoutParametersOptionsDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersSSLKey       = change(all, func(t *T) { t.Parameters.Options = nil; t.Parameters.SSLKey = "" })
		withoutParametersSSLKeyDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&sslcert=path%2Fto%2Fssl%2Fcert&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersSSLCert       = change(all, func(t *T) { t.Parameters.Options = nil; t.Parameters.SSLCert = "" })
		withoutParametersSSLCertDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersSSLRoot       = change(all, func(t *T) { t.Parameters.Options = nil; t.Parameters.SSLRootCert = "" })
		withoutParametersSSLRootDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable"

		withoutParametersSSLMod       = change(all, func(t *T) { t.Parameters.Options = nil; t.Parameters.Mode = 0 })
		withoutParametersSSLModDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?application_name=application_name&sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParametersApplicationName       = change(all, func(t *T) { t.Parameters.Options = nil; t.Parameters.ApplicationName = "" })
		withoutParametersApplicationNameDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb?sslcert=path%2Fto%2Fssl%2Fcert&sslkey=path%2Fto%2Fssl%2Fkey&sslmode=disable&sslrootcert=path%2Fto%2Fssl%2Froot"

		withoutParameters       = change(all, func(t *T) { t.Parameters = nil })
		withoutParametersDSNExp = "postgres://user:password@127.0.0.1:26257/defaultdb"
	)

	testCases := map[string]struct {
		cfg T
		exp string
	}{
		"all":                                 {all, allDSNExp},
		"without_parameters_options_variable": {withoutParametersOptionsVariable, withoutParametersOptionsVariableDSNExp},
		"without_parameters_options_cluster":  {withoutParametersOptionsCluster, withoutParametersOptionsClusterDSNExp},
		"without_parameters_options":          {withoutParametersOptions, withoutParametersOptionsDSNExp},
		"without_parameters_ssl_key":          {withoutParametersSSLKey, withoutParametersSSLKeyDSNExp},
		"without_parameters_ssl_cert":         {withoutParametersSSLCert, withoutParametersSSLCertDSNExp},
		"without_parameters_ssl_root":         {withoutParametersSSLRoot, withoutParametersSSLRootDSNExp},
		"without_parameters_ssl_mod":          {withoutParametersSSLMod, withoutParametersSSLModDSNExp},
		"without_parameters_application_name": {withoutParametersApplicationName, withoutParametersApplicationNameDSNExp},
		"without_parameters":                  {withoutParameters, withoutParametersDSNExp},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			dsn, err := tc.cfg.DSN()
			r.NoError(err)
			r.Equal(tc.exp, dsn)
		})
	}
}
