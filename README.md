- CREATE USER: Creates a database user.
- shoe_shop_user: Our application's database username.
- PASSWORD: Sets the user's password.

go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get github.com/go-playground/validator/v10
go get github.com/stretchr/testify


pg_lsclusters
sudo systemctl status postgresql

sudo systemctl start postgresql

psql -h localhost -p 5433 -U shoe_shop_user -d shoe_shop -W