# go-health

Health checks for the go programming language.

## Usage

```go
checker := health.NewHTTPHealthCheckBuilder("http://localhost:8080/healthy").
    WithMethod("POST").
    WithHeader("Content-Type", "nothing").
    WithBasicAuth("trillian", "tricia123").
    WithExpectedStatusCode(204).
    Build()

watcher := health.NewWatcher()
err := watcher.WaitUntilHealthy(checker)
if err != nil {
  // check does not get healthy
}
```