## ホストとWSL2のポート接続方法

1. WSL2で実行
```console
ip addr show eth0
```

アドレスを取得する
```console
# XXX.XX.XX.X/XX の部分

3: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether 00:15:5d:00:4a:41 brd ff:ff:ff:ff:ff:ff
    inet XXX.XX.XX.X/XX brd 172.18.63.255 scope global eth0
       valid_lft forever preferred_lft forever
```

2. ポート接続
powershellを管理者権限で立ち上げて実行
```console
netsh interface portproxy add v4tov4 listenport=8081 listenaddress=0.0.0.0 connectport=8081 connectaddress=<WSL2_IP_ADDRESS>
```

3. ファイアーウォール許可

```console
netsh advfirewall firewall add rule name="WSL2 port 8081" dir=in action=allow protocol=TCP localport=8081
```
