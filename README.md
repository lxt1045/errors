# errors 包

## 设计思路

1、记录报错现场的调用栈  
2、报错后，返回路径上可以添加额外的日志信息，但是不会记录该位置的调用栈，只会记录添加信息位置的函数信息  
3、在最终处理错误的地方打印错误信息，err.Err()返回json格式日志信息  
  

## 使用方式

1、在产生错的地方调用：
```go
// code,msg,errlog:=<http需要返回的业务错误码>, <http需要返回的业务错误信息>, <需要打印的日志信息，比如：上下文变量，提示log>
err := errors.NewErr(code,msg,errlog)
```

2、在得到函数报错时，如果是中间处理层，可以不输出日志，调用：
```go
if x,err:=pkg.DoFunc(x,y){
    // errlog:= <需要打印的日志信息，比如：上下文变量，提示log>
    err = errors.Trace(err, errlog)
    return err
}
if x,err:=pkg.DoFunc(x,y){
    // errlog:= <需要打印的日志信息，比如：上下文变量，提示log>
    err = errors.Trace(err, errors.Fields("req", req, "order", order))
    return err
}
```

3、在得到函数报错时，如果是顶层调用者，输出日志，调用：
```go
x,err:=pkg.DoFunc(x,y)
if err!=nil {
    err = errors.Trace(err, errors.Fields("x", x, "y", y))
    log.Warn(err.Err())
    return err
}
```

4、日志输出内容：   
err.Err() 会返回一个json字符串，包含一下信息：
```json
{
    "code": 21200017,
    "message": "number of cards exceeds the limit",SetTraceID(uuid.New().String()) 
    "logs": [ // 函数调用栈
        {
            "ts": "2021-05-11T10:08:37.791737+08:00", // errors.NewErr 的时间 errors.NewErr，会输出 stack 
            "stack": [
                "(xxxxx/models/card.go:123) models.(*Account).GetOrCreate",
                "(xxxxx/handler.go:119) main.(*Handler).bindAccount"
            ],
            "msg": { // msg 就是 errors.NewErr(code,msg,errlog) 中的errlog，interface{}格式
                "a": {
                    "id": 0,
                    "user_id": 1030
                }
            }
        },
        {
            "ts": "2021-05-11T10:08:37.832787+08:00", //errors.Trace 的时间，errors.Trace 只输出 func
            "func": "(xxxxx/handler.go:122) main.(*Handler).bindAccount",
            "msg": {
                "pAccount": null,
                "req": {
                    "user_id": 1030,
                    "account_type": 2,
                    "account_number": ""
                }
            }
        }
    ]
}
```