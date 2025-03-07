package rest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRestConfig(t *testing.T) {
	tests := []struct {
		name       string
		kubeConfig string
		overrides  clientcmd.ConfigOverrides
		wantErr    bool
		want       *rest.Config
	}{{
		name:       "no cluster",
		kubeConfig: ".",
		wantErr:    true,
	}, {
		name:       "empty",
		kubeConfig: "../../../testdata/.kube/config",
		want: &rest.Config{
			Host:  "https://127.0.0.1:53742",
			QPS:   300,
			Burst: 300,
		},
	}, {
		name:       "context override",
		kubeConfig: "../../../testdata/.kube/config",
		overrides: clientcmd.ConfigOverrides{
			CurrentContext: "foo",
		},
		want: &rest.Config{
			Host:  "https://127.0.0.1:1234",
			QPS:   300,
			Burst: 300,
		},
	}, {
		name:       "timeout override",
		kubeConfig: "../../../testdata/.kube/config",
		overrides: clientcmd.ConfigOverrides{
			Timeout: "30s",
		},
		want: &rest.Config{
			Host:    "https://127.0.0.1:53742",
			QPS:     300,
			Burst:   300,
			Timeout: 30 * time.Second,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("KUBECONFIG", tt.kubeConfig)
			got, err := DefaultConfig(tt.overrides)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, got) {
					assert.Equal(t, tt.want.Host, got.Host)
					assert.Equal(t, tt.want.APIPath, got.APIPath)
					assert.Equal(t, tt.want.Username, got.Username)
					assert.Equal(t, tt.want.Password, got.Password)
					assert.Equal(t, tt.want.BearerToken, got.BearerToken)
					assert.Equal(t, tt.want.BearerTokenFile, got.BearerTokenFile)
					assert.Equal(t, tt.want.UserAgent, got.UserAgent)
					assert.Equal(t, tt.want.DisableCompression, got.DisableCompression)
					assert.Equal(t, tt.want.QPS, got.QPS)
					assert.Equal(t, tt.want.Burst, got.Burst)
					assert.Equal(t, tt.want.Timeout, got.Timeout)
				}
			}
		})
	}
}
