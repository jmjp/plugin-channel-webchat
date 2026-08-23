# plugin-channel-webchat

Plugin de canal (channel) do Fluxa que implementa `ChannelPluginServer`.

Recebe webhooks HTTP de clientes webchat (payload JSON simples:
`{"customer_ref": "...", "text": "..."}`) e expõe o contrato
normalizado `ChannelPluginServer` via gRPC Unix socket para o core.

## Build

```bash
go build -o bin/plugin-channel-webchat ./cmd/plugin-channel-webchat
```

## Capability

`channel.webchat` — expõe um endpoint HTTP local para clientes webchat
e um Unix socket gRPC para o core Fluxa.
