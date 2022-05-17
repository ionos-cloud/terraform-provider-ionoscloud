# Multiple NICs under the same IP Failover

This example aims to exemplify the way in which a secondary NIC can be added to an IP Failover, following these steps: 
1) Creating NIC A with failover IP on LAN 1 
2) Create the IP Failover on LAN 1 with NIC A and failover IP of NIC A (A becomes now "master", no slaves)
3) Create NIC B unde the same LAN but with the failover IP, that depends on the creation of the IP Failover ( B becomes now a slave, A remains master)

To test this please run [the main.tf](main.tf) plan.

