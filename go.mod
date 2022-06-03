module github.com/lxt1045/errors

go 1.16

require (
	github.com/bytedance/gopkg v0.0.0-20220531084716-665b4f21126f
	github.com/petermattis/goid v0.0.0-20220526132513-07eaf5d0b9f4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.1
)

// replace github.com/petermattis/goid => github.com/lxt1045/goid v0.0.0-20220606075709-6d67c0e3a5ea
replace github.com/petermattis/goid => /Users/bytedance/go/src/github.com/lxt1045/goid
