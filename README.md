# QIAO: 一个轻量级的基于高级语言的路由器软件

##### 2023.6.5 更新
根目录新增目录 dhcpv6，/setup 拷贝了 dhcpv6 虚拟网络搭建和测试脚本

##### 2023.6.9 更新
1. 迁移到 vmware ubuntu 成功
2. make build && make run 成功
> mininet h1 ping h2 成功
> mininet h2 ping h1 成功
3. 创建新分支 dhcpv6
4. 尝试只 merge 某个文件


# How to play with it?
1. 需要在ubuntu系统下运行.
2. 本项目最初尝试使用mininet进行拓扑的搭建和测试.
3. 但是Qiao是基于IPv6寻址的, 由于mininet不支持IPv6, 我们使用了ipmininet作为mininet的扩展，进行环境的搭建。
4. 安装mininet，请参考http://mininet.org/download/
5. 安装ipmininet，请参考https://ipmininet.readthedocs.io/en/latest/install.html#manual-installation
6. 目前支持go1.18版本, 不确定在其他版本上是否适配
7. 运行go mod tidy安装基础依赖
8. 运行make run, 运行ipmininet搭建拓扑
9. 在mininet> 终端下输入h1 ping h2, 若能够ping通，说明环境搭建成功


## 关于环境搭建
可以使用setup里的sh文件，用ip命令搭建拓扑;
也可以使用ipmininet.

关于ipmininet的基本使用:
在typo.py文件里，有关于网络拓扑的配置，初始化的时候，r1、r2、r3会在后台运行qiao
```
make run
```
执行make run之后，出现mininet>
可以与之交互。

### 使用go1.18以上的环境
ubuntu如果用apt install go安装的话，默认版本是1.13。
可以使用gvm来对go版本进行管理:https://github.com/moovweb/gvm
```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /root/.gvm/scripts/gvm
gvm use go1.18
```

### 关于typo.py
在typo.py里，我们给出了两个主机和三个路由器：
```
        h1 = self.addHost("h1", use_v4=False)
        r1 = self.addHost("r1", use_v4=False)
        r2 = self.addHost("r2", use_v4=False)
        r3 = self.addHost("r3", use_v4=False)
        h2 = self.addHost("h2", use_v4=False)
```

如果要在特定的主机下执行命令，在前面加主机的前缀就好，例如h1要ping h2，可以执行h1 ping h2。
另一个基础用法是，给这个主机打开一个新终端，在这个终端下执行命令并进行观测，可以执行
```
xterm h1 r1 r2 r3 h2
```

另一个调试常用方法是，注释掉以下typo.py里的三行代码
```
    net['r1'].cmd('nohup ./qiao &')
    net['r2'].cmd('nohup ./qiao &')
    net['r3'].cmd('nohup ./qiao &')
```
使用xterm手动让路由器运行qiao软件，并观察qiao的调试信息，与此同时还可以打开wireshark对机器收到的数据包进行分析。
