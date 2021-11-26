module wake_up_backend

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/cors v1.2.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/stong1994/kit_golang v0.0.0-20211003064513-7826fbb57f46
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.8.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
)

replace github.com/stong1994/kit_golang v0.0.0-20211003064513-7826fbb57f46 => ../kit_golang
