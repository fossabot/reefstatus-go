FROM scratch

COPY main  /server
COPY frontend/ReefStatus/dist  /server/frontend/ReefStatus/dist

WORKDIR  /server

CMD ["/server/main"]
