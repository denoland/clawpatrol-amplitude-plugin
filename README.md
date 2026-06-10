# Claw Patrol Amplitude credential plugin

External Claw Patrol credential plugin for Amplitude's hosted MCP OAuth flow.

This is a proof of concept for Claw Patrol's external credential metadata: the plugin declares a dynamic MCP OAuth flow, pushes Amplitude CLI placeholders into wrapped agent environments, and injects the gateway-held OAuth access token into intercepted Amplitude MCP HTTPS requests.

## Credential type

```hcl
plugin "amplitude" {
  source = "/opt/clawpatrol/plugins/clawpatrol-amplitude-plugin"
}

endpoint "https" "amplitude_mcp" {
  hosts = ["mcp.eu.amplitude.com"]
}

credential "amplitude_oauth" "amplitude" {
  endpoint    = https.amplitude_mcp
  region      = "eu"
  placeholder = "PH_amplitude_oauth"
}
```

See `examples/amplitude.hcl` for a minimal gateway snippet.

## Pushed environment variables

- `AMPLITUDE_ACCESS_TOKEN=PH_amplitude_oauth`
- `AMPLITUDE_OAUTH_TOKEN=PH_amplitude_oauth`
- `AMPLITUDE_REGION=<region>`

## Build

Until the external credential runtime lands upstream, build this next to a Claw Patrol checkout that contains the `pluginsdk` changes:

```bash
git clone https://github.com/dhruvkelawala/clawpatrol.git
cd clawpatrol
git checkout external-http-credential-runtime
cd ..

git clone https://github.com/dhruvkelawala/clawpatrol-amplitude-plugin.git
cd clawpatrol-amplitude-plugin
go test ./...
go build -o clawpatrol-amplitude-plugin .
```

The temporary `replace github.com/denoland/clawpatrol => ../clawpatrol` in `go.mod` is intentional for this POC.

## Gateway install sketch

```bash
sudo install -m 0755 clawpatrol-amplitude-plugin /opt/clawpatrol/plugins/clawpatrol-amplitude-plugin
```

Then connect the OAuth credential through the Claw Patrol dashboard. Do not put OAuth access tokens, refresh tokens, or copied callback URLs in HCL or environment files.
