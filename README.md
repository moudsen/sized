# Linux Daemon "sized"
Golang Linux daemon to serve directory size requests as input for example for Zabbix Agent
# Origin
Zabbix has the ability to determine directory sizes with vfs.dir.size but requires Linux rights to be able to at least read monitored directories. This results in a complex setup in case of directories owned with other Linux users (like vmail, mysql, ...).
Although a simple "du" could help, this would require sudo rights (and clutters my sudo with Zabbix related messages) ...
Hence I have decided to write a simple "sized" deamon that can help a GET request with a 'dir' parameter and simply returns the size of the requested directory.
# sized.go (daemon)
The sized daemon depends on "github.com/takama/daemon". Add to Go library before compiling sized.go.
Once compiled copy "sized" to /usr/local/sbin and install/use the service (note: must be root or able to use sudo):
- Install with "sized install"
- Start with "sized start"
- Stop with "sized stop"
- Status with "sized status"
- Uninstall with "sized uninstall"
The daemon listens on port 7007 for http GET requests:
- http://127.0.0.1/size?dir=<dir>
# Command line client dirsizeclient.go
Once compiled, copy to /usr/local/bin and use:
  - "dirsizeclient <dir>"
Obviously the daemon must be running for this command to work.
# Configuration
_Daemon Logging consideration ("/var/log/sized.log")_
- The sized daemon logs daemon operations and sizing requests to syslog (including timing). In case of rsyslogd you can log the messages to a specific file
- /etc/rsyslog.d/sized.conf: ":programname, isequal, "sized" /var/log/sized.log"
_Zabbix Agent configuration_
- /etc/zabbix/zabbix_agent.d/dirsize.conf: "UserParameter=dirsize.size[*], /usr/local/bin/dirsizeclient $1"
- Add an item in Zabbix to the concerning host, key: "dirsize.size(<directory>)", suggest to set to 15m updates
- Add a graph to the same host (set name and "add" the item)
# Experiences so far
Very fast and simple (even with large number of files works perfectly in seconds!) for my own use cases, but happy to learn from your experiences!
