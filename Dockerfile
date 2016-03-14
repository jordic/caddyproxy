
FROM scratch
MAINTAINER Jordi Collell <j@tmpo.io>

COPY main caddy-gen
COPY Caddyfile Caddyfile

VOLUME = ["/Caddyfile"]
ENTRYPOINT = ["/caddy-gen"]

