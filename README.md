# DNS C2 Server

This is a rough try of a c2 dns server setup, including an agent as well. 

Simple example of a command output `cat ./main.go` transmitted through TXT RRs. The agent is on the right, and it is replying with FQDNs with the data encoded with base64.
![poc](/assets/images/poc_scr.jpg)
