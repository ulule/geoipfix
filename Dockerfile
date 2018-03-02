FROM stretch-slim

ADD bin/ipfix /ipfix

CMD ["/ipfix"]
