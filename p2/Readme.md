### Normal Usage

Assume you are in `code` dirctory

1.make configuration in config file, which is in /code/src/config/config.go

if you want to test in different machine, you should replace 127.0.0.1 to other ips.

```golang
var Servers = []string{"127.0.0.1:1230", "127.0.0.1:1231", "127.0.0.1:1232", "127.0.0.1:1233", "127.0.0.1:1234"}
```

2.assmue you are under `code` directory, open 5 servers:

```go
go run src/entry/serverMain.go
please input serverId, which is corresponding to the config file, indexed from 0

// then you should input the serverId which is corresponding to config.Servers's index
```

3.assmue you are under `code` directory, open client, and try to new b = 1 on server0 and server0 would broadcast it to all of the servers

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
go test -v -run Test2
```

Note that I just use the `testing` package to execute test, the 'PASS' is not means it is correct. 

**Using docker:**

Or use docker test it automatically:

```shell
① open terminal under "p2" directory

② docker image build -f Dockerfile -t cs2510_dk .

③ docker run  cs2510_dk
```

**!!!!!You can find the test log under `p1/code/src/test1.txt or test2.txt`**

#### (1) No conflict

① open 0 ~ 3 servers, but not 4

② client new b = 1 on server 0

③ server 4 online

③ client update  b = 2 on server 1

④ client retrieve all of the value of b in each servers, print the corresponding key, value and vector

Sample log:

```
2022/03/11 08:57:22 client new b = 1 to server0, in this time, 0 ~ 3 are online, but 4 offline
2022/03/11 08:57:26 4 online now
2022/03/11 08:57:28 client update b = 2 to server0, in this time, 0 ~ 4 are online
2022/03/11 08:57:30 client retrieve all of the servers b = ?
2022/03/11 08:57:30 server0:
2022/03/11 08:57:30 the value:[2], the vecotor:[[2 2 0 0 0]] 
2022/03/11 08:57:30 server1:
2022/03/11 08:57:30 the value:[2], the vecotor:[[2 2 0 0 0]] 
2022/03/11 08:57:30 server2:
2022/03/11 08:57:30 the value:[2], the vecotor:[[2 2 0 0 0]] 
2022/03/11 08:57:30 server3:
2022/03/11 08:57:30 the value:[2], the vecotor:[[2 2 0 0 0]] 
2022/03/11 08:57:30 server4:
2022/03/11 08:57:30 the value:[2], the vecotor:[[2 2 0 0 0]] 
```

#### (2)  conflict

①  open 0 ~ 4 (all) servers

② client new b = 1 on server 0    and  client new b = 2 on server 3

In this case, server 1, 2, 4 will get (1,0,0,0,0), (0, 0, 0, 1, 0). Confliction happens. However, because it is concurrent situation. Maybe need test several times.

Sample log:

```
2022/03/11 09:01:04 client new b = 1 to server0, in this time, 0 ~ 3 are online, but 4 offline
2022/03/11 09:01:06 client update b = 1 to server0, and send b = 2 to server 4 simutaneously, 0 ~ 4 are online
2022/03/11 09:01:08 client retrieve all of the servers b = ?
2022/03/11 09:01:08 server0:
2022/03/11 09:01:08 the value:[1 2], the vecotor:[[2 0 0 0 0] [0 0 0 2 0]] 
2022/03/11 09:01:08 conflict happens !!!!!
2022/03/11 09:01:08 server1:
2022/03/11 09:01:08 the value:[1 2], the vecotor:[[2 0 0 0 0] [0 0 0 2 0]] 
2022/03/11 09:01:08 conflict happens !!!!!
2022/03/11 09:01:08 server2:
2022/03/11 09:01:08 the value:[1 2], the vecotor:[[2 0 0 0 0] [0 0 0 2 0]] 
2022/03/11 09:01:08 conflict happens !!!!!
2022/03/11 09:01:08 server3:
2022/03/11 09:01:08 the value:[2 1], the vecotor:[[0 0 0 2 0] [2 0 0 0 0]] 
2022/03/11 09:01:08 conflict happens !!!!!
2022/03/11 09:01:08 server4:
2022/03/11 09:01:08 the value:[2 1], the vecotor:[[0 0 0 2 0] [2 0 0 0 0]] 
2022/03/11 09:01:08 conflict happens !!!!!
```
