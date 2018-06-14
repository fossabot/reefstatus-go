FROM scratch

COPY frontend/ReefStatus/dist/.  /server/frontend/ReefStatus/dist
COPY main  /server

WORKDIR  /server

CMD ["/server/main"]
