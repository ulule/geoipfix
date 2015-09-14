FROM scratch

ADD bin/ipfix /ulule-api

CMD ["/ipfix"]
