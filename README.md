# AmpFilter
Filters lists containing various IPs that responded with a x. bytes response by sending an UDP packet. The Filter sends a Payload to every IP in a text file parsed by a newline and checks if the IP responds with a greater size of const declared bytes argument.

- This is intented to scan (rescan) for various DDoS reflectors such as (NTP, DVR, DNS, LDAP)
- Use zmap to scan the world ;)

# How it works
![unknown](https://user-images.githubusercontent.com/65712074/156186925-99709688-05ad-41f0-a06e-57ffbdaea5b1.png)
