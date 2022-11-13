## Summary

There's a [proposal](https://github.com/golang/go/issues/56345) to add structured logging to the standard library. See 
[here](https://pkg.go.dev/golang.org/x/exp/slog) for implementation. From here on out, I will use the name slog to refer 
to this proposed implementation.

This repo is a proof of concept for kicking the tires on slog library. I wanted to see:

- how easy it was to create handler that emits syslog messages. turns out it was pretty easy
- provide a frontend shim for slog that supports a fluent api for adding attributes. The slog implementation
  avoids fluent api because of the performance implications around memory allocation. I wanted 
  to see what a fluent api would look like and what the performance penalty is. Based on some rough benchmarks,
  it looks like the fluent api adds 5 additional allocations.

This software is alpha quality software, is intentionally tagged at pre 1.0.0 semver, and provides zero compatibility
guarantees. Use at your own risk.


