# Twilio Cost Exporter

A Prometheus exporter for Twilio usage costs.

## Configuration

The exporter is configured via environment variables:

| Variable | Description |
|Str|Str|
| `TWILIO_ACCOUNT_ID` | Your Twilio Account SID (starts with AC) |
| `TWILIO_SID` | Your Twilio API Key SID (starts with SK) |
| `TWILIO_SECRET` | Your Twilio API Key Secret |
| `PORT` | The port to listen on (default: 8080) |

## Metrics

The exporter exposes a single metric `twilio_usage` with the following labels:

- `category`: The usage category (e.g., `sms-outbound`, `calls-inbound`).
- `subresource`: The time period (`Today`, `Yesterday`, `ThisMonth`, `LastMonth`).
- `unit`: The unit of measurement (`usd`, `messages`, `minutes`, etc.).

Example:

```
twilio_usage{category="sms-outbound",subresource="Today",unit="usd"} 0.05
```

## Running

```bash
# Build
make build

# Run
make run
```

## License

MIT
