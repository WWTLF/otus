export POSTGRESQL_URL="postgres://postgres:$PGPASSWORD@$POSTGRESQL_HOSTNAME:$POSTGRESQL_PORT_NUMBER/userlist?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path migrations up