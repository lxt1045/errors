# errors 包

## 原理
一句话：通过汇编，从调用栈中获取 pc 和 pc 列表。\
性能提升的具体细节，这两篇技术文章中有详细说明：
1. [关于 golang 错误处理的一些优化想法](https://juejin.cn/post/7121929424148103198)
2. [golang文件行号探索](https://juejin.cn/post/7124334239692095501)

此库下有两个功能模块：
1. errors：功能和 [pkg/errors](https://github.com/pkg/errors) 类似，性能比后者高一个数量级以上。
2. errors/logrus 和 errors/zap：分别包装了 [sirupsen/logrus](https://github.com/sirupsen/logrus) 和 [go.uber.org/zap](https://github.com/uber-go/zap) ，目标是 100% 兼容后者。 利用了 errors 获取行号的接口，能减少 1300ns ~ 2500ns 的时间损耗，会持续优化。如有不兼容的地方，欢迎吐槽。

## 模拟 Go2 错误处理方式
参考 Go2 的 check与hanle关键字，实现了类是的错误处理逻辑：
```go
func TestHandlerCheck(t *testing.T) {
	defer func() {
		fmt.Printf("1 -> ")
	}()

	handler, err := NewHandler() // 当 tag.Try(err) 时，跳转此处并返回 err1
	fmt.Printf("2 -> ")
	if err != nil {
		fmt.Printf("3 -> ")
		return
	}

	fmt.Printf("5 -> ")
	handler.Check(errors.New("err"))

	fmt.Printf("6 -> ")
	return
}
```
以上代码将输出：
```log
2 -> 5 -> 2 -> 3 -> 1 ->
```
当然，如果使用 defer + panic 实现相关功能也可以。
不过如果忘了 defer recover 有可能会早成程序退出，而且很多公司都禁用这种方式。

## 性能基准测试

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


2. errors/logrus 和 errors/zap 性能提升

errors/logrus 和 errors/zap 分别替换了 [sirupsen/logrus](https://github.com/sirupsen/logrus) 和 [go.uber.org/zap](https://github.com/uber-go/zap) 的代码行号获取逻辑。

由结果可知，能减少 1300ns ~ 2500ns 的损耗，而且是兼容性升级，非常值得尝试。

[BenchmarkLog](https://github.com/lxt1045/errors/blob/main/zap/zap_test.go#L85)
```go
func BenchmarkLog(b *testing.B) {
	b.Run("logrus", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		// logrus.SetReportCaller(true)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.WithFields(logrus.Fields{
				"string": "some string format log information",
				"int":    3,
			}).Info("some log messages")
		}
	})
	b.Run("logrus+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		logger.SetReportCaller(true)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.WithFields(logrus.Fields{
				"string": "some string format log information",
				"int":    3,
			}).Info("some log messages")
		}
	})
	b.Run("logrus+lxt caller", func(b *testing.B) {
		// logrus.SetReportCaller(false)
		b.StopTimer()
		b.ReportAllocs()
		logger := lxtlog.New()
		logger.SetOutput(io.Discard)
		// logrus.SetReportCaller(true)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.WithFields(lxtlog.Fields{
				"string": "some string format log information",
				"int":    3,
			}).Info("some log messages")
		}
	})

	b.Run("zap", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("some log messages",
				zap.String("string", `some string format log information`),
				zap.Int("int", 3),
			)
		}
	})
	b.Run("zap+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core, zap.WithCaller(true))
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("some log messages",
				zap.String("string", `some string format log information`),
				zap.Int("int", 3),
			)
		}
	})
	b.Run("zap+lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := New(core, zap.WithCaller(false))
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("some log messages",
				zap.String("string", `some string format log information`),
				zap.Int("int", 3),
			)
		}
	})

	b.Run("zap-sugar", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core)
		sugar := logger.Sugar()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sugar.Info("some log messages",
				"string", `some string format log information`,
				"int", 3,
				"backoff", time.Second,
			)
		}
	})
	b.Run("zap-sugar+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core, zap.WithCaller(true))
		sugar := logger.Sugar()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sugar.Info("some log messages",
				"string", `some string format log information`,
				"int", 3,
				"backoff", time.Second,
			)
		}
	})

	b.Run("zap-sugar+lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := New(core, zap.WithCaller(false))
		sugar := logger.Sugar()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sugar.Info("some log messages",
				"string", `some string format log information`,
				"int", 3,
				"backoff", time.Second,
			)
		}
	})

	b.Run("lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c := CallerFrame(errors.GetPC())
			io.Discard.Write([]byte(zap.String("caller", c.File).String))
		}
	})
}
```
测试结果如下：
```sh
BenchmarkLog/logrus
BenchmarkLog/logrus-12         	  382123	      3144 ns/op	    1361 B/op	      23 allocs/op
BenchmarkLog/logrus+caller
BenchmarkLog/logrus+caller-12  	  173437	      6745 ns/op	    2355 B/op	      34 allocs/op
BenchmarkLog/logrus+lxt_caller
BenchmarkLog/logrus+lxt_caller-12         	  239078	      4836 ns/op	    2082 B/op	      31 allocs/op
BenchmarkLog/zap
BenchmarkLog/zap-12                       	 1457443	       812.7 ns/op	     152 B/op	       3 allocs/op
BenchmarkLog/zap+caller
BenchmarkLog/zap+caller-12                	  461288	      2391 ns/op	     401 B/op	       6 allocs/op
BenchmarkLog/zap+lxt_caller
BenchmarkLog/zap+lxt_caller-12            	 1000000	      1053 ns/op	     409 B/op	       4 allocs/op
BenchmarkLog/zap-sugar
BenchmarkLog/zap-sugar-12                 	 1411080	       848.2 ns/op	     112 B/op	       4 allocs/op
BenchmarkLog/zap-sugar+caller
BenchmarkLog/zap-sugar+caller-12          	  388030	      3542 ns/op	     361 B/op	       7 allocs/op
BenchmarkLog/zap-sugar+lxt_caller
BenchmarkLog/zap-sugar+lxt_caller-12      	 1171387	      1015 ns/op	     176 B/op	       5 allocs/op
BenchmarkLog/lxt_caller
BenchmarkLog/lxt_caller-12                	35303271	        32.32 ns/op	      24 B/op	       1 allocs/op
```

## 设计思路



# 交流学习
![扫码加微信好友](https://github.com/lxt1045/wechatbot/blob/main/resource/Wechat-lxt.png "微信")