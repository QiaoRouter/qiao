## Update on June 5, 2023
A new directory named 'dhcpv6' is added in the root directory. The /setup has copied the scripts for building and testing the dhcpv6 virtual network.

## Update on June 9, 2023
1. Successfully migrated to vmware ubuntu.
2. Successfully built and run the program with 'make build' and 'make run'.
> 'mininet h1 ping h2' successfully pings.
> 
> 'mininet h2 ping h1' successfully pings.

3. Created a new branch called 'dhcpv6'.
4. Attempted to merge only a certain file.

# How to play with it?
1. It needs to be run under the Ubuntu system.
2. This project initially attempted to use mininet for topology construction and testing.
3. But Qiao is based on IPv6 addressing, and since mininet does not support IPv6, we used ipmininet as an extension of mininet for environment construction.
4. To install mininet, please refer to http://mininet.org/download/
5. To install ipmininet, please refer to https://ipmininet.readthedocs.io/en/latest/install.html#manual-installation
6. Currently supports go1.18 version, it is uncertain whether it is adapted to other versions.
7. Run 'go mod tidy' to install basic dependencies.
8. Run 'make run', run ipmininet to build topology.
9. Input 'h1 ping h2' in the 'mininet>' terminal, if it can ping successfully, it indicates that the environment is successfully built.

# About environment construction
You can use the sh file in the setup to build the topology with the ip command;
You can also use ipmininet.

For basic use of ipmininet:
In the typo.py file, there is the configuration of network topology. When initialized, r1, r2, and r3 will run qiao in the background.

```
make run
```

After executing 'make run', 'mininet>' appears.
You can interact with it.

# Use an environment above go1.18

If you install go with 'apt install go' in ubuntu, the default version is 1.13.
You can use gvm to manage the go version: https://github.com/moovweb/gvm

```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /root/.gvm/scripts/gvm
gvm use go1.18
```

# About typo.py
In typo.py, we provided two hosts and three routers:

```
        h1 = self.addHost("h1", use_v4=False)
        r1 = self.addHost("r1", use_v4=False)
        r2 = self.addHost("r2", use_v4=False)
        r3 = self.addHost("r3", use_v4=False)
        h2 = self.addHost("h2", use_v4=False)
```

If you want to execute a command under a specific host, just add the prefix of the host. For example, if h1 wants to ping h2, you can execute 'h1 ping h2'.
Another basic usage is to open a new terminal for this host, execute the command and observe in this terminal, you can execute:

```
xterm h1 r1 r2 r3 h2
```

Another common method for debugging is to comment out the following three lines of code in typo.py:

```
    net['r1'].cmd('nohup ./qiao &')
    net['r2'].cmd('nohup ./qiao &')
    net['r3'].cmd('nohup ./qiao &')
```

Use xterm to manually run the Qiao software on the router and observe the debug information from Qiao. Meanwhile, you can also open Wireshark to analyze the packets received by the machine.




