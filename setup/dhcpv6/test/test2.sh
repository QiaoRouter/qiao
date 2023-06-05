#!/bin/sh
# Test 2
# Check router advertisement


# 这个命令是在指定的网络命名空间 PC1 中执行 rdisc6 pc1r1 命令。
# 具体来说，它会启动 rdisc6 工具，以便在 PC1 命名空间中向 pc1r1 路由器发送 Router Advertisement 消息。

# rdisc6 命令用于发送 IPv6 Router Advertisement 消息。

# 在 IPv6 网络中，路由器通过发送 Router Advertisement 消息来通知其他主机连接到该网络的路由器信息。
# 这些信息包括路由器的 IPv6 地址、网络前缀和其他可选信息。

# 通过在 PC1 命名空间中执行 rdisc6 pc1r1 命令，可以模拟 PC1 主机接收来自 pc1r1 路由器的 Router Advertisement 消息，
# 以便 PC1 主机可以了解到 pc1r1 路由器的信息并配置其 IPv6 地址。

ip netns exec PC1 rdisc6 pc1r1
