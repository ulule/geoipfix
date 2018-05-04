FROM debian:stretch-slim

RUN apt-get update \
	&& apt-get install -y wget

RUN mkdir -p /usr/share/geoip \
	&& wget -O /tmp/GeoLite2-City.tar.gz http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz \
	&& tar xf /tmp/GeoLite2-City.tar.gz -C /usr/share/geoip --strip 1 \
	&& gzip /usr/share/geoip/GeoLite2-City.mmdb \
	&& ls -al /usr/share/geoip/

ADD bin/geoipfix /geoipfix

CMD ["/geoipfix"]
