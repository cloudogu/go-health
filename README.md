# go-health

Health checks for the go programming language.

## Usage

First we have to create a health check:

### http

```go
checker := health.NewHTTPHealthCheckBuilder("http://localhost:8080/healthy").
    WithMethod("POST").
    WithHeader("Content-Type", "text/plain").
    WithBasicAuth("trillian", "tricia123").
    WithExpectedStatusCode(204).
    Build()
```

### tcp

```go
checker := health.NewTCPHealthCheckBuilder(6379).
  WithHostname("redis.hitchhiker.com").
  WithTimeout(1 * time.Second).
  Build()
```

Than we can wait until the check becomes healthy:

```go
watcher := health.NewWatcher()
err := watcher.WaitUntilHealthy(checker)
if err != nil {
  // check does not become healthy and the timeout is reached
}
```
