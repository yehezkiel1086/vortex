package virustotal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVirusTotalClient(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-apikey") != "test-key" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"attributes": {
					"reputation": -35,
					"last_analysis_stats": {
						"harmless": 2,
						"malicious": 58,
						"suspicious": 1,
						"undetected": 10
					},
					"tags": ["trojan", "stealer"],
					"popular_threat_classification": {
						"suggested_threat_label": "trojan.redline"
					}
				}
			}
		}`))
	}))
	defer mockServer.Close()

	client := New("test-key", mockServer.URL, 2*time.Second)
	ctx := context.Background()

	t.Run("LookupHash", func(t *testing.T) {
		data, err := client.LookupHash(ctx, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.MaliciousVotes != 58 {
			t.Errorf("expected 58 malicious votes, got %d", data.MaliciousVotes)
		}
		if data.Reputation != -35 {
			t.Errorf("expected reputation -35, got %d", data.Reputation)
		}
		if data.MalwareFamily != "trojan.redline" {
			t.Errorf("expected malware family trojan.redline, got %s", data.MalwareFamily)
		}
	})
}
