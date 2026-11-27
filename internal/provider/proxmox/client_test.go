package proxmox

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListNodes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"success", 200, `{"data":[{"node":"lab","status":"online"}]}`},
		{"denied", 403, `{}`},
		{"invalid", 200, `{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api2/json/nodes" || r.Method != http.MethodGet || r.Header.Get("Authorization") != "PVEAPIToken=user@pve!test=secret" {
					t.Error("unexpected request")
				}
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprint(w, tc.body)
			}))
			defer s.Close()
			c, err := New(s.URL, "user@pve!test", "secret")
			if err != nil {
				t.Fatal(err)
			}
			c.http = s.Client()
			nodes, err := c.ListNodes(context.Background())
			if tc.name == "success" {
				if err != nil || len(nodes) != 1 || nodes[0].Name != "lab" {
					t.Fatalf("nodes=%v error=%v", nodes, err)
				}
			} else if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestRejectHTTP(t *testing.T) {
	if _, err := New("http://localhost", "id", "secret"); err == nil {
		t.Fatal("expected HTTPS requirement")
	}
}
