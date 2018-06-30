# http-proxy-dev

http-proxy-dev works as a proxy recieve http request and send towards any new target. 

u can dump full of request or response out in http message as u wish.

```
Usage of ./http-proxy-dev:
  -l string
        path of log file (default "/dev/stdout")
  -m int
        switch if print full http message (default 1)
  -p int
        port which server listened (default 80)
  -t string
        port which proxy redirect to (default "http://127.0.0.1:8080")
```
