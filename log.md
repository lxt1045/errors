# 日志库的性能提升


# 1. Go 下的性能追踪
由于笔者能力有限，这里并不会讲 eBPF、SystemTap 和 DTrace 的使用，仅根据自己经验简单介绍 Go 下的 pprof 工具的一种使用方式。
## 1.1 pprof 工具使用：
参考文档:\
[golang 性能优化分析工具 pprof (上) - 基础使用介绍](https://www.cnblogs.com/jiujuan/p/14588185.html)\
[《Go 语言编程之旅》6.1 Go 大杀器之性能剖析 PProf（上）](https://golang2.eddycjy.com/posts/ch6/01-pprof-1/)

我们知道 Go 的 pprof 工具有三种使用方式：
runtime/pprof：采集程序（非 Server）的指定区块的运行数据进行分析。\
net/http/pprof：基于 HTTP Server 运行，并且可以采集运行时数据进行分析。\
go test：通过运行测试用例，并指定所需标识来进行采集。

这里仅介绍 go test 方式，分为这几个步骤：\
1. 这种方式需要先创建一个 Benchmark 测试函数。
```go
func Benchmark_Sample(b *testing.B) {
    for i := 0; i < b.N; i++ {
        f()
    }
}
```
2. 执行测试，生成 cpu.prof 测试文档
```sh
go test -benchmem -run=^$ -bench ^BenchmarkZapCaller$ github.com/lxt1045/errors/zap -count=1 -v -cpuprofile cpu.prof
```
3. 执行编译命令生成二进制
```sh
go test -benchmem -run=^$ -bench ^BenchmarkZapCaller$ github.com/lxt1045/errors/zap -c -o test.bin 
```
4. 使用 go tool 命令解析 cpu.prof 测试文档
```sh
go tool pprof ./test.bin cpu.prof
```
5. 使用以下命令查看:\
5.1 查看 graph 图
```sh
web 
```
5.2 查看排行
```sh
top n 
```
输出例子：
```
Showing nodes accounting for 9270ms, 67.91% of 13650ms total
Dropped 172 nodes (cum <= 68.25ms)
Showing top 10 nodes out of 116
      flat  flat%   sum%        cum   cum%
    2490ms 18.24% 18.24%     2490ms 18.24%  runtime.madvise
    1880ms 13.77% 32.01%     1880ms 13.77%  runtime.pthread_cond_signal
     970ms  7.11% 39.12%     1230ms  9.01%  [test.bin]
     920ms  6.74% 45.86%      940ms  6.89%  runtime.pthread_cond_wait
     720ms  5.27% 51.14%     2570ms 18.83%  github.com/lxt1045/json.parseObj
     640ms  4.69% 55.82%      640ms  4.69%  github.com/bytedance/sonic/internal/native/avx2.__native_entry__
     550ms  4.03% 59.85%      740ms  5.42%  github.com/lxt1045/json.(*tireTree).Get
     530ms  3.88% 63.74%      710ms  5.20%  runtime.scanobject
     300ms  2.20% 65.93%      300ms  2.20%  runtime.memmove
     270ms  1.98% 67.91%      270ms  1.98%  runtime.kevent
```
5.3 查看函数内每行代码开销
```sh
list func_name 
```
输出例子：
```sh
Total: 13.65s
ROUTINE ======================== github.com/lxt1045/json.(*tireTree).Get in /Users/bytedance/go/src/github.com/lxt1045/json/tire_tree.go
     550ms      740ms (flat, cum)  5.42% of Total
         .          .    282:           return nil
         .          .    283:   }
         .          .    284:
         .          .    285:   return nil
         .          .    286:}
      30ms       30ms    287:func (root *tireTree) Get(key string) *TagInfo {
      10ms       10ms    288:   status := &root.tree[0]
         .          .    289:   // for _, c := range []byte(key) {
      20ms       20ms    290:   for i := 0; i < len(key); i++ {
      10ms       10ms    291:           c := key[i]
         .          .    292:           k := c & 0x7f
     160ms      160ms    293:           next := status[k]
         .          .    294:           if next.next >= 0 {
         .          .    295:                   i += int(next.skip)
      10ms       10ms    296:                   status = &root.tree[next.next]
         .          .    297:                   continue
         .          .    298:           }
      10ms       10ms    299:           if next.idx >= 0 {
      40ms       40ms    300:                   tag := root.tags[next.idx]
     250ms      440ms    301:                   if len(key) > len(tag.TagName) && key[len(tag.TagName)] == '"' && tag.TagName == key[:len(tag.TagName)] {
      10ms       10ms    302:                           return tag
         .          .    303:                   }
         .          .    304:           }
         .          .    305:           return nil
         .          .    306:   }
         .          .    307:
```
5.4 通过浏览器查看测试结果（火焰图、graph 图等）
```sh
go tool pprof -http=:8080 cpu.prof
```
执行后，通过浏览器打开 http://localhost:8080/ 链接就可以查看了。

# 2. logrus 和 zap 中获取代码行号损耗

logrus 的占 17% 左右：
![logrus](https://github.com/lxt1045/errors/blob/main/resource/logrus_flamegraph.jpg)

zap 的占 28% 左右：
![logrus](https://github.com/lxt1045/errors/blob/main/resource/zap_flamegraph.jpg)

所以，我们就拿 "获取代码行号" 逻辑开到，做优化。

经过优化，logrus 和 zap 能减少 1300ns ~ 2500ns 的损耗。\
结果如下：
[BenchmarkLog](https://github.com/lxt1045/errors/blob/main/zap/zap_test.go#L85)
```sh
BenchmarkLog/logrus
BenchmarkLog/logrus-12         	  387568	      3080 ns/op	    1361 B/op	      23 allocs/op
BenchmarkLog/logrus+caller
BenchmarkLog/logrus+caller-12  	  178818	      6768 ns/op	    2355 B/op	      34 allocs/op
BenchmarkLog/logrus+lxt_caller
BenchmarkLog/logrus+lxt_caller-12         	  237561	      4775 ns/op	    2082 B/op	      31 allocs/op
BenchmarkLog/zap
BenchmarkLog/zap-12                       	 1510062	       795.7 ns/op	     152 B/op	       3 allocs/op
BenchmarkLog/zap+caller
BenchmarkLog/zap+caller-12                	  468783	      2354 ns/op	     401 B/op	       6 allocs/op
BenchmarkLog/zap+lxt_caller
BenchmarkLog/zap+lxt_caller-12            	 1000000	      1054 ns/op	     409 B/op	       4 allocs/op
BenchmarkLog/zap-sugar
BenchmarkLog/zap-sugar-12                 	 1301520	       844.2 ns/op	     112 B/op	       4 allocs/op
BenchmarkLog/zap-sugar+caller
BenchmarkLog/zap-sugar+caller-12          	  386512	      2919 ns/op	     361 B/op	       7 allocs/op
BenchmarkLog/zap-sugar+lxt_caller
BenchmarkLog/zap-sugar+lxt_caller-12      	 1000000	      1020 ns/op	     176 B/op	       5 allocs/op
BenchmarkLog/lxt_caller
BenchmarkLog/lxt_caller-12                	34919370	        34.30 ns/op	      24 B/op	       1 allocs/op
```

# 3. 优化原理
