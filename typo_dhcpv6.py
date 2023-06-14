from ipmininet.iptopo import IPTopo
from ipmininet.ipnet import IPNet
from ipmininet.cli import IPCLI
import os

class MyTopology(IPTopo):

    def build(self, *args, **kwargs):

        h1 = self.addHost("h1", use_v4=False)
        r1 = self.addHost("r1", use_v4=False)
        r2 = self.addHost("r2", use_v4=False)

        h1r1 = self.addLink(h1, r1)
        r1r2 = self.addLink(r1, r2)
        r1r2[r2].addParams(ip="fd00::3:2/112")

        super().build(*args, **kwargs)

def close():
    try:
        os.remove('./if_name_to_ipv6')
        os.remove('./nohup.out')
        os.remove('./qiao')
    finally:
        pass

if __name__ == '__main__':
    net = IPNet(topo=MyTopology())
    # net['h1'].cmd('ip -6 r add default via fd00::1:1')
    net['h1'].cmd('ethtool -K h1-eth0 tx off')

    net['r1'].cmd('ethtool -K r1-eth0 tx off')
    net['r1'].cmd('ethtool -K r1-eth1 tx off')

    net['r2'].cmd('ethtool -K r2-eth0 tx off')
    net['r2'].cmd('ethtool -K r2-eth1 tx off')

    net['r2'].cmd('ip r add fd00::1:0/112 via fd00::3:1 dev r2r1')

    net['r1'].cmd('nohup ./qiao --dhcpv6 &')
    net['r2'].cmd('nohup ./qiao --dhcpv6 &')

    try:
        net.start()
        IPCLI(net)
    finally:
        net.stop()
        close()