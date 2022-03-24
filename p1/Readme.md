### Normal Usage

Assume you are in `code` dirctory

1.firstly, open server:

```go
go run ./src/entry/serverMain.go

// then, the Ip address would appear
// eg. the server's ip is :172.26.19.135:1234
```

2.Second, open a client, you should enter the related parameter one by one

```go
go run src/entry/clientMain.go
//eg:
//please input your name:
//human
//please input your group:
//1
//please input <serverIp:severPort:>
//172.26.19.135:1234
//please input the port you want to use this app, if you want to open multiple apps on the same machine, the ports should be differenteg: 1200 ~ 1230
//1220

Then you can input your message, the message would be forward to the user who in the same group
```

### Test

#### 1）How to test it

**Test it manually:**

There are three logs under `code/src/testEntry/test{textnumber}/log/Alice.txt`: Alice.txt,  Bob.txt and Chad txt

Assuming you are in `code` directory, 

You can run the test by:

```bash
cd ./src
go test -v -run Test1 
go test -v -run Test2
go test -v -run Test3 

or go test -v
```

Note that I just use the `testing` package to execute test, the 'PASS' is not means it is correct. 

**Using docker:**

Or use docker test it automatically:

```shell
① open terminal under "p1" directory

② docker image build -f Dockerfile -t cs2510_dk .

③ docker run  cs2510_dk
```

**!!!!!You can find the test log under `p1/code/src/testEntry`**

#### 2) test1: get history message test

In this test,  the {testnumber} = 1

Alice online first, and send the "I am Alice"
Then, after 5 seconds, Bob and Chad online, the would receive the Alice's history message: "I am Alice"
**The log example:**
**Alice.txt**:15:29:58  (2022-02-13 15:28:58 ,group:1,name:Alice):I am Alice, nice to meet you 

**Bob.txt**: 15:29:03 (2022-02-13 15:28:58 ,group:1,name:Alice):I am Alice, nice to meet you

**Chad.txt**: 15:29:03 (2022-02-13 15:28:58 ,group:1,name:Alice):I am Alice, nice to meet you

15:29:03 is the time Alice send the message

15:29:03 is the time Bob and Chad online and receive this message

#### 3) test2: group test

In this test,  the {testnumber} = 2

Alice, Bob, and Chad are online. Bob sends a message to all, Chad and Alice receives the message (The sender Bob doesnt receive the message from the server). Alice sends a message to all, Bob and Chad receives it (but not Alice). Doug, not part of the group, joins the server but receives no message.

**The log example:**
**Alice.txt**:17:09:53 (2022-02-13 17:09:52 ,group:1,name:Bob):I am Bob, nice to meet you

**Bob.txt**: 17:09:52 (2022-02-13 17:09:52 ,group:1,name:Bob):I am Bob, nice to meet you 

**Chad.txt**: 17:09:53 (2022-02-13 17:09:52 ,group:1,name:Bob):I am Bob, nice to meet you

**Doug.txt:** null



#### 4) test3: offline and re-online history message test

In this test,  the {testnumber} = 3

Three member Alice, Bob, Chad. All of them in  the same group

Alice send a message "I am Alice, nice to meet you" to the group, Bob would receive it.

Then Bob offline.

Chad  online and send the message "I am Chad, nice to meet you"to the group

Then Bob online, he only would receive the "I am Chad, nice to meet you", because  "I am Alice, nice to meet you" he have already read. . And since the log is recreate, the log only has "I am Chad, nice to meet you"

**The log example:**

**Alice.txt**:

2022/02/13 22:08:21 (2022-02-13 22:08:21 ,group:1,name:Alice):I am Alice, nice to meet you
2022/02/13 22:08:25 (2022-02-13 22:08:25 ,group:1,name:Chad):I am Chad, nice to meet you

**Bob.txt**:  
2022/02/13 22:08:26 (2022-02-13 22:08:25 ,group:1,name:Chad):I am Chad, nice to meet you

**Chad.txt**: 
2022/02/13 22:08:25 (2022-02-13 22:08:21 ,group:1,name:Alice):I am Alice, nice to meet you
2022/02/13 22:08:25 (2022-02-13 22:08:25 ,group:1,name:Chad):I am Chad, nice to meet you



