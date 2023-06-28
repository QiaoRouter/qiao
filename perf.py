from ipmininet.iptopo import IPTopo
from ipmininet.ipnet import IPNet
from ipmininet.cli import IPCLI
import os


class MyTopology(IPTopo):

    def build(self, *args, **kwargs):
        hl_lst = []
        hr_lst = []
        for i in range(1, 1000):
            hl = self.addHost("h" + str(i), user_v4=False)
            hl_lst.append(hl)
            hr = self.addHost("h" + str(i + 1000), user_v4=False)
            hr_lst.append(hr)

        r1 = self.addHost("r1", use_v4=False)
        r2 = self.addHost("r2", use_v4=False)
        r3 = self.addHost("r3", use_v4=False)

        for i in range(1, 1000):
            linkl = self.addLink(hl_lst[i], r1)
            linkr = self.addLink(hr_lst[i], r3)

        r1r2 = self.addLink(r1, r2)
        r1r2[r1].addParams(ip="fd00::3:1/112")
        r1r2[r2].addParams(ip="fd00::3:2/112")

        r2r3 = self.addLink(r2, r3)
        r2r3[r2].addParams(ip="fd00::4:1/112")
        r2r3[r3].addParams(ip="fd00::4:2/112")

        super().build(*args, **kwargs)


def close():
    try:
        os.remove('./if_name_to_ipv6')
        os.remove('./nohup.out')
    finally:
        pass


if __name__ == '__main__':
    net = IPNet(topo=MyTopology())
    net['h1'].cmd('ip -6 r add default via fd00::1:1')
    net['h1'].cmd('ethtool -K h1-eth0 tx off')
    #
    net['h2'].cmd('ip -6 r add default via fd00::5:2')
    net['h2'].cmd('ethtool -K h2-eth0 tx off')

    net['r1'].cmd('sysctl -w net.ipv6.conf.all.forwarding=1')
    net['r1'].cmd('ethtool -K r1-eth0 tx off')
    net['r1'].cmd('ethtool -K r1-eth1 tx off')

    net['r2'].cmd('ethtool -K r2-eth0 tx off')
    net['r2'].cmd('ethtool -K r2-eth1 tx off')
    net['r3'].cmd('sysctl -w net.ipv6.conf.all.forwarding=1')

    net['r3'].cmd('sysctl -w net.ipv6.conf.all.forwarding=1')
    net['r3'].cmd('ethtool -K r3-eth0 tx off')
    net['r3'].cmd('ethtool -K r3-eth1 tx off')

    net['r1'].cmd('nohup ./qiao &')
    net['r2'].cmd('nohup ./qiao &')
    net['r3'].cmd('nohup ./qiao &')
    try:
        net.start()
        IPCLI(net)
    finally:
        net.stop()
        close()
