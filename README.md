# Linux Daemon "sized"
Golang Linux daemon to serve directory size requests as input for example for Zabbix Agent
# Origin
Zabbix has the ability to determine directory sizes with vfs.dir.size but requires Linux rights to be able to at least read monitored directories. This results in a complex setup in case of directories owned with other Linux users (like vmail, mysql, ...).
Although a simple "du" could help, this would require sudo rights (and clutters my sudo with Zabbix related messages) ...
Hence I have decided to write a simple "sized" deamon that can help a GET request with a 'dir' parameter and simply returns the size of the requested directory.
# Installation
Under construction, basically Go default installation and 1 additional library for the 'deamon' functions suffices.
# Configuration
Under construction
- Daemon Logging consideration ("/var/log/sized.log")
- Zabbix Agent configuration
# Experiences so far
Very fast and simple (even with large number of files works perfectly in seconds!) for my own use cases, but happy to learn from your experiences!
