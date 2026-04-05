# Stockyard Booking

**Self-hosted appointment booking and scheduling**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted tools.

## Quick Start

```bash
curl -fsSL https://stockyard.dev/tools/booking/install.sh | sh
```

Or with Docker:

```bash
docker run -p 9800:9800 -v booking_data:/data ghcr.io/stockyard-dev/stockyard-booking
```

Open `http://localhost:9800` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9800` | HTTP port |
| `DATA_DIR` | `./booking-data` | SQLite database directory |
| `STOCKYARD_LICENSE_KEY` | *(empty)* | License key for unlimited use |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 5 records | Unlimited |
| Price | Free | Included in bundle or $29.99/mo individual |

Get a license at [stockyard.dev](https://stockyard.dev).

## License

Apache 2.0
