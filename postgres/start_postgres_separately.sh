docker build -t postgresdb .
docker run -d -P -p 6432:5432 -v ../_postgres --name postgrescontainer postgresdb
