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

# Add this credential to the profile used by your agent, for example:
# profile "mario" {
#   credentials = [amplitude_oauth.amplitude]
# }

rule "amplitude-mcp" {
  endpoint = https.amplitude_mcp
  verdict  = "allow"
}
