FROM scrach

COPY .  /server

WORKDIR  /server

CMD ["/server/main"]
