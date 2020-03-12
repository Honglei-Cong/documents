
### process of migration

和Bison Trails确认：

1. 如何重启节点
2. 开ticket的反应时间
3. 重置google auth code

操作步骤：

1. 开始 google meet
2. bison 通过 1password，获得节点钱包
3. bison 通过 telegram，获得密码
4. 登录bison trails网站
5. 销毁原bison validator节点
6. 确认bison trails 节点启动命令
7. 在bison trails网站上，创建新的节点（AWS West-1, California/Oregon）
   1. bison trails节点启动
8. 获得新节点的IP
   1. peers.recent
   2. peers.rsv
9. 更新所有节点的IP列表（防火墙）(只更新共识节点配置，不重启节点)
   1. 我们更新seed/dapp/consensus节点的IP列表
   2. 火币更新节点的IP列表
   3. OK更新节点的IP列表
10. 停止我们的Phecda节点
11. 依次重启 我们/火币/ OK 的节点 
    1.  火币OK的节点在我们重启完成后，通知他们依次重启
12. 完成后，通知Bison Trails，其他的节点firewall已经配置好，
13. bison trails 更新 peers.recent文件，并重启节点
14. 监控Bison trails节点高度
15. bison trails同步完成后，通过浏览器确认参与出块
16. 所有节点确认和bison trails节点网络链接正常 
    1.  (netstat -antp | grep ontology) 查看ontology进程的网络链接

