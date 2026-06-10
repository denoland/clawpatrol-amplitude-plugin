package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/denoland/clawpatrol/pluginsdk"
)

const (
	pluginName        = "amplitude"
	pluginVersion     = "0.1.0"
	credentialType    = "amplitude_oauth"
	phAmplitudeOAuth  = "PH_amplitude_oauth"
	amplitudeRedirect = "http://localhost:8900/callback"
)

var amplitudeScopes = []string{"mcp:read", "mcp:write", "offline_access"}

type amplitudeConfigInput struct {
	Region string `json:"region"`
}

type amplitudeConfig struct {
	Region  string `json:"region"`
	BaseURL string `json:"base_url"`
}

func main() {
	pluginsdk.Run(&pluginsdk.Plugin{
		Name:        pluginName,
		Version:     pluginVersion,
		Credentials: []pluginsdk.CredentialDef{amplitudeOAuthDef()},
	})
}

func amplitudeOAuthDef() pluginsdk.CredentialDef {
	return pluginsdk.CredentialDef{
		TypeName:       credentialType,
		Disambiguators: []string{"placeholder"},
		HTTPInject:     true,
		Schema: pluginsdk.Schema{Fields: []pluginsdk.SchemaField{
			{Name: "region", TypeString: "string"},
		}},
		Build: func(req pluginsdk.BuildRequest) (any, error) {
			cfg, err := buildAmplitudeConfig(req.ConfigJSON)
			if err != nil {
				return nil, err
			}
			return pluginsdk.CredentialBuildResult{
				Canonical: cfg,
				Metadata: pluginsdk.CredentialMetadata{
					Disambiguators: []string{"placeholder"},
					EnvVars:        amplitudeEnvVars(cfg),
					OAuth:          amplitudeOAuthFlow(cfg),
					HTTPInject:     true,
				},
			}, nil
		},
		InjectHTTP: injectAmplitudeHTTP,
	}
}

func buildAmplitudeConfig(raw []byte) (amplitudeConfig, error) {
	var in amplitudeConfigInput
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &in); err != nil {
			return amplitudeConfig{}, err
		}
	}
	region := strings.ToLower(strings.TrimSpace(in.Region))
	if region == "" {
		region = "us"
	}
	if region != "us" && region != "eu" {
		return amplitudeConfig{}, fmt.Errorf(`region must be "us" or "eu"`)
	}
	return amplitudeConfig{Region: region, BaseURL: amplitudeMCPBaseURL(region)}, nil
}

func amplitudeMCPBaseURL(region string) string {
	if region == "eu" {
		return "https://mcp.eu.amplitude.com"
	}
	return "https://mcp.amplitude.com"
}

func amplitudeEnvVars(cfg amplitudeConfig) []pluginsdk.EnvVar {
	return []pluginsdk.EnvVar{
		{Name: "AMPLITUDE_ACCESS_TOKEN", Value: phAmplitudeOAuth, Description: "Amplitude CLI OAuth access token placeholder"},
		{Name: "AMPLITUDE_OAUTH_TOKEN", Value: phAmplitudeOAuth, Description: "Amplitude CLI OAuth access token placeholder"},
		{Name: "AMPLITUDE_REGION", Value: cfg.Region, Description: "Amplitude MCP region"},
	}
}

func amplitudeOAuthFlow(cfg amplitudeConfig) *pluginsdk.OAuthIntegration {
	return &pluginsdk.OAuthIntegration{
		Type:   "oauth2",
		Header: "Authorization",
		Prefix: "Bearer ",
		Flow:   "dynamic_mcp",
		OAuth: pluginsdk.OAuthConfig{
			AuthURL:     cfg.BaseURL + "/authorize",
			TokenURL:    cfg.BaseURL + "/token",
			RegisterURL: cfg.BaseURL + "/register",
			RedirectURI: amplitudeRedirect,
			Scopes:      append([]string(nil), amplitudeScopes...),
		},
	}
}

func injectAmplitudeHTTP(_ context.Context, req pluginsdk.HTTPInjectRequest) (*pluginsdk.HTTPInjectResponse, error) {
	accessToken := strings.TrimSpace(string(req.CredentialSecret))
	if accessToken == "" {
		return &pluginsdk.HTTPInjectResponse{}, nil
	}
	return &pluginsdk.HTTPInjectResponse{
		Headers: []pluginsdk.HeaderMutation{{
			Op:     pluginsdk.HeaderSet,
			Name:   "Authorization",
			Values: []string{"Bearer " + accessToken},
		}},
		Redactions: []string{accessToken},
	}, nil
}
