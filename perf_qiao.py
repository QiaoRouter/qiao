from ipmininet.iptopo import IPTopo
from ipmininet.ipnet import IPNet
from ipmininet.cli import IPCLI
import os
import time

N = 10  # number of hosts on each side of the network

class MyTopology(IPTopo):

    def build(self, *args, **kwargs):
        r1 = self.addHost('r1', use_v4=False)
        r2 = self.addHost('r2', use_v4=False)
        r3 = self.addHost('r3', use_v4=False)

        for i in range(N):
            h1 = self.addHost(f'h1_{i}', use_v4=False)
            h2 = self.addHost(f'h2_{i}', use_v4=False)

            h1r1 = self.addLink(h1, r1)
            h2r3 = self.addLink(h2, r3)

            h1r1[h1].addParams(ip=f"2001:1:{i+1}::1/64")
            h1r1[r1].addParams(ip=f"2001:1:{i+1}::2/64")

            h2r3[h2].addParams(ip=f"2001:3:{i+1}::1/64")
            h2r3[r3].addParams(ip=f"2001:3:{i+1}::2/64")

        self.addLink(r1, r2, params1={"ip": "2001:12::1/64"}, params2={"ip": "2001:12::2/64"})
        self.addLink(r2, r3, params1={"ip": "2001:23::1/64"}, params2={"ip": "2001:23::2/64"})

        super().build(*args, **kwargs)

if __name__ == "__main__":
    net = IPNet(topo=MyTopology())
    net['r1'].cmd('sysctl -w net.ipv6.conf.all.forwarding=1')
    net['r3'].cmd('sysctl -w net.ipv6.conf.all.forwarding=1')
    for i in range(N):
        net[f'h1_{i}'].cmd('ip -6 route add default via 2001:1:{}::2'.format(i+1))
        net[f'h1_{i}'].cmd('ethtool -K h1_{}-eth0 tx off'.format(i+1))
        net[f'h2_{i}'].cmd('ip -6 route add default via 2001:3:{}::2'.format(i+1))
        net[f'h2_{i}'].cmd('ethtool -K h2_{}-eth0 tx off'.format(i+1))

        net['r1'].cmd('ethtool -K r1-eth{} tx off'.format(i+1))
        net['r3'].cmd('ethtool -K r2-eth{} tx off'.format(i+1))


    net['r3'].cmd('sysctl -w net.ipv6.conf.all.forwarding=1')
    net['r3'].cmd('ethtool -K r3-eth0 tx off')
    net['r3'].cmd('ethtool -K r3-eth1 tx off')

    net['r1'].cmd('nohup ./qiao &')
    net['r2'].cmd('nohup ./qiao &')
    net['r3'].cmd('nohup ./qiao &')
    try:
        net.start()
        time.sleep(5)  # wait for network to start
        for i in range(N):
            net[f'h1_{i}'].cmd('nohup ping h2_{} &'.format(i))
            net[f'h2_{i}'].cmd('nohup ping h1_{} &'.format(i))
        IPCLI(net)
    finally:
        net.stop()
        close()

