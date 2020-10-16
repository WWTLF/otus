FROM migrate/migrate
WORKDIR /db
COPY migrate.sh .
WORKDIR /db/migrations
COPY . .
WORKDIR /db
