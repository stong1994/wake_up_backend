module wake_up_backend

go 1.16

require (
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/cors v1.2.0
	github.com/stong1994/kit_golang v0.0.0-20211003064513-7826fbb57f46
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
)

replace github.com/stong1994/kit_golang v0.0.0-20211003064513-7826fbb57f46 => ../kit_golang
