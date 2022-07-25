module github.com/warent/calzone/service

go 1.18

require github.com/google/go-github/v45 v45.2.0

require (
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/warent/calzone/service/structures/args v0.0.0 => ./structures/args

replace github.com/warent/calzone/service/structures/responses v0.0.0 => ./structures/responses
