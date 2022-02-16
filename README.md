# cs2510_distributed_system



坑：

Golang rpc 必须关闭 conn,listener .Close

如果用ctl + c去关闭程序 ,defer是不会执行的

