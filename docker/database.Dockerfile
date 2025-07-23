FROM postgres:latest

ENV POSTGRES_DB=exercise
ENV POSTGRES_USER=super_secure_username
ENV POSTGRES_PASSWORD=super_secure_password

COPY database/schema.sql /docker-entrypoint-initdb.d/schema.sql