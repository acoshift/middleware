# middleware

net/http middleware collection

## Chaining

```go
middleware.Chain(
    middleware.HSTSPreload(),
    middleware.NonWWWRedirect(),
    middleware.Compress(middleware.DeflateCompressor),
    middleware.Compress(middleware.GzipCompressor),
)
```

## HSTS

```go
middleware.HSTS(HSTSConfig{
    MaxAge:            3600 * time.Second,
    IncludeSubDomains: true,
    Preload:           false,
})
```

```go
middleware.HSTS(middleware.DefaultHSTS)
```

```go
middleware.HSTS(middleware.PreloadHSTS)
```

```go
middleware.HSTSPreload()
```

## Compressor

```go
middleware.Compress(middleware.CompressConfig{
    New: func() Compressor {
        g, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
        if err != nil {
            panic(err)
        }
        return g
    },
    Encoding:  "gzip",
    Vary:      true,
    Types:     "text/plain text/html",
    MinLength: 1000,
})
```

```go
middleware.Compress(middleware.GzipCompressor)
```

## Redirector

```go
middleware.NonWWWRedirect()
```

```go
middleware.WWWRedirect()
```