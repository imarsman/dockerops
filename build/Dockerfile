# Base docker image for running NanoVMS

FROM amd64/alpine:latest

# RUN apk upgrade && apk update && apk add qemu-img qemu-system-x86_64 bash bash-completion
RUN apk upgrade && apk update && apk add qemu curl bash bash-completion

COPY ./assets/run.sh /app/run.sh
RUN chmod 755 /app/run.sh

ADD https://storage.googleapis.com/cli/linux/ops /app/ops

RUN chmod 755 /app/run.sh
RUN chmod 755 /app/ops

#WORKDIR /app

ENTRYPOINT ["/app/run.sh"]
#RUN /bin/ls -lah