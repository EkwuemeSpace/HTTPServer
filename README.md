# Go HTTP Server — 7 Endpoints

A minimal HTTP server built from scratch using only Go's standard library (`net/http`), covering routing, query parameters, request bodies, validation, headers, simple auth, and redirects.

## Requirements

- Go 1.18+ (no external dependencies)

## Running the server

```bash
go run main.go
```

Server starts on `http://localhost:8080`.

## Endpoints

| Method | Path         | Description                                                                 |
|--------|--------------|-------------------------------------------------------------------------------|
| GET    | `/ping`      | Returns `pong`                                                               |
| GET    | `/hello`     | Returns `Hello, {name}!` using `?name=` query param, defaults to `Guest`      |
| GET    | `/count`     | Returns instructions to POST text                                            |
| POST   | `/count`     | Returns character count of the request body                                 |
| GET    | `/calculate` | Computes `a` `op` `b` (`add`/`subtract`/`multiply`) via query params         |
| GET    | `/agent`     | Echoes back the client's `User-Agent` header                                |
| GET    | `/dashboard` | Protected route, requires header `X-Api-Key: secret123`                      |
| GET    | `/legacy`    | 301 redirect to `/v2`                                                       |
| GET    | `/v2`        | Returns a welcome message                                                   |

## Status codes used

| Code | Meaning              | Where                                                       |
|------|----------------------|--------------------------------------------------------------|
| 200  | OK                   | Default success response                                    |
| 301  | Moved Permanently    | `/legacy` → `/v2`                                            |
| 400  | Bad Request          | `/calculate`: invalid number or unknown operation            |
| 401  | Unauthorized         | `/dashboard`: missing or incorrect `X-Api-Key`                |
| 404  | Not Found            | Any unregistered path (handled automatically by `net/http`)  |
| 405  | Method Not Allowed   | Registered path hit with an unsupported HTTP method           |

## Example requests

```bash
curl http://localhost:8080/ping
curl "http://localhost:8080/hello?name=Alice"
curl -X POST -d "Golang" http://localhost:8080/count
curl "http://localhost:8080/calculate?op=add&a=12&b=8"
curl -H "User-Agent: CustomTester/1.0" http://localhost:8080/agent
curl -H "X-Api-Key: secret123" http://localhost:8080/dashboard
curl -iL http://localhost:8080/legacy
```

## Testing

An automated test script (`test_endpoints.sh`) verifies all 7 endpoints:

```bash
chmod +x test_endpoints.sh
./test_endpoints.sh
```

All 14 checks currently pass.

## Notes / things learned building this

- 404 means the path itself was never registered. 405 means the path exists but the method used isn't allowed on it — these are commonly confused but mean different things.
- Header *names* are case-insensitive (`X-Api-Key` == `x-api-key`); header *values* are not.
- Query parameters are always strings — a URL is just text, so numeric-looking params still need `strconv.Atoi` to convert.
- A redirect (301/302) is a separate response containing only a status code and a `Location` header — the client makes a second request to actually fetch the new content.
