FROM migrate/migrate

COPY ./migrations /migrations

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["migrate -path /migrations -database postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up"]
