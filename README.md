# Web Proxy
[![Docker Automated buil](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/sommereng/web-proxy/)

This is a simple web proxy e.g. for Docker setups. The purpose is to run this proxy on port `80` and additional web projects on several other ports e.g. `50001`, `50002`, etc. Through the web proxy, all web projects are available by port `80`. The differentiation is possible by different domains for these web projects.

# Configuration
Append an arbitrary number of hosts and destinations to the program name, e.g. for Linux and macOS `./WebProxy myhost1=>http://www.another-domain.com myhost2=>http://www.test.com`. In this example, any request where the host contains `myhost1` will be fed by `http://www.another-domain.com`, etc. Instead of a domain name, the destination could be an IP address instead: `./WebProxy myhost1=>http://www.another-domain.com myhost2=>http://127.0.0.1:50001`. The matching of the host part allows partial matching. Thus, `myhost1` will match e.g. `http://myhost1.my-domain.com` as well as `http://www.myhost1-cool.com`, etc. pp.