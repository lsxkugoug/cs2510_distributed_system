### Normal Usage

Assume you are in `code` dirctory

1.make configuration in config file, which is in /code/src/config/config.go

if you want to test in different machine, you should replace 127.0.0.1 to other ips. and you should set Nr and Nw, if N == 5, you can set Nr = 3, Nw = 3

```golang
var Servers = []string{"127.0.0.1:1230", "127.0.0.1:1231", "127.0.0.1:1232", "127.0.0.1:1233", "127.0.0.1:1234"}
Nr = 3
Nw = 3
```

2.assmue you are under `code` directory, open 5 servers:

```go
go run src/entry/serverMain.go
please input serverId, which is corresponding to the config file, indexed from 0

// then you should input the serverId which is corresponding to config.Servers's index
```

3.assmue you are under `code` directory, open client, and try to new b = 1 on server0~4

```go
go run src/entry/clientMain.go

//(base) shixianglong@lipingde-iPhone code % go run src/entry/clientMain.go
//please input the operation you want, 0:update or new 1: retrieve data
//0
//please input the serverId you want to make operation, you can input 0 ~ 4
//0
//please input the key you want update or create
//b
//please input the corresponding value
//1
```

### Test

#### 1）How to test it

**Test it manually:**

Assuming you are in `code` directory, 

You can run the test by:

```bash
cd ./src
go test -v -run Test1 
```

Note that I just use the `testing` package to execute test, the 'PASS' is not means it is correct. 

**Using docker:**

Or use docker test it automatically:

```shell
① open terminal under "p3" directory

② docker image build -f Dockerfile -t cs2510_dk .

③ docker run  cs2510_dk
```

**!!!!!You can find the test log under `p1/code/src/testEntry/logs/server0~server4.txt`**

The test design is writen in the report document.
