# ipv4range

Fast IPv4 address range matcher.

## Install

```
go get -u github.com/okzk/ipv4range
```

## Example

```go
matcher, _ := ipv4range.NewMatcher("10.10.0.0/16", "10.20.0.0/16", "10.30.0.0/16")

if matcher.Match("10.10.10.10") {
    // do something	
}
```

## License

MIT
