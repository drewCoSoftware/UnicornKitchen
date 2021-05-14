docker pull postgres
docker pull dpage/pgadmin4

docker stop /dev-postgres
docker rm /dev-postgres

docker stop /dev-pgadmin
docker rm /dev-pgadmin

docker run -d --name dev-postgres -e POSTGRES_PASSWORD=abc123 -v ${HOME}/postgres-data/:/var/lib/postgresql/data -p 5432:5432 postgres
docker run -p 80:80 -e 'PGADMIN_DEFAULT_EMAIL=user@domain.local' -e 'PGADMIN_DEFAULT_PASSWORD=SuperSecret' --name dev-pgadmin -d dpage/pgadmin4