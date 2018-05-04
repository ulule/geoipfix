#!/bin/sh
while true; do
  echo 'Updating GeoIP database...'
  mkdir -p /usr/share/geoip
  wget -O /tmp/GeoLite2-City.tar.gz http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz
  tar xf /tmp/GeoLite2-City.tar.gz -C /usr/share/geoip --strip 1
  mv /tmp/GeoLite2-City.tar.gz /usr/share/geoip
  ls -al /usr/share/geoip/

  echo 'Sleeping for a day...';
  sleep 86400;
done
