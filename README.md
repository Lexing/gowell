#### gowell 

Gowell provides basic utilities for go servers. Specifically, it provides 
following HTTP handler automatically: 

1. `/flagz` handler, which outputs all flags defined on the server and its value.
2. `/healthz` handler, which outputs "ok" when the server is ready. (TODO: health channel)
3. `/metrics` handler, which displays the variables/status reported in go server, with [prometheus](https://prometheus.io/) monitoring solution.

Example:

In your server, add (assuming port 8080)
```
InitializeHTTPService(":8080")
```

And when initialization is done, e.g. loading. Call
```
NoteHealthy()
```