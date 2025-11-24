# Heroku Cost Exporter

**Generated using [Antigravity](https://antigravity.dev)**

A Prometheus exporter that collects cost and usage data from Heroku, providing metrics for invoices and team usage.

## Features

- **Invoice Metrics**: Latest invoice details including total, charges, credits, and add-ons.
- **Team Usage Metrics**: Monthly team usage costs broken down by app and expense type (dynos, data, addons, etc.).
- **Benchmarking**: Separate metrics for the **current month** (ramping up) and **previous month** (final) to allow for easy comparison.
- **Prometheus Integration**: Exposes metrics at `/metrics` endpoint.
- **Health Check**: Health endpoint at `/health`.
- **Automatic Updates**: Metrics refresh every hour.

## Metrics

### Invoice Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `heroku_invoice_current_total_usd` | Gauge | Current (partial) Heroku invoice total in USD | `payment_status`, `state`, `period_end` |
| `heroku_invoice_current_charges_total_usd` | Gauge | Current Heroku invoice charges total in USD | `period_end` |
| `heroku_invoice_current_credits_total_usd` | Gauge | Current Heroku invoice credits total in USD | `period_end` |
| `heroku_invoice_current_addons_total_usd` | Gauge | Current Heroku invoice addons total in USD | `period_end` |
| `heroku_invoice_current_database_total_usd` | Gauge | Current Heroku invoice database total in USD | `period_end` |
| `heroku_invoice_current_platform_total_usd` | Gauge | Current Heroku invoice platform total in USD | `period_end` |
| `heroku_invoice_current_dyno_units` | Gauge | Current Heroku invoice dyno units | `period_end` |
| `heroku_invoice_current_weighted_dyno_hours` | Gauge | Current Heroku invoice weighted dyno hours | `period_end` |
| `heroku_invoice_previous_total_usd` | Gauge | Previous (complete) Heroku invoice total in USD | `payment_status`, `state`, `period_end` |
| `heroku_invoice_previous_charges_total_usd` | Gauge | Previous Heroku invoice charges total in USD | `period_end` |
| `heroku_invoice_previous_credits_total_usd` | Gauge | Previous Heroku invoice credits total in USD | `period_end` |
| `heroku_invoice_previous_addons_total_usd` | Gauge | Previous Heroku invoice addons total in USD | `period_end` |
| `heroku_invoice_previous_database_total_usd` | Gauge | Previous Heroku invoice database total in USD | `period_end` |
| `heroku_invoice_previous_platform_total_usd` | Gauge | Previous Heroku invoice platform total in USD | `period_end` |
| `heroku_invoice_previous_dyno_units` | Gauge | Previous Heroku invoice dyno units | `period_end` |
| `heroku_invoice_previous_weighted_dyno_hours` | Gauge | Previous Heroku invoice weighted dyno hours | `period_end` |

### Team Usage Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `heroku_team_usage_current_usd` | Gauge | Current month Heroku team usage cost in USD | `team_name`, `app_name`, `type`, `month` |
| `heroku_team_usage_previous_usd` | Gauge | Previous month Heroku team usage cost in USD | `team_name`, `app_name`, `type`, `month` |
| `heroku_team_usage_current_dyno_units` | Gauge | Current month Heroku team usage dyno units | `team_name`, `app_name`, `month` |
| `heroku_team_usage_previous_dyno_units` | Gauge | Previous month Heroku team usage dyno units | `team_name`, `app_name`, `month` |

**Note**: 
- For USD metrics: `app_name="_team_total"` represents the total cost for the team for that specific type.
- For USD metrics: `type` can be `data`, `addons`, `connect`, or `partner` (dynos are tracked separately in dyno_units metrics).
- For dyno_units metrics: Values represent dyno-months (e.g., 0.368 = ~11 days of a single dyno).

## Prerequisites

- Heroku API Key
- Heroku Team ID
- Go 1.23+ (for building from source)

## Installation

### Using Docker

```bash
docker run -p 8080:8080 \
  -e HEROKU_API_KEY=your_api_key \
  -e HEROKU_TEAM_ID=your_team_id \
  ghcr.io/loke/heroku-cost-exporter:latest
```

### Building from Source

```bash
git clone https://github.com/LOKE/heroku-cost-exporter.git
cd heroku-cost-exporter
go build -o heroku-cost-exporter .
./heroku-cost-exporter
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | HTTP server port (default: 8080) | No |
| `HEROKU_API_KEY` | Your Heroku API Key | Yes |
| `HEROKU_TEAM_ID` | The ID of the Heroku Team to track | Yes |

## Caveats

- **Invoice Selection**: The exporter tracks both the **current** (partial/open) invoice and the **previous** (complete) invoice. It assumes the most recent invoice in the API response is the current one, and the second-to-last is the previous one.
- **Current Invoice**: Represents the ongoing billing cycle and will change throughout the month as usage accumulates.
- **Previous Invoice**: Represents the finalized invoice from the last billing cycle, providing stable data for historical tracking.

## Usage

1. Start the exporter:
   ```bash
   export HEROKU_API_KEY="your-key"
   export HEROKU_TEAM_ID="your-team-id"
   ./heroku-cost-exporter
   ```

2. Check health:
   ```bash
   curl http://localhost:8080/health
   ```

3. View metrics:
   ```bash
   curl http://localhost:8080/metrics
   ```

4. Configure Prometheus to scrape the metrics:
   ```yaml
   scrape_configs:
     - job_name: 'heroku-cost-exporter'
       static_configs:
         - targets: ['localhost:8080']
   ```

## Development

### Make Targets

```bash
make build    # Build the application
make run      # Run the application
make test     # Run tests
make tidy     # Clean up dependencies
make clean    # Clean build artifacts
make help     # Show available targets
```

## License

This project is licensed under the MIT License.
