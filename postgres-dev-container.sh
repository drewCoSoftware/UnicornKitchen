docker pull postgres
docker run -d --name dev-postgres -e POSTGRES_PASSWORD=abc123 -v ${HOME}/postgres-data/:/var/lib/postgresql/data -p 5432:5432 postgres


docker pull dpage/pgadmin4
docker run -p 80:80 -e 'PGADMIN_DEFAULT_EMAIL=user@domain.local' -e 'PGADMIN_DEFAULT_PASSWORD=SuperSecret' --name dev-pgadmin -d dpage/pgadmin4