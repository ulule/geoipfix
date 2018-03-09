FROM stretch-slim

ADD bin/geoipfix /geoipfix

CMD ["/geoipfix"]
