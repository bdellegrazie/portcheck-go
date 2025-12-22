# Portcheck

Simple netstat like port check for containers, intended to be used as a dumb readiness check for
TCP or UDP services.

# Usage

```
Usage of portcheck:
  -addr string
    	IP address to check against (default "0.0.0.0")
  -port uint
    	Port to check (0 > n < 65536) (default 8000)
  -protocol string
    	Protocol, one of: tcp, tcp6, udp, udp6 (default "tcp")
  -silent
    	Be silent
  -version
    	Get application version
```

| Exit Code | Meaning               |
| --------- | --------------------- |
| 0         | Port is listening     |
| 1         | Error occurred        |
| 10        | Port is not listening |

# Contributors

@bdellegrazie (author)

<!-- References -->
