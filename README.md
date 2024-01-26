# Gofresher

Scheduled and just-in-time state refresher, Gofresher allows to manage a central state that is refreshed both periodically and opportunistically. In all cases, refresh actions are deduplicated.

```go
// Create a state holder where a value is does not need refreshing for 50 minutes
tokenHolder := gofresher.NewGofresher[string](50*time.Minute, func(state *int) (*string, error) {
  token, err := fetchHourLongToken()

  return token, err
})

// Start a background task to check if the token needs refreshing every minute
tokenHolder.Start(time.Minute)

// Retrieve the current state without refreshing
token, err, refreshed := tokenHolder.State(false)

// Retrieve the current state, refreshing if needed
token, err, refreshed := tokenHolder.State(true)

// Force refreshing the current state
token, err := tokenHolder.ForceRefresh()
```
