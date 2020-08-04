## development #################################################################

FROM golang:1.14 AS development

RUN apt-get update && apt-get -y install curl software-properties-common && curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get update && apt-get -y install git nodejs unzip

RUN curl -Ls https://github.com/mattgreen/watchexec/releases/download/1.8.6/watchexec-1.8.6-x86_64-unknown-linux-gnu.tar.gz | \
  tar -C /usr/bin --strip-components 1 -xz

# RUN apt-get update && apt-get -y install python python-pip && pip install awscli

# RUN curl -Ls https://github.com/convox/convox/releases/download/3.0.9/convox-linux -o /usr/bin/convox && \
#   chmod +x /usr/bin/convox

# RUN curl -s https://convox.s3.amazonaws.com/release/20200302115619/cli/linux/convox -o /usr/bin/convox2 && \
#   chmod +x /usr/bin/convox2

# RUN curl -L https://releases.hashicorp.com/terraform/0.12.21/terraform_0.12.21_linux_amd64.zip -o terraform.zip && \
#   unzip terraform.zip -d /tmp && mv /tmp/terraform /usr/bin/terraform && rm terraform.zip

ENV MODE=development
ENV VERSION=dev

WORKDIR /usr/src/console/web
COPY web/node_modules ./node_modules
RUN npm rebuild

WORKDIR /usr/src/console
COPY go.mod go.sum ./
COPY vendor vendor
RUN go build -mod=vendor --ldflags="-s -w" $(go list -mod=vendor ./vendor/...)

COPY . .

RUN make build

CMD ["bash"]

## package #####################################################################

FROM golang:1.13 AS package

RUN apt-get update && apt-get -y install curl software-properties-common && curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get update && apt-get -y install git nodejs upx-ucl

WORKDIR /usr/src/console

COPY --from=development /usr/src/console .

RUN make dist package build compress

## production ##################################################################

FROM ubuntu:18.04 AS production

RUN apt-get update && apt-get install -y curl git unzip

# COPY --from=development /usr/bin/convox /usr/bin/
# COPY --from=development /usr/bin/convox2 /usr/bin/
# COPY --from=development /usr/bin/terraform /usr/bin/

ARG VERSION=dev

ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

WORKDIR /

COPY bin/web /bin/
# COPY bin/worker /bin/

# COPY --from=package /go/bin/job /go/bin/
# COPY --from=package /go/bin/rack /go/bin/
# COPY --from=package /go/bin/task /go/bin/
COPY --from=package /go/bin/web /go/bin/
# COPY --from=package /go/bin/worker /go/bin/

RUN groupadd -r console && useradd -r -g console console
RUN mkdir -p /home/console && chown -R console:console /home/console

USER console

CMD ["bash"]
