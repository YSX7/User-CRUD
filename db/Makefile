DATABASE_URL=postgres://postgres:admin@localhost:5432/vuegolang?sslmode=disable

up:
	migrate -database "$(DATABASE_URL)" -path migrations up

down:
	migrate -database "$(DATABASE_URL)" -path migrations down

create:
	migrate create -ext sql -dir migrations -seq $(NAME)