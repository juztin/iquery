# minty/iquery
#

FROM minty/db2
MAINTAINER Justin Wilson <iquery@minty.io>


WORKDIR /opt/iquery

COPY iquery .
COPY templates/ templates/
COPY static/ static/

# Disable cgo pointer validations (for 'bitbucket.org/phiggins/db2cli')
#  - https://golang.org/doc/go1.6#cgo
ENV GODEBUG="cgocheck=0"

CMD ["./iquery"]
