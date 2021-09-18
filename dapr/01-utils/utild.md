##### host.go

```go
if val, ok := os.LookupEnv(HostIPEnvVar); ok && val != "" {
		return val, nil
	}
```

首先读取环境变量，读取到直接返回

```go
conn, err := net.Dial("udp", "8.8.8.8:80")
defer conn.Close()
return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
```

通过访问谷歌dns的80端口查询本机IP地址

```go
addrs, err := net.InterfaceAddrs()
for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
```

当网络故障时查询本地网络地址，获取全部地址，然后去除本地回环地址并且可以转换为ipv4地址，最终确定本地的地址。

最终函数GetHostAddress返回本机的出口IP地址。

##### resolvconf.go

用来获取本pod在k8s集群中的域名

```go
func getResolvContent(resolvPath string) ([]byte, error) {
	return ioutil.ReadFile(resolvPath)
}
```

读取dns配置文件"/etc/resolv.conf"

```go
func getResolvSearchDomains(resolvConf []byte) []string {
	var (
		domains []string
		lines   [][]byte
	)
	# 将文件内容由[]byte类型转换为scanner类型
	scanner := bufio.NewScanner(bytes.NewReader(resolvConf))
	for scanner.Scan() {
		line := scanner.Bytes()
		# 过滤掉注释内容
		commentIndex := bytes.Index(line, []byte(commentMarker))
		if commentIndex == -1 {
			lines = append(lines, line)
		} else {
			lines = append(lines, line[:commentIndex])
		}
	}

	for _, line := range lines {
		match := searchRegexp.FindSubmatch(line)
		if match == nil {
			continue
		}
		domains = strings.Fields(string(match[1]))
	}

	return domains
}
```

解析配置文件，返回本机的域名

```go
func getClusterDomain(resolvConf []byte) (string, error) {
	var kubeClusterDomian string
	searchDomains := getResolvSearchDomains(resolvConf)
	sort.Strings(searchDomains)
	if len(searchDomains) == 0 || searchDomains[0] == "" {
		kubeClusterDomian = DefaultKubeClusterDomain
	} else {
		kubeClusterDomian = searchDomains[0]
	}
	return kubeClusterDomian, nil
}
```

返回第一个域名，没有则返回默认域名"cluster.local"

```go
func GetKubeClusterDomain() (string, error) {
	resolvContent, err := getResolvContent(defaultResolvPath)
	if err != nil {
		return "", err
	}
	return getClusterDomain(resolvContent)
}
```

对外暴露的函数

##### utils.go