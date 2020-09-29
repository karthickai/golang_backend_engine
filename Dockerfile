FROM ubuntu:18.04

MAINTAINER karthick.pannerselvam@humetis.in

RUN apt-get update
RUN apt-get update --fix-missing
RUN apt-get install -y wkhtmltopdf
COPY bin/c2e_backend /usr/local/bin

ENV DEBUG true
EXPOSE 3000

CMD ["/usr/local/bin/c2e_backend"]

