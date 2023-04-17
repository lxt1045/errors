# errors 包

## 原理
一句话：通过汇编，从调用栈中获取 pc 和 pc 列表。\
性能提升的具体细节，这两篇技术文章中有详细说明：
1. [关于 golang 错误处理的一些优化想法](https://juejin.cn/post/7121929424148103198)
2. [golang文件行号探索](https://juejin.cn/post/7124334239692095501)

此库下有两个功能模块：
1. errors：功能和 [pkg/errors](https://github.com/pkg/errors) 类似，性能比后者高一个数量级以上。
2. errors/logrus：功能和 [sirupsen/logrus](https://github.com/sirupsen/logrus) 一样，目标是 100% 兼容。 利用了 errors 获取行号的接口，性能比后者高 35% 以上，会持续优化。如有不兼容的地方，欢迎吐槽。

## Go2 Error Handling
参考 Go2 的 check与hanle关键字，实现了类是的错误处理逻辑：
```go
func TestTagTry(t *testing.T) {
	defer func() {
		fmt.Printf("1 -> ")
	}()

	tag, err := NewTag() // 当 tag.Try(err) 时，跳转此处并返回 err1
	fmt.Printf("2 -> ")
	if err != nil {
		fmt.Printf("3 -> ")
		// 参考： https://github.com/golang/proposal/blob/master/design/34481-opencoded-defers.md
		// defer 在 loop 中，导致编译器对 defer 内联优化策略的改变：不再使用 "open-coded defers" 策略，否则 "4 -> " 将不输出。
		for false {
			defer func() {}()
		}
		return
	}

	defer func() {
		fmt.Printf("4 -> ")
	}()

	fmt.Printf("5 -> ")
	tag.Try(errors.New("err"))

	fmt.Printf("6 -> ")
	return
}
```
以上代码将输出：
```log
2 -> 5 -> 2 -> 3 -> 4 -> 1 ->
```
当然，如果使用 defer + panic 实现相关功能也可以。
不过如果忘了 defer recover 有可能会早成程序退出，而且很多公司都禁用这种方式。

## 性能测试

1. errors 和 [pkg/errors](https://github.com/pkg/errors) 比较

由测试结果可知，在性能上，这个 errors 库比 Go 语言官方 errors 库（不带调用栈）的性能还要好。
和 [pkg/errors](https://github.com/pkg/errors) 已经拉开了一个数量级以上的差距。

[BenchmarkNewAndFormatting](https://github.com/lxt1045/errors/blob/main/formatter_test.go#L363)
```go
func BenchmarkNewAndFormatting(b *testing.B) {
    depths := []int{1, 10} //嵌套深度
    std, pkg, lxt := "std", "pkg", "lxt"

    stdText := func(err error) []byte {
        buf := bytes.NewBuffer(make([]byte, 0, 1024))
        for ; err != nil; err = errors.Unwrap(err) {
            buf.WriteString(err.Error())
        }
        return buf.Bytes()
    }

    runs := []struct {
        t    string          //函数名字
        name string          //函数名字
        f    func(depth int) //调用方法
    }{
        {std, "text", func(depth int) {
            err := errors.New(errMsg)
            for j := 0; j < depth; j++ {
                err = fmt.Errorf("%w; %s", err, errTrace)
            }
            stdText(err)
        }},
        {lxt, "text", func(depth int) {
            var err error = NewCode(0, errCode, errMsg)
            for j := 0; j < depth; j++ {
                err = Wrap(err, errTrace)
            }
            MarshalText(err)
        }},
        {lxt, "json", func(depth int) {
            var err error = NewCode(0, errCode, errMsg)
            for j := 0; j < depth; j++ {
                err = Wrap(err, errTrace)
            }
            MarshalJSON(err)
        }},
        {pkg, "text.%+v", func(depth int) {
            err := pkgerrs.New(errMsg)
            for j := 0; j < depth; j++ {
                err = pkgerrs.Wrap(err, errTrace)
            }
            _ = fmt.Sprintf("%+v", err)
        }},
        {pkg, "text.%v", func(depth int) {
            err := pkgerrs.New(errMsg)
            for j := 0; j < depth; j++ {
                err = pkgerrs.Wrap(err, errTrace)
            }
            _ = fmt.Sprintf("%v", err)
        }},
    }

    for _, run := range runs {
        for _, depth := range depths {
            name := fmt.Sprintf("%s.%s-%d", run.t, run.name, depth)
            b.Run(name, func(b *testing.B) {
                b.ReportAllocs()
                b.ResetTimer()
                for i := 0; i < b.N; i++ {
                    run.f(depth)
                }
                b.StopTimer()
            })
        }
    }
}


```
测试结果
```sh
BenchmarkNewAndFormatting/std.text-1
BenchmarkNewAndFormatting/std.text-1-12      1963789   658.9 ns/op  1088 B/op   4 allocs/op
BenchmarkNewAndFormatting/std.text-10
BenchmarkNewAndFormatting/std.text-10-12     452484    2575 ns/op  1913 B/op  22 allocs/op
BenchmarkNewAndFormatting/lxt.text-1
BenchmarkNewAndFormatting/lxt.text-1-12       2825418   429.8 ns/op   421 B/op   4 allocs/op
BenchmarkNewAndFormatting/lxt.text-10
BenchmarkNewAndFormatting/lxt.text-10-12      831126    1529 ns/op  1814 B/op  13 allocs/op
BenchmarkNewAndFormatting/lxt.json-1
BenchmarkNewAndFormatting/lxt.json-1-12       2325892   650.5 ns/op   485 B/op   4 allocs/op
BenchmarkNewAndFormatting/lxt.json-10
BenchmarkNewAndFormatting/lxt.json-10-12       570873   1912 ns/op  2071 B/op  13 allocs/op
BenchmarkNewAndFormatting/pkg.text.%+v-1
BenchmarkNewAndFormatting/pkg.text.%+v-1-12    110577   9163 ns/op  1827 B/op  28 allocs/op
BenchmarkNewAndFormatting/pkg.text.%+v-10
BenchmarkNewAndFormatting/pkg.text.%+v-10-12   24076    52849 ns/op  9980 B/op   154 allocs/op
BenchmarkNewAndFormatting/pkg.text.%v-1
BenchmarkNewAndFormatting/pkg.text.%v-1-12     534991   2099 ns/op   672 B/op   9 allocs/op
BenchmarkNewAndFormatting/pkg.text.%v-10
BenchmarkNewAndFormatting/pkg.text.%v-10-12    95394    11317 ns/op  4315 B/op  54 allocs/op
```


2. errors/logrus 和 [sirupsen/logrus](https://github.com/sirupsen/logrus)

由结果可知，性能提升了 35% 以上。

[BenchmarkLog](https://github.com/lxt1045/errors/blob/main/logrus/sample_test.go#L96)
```go

func BenchmarkLog(b *testing.B) {
    bs := make([]byte, 1<<20)
    w := bytes.NewBuffer(bs)
    logrus.SetReportCaller(true)
    logrus.SetOutput(w)
    // logrus.SetLevel(logrus.DebugLevel)
    logrus.SetFormatter(&logrus.JSONFormatter{})
    // h := &Hook{AppName: "awesome-web"}
    // logrus.AddHook(h)
    logrus.Info("info msg")
    // b.Log(w.String())

    ctx := context.TODO()

    b.Run("logrus+caller", func(b *testing.B) {
        logrus.SetReportCaller(true)
        for i := 0; i < b.N; i++ {
            logrus.WithContext(ctx).Info("info msg")
            if w.Len() > len(bs)-64 {
                w.Reset()
            }
        }
    })

    b.Run("logrus", func(b *testing.B) {
        logrus.SetReportCaller(false)
        for i := 0; i < b.N; i++ {
            WithContext(ctx).Info("info msg")
            if w.Len() > len(bs)-64 {
                w.Reset()
            }
        }
    })

    b.Run("logrus+lxt caller", func(b *testing.B) {
        logrus.SetReportCaller(false)
        for i := 0; i < b.N; i++ {
            logrus.WithContext(ctx).Info("info msg")
            if w.Len() > len(bs)-64 {
                w.Reset()
            }
        }
    })
}
```
测试结果如下：
```sh
BenchmarkLog/logrus+caller
BenchmarkLog/logrus+caller-12      169206    6166 ns/op    2172 B/op    36 allocs/op
BenchmarkLog/logrus
BenchmarkLog/logrus-12             265323    3942 ns/op    2317 B/op    35 allocs/op
BenchmarkLog/logrus+lxt_caller
BenchmarkLog/logrus+lxt_caller-12  422571    2413 ns/op    1354 B/op    25 allocs/op
```

## 设计思路

