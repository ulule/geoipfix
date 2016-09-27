FROM scratch

ADD bin/ipfix /ipfix

CMD ["/ipfix"]
