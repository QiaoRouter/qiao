# Qiao: 一个轻量级的基于高级语言的路由器软件

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