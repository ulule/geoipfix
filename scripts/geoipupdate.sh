#!/bin/sh
while true; do
  echo 'Updating GeoIP database...'
  mkdir -p /usr/share/geoip
  wget -O /tmp/country.tar.gz http://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz
  tar xf /tmp/country.tar.gz -C /usr/share/geoip --strip 1
  ls -al /usr/share/geoip/

  echo 'Sleeping for a day...';
  sleep 86400;
done
