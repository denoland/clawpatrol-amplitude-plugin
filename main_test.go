package main

import (
	"context"
	"testing"

	"github.com/denoland/clawpatrol/pluginsdk"
)

func TestAmplitudeBuildEUOAuthMetadata(t *testing.T) {
	out, err := amplitudeOAuthDef().Build(pluginsdk.BuildRequest{ConfigJSON: []byte(`{"region":"eu"}`)})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	res := out.(pluginsdk.CredentialBuildResult)
	cfg := res.Canonical.(amplitudeConfig)
	if cfg.Region != "eu" || cfg.BaseURL != "https://mcp.eu.amplitude.com" {
		t.Fatalf("cfg = %#v", cfg)
	}
	if !res.Metadata.HTTPInject {
		t.Fatalf("HTTPInject = false")
	}
	if res.Metadata.OAuth == nil {
		t.Fatalf("OAuth metadata missing")
	}
	flow := res.Metadata.OAuth
	if flow.Flow != "dynamic_mcp" || flow.Type != "oauth2" || flow.Header != "Authorization" || flow.Prefix != "Bearer " {
		t.Fatalf("flow = %#v", flow)
	}
	if flow.OAuth.AuthURL != "https://mcp.eu.amplitude.com/authorize" {
		t.Fatalf("AuthURL = %q", flow.OAuth.AuthURL)
	}
	if flow.OAuth.TokenURL != "https://mcp.eu.amplitude.com/token" {
		t.Fatalf("TokenURL = %q", flow.OAuth.TokenURL)
	}
	if flow.OAuth.RegisterURL != "https://mcp.eu.amplitude.com/register" {
		t.Fatalf("RegisterURL = %q", flow.OAuth.RegisterURL)
	}
	if flow.OAuth.RedirectURI != amplitudeRedirect {
		t.Fatalf("RedirectURI = %q", flow.OAuth.RedirectURI)
	}
	wantScopes := []string{"mcp:read", "mcp:write", "offline_access"}
	if len(flow.OAuth.Scopes) != len(wantScopes) {
		t.Fatalf("scopes = %#v", flow.OAuth.Scopes)
	}
	for i, want := range wantScopes {
		if flow.OAuth.Scopes[i] != want {
			t.Fatalf("scope[%d] = %q, want %q", i, flow.OAuth.Scopes[i], want)
		}
	}
	vars := map[string]string{}
	for _, ev := range res.Metadata.EnvVars {
		vars[ev.Name] = ev.Value
	}
	if vars["AMPLITUDE_ACCESS_TOKEN"] != phAmplitudeOAuth || vars["AMPLITUDE_OAUTH_TOKEN"] != phAmplitudeOAuth || vars["AMPLITUDE_REGION"] != "eu" {
		t.Fatalf("env vars = %#v", vars)
	}
}

func TestAmplitudeDefaultsToUS(t *testing.T) {
	cfg, err := buildAmplitudeConfig(nil)
	if err != nil {
		t.Fatalf("buildAmplitudeConfig: %v", err)
	}
	if cfg.Region != "us" || cfg.BaseURL != "https://mcp.amplitude.com" {
		t.Fatalf("cfg = %#v", cfg)
	}
}

func TestAmplitudeRejectsInvalidRegion(t *testing.T) {
	_, err := buildAmplitudeConfig([]byte(`{"region":"moon"}`))
	if err == nil {
		t.Fatalf("expected invalid region error")
	}
}

func TestAmplitudeInjectHTTPSetsBearer(t *testing.T) {
	out, err := injectAmplitudeHTTP(context.Background(), pluginsdk.HTTPInjectRequest{
		CredentialSecret: []byte("amp-access-token"),
	})
	if err != nil {
		t.Fatalf("InjectHTTP: %v", err)
	}
	if len(out.Headers) != 1 || out.Headers[0].Name != "Authorization" || out.Headers[0].Values[0] != "Bearer amp-access-token" {
		t.Fatalf("headers = %#v", out.Headers)
	}
	if len(out.Redactions) != 1 || out.Redactions[0] != "amp-access-token" {
		t.Fatalf("redactions = %#v", out.Redactions)
	}
}
