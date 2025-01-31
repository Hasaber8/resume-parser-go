# Resume parser

Simple attempt at understanding how career websites parse your resume pdf files and put them automatically into sections.

#### Build locally
```
go build -o parser ./cmd/parser
```

#### Usage of parser:-

Requires java8+ installed locally.
(there was no reliable go - pdf to text - utility)

```
-debug
        Enable debug output
-format=string
        Output format (json or text) (default "json")
-timeout=duration
        Processing timeout (default 30s)
```        
