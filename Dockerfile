FROM golang:1.14.2

WORKDIR /usr/src/app

COPY ./install_prerequisites.sh ./install_prerequisites.sh

RUN chmod +x ./install_prerequisites.sh&& /bin/bash ./install_prerequisites.sh

COPY ./ ./

CMD ["/bin/bash", "./start.sh"]