# layer-four-forwarder

Created by me when I get f**ked by macOS's firewall. Can't accept inbound connection with dlv / air / etc. 

To bypass the limitation, a layer 4 forwarder is needed.

Currently only tcp

# Install 
```bash
go install github.com/bookqaq/layer-four-forwarder@latest
```

# Usage
```bash
layer-four-forwarder -src 0.0.0.0:8080 -dst 127.0.0.1:8081
```
