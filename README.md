# p2pssh
## p2pssh is based on libp2p
## 1. set p2pssh deamon
### It is a daemon process just like sshd
```
wany@WANY:~$ ./p2pssh daemon
Your PeerID is :Qmb3Tn7SPGxXn7ENagQUY9mVwhNqhr8Ac4C9mWiDqrrbST
Listen:[/ip4/127.0.0.1/udp/9000/quic /ip4/192.168.0.133/udp/9000/quic /ip4/172.17.0.1/udp/9000/quic]
```
## 2. In another computer, try to ping this peer PeerID

```
wany@WANY:~$ ./p2pssh ping Qmb3Tn7SPGxXn7ENagQUY9mVwhNqhr8Ac4C9mWiDqrrbST
ping took: 35.50429ms
ping took: 33.223504ms
ping took: 34.923255ms
ping took: 33.357161ms
ping took: 38.924518ms
ping took: 34.280522ms
ping took: 32.667023ms
ping took: 34.534543ms
```
## 3. login it
```
wany@WANY:~$ ./p2pssh login wany@Qmb3Tn7SPGxXn7ENagQUY9mVwhNqhr8Ac4C9mWiDqrrbST
Password:  //input your password
wany@WANY:~$
```