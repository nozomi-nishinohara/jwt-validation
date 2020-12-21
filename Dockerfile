FROM scratch as client
ARG ARG_DOCKER_CLIENT_VERSION=19.03.8
ENV DOCKER_CLIENT_VERSION=${ARG_DOCKER_CLIENT_VERSION}
ENV DOCKER_API_VERSION=1.40
ADD https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_CLIENT_VERSION}.tgz .

FROM golang:1.14-alpine
ENV GOPRIVATE="bitbucket.org/dxgogo,github.com/belldata"
ENV TZ=Asia/Tokyo
ARG ARG_DOCKER_CLIENT_VERSION=19.03.8

ENV DOCKER_CLIENT_VERSION=${ARG_DOCKER_CLIENT_VERSION}
ENV DOCKER_API_VERSION=1.40
COPY --from=client docker-${DOCKER_CLIENT_VERSION}.tgz .

RUN apk --no-cache add tzdata gcc libc-dev git make bash \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && apk del tzdata \
    && tar xzvf docker-${DOCKER_CLIENT_VERSION}.tgz \
    && mv docker/* /usr/bin/ \
    && rm -rf docker-${DOCKER_CLIENT_VERSION}.tgz \
    && mkdir /src \
    && rm  -rf /tmp/* /var/cache/apk/* \
    && cd /tmp \
    && git clone https://github.com/awslabs/git-secrets.git \
    && cd git-secrets \
    && make install \
    && git secrets --register-aws --global \
    && cd ../ \
    && rm -rf git-secrets

WORKDIR /src/

CMD [ "sh" ]