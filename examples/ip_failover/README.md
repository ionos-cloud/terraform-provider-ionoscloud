# Multiple NICs under the same IP Failover

This example aims to exemplify the way in which a secondary NIC can be added to an IP Failover, following these steps: 
1) Creating NIC A with failover IP on LAN 1 
2) Create NIC B unde the same LAN but with a different IP 
3) Create the IP Failover on LAN 1 with NIC A and failover IP of NIC A (A becomes now "master", no slaves)
4) Update NIC B IP to be the failover IP ( B becomes now a slave, A remains master)

To test this please run firstly [the main.tf from inital_structure folder](initial_infrastructure/main.tf), followed by [the one in update_nic folder](update_nic/main.tf).

After this you can create a new NIC C, NIC D and so on, in LAN 1, directly with the failover IP.