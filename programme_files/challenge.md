# Go HTTP Server — Exercises & Cheat Sheet

## Exercises

### Exercise 1: Basic Ping-Pong Server

**Goal:** Build a minimal web server that listens on port 8080 and responds with `pong` when a user visits the `/ping` route.

**Tasks:**
- Create a route handler for `/ping` using `http.HandleFunc`.
- Use `w.Write()` or `fmt.Fprint()` to send a plain text response.
- Start the server on port `:8080` using `http.ListenAndServe`.

---

### Exercise 2: Query Parameters & Path Validation

**Goal:** Create a `/hello` endpoint that reads a `name` query parameter (e.g., `/hello?name=Alice`) and responds with `Hello, Alice!`. If the parameter is missing, default to `Hello, Guest!`.

**Tasks:**
- Extract query parameters using `r.URL.Query().Get("name")`.
- Reject any HTTP method that is not GET by returning an `http.StatusMethodNotAllowed` status code.

---

### Exercise 3: Text Counter (URL Variables & Methods)

**Goal:** Build a server with a `/count` route. If a user sends a GET request, return the text `Send a POST request with text to count words`. If they send a POST request, read the text body and return the number of characters.

**Key Tasks:**
- Differentiate between GET and POST methods using `r.Method`.
- Read the entire request body using `io.ReadAll(r.Body)`.
- Return the character length as a string.

---

### Exercise 4: Basic Math API (Multiple Query Parameters)

**Goal:** Create a `/calculate` route that accepts three query parameters: `op` (operation), `a`, and `b`. For example, `/calculate?op=add&a=10&b=5` should respond with `Result: 15`.

**Key Tasks:**
- Parse string query variables using `r.URL.Query().Get()`.
- Convert string inputs to integers using `strconv.Atoi()`.
- Support `add`, `subtract`, and `multiply`. Return an HTTP 400 Bad Request if the operation is unknown or parsing fails.

---

### Exercise 5: User-Agent Echo

**Goal:** Create an `/agent` route that reads the client's `User-Agent` header and echoes it back in the response.

---

### Exercise 6: Secure Dashboard (Simple Authorization Headers)

**Goal:** Create a `/dashboard` route that acts as a protected page. If the client does not provide a specific API key in their headers, block them.

**Key Tasks:**
- Read custom headers using `r.Header.Get("X-API-Key")`.
- Match it against a hardcoded value (e.g., `secret123`).
- Use `http.StatusUnauthorized` (401) to reject bad keys.

---

### Exercise 7: Simple Redirector (Status Codes)

**Goal:** Create a `/legacy` route. Whenever a user hits this endpoint, permanently redirect them to a new route `/v2` with a friendly `Welcome to version 2` message.

**Key Tasks:**
- Redirect traffic using the `http.Redirect` helper function.
- Use the proper status code for a permanent move (`http.StatusMovedPermanently`).

---

## Go HTTP Server Cheat Sheet

### 1. Routing & Server Management (`net/http`)

| Function | Description |
|---|---|
| `http.HandleFunc(pattern string, handler func(ResponseWriter, *Request))` | Registers a handler function for a specific URL path pattern. |
| `http.ListenAndServe(addr string, handler Handler) error` | Starts an HTTP server on the specified address (e.g., `:8080`). Passing `nil` uses the default system router (`DefaultServeMux`). |
| `http.Error(w ResponseWriter, error string, code int)` | A helper that sends a specific error message string and numeric HTTP status code back to the client, automatically ending the request lifecycle safely. |
| `http.Redirect(w ResponseWriter, r *Request, url string, code int)` | Sends an HTTP redirect status code (like 301 or 302) forcing the client's browser to jump to a new target URL path. |

### 2. Reading Inputs & Formatting (`io`, `strconv`, `fmt`)

| Function | Description |
|---|---|
| `io.ReadAll(r io.Reader) ([]byte, error)` | Reads all remaining data from an input stream (like `r.Body`) until it hits the end of the file/stream (EOF). Returns a byte slice. |
| `strconv.Atoi(s string) (int, error)` | Stands for "ASCII to Integer". Converts a text string into a native numeric `int`. Returns an error if the text contains non-numeric characters. |
| `fmt.Fprintf(w io.Writer, format string, a ...any)` | Formats data according to a template string and writes the output directly into an open network socket stream or file (`w`). |

### Request Context Breakdown (`*http.Request`)

When a client hits your server, all incoming information is bundled inside the pointer to the `http.Request` struct (usually named `r`):

| Struct Field / Method | Purpose / Explanation | Example Usage |
|---|---|---|
| `r.Method` | A string representing the incoming HTTP type. Always use standard library constants for comparisons. | `if r.Method == http.MethodPost` |
| `r.Body` | An open stream containing data uploaded via POST/PUT requests. Always close it after reading to prevent memory leaks. | `defer r.Body.Close()` |
| `r.URL.Query()` | Parses the raw URL query string (everything after the `?`) and returns a map-like structure (`Values`). | `queryMap := r.URL.Query()` |
| `r.URL.Query().Get(key)` | Fetches the value of a specific query parameter. Returns an empty string `""` if the key does not exist. | `name := r.URL.Query().Get("name")` |
| `r.Header.Get(key)` | Fetches the value of a specific HTTP request header (case-insensitive). Returns `""` if missing. | `token := r.Header.Get("X-API-Key")` |

### Response Management (`http.ResponseWriter`)

The `http.ResponseWriter` interface (usually named `w`) is your pipeline to send data back to the client. Keep this execution order in mind:

1. **Modify Headers First** — Call `w.Header().Set("Key", "Value")` before anything else if changing content types or metadata.
2. **Write Status Code Second** — Call `w.WriteHeader(http.StatusCreated)` if returning something other than a standard 200 OK.
3. **Write Body Content Last** — Call `w.Write([]byte)` or `fmt.Fprintf(w)` to push actual visual content to the user. Writing to the body locks your headers and status code automatically.

---

## Terminal Test Script (`test_endpoints.sh`)

This shell script acts as an automated collection of `curl` requests. Run it directly from your terminal to verify your 7 endpoints meet the requirements.

Create a file named `test_endpoints.sh`, paste the code below, and make it executable:

```bash
chmod +x test_endpoints.sh
./test_endpoints.sh
```

```bash
#!/bin/bash

# Configuration
SERVER_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0;0m' # No Color

echo -e "${BLUE}=== Starting Go HTTP Exercise Verification Script ===${NC}\n"

# Exercise 1: Basic Ping-Pong Server
echo -e "${BLUE}[Exercise 1: /ping]${NC}"
RESPONSE=$(curl -s "$SERVER_URL/ping")
if [ "$RESPONSE" == "pong" ]; then
    echo -e "${GREEN}✔ PASS: Got 'pong'${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 'pong', got '$RESPONSE'${NC}"
fi
echo ""

# Exercise 2: Query Parameters & Path Validation
echo -e "${BLUE}[Exercise 2: /hello]${NC}"
# Test with name
RESP_NAME=$(curl -s "$SERVER_URL/hello?name=Alice")
if [[ "$RESP_NAME" == *"Hello, Alice!"* ]]; then
    echo -e "${GREEN}✔ PASS: Query param parsed successfully ('Hello, Alice!')${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 'Hello, Alice!', got '$RESP_NAME'${NC}"
fi

# Test default guest
RESP_GUEST=$(curl -s "$SERVER_URL/hello")
if [[ "$RESP_GUEST" == *"Hello, Guest!"* ]]; then
    echo -e "${GREEN}✔ PASS: Default fallback working ('Hello, Guest!')${NC}"
else
    echo -e "${RED}✘ FAIL: Expected fallback, got '$RESP_GUEST'${NC}"
fi

# Test invalid method validation
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$SERVER_URL/hello")
if [ "$STATUS_CODE" == "405" ]; then
    echo -e "${GREEN}✔ PASS: Blocked POST request with Status 405 Method Not Allowed${NC}"
else
    echo -e "${RED}✘ FAIL: POST request gave status $STATUS_CODE instead of 405${NC}"
fi
echo ""

# Exercise 3: Text Counter
echo -e "${BLUE}[Exercise 3: /count]${NC}"
# Test GET
RESP_GET=$(curl -s "$SERVER_URL/count")
if [[ "$RESP_GET" == *"Send a POST request"* ]]; then
    echo -e "${GREEN}✔ PASS: GET request displays instruction text${NC}"
else
    echo -e "${RED}✘ FAIL: Unexpected GET response '$RESP_GET'${NC}"
fi

# Test POST
RESP_POST=$(curl -s -X POST -d "Golang" "$SERVER_URL/count")
if [[ "$RESP_POST" == *"6"* ]]; then
    echo -e "${GREEN}✔ PASS: POST request calculated length correctly ('Golang' = 6)${NC}"
else
    echo -e "${RED}✘ FAIL: Expected length 6, got '$RESP_POST'${NC}"
fi
echo ""

# Exercise 4: Basic Math API
echo -e "${BLUE}[Exercise 4: /calculate]${NC}"
# Test valid math
RESP_MATH=$(curl -s "$SERVER_URL/calculate?op=add&a=12&b=8")
if [[ "$RESP_MATH" == *"20"* ]]; then
    echo -e "${GREEN}✔ PASS: 12 + 8 = 20 handled successfully${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 20, got '$RESP_MATH'${NC}"
fi

# Test invalid input validation
STATUS_MATH=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/calculate?op=multiply&a=abc&b=5")
if [ "$STATUS_MATH" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Rejected non-integer values with Status 400 Bad Request${NC}"
else
    echo -e "${RED}✘ FAIL: Expected status 400 for bad parameters, got $STATUS_MATH${NC}"
fi
echo ""

# Exercise 5: User-Agent Echo
echo -e "${BLUE}[Exercise 5: /agent]${NC}"
RESP_AGENT=$(curl -s -H "User-Agent: CustomTester/1.0" "$SERVER_URL/agent")
if [[ "$RESP_AGENT" == *"CustomTester/1.0"* ]]; then
    echo -e "${GREEN}✔ PASS: Header value extracted and echoed back${NC}"
else
    echo -e "${RED}✘ FAIL: Expected agent info to be visible, got '$RESP_AGENT'${NC}"
fi
echo ""

# Exercise 6: Secure Dashboard
echo -e "${BLUE}[Exercise 6: /dashboard]${NC}"
# Test unauthorized access
STATUS_DASH_BAD=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/dashboard")
if [ "$STATUS_DASH_BAD" == "401" ]; then
    echo -e "${GREEN}✔ PASS: Missing API key blocked with Status 401 Unauthorized${NC}"
else
    echo -e "${RED}✘ FAIL: Expected status 401 for unauthorized traffic, got $STATUS_DASH_BAD${NC}"
fi

# Test authorized access
RESP_DASH_GOOD=$(curl -s -H "X-API-Key: secret123" "$SERVER_URL/dashboard")
if [[ "$RESP_DASH_GOOD" == *"Welcome"* ]]; then
    echo -e "${GREEN}✔ PASS: Access granted with correct token header${NC}"
else
    echo -e "${RED}✘ FAIL: Correct token rejected. Response: '$RESP_DASH_GOOD'${NC}"
fi
echo ""

# Exercise 7: Simple Redirector
echo -e "${BLUE}[Exercise 7: /legacy -> /v2]${NC}"
# Test redirect status
REDIRECT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/legacy")
if [ "$REDIRECT_STATUS" == "301" ]; then
    echo -e "${GREEN}✔ PASS: Route /legacy issues a 301 Permanent Redirect${NC}"
else
    echo -e "${RED}✘ FAIL: Expected redirect status 301, got $REDIRECT_STATUS${NC}"
fi

# Test following the location redirect
RESP_REDIRECT=$(curl -s -L "$SERVER_URL/legacy")
if [[ "$RESP_REDIRECT" == *"version 2"* ]]; then
    echo -e "${GREEN}✔ PASS: Followed redirect pipeline to /v2 successfully${NC}"
else
    echo -e "${RED}✘ FAIL: Target redirection payload path failed. Got: '$RESP_REDIRECT'${NC}"
fi

echo -e "\n${BLUE}=== Testing Complete ===${NC}"
```
