package config

/**
* @program: code
* @description: used to config the serverIp, client and servers would read it
* @author: Shixiang Long
* @create: 2022-03-05 21:40
**/

var Servers = []string{"127.0.0.1:1230", "127.0.0.1:1231", "127.0.0.1:1232", "127.0.0.1:1233", "127.0.0.1:1234"}

func IndexOf(value string, slice []string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
