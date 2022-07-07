module github.com/lxt1045/errors

go 1.18

require (
	github.com/petermattis/goid v0.0.0-20220526132513-07eaf5d0b9f4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

// replace github.com/petermattis/goid => github.com/lxt1045/goid v0.0.0-20220606075709-6d67c0e3a5ea
replace github.com/petermattis/goid => /Users/bytedance/go/src/github.com/lxt1045/goid
