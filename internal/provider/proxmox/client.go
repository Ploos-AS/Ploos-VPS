package proxmox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Ploos-AS/Ploos-VPS/internal/provider"
)

type Client struct {
	base  string
	token string
	http  *http.Client
}

var _ provider.Inventory = (*Client)(nil)

func New(endpoint, tokenID, secret string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return nil, errors.New("expected HTTPS Proxmox origin")
	}
	if tokenID == "" || secret == "" || strings.ContainsAny(tokenID+secret, "\r\n") {
		return nil, errors.New("invalid API token")
	}
	return &Client{
		base:  strings.TrimRight(endpoint, "/") + "/api2/json",
		token: "PVEAPIToken=" + tokenID + "=" + secret,
		http: &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (c *Client) ListNodes(ctx context.Context) ([]provider.Node, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/nodes", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Proxmox request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Proxmox returned HTTP %d", resp.StatusCode)
	}
	var payload struct {
		Data []struct {
			Node   string `json:"node"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&payload); err != nil {
		return nil, errors.New("invalid Proxmox response")
	}
	if payload.Data == nil {
		return nil, errors.New("missing Proxmox node data")
	}
	nodes := make([]provider.Node, 0, len(payload.Data))
	for _, n := range payload.Data {
		if n.Node == "" {
			return nil, errors.New("missing node name")
		}
		nodes = append(nodes, provider.Node{Name: n.Node, Status: n.Status})
	}
	return nodes, nil
}
