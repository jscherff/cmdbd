# CMDBd
The _**Configuration Management Database Daemon**_ is a lightweight HTTP server that provides a RESTful JSON API for workstations to register and manage information about attached devices. The _**Configuration Management Database Client**_ or [**CMDBc**](https://github.com/jscherff/cmdbc/blob/master/README.md) is the complementary component that collects configuration information for attached devices and reports that information to the server for storage in the database. **CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis. **CMDBd** stores the information in a back-end database.

### System Requirements
**CMDBd** is written in **Go** and can be compiled for any operating system and architecture. This document assumes **CMDBd** will be installed on **Red Hat Enterprise Linux** or **CentOS** release 7 -- or an equivalent operating system that supports **RPM** package management and the **SystemD** initialization system. It requires **MySQL 5.7** or **MariaDB 10.2** or higher for the back-end database.

### Installation
You can build the RPM package with only the RPM spec file, [`cmdbd.spec`](https://github.com/jscherff/cmdbd/blob/master/rpm/cmdbd.spec), using the following commands:
```sh
wget https://raw.githubusercontent.com/jscherff/cmdbd/master/rpm/cmdbd.spec
rpmbuild -bb --clean cmdbd.spec
```
You will need to install the `git`, `golang`, and `rpm-build` packages and their dependencies in order to perform the build. Once you've built the RPM, you can install it with this command:
```sh
rpm -i ${HOME}/rpmbuild/RPMS/{arch}/cmdbd-{version}-{release}.{arch}.rpm
```
Where `{arch}` is your system architecture (e.g. `x86_64`), `{version}` is the package version, (e.g. `1.0.0`), and `{release}` is the package release (e.g. `1.el7.centos`). The package will install the following files:
* **`/usr/sbin/cmdbd`** is the CMDBd daemon.
* **`/etc/cmdbd/config.json`** is the CMDBd configuration file.
* **`/usr/lib/systemd/system/cmdbd.service`** is the SystemD service configuration.
* **`/usr/share/doc/cmdbd-x.y.z/LICENSE`** is the Apache 2.0 license.
* **`/usr/share/doc/cmdbd-x.y.z/README.md`** is this documentation file.
* **`/usr/share/doc/cmdbd-x.y.z/cmdbd.sql`** is the database creation SQL.
* **`/usr/share/doc/cmdbd-x.y.z/users.sql`** is the application user creation SQL.
* **`/var/log/cmdbd`** is the directory where CMDBd writes its log files.

Once the package is installed, you must create the database schema, objects, and user account on the target database server using the provided SQL, `cmdb.sql` and `users.sql`. You must also modify `config.json` configuration file to reflect the correct database server hostname and port, database user and password, application listener port, and other preferences (see below). By default, the `config.json` file is owned by the daemon user account and is not world-readable. You should not relax the permissions mode of this file as it contains the database password.

### Configuration
The JSON configuration file, [`config.json`](https://github.com/jscherff/cmdbd/blob/master/config.json), is mostly self-explanatory. The default settings are sane and you should not have to change them in most use cases.

##### Server Settings
Parameters that affect the behavior of the HTTP server.
```json
"Server": {
    "Addr": ":8080",
    "ReadTimeout": 10,
    "WriteTimeout": 10,
    "MaxHeaderBytes": 1048576,
    "HttpBodySizeLimit": 1048576,
    "AllowedContentTypes": ["application/json"]
}
```
* **`Addr`** is the hostname or IP address and port of the listener, separated by a colon. If blank, the daemon will listen on all network interfaces.
* **`ReadTimeout`** is the maximum duration in seconds for reading the entire HTTP request, including the body.
* **`WriteTimeout`** is the maximum duration in seconds before timing out writes of the response.
* **`MaxHeaderBytes`** is the maximum size in bytes of the request header.
* **`HttpBodySizeLimit`** is the maximum size in bytes of the request body.
* **`AllowedContentTypes`** is a comma-separated list of allowed media types.

##### Database Settings
Parameters for communicating with the database server.
```json
"Database": {
    "Driver": "mysql",
    "Config": {
        "User": "cmdbd",
        "Passwd": "K2Cvg3NeyR",
        "Net": "",
        "Addr": "localhost",
        "DBName": "gocmdb",
        "Params": null
    },
    ...
}
```
* **`Driver`** is the database driver. Only `mysql` is supported.
* **`User`** is the database user the daemon uses to access the database.
* **`Passwd`** is the database user password. The default, shown, should be changed in production.
* **`Net`** is the port on which the database is listening. If blank, the daemon will use the MySQL default port, 3306.
* **`Addr`** is the database hostname or IP address.
* **`DBName`** is the database schema used by the application.
* **`Params`** are additional parameters to pass to the driver (advanced).

##### Logger Settings
Parameters that determine log file names and logging behavior.
```json
"Loggers": {
    "system": {
        "LogFile": "system.log",
        "LogFlags": ["date","time","shortfile"],
        "Stdout": false,
        "Stderr": false,
        "Syslog": false
    },
    "access": {
        "LogFile": "access.log",
        "LogFlags": [],
        "Stdout": false,
        "Stderr": false,
        "Syslog": true
    },
    "error": {
        "LogFile": "error.log",
        "LogFlags": ["date","time","shortfile"],
        "Stdout": false,
        "Stderr": false,
        "Syslog": false
    }
}
```
* **`LogFile`** is the filename of the log file.
* **`LogFlags`** specifies information to include in the prefix of each log entry. The following [case-sensitive] flags are supported:
    * **`date`** includes date of the event in `YYYY/MM/DD` format.
    * **`time`** includes local time of the event in `HH:MM:SS` 24-hour clock format.
    * **`utc`** includes time in UTC rather than local time.
    * **`standard`** is shorthand for `date` and `time`.
    * **`longfile`** includes the long filename of the source file of the code that generated the event.
    * **`shortfile`** includes the short filename of the source file of the code that generated the event.
* **`Stdout`** causes the daemon to write log entries to standard output (console) in addition to other destinations.
* **`Stderr`** causes the daemon to write log entries to standard error in addition to other destinations.
* **`Syslog`** causes the daemon to write log entries to a local or remote syslog daemon using the `Syslog` configuration settings, below.

##### Syslog Settings
Parameters for communicating with a local or remote syslog server.
```json
"Syslog": {
    "Protocol": "tcp",
    "Port": "1468",
    "Host": "localhost",
    "Tag": "cmdbd",
    "Facility": "LOG_LOCAL7",
    "Severity": "LOG_INFO"
}
```
* **`Protocol`** is the transport-layer protocol used by the syslog daemon (blank for local).
* **`Port`** is the port used by the syslog daemon (blank for local).
* **`Host`** is the hostname or IP address of the syslog daemon (blank for local).
* **`Tag`** is an arbitrary string to add to the event.
* **`Facility`** specifies the type of program that is logging the message (see [RFC 5424](https://tools.ietf.org/html/rfc5424)):
    * **`LOG_KERN`** -- kernel messages
    * **`LOG_USER`** -- user-level messages
    * **`LOG_MAIL`** -- mail system
    * **`LOG_DAEMON`** -- system daemons
    * **`LOG_AUTH`** -- security/authorization messages
    * **`LOG_SYSLOG`** -- messages generated internally by syslogd
    * **`LOG_LPR`** -- line printer subsystem
    * **`LOG_NEWS`** -- network news subsystem
    * **`LOG_UUCP`** -- UUCP subsystem
    * **`LOG_CRON`** -- security/authorization messages
    * **`LOG_AUTHPRIV`** -- FTP daemon
    * **`LOG_FTP`** -- scheduling daemon
    * **`LOG_LOCAL0`** -- local use 0
    * **`LOG_LOCAL1`** -- local use 1
    * **`LOG_LOCAL2`** -- local use 2
    * **`LOG_LOCAL3`** -- local use 3
    * **`LOG_LOCAL4`** -- local use 4
    * **`LOG_LOCAL5`** -- local use 5
    * **`LOG_LOCAL6`** -- local use 6
    * **`LOG_LOCAL7`** -- local use 7
* **`Severity`** specifies the severity of the event (see [RFC 5424](https://tools.ietf.org/html/rfc5424)):
    * **`LOG_EMERG`** -- system is unusable
    * **`LOG_ALERT`** -- action must be taken immediately
    * **`LOG_CRIT`** -- critical conditions
    * **`LOG_ERR`** -- error conditions
    * **`LOG_WARNING`** -- warning conditions
    * **`LOG_NOTICE`** -- normal but significant conditions
    * **`LOG_INFO`** -- informational messages
    * **`LOG_DEBUG`** -- debug-level messages

##### Log Directory Settings
Directory where log files are written.
```json
"LogDir": {
    "Windows": "log",
    "Linux": "/var/log/cmdbd"
}
```
* **`Windows`** is the log directory to use for Windows installations.
* **`Linux`** is the log directory to use for Linux installations.

##### Global Options
System-wide parameters.
```json
"Options": {
    "Stdout": false,
    "Stderr": false,
    "Syslog": false,
    "RecoveryStack": false
}
```
* **`Stdout`** causes _all logs_ to be written to standard output; it overrides `Stdout` setting for individual logs.
* **`Stderr`** causes all logs to be written to standard error; it overrides `Stderr` setting for individual logs.
* **`Syslog`** causes all logs to be written to the configured syslog daemon; it overrides `Syslog` settings for individual logs.
* **`RecoveryStack`** enables or suppresses writing of the stack track to the error log on panic conditions.

### Startup
Once all configuration tasks are complete, the daemon can be started with the following command:
```sh
systemctl start cmdbd
```
Service access, system events, and errors are written to the following log files:
* **`system.log`** records significant, non-error events.
* **`access.log`** records client activity in Apache Combined Log Format.
* **`error.log`** records service and database errors.

The daemon can also be started from the command line. The following command-line options are available:
* **`-config`** specifies an alternate JSON configuration file; the default is `/etc/cmdbd/config.json`.
* **`-stdout`** causes _all logs_ to be written to standard output; it overrides `Stdout` setting for individual logs.
* **`-stderr`** causes _all logs_ to be written to standard error; it overrides `Stderr` setting for individual logs.
* **`-syslog`** causes _all logs_ to be written to the configured syslog daemon; it overrides `Syslog` setting for individual logs.
* **`-help`** displays the above options with a short description.

Starting the daemon manually with console logging (`-stdout` option) is good for troubleshooting. Use the `su` command as `root` (or `sudo su` as a non-privileged user) to start the daemon with the `cmdbd` account or it will not be able to write to its log files:
```sh
su - cmdbd -c '/usr/sbin/cmdbd -stdout'
```
You can also start the daemon manually as `root`, but this can hide permissions-base issues when troubleshooting. (_For security reasons, the daemon should never run as `root` in production; it should always run in the context of a nonprivileged account._) Manual startup example:
```sh
[root@sysadm-dev-01 ~]# su - cmdbd  -c '/usr/sbin/cmdbd -help'
Usage of /usr/sbin/cmdbd:
  -config file
        Web server configuration file (default "/etc/cmdbd/config.json")
  -stderr
        Enable logging to stderr
  -stdout
        Enable logging to stdout
  -syslog
        Enable logging to syslog

[root@sysadm-dev-01 ~]# su - cmdbd  -c '/usr/sbin/cmdbd -stdout'
system 2017/09/30 09:55:38 main.go:62: Database "10.2.9-MariaDB" (cmdbd@localhost/gocmdb) using "mysql" driver
system 2017/09/30 09:55:38 main.go:63: Server started and listening on ":8080"
```

### API Endpoints
| Endpoint | Method | Purpose
| :------ | :------ | :------ |
| **`/usbci/checkin/{host}/{vid}/{pid}`** | POST | Submit configuration information for a new device or update information for an existing device. |
| **`/usbci/checkout/{host}/{vid}/{pid}/{sn}`** | GET | Obtain configuration information for a previously-registered, serialized device in order to perform a change audit. |
| **`/usbci/audit/{host}/{vid}/{pid}/{sn}`** | POST | Submit the results of a change audit on a serialized device. Results include the attribute name, previous value, and new value for each modified attribute.
| **`/usbci/newsn/{host}/{vid}/{pid}`** | POST | Obtain a new unique serial number from the server for assignment to the attached device. |

### API Parameters
* **`host`** is the _hostname_ of the workstation to which the device is attached.
* **`vid`** is the _vendor ID_ of the device.
* **`pid`** is the _product ID_ of the device.
* **`sn`** is the _serial number_ of the device.
