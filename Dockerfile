FROM scratch

COPY .  /server

WORKDIR  /server

CMD ["/server/main"]
