FROM golang:latest AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main main.go

RUN apt-get -y update && apt-get install -y tzdata

ENV PGVER 12
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "ALTER USER postgres WITH PASSWORD 'Qwerty123';" &&\
    createdb -O postgres forum &&\
    /etc/init.d/postgresql stop

RUN echo "listen_addresses='*'\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /app
COPY --from=builder /app .

EXPOSE 5000

ENV PGPASSWORD root
CMD service postgresql start && psql -h localhost -d forum -U postgres -p 5432 -a -q -f ./migrations/init.sql && ./main