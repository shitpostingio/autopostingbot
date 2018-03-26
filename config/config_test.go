package config

import (
	"reflect"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    Config
		wantErr bool
	}{
		{
			"a nonexistent configuration file",
			"/this/path/does/not.exists",
			Config{},
			true,
		},
		{
			"a well-formed configuration file",
			"./testing/complete_config.toml",
			Config{
				BotToken:         "000000000:abcdefghijklmnopqrstuvwxyz123456789",
				Port:             88,
				ReverseProxy:     false,
				ReverseProxyPort: 1234,
				IP:               "127.0.0.1",
				Domain:           "shitposting.io",
				TLS:              false,
				TLSCertPath:      "/path/to/the/tls/certificate",
				TLSKeyPath:       "/path/to/the/tls/key",
			},
			false,
		},
		{
			"a configuration file where ReverseProxy is false, and Port is invalid",
			"./testing/reverseproxy_false_port_malformed.toml",
			Config{},
			true,
		},
		{
			"an empty config file",
			"./testing/empty.toml",
			Config{},
			true,
		},
		{
			"missing BotToken",
			"./testing/missing_bottoken.toml",
			Config{},
			true,
		},
		{
			"ReverseProxy is true, but no ReverseProxyPort has been declared",
			"./testing/reverseproxy_true_missing_reverseproxyport.toml",
			Config{},
			true,
		},
		{
			"missing Domain",
			"./testing/missing_domain.toml",
			Config{},
			true,
		},
		{
			"Domain has http://",
			"./testing/domain_has_http.toml",
			Config{},
			true,
		},
		{
			"Domain has https://",
			"./testing/domain_has_https.toml",
			Config{},
			true,
		},
		{
			"wants TLS, but certificate is missing",
			"./testing/has_tls_missing_cert.toml",
			Config{},
			true,
		},
		{
			"wants TLS, but key is missing",
			"./testing/has_tls_missing_key.toml",
			Config{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfigFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_BindString(t *testing.T) {
	type fields struct {
		BotToken         string
		Port             int
		ReverseProxy     bool
		ReverseProxyPort int
		IP               string
		Domain           string
		TLS              bool
		TLSCertPath      string
		TLSKeyPath       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "a complete, working Config",
			fields: fields{
				BotToken:     "thisisatoken",
				Port:         88,
				ReverseProxy: false,
				IP:           "127.0.0.2",
				Domain:       "memes.inc",
				TLS:          false,
			},
			want: "127.0.0.2:88",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				BotToken:         tt.fields.BotToken,
				Port:             tt.fields.Port,
				ReverseProxy:     tt.fields.ReverseProxy,
				ReverseProxyPort: tt.fields.ReverseProxyPort,
				IP:               tt.fields.IP,
				Domain:           tt.fields.Domain,
				TLS:              tt.fields.TLS,
				TLSCertPath:      tt.fields.TLSCertPath,
				TLSKeyPath:       tt.fields.TLSKeyPath,
			}
			if got := c.BindString(); got != tt.want {
				t.Errorf("Config.BindString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WebHookURL(t *testing.T) {
	// mainly a placeholder, better testing should be written
	type fields struct {
		BotToken         string
		Port             int
		ReverseProxy     bool
		ReverseProxyPort int
		IP               string
		Domain           string
		TLS              bool
		TLSCertPath      string
		TLSKeyPath       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "a complete, working Config",
			fields: fields{
				BotToken:     "thisisatoken",
				Port:         88,
				ReverseProxy: false,
				IP:           "127.0.0.2",
				Domain:       "memes.inc",
				TLS:          false,
			},
			want: "https://memes.inc:88/thisisatoken/updates",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				BotToken:         tt.fields.BotToken,
				Port:             tt.fields.Port,
				ReverseProxy:     tt.fields.ReverseProxy,
				ReverseProxyPort: tt.fields.ReverseProxyPort,
				IP:               tt.fields.IP,
				Domain:           tt.fields.Domain,
				TLS:              tt.fields.TLS,
				TLSCertPath:      tt.fields.TLSCertPath,
				TLSKeyPath:       tt.fields.TLSKeyPath,
			}
			if got := c.WebHookURL(); got != tt.want {
				t.Errorf("Config.WebHookURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isStandardPort(t *testing.T) {
	tests := []struct {
		name string
		port int
		want bool
	}{
		{
			"port 443",
			443,
			true,
		},
		{
			"port 80",
			80,
			true,
		},
		{
			"port 88",
			88,
			true,
		},
		{
			"port 8443",
			8443,
			true,
		},
		{
			"port that is not standard: 60444",
			60444,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStandardPort(tt.port); got != tt.want {
				t.Errorf("isStandardPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
