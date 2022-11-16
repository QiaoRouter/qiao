from ipmininet.iptopo import IPTopo
from ipmininet.ipnet import IPNet
from ipmininet.cli import IPCLI

class MyTopology(IPTopo):

    def build(self, *args, **kwargs):

        h1 = self.addHost("h1", use_v4=False)
        r1 = self.addHost("r1", use_v4=False)
        r2 = self.addHost("r2", use_v4=False)
        r3 = self.addHost("r3", use_v4=False)
        h2 = self.addHost("h2", use_v4=False)

        h1r1 = self.addLink(h1, r1)
        h1r1[h1].addParams(ip="fd00::1:2/112")
        h1r1[r1].addParams(ip="fd00::1:1/112")

        r1r2 = self.addLink(r1, r2)
        r1r2[r1].addParams(ip="fd00::3:1/112")
        r1r2[r2].addParams(ip="fd00::3:2/112")

        r2r3 = self.addLink(r2, r3)
        r2r3[r2].addParams(ip="fd00::4:1/112")
        r2r3[r3].addParams(ip="fd00::4:2/112")

        r3h2 = self.addLink(r3, h2)
        r3h2[r3].addParams(ip="fd00::5:2/112")
        r3h2[h2].addParams(ip="fd00::5:1/112")

        super().build(*args, **kwargs)


if __name__ == '__main__':
    net = IPNet(topo=MyTopology())
    net['h1'].cmd('ip -6 r add default via fd00::1:1')
    net['h2'].cmd('ip -6 r add default via fd00::5:2')

    try:
        net.start()
        IPCLI(net)
    finally:
        net.stop()