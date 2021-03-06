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
- http://127.0.0.1/size?dir=(directory)
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
- /etc/zabbix/zabbix_agent.d/dircount.conf: "UserParameter=dircount.size[*], /usr/local/bin/dircountclient $1"
- Add an item in Zabbix to the concerning host, key: "dirsize.size(directory)" or "dirsize.count(directory)", suggest to set to 15m updates for balanced system load
- Add a graph to the same host (set name and "add" the items size(left) and count (right))
# Experiences so far
Very fast and simple (even with large number of files works perfectly in seconds!) for my own use cases, but happy to learn from your experiences!
# Todo
- Combine into single Zabbix configuration file
- Use MACRO or otherwise to set folder and create Template for usage across systems (instead of individual items per folder currently)
# Future ideas/extensions
I'm currently writing a lot of Go routines to handle/analyze database, mail, security and now more system related information that end up providing information for systems like Zabbix and Grafana. I'm considering to create a single (but extendable) daemon in Go (routines) with an associated library of functions that will be able to handle and store all these systems information processing and analysis instead of introducing many individual daemons (including automatic configuration where possible).
- Mail Log Processor (Dovecot, Postfix) including "jail processor" (very detailed mail reporting and analysis facility; not released to public yet)
- Log display/filtering - non-root accessible (Logdaemon, see https://github.com/moudsen/logdaemon)
- Backup analysis and reporting (Hashbackup, Backblaze buckets)
- Information parser from e-mail messages (i.e. Service invoices)
- Sizing information - non-root accessible (MariaDB, Directories, vMail, Logs, etc)
- Traceroute analysis
- Slack (and others) integration
