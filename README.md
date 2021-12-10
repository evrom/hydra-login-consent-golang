# Hydra Login Consent Golang

This is a WIP re-implementation of https://github.com/ory/hydra-login-consent-node in golang. 

## How to run

### starting Hydra
- clone https://github.com/ory/hydra
- edit https://github.com/ory/hydra/blob/21b470dce2df5495484701d009b1aabe136d4c28/quickstart.yml#L68 to  `- "3001:3000"`, so the Login Consent Golang can run on port 3000 instead of the Login Consent Node
- Start Hydra as explained in the quickstart guide here: https://www.ory.sh/hydra/docs/5min-tutorial
- then run `go run main.go` in the VC root of this repository
- also run the example server from https://github.com/layeh/radius
