#! /bin/sh
cd /api
./pgweb --user=$DB_USER --pass=$DB_PASSWORD --host=$DB_HOST --port=$DB_PORT --db=$DB_NAME --bind=0.0.0.0 -d --listen=$PGWEB_PORT --ssl=disable &

./app
