module github.com/araquach/apiFinance23

go 1.18

require (
	github.com/araquach/apiHelpers v0.0.1
	github.com/araquach/dbService v0.0.1
	github.com/rs/cors v1.9.0
)

require (
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
	gorm.io/gorm v1.25.1 // indirect
)

replace (
	github.com/araquach/apiHelpers => ../apiHelpers
	github.com/araquach/dbService => ../dbService
)
