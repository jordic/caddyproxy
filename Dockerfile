
FROM scratch
MAINTAINER Jordi Collell <j@tmpo.io>

ADD main caddygen
ADD Caddyfile Caddyfile

VOLUME = ["/Caddyfile"]
CMD = ["/caddygen"]
