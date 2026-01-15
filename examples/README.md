# FindLabs API Go SDK Examples

This directory contains practical examples demonstrating how to use the FindLabs API Go SDK.

## Directory Structure

```
examples/
├── auth/
│   └── generate-token/     # Generate JWT tokens with Basic Auth
└── simple/
    ├── get-blocks/          # Fetch blocks by height
    ├── get-events/          # Query events by name and block range
    ├── get-transaction/     # Get transaction details by ID
    └── get-transaction-events/ # Get events for a specific transaction
```

## Prerequisites

Before running any examples, you need to set your FindLabs API credentials as environment variables:

```bash
export FINDAPI_USERNAME="your-username"
export FINDAPI_PASSWORD="your-password"
```

## Running Examples

Each example is a standalone Go program. Navigate to the example directory and run it:

```bash
# Example: Generate JWT token
cd examples/auth/generate-token
go run main.go

# Example: Get blocks
cd examples/simple/get-blocks
go run main.go

# Example: Get events
cd examples/simple/get-events
go run main.go

# Example: Get transaction details
cd examples/simple/get-transaction
go run main.go

# Example: Get transaction events
cd examples/simple/get-transaction-events
go run main.go
```

## Examples Overview

### Auth Service

#### generate-token
Demonstrates how to manually generate JWT tokens with different expiry durations.

**Key Features:**
- Generate tokens with custom expiry (1 hour, 24 hours)
- Display token metadata (expiry time, issued time)
- Format Unix timestamps into human-readable dates

### Simple Service

#### get-blocks
Fetch blocks at a specific height with optional pagination offset.

**Key Features:**
- Retrieve blocks by height
- Use pagination with offset parameter
- Display block metadata and transaction IDs

#### get-events
Query events by name within a block height range.

**Key Features:**
- Filter events by event name
- Specify block height range
- Paginate through large result sets
- Pretty-print event field data

#### get-transaction
Retrieve detailed transaction information by transaction ID.

**Key Features:**
- Get complete transaction details
- Display gas usage and fees
- Show transaction participants (payer, proposer, authorizers)
- View transaction events
- Access transaction body (Cadence script)

#### get-transaction-events
Get all events emitted by a specific transaction.

**Key Features:**
- Fetch events for a transaction
- Pagination support
- Filter and count events by type
- Pretty-print event fields

## Common Patterns

### Authentication
All examples use automatic JWT authentication. The client handles:
- Initial token generation
- Token caching
- Automatic token refresh before expiry
- Thread-safe token management

```go
client := findapi.NewClient(username, password)
// Token management is automatic from here
```

### Error Handling
Examples demonstrate proper error handling:

```go
result, err := client.Simple.GetBlocks(ctx, height, nil)
if err != nil {
    log.Fatalf("Failed to get blocks: %v", err)
}
```

### Pagination
Many endpoints support pagination for large datasets:

```go
offset := 0
for {
    results, err := client.Simple.GetEvents(ctx, eventName, from, to, &offset)
    if err != nil {
        log.Fatal(err)
    }

    if len(results.Events) == 0 {
        break
    }

    // Process results...
    offset += len(results.Events)
}
```

### Context Usage
All API calls accept a context for timeout and cancellation:

```go
ctx := context.Background()
// Or with timeout:
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.Simple.GetBlocks(ctx, height, nil)
```

## Customization

You can modify the examples to:
- Query different block heights or transaction IDs
- Change event names and block ranges
- Adjust pagination parameters
- Add custom error handling or retry logic
- Format output differently

## Rate Limiting

The SDK automatically handles rate limiting (HTTP 429 responses):
- Respects `Retry-After` headers
- Retries up to 3 times with appropriate delays
- Returns `RateLimitError` if retries are exhausted

## Tips

1. **Start with small ranges**: When querying events, start with small block ranges to avoid timeouts
2. **Use pagination**: For large datasets, use offset-based pagination to fetch results in chunks
3. **Cache tokens**: The client automatically caches JWT tokens, but you can also generate and store them manually for use in other applications
4. **Monitor rate limits**: If you're making many requests, implement exponential backoff on `RateLimitError`

## Next Steps

After trying these examples:
1. Explore the SDK documentation in the main [README.md](../../README.md)
2. Check the [API reference](https://docs.findlabs.com) for available endpoints
3. Build your own applications using these patterns as a foundation

## Troubleshooting

**Authentication errors:**
- Verify your credentials are correct
- Check that environment variables are set properly
- Ensure your API account has necessary permissions

**Rate limit errors:**
- Reduce request frequency
- Implement exponential backoff
- Consider caching responses when appropriate

**Connection timeouts:**
- Check your network connectivity
- Verify the API endpoint is accessible
- Consider increasing context timeout for large queries
