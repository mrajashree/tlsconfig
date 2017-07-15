For running in an ubuntu container, you can add the ca-certificates to `/etc/ssl/certs/ca-certificates.crt`
1. Replace [line](https://github.com/mrajashree/tlsconfig/blob/master/main.go#L36) `var certFile = "ca-certificates.crt"` with `var certFile = "/etc/ssl/certs/ca-certificates.crt"`
2. Within the container, `cp ca-certificates.crt /usr/local/share/ca-certificates/` and `update-ca-certificate`
