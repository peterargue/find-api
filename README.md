# Find API Go SDK

A Go client library for the Find Flow blockchain data API. This SDK provides a clean, idiomatic Go interface for accessing Flow blockchain data with built-in authentication, rate limiting, and error handling.

## Features

- **Fluent Builder API**: Intuitive builder pattern for constructing requests
- **Automatic JWT Authentication**: Seamless username/password authentication with automatic token refresh
- **Rate Limiting**: Built-in detection and handling of rate limits (HTTP 429) with automatic retries
- **Context Support**: All methods accept `context.Context` for cancellation and timeout control
- **Typed Responses**: Strongly typed response structures matching Flow blockchain data
- **Service-based Organization**: Clean API with dedicated services (Simple, Auth, etc.)
- **Comprehensive Error Handling**: Typed errors with blockchain-specific error codes

## Installation

```bash
go get github.com/peterargue/find-api
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    findapi "github.com/peterargue/find-api"
)

func main() {
    // Create a new client
    client := findapi.NewClient("your-username", "your-password")

    // Create a context
    ctx := context.Background()

    // Get blocks at a specific height using the fluent builder API
    blocks, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, block := range blocks.Blocks {
        fmt.Printf("Block %d has %d transactions\n", block.Height, block.TxCount)
    }
}
```

For more detailed examples, see the [examples/](examples/) directory.

## Examples

The [examples/](examples/) directory contains runnable examples for all API endpoints:

- **Auth**: [Generate JWT tokens](examples/auth/generate-token/)
- **Simple**:
  - [Get blocks by height](examples/simple/get-blocks/)
  - [Query events by name](examples/simple/get-events/)
  - [Get transaction details](examples/simple/get-transaction/)
  - [Get transaction events](examples/simple/get-transaction-events/)

Each example includes detailed comments and demonstrates best practices.

## Configuration

### Custom Base URL

```go
client := findapi.NewClient(
    "username",
    "password",
    findapi.WithBaseURL("https://custom-api.example.com"),
)
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

client := findapi.NewClient(
    "username",
    "password",
    findapi.WithHTTPClient(httpClient),
)
```

## Simple API Endpoints

The Simple API uses a fluent builder pattern for constructing requests. All builders have a `Do(ctx)` method to execute the request.

### Get Blocks

Retrieve blocks based on height and offset:

```go
// Get blocks at a specific height
blocks, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
if err != nil {
    log.Fatal(err)
}

// With offset for pagination
blocks, err := client.Simple.GetBlocks().Height(96708412).Offset(10).Do(ctx)
```

### Get Events

Retrieve events of a specific name within a block height range:

```go
events, err := client.Simple.GetEvents().
    Name("A.921ea449dffec68a.FlovatarMarketplace.FlovatarPriceChanged").
    FromHeight(102968960).
    ToHeight(103850311).
    Do(ctx)
if err != nil {
    log.Fatal(err)
}

for _, event := range events.Events {
    fmt.Printf("Event at block %d: %s\n", event.BlockHeight, event.Name)
}

// With offset for pagination
events, err := client.Simple.GetEvents().
    Name("A.921ea449dffec68a.FlovatarMarketplace.FlovatarPriceChanged").
    FromHeight(102968960).
    ToHeight(103850311).
    Offset(100).
    Do(ctx)
```

### Get Transaction

Retrieve a transaction by its ID:

```go
tx, err := client.Simple.GetTransaction().
    ID("b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562").
    Do(ctx)
if err != nil {
    log.Fatal(err)
}

for _, t := range tx.Transactions {
    fmt.Printf("Status: %s, Fee: %.8f\n", t.Status, t.Fee)
}
```

### Get Transaction Events

Retrieve events for a specific transaction:

```go
events, err := client.Simple.GetTransactionEvents().
    TransactionID("b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562").
    Do(ctx)
if err != nil {
    log.Fatal(err)
}

for _, event := range events.Events {
    fmt.Printf("Event: %s\n", event.Name)
}

// With offset for pagination
events, err := client.Simple.GetTransactionEvents().
    TransactionID("b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562").
    Offset(10).
    Do(ctx)
```

## Error Handling

The SDK provides typed errors for better error handling:

```go
blocks, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
if err != nil {
    // Check for rate limit errors
    if findapi.IsRateLimitError(err) {
        log.Printf("Rate limited: %v", err)
        // Implement custom backoff logic
        return
    }

    // Check for API errors
    if findapi.IsAPIError(err) {
        apiErr := err.(*findapi.APIError)
        log.Printf("API error (status %d): %s", apiErr.StatusCode, apiErr.Message)
        return
    }

    log.Fatal(err)
}
```

## Rate Limiting

The SDK automatically handles rate limiting:

- Detects HTTP 429 responses
- Respects `Retry-After` headers
- Automatically retries up to 3 times with appropriate delays
- Returns a `RateLimitError` if all retries are exhausted

```go
blocks, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
if err != nil {
    if rateLimitErr, ok := err.(*findapi.RateLimitError); ok {
        // Wait for the suggested duration before retrying
        time.Sleep(rateLimitErr.RetryAfter)
    }
}
```

## Pagination

For endpoints that support pagination, use the `Offset()` builder method:

```go
offset := 0
for {
    events, err := client.Simple.GetEvents().
        Name("EventName").
        FromHeight(100).
        ToHeight(200).
        Offset(offset).
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if len(events.Events) == 0 {
        break
    }

    // Process events...

    offset += len(events.Events)

    // Check if we got less than max (100), indicating last page
    if len(events.Events) < 100 {
        break
    }
}
```

## Authentication

JWT authentication is handled automatically:

- Tokens are generated on first request using Basic Auth (username/password)
- Tokens are cached and reused for subsequent requests
- Tokens are automatically refreshed before expiration (1-minute buffer)
- Token refresh is thread-safe with mutex locking

You can also manually generate a token:

```go
token, err := client.Auth.GenerateToken(ctx, time.Hour)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Token expires at: %d\n", token.Exp)
```

## Testing

Run the test suite:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

## Project Structure

```
.
├── client.go           # Main client with JWT token management
├── client_test.go      # Client tests (rate limiting, etc.)
├── errors.go          # Error types
├── example_test.go    # Usage examples
├── auth/              # Auth API module
│   ├── auth.go        # Auth API service (token generation)
│   └── auth_test.go   # Unit tests
└── simple/            # Simple API module
    ├── simple.go      # Simple API service
    └── simple_test.go # Unit tests with mocked responses
```

The SDK uses a modular structure where each API group (Auth, Simple, Accounting, Staking, etc.) is in its own subpackage. This makes it easy to:
- Add new API groups without cluttering the main package
- Test each service independently
- Maintain clear separation of concerns
- Keep authentication logic (JWT management) separate from auth API endpoints

## Adding New API Groups

The SDK is designed to make adding new API groups straightforward. To add a new API group (e.g., Accounting):

1. Create a new subdirectory: `accounting/`
2. Create `accounting/accounting.go` with the service implementation
3. Implement the `Client` interface for making requests
4. Add the service to the main `Client` struct in `client.go`
5. Initialize it in `NewClient()`

Example structure:
```go
// accounting/accounting.go
package accounting

import (
    "context"
    "net/http"
    "net/url"
)

// Client interface for making HTTP requests
type Client interface {
    DoRequest(ctx context.Context, method, path string, query url.Values) (*http.Response, error)
    DecodeResponse(resp *http.Response, v any) error
}

type Service struct {
    client Client
}

func NewService(client Client) *Service {
    return &Service{client: client}
}

// Add methods for accounting endpoints...
func (s *Service) GetAccount(ctx context.Context, address string) (*AccountResponse, error) {
    // Implementation...
}
```

Then in `client.go`:
```go
import "github.com/peterargue/find-api/accounting"

type Client struct {
    // ...
    Simple     *simple.Service
    Auth       *auth.Service
    Accounting *accounting.Service  // Add this
}

func NewClient(...) *Client {
    // ...
    c.Accounting = accounting.NewService(c)  // Add this
    // ...
}
```

## Future Enhancements

- Support for additional API groups (Accounting, Staking, Bulk, etc.)
- Webhook support
- Streaming event subscriptions
- Enhanced retry strategies with exponential backoff
- Metrics and logging hooks

## License

Apache 2.0

## Support

For issues and questions, please open an issue on the GitHub repository.
