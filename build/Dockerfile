# Base docker image for running NanoVMS

FROM amd64/alpine:latest

# RUN apk upgrade && apk update && apk add qemu-img qemu-system-x86_64 bash bash-completion
RUN apk upgrade && apk update && apk add curl bash bash-completion

# COPY ./assets/ops /app/ops
ADD https://storage.googleapis.com/cli/linux/ops /app/ops

RUN chmod 755 /app/ops

#RUN /bin/ls -lah