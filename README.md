# CMDBd
The _Configuration Management Database Daemon_ is a lightweight HTTP server that provides a RESTful JSON API for workstations to register and manage information about attached devices. The [_Configuration Management Database Client_](https://github.com/jscherff/cmdbc/blob/master/README.md) or **CMDBc** is the complementary component that collects configuration information for attached devices and reports that information to the server for storage in the database. **CMDBc** can register attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis. **CMDBd** stores the information in a back-end database. It requires **MySQL 5.7** or **MariaDB 10.2** or higher.

### System Requirements
**CMDBd** is written in **Go** and can be compiled for any operating system and architecture. This document assumes **CMDBd** will be installed on **Red Hat Enterprise Linux** or **CentOS** release 7 or equivalent operating system that supports the RPM package management and the SystemD init system.

### Installation
The RPM package can be built using only the RPM spec file, [**`cmdbd.spec`**](https://github.com/jscherff/cmdbd/blob/master/rpm/cmdbd.spec), using the following commands:
```sh
wget https://raw.githubusercontent.com/jscherff/cmdbd/master/rpm/cmdbd.spec
rpmbuild -bb --clean cmdbd.spec
```
The package will install the following files:
* **`/usr/sbin/cmdbd`** is the CMDBd daemon.
* **`/etc/cmdbd/config.json`** is the CMDBd configuration file.
* **`/usr/lib/systemd/system/cmdbd.service`** is the SystemD service configuration.
* **`/usr/share/doc/cmdbd-x.y.z/LICENSE`** is the Apache 2.0 license.
* **`/usr/share/doc/cmdbd-x.y.z/README.md`** is this documentation file.
* **`/usr/share/doc/cmdbd-x.y.z/cmdbd.sql`** is the database creation SQL.
* **`/usr/share/doc/cmdbd-x.y.z/users.sql`** is the application user creation SQL.
* **`/var/log/cmdbd`** is the directory where CMDBd writes its log files.

Once the package is installed, the database schema, objects, and user account must be created on the target database server using the provided SQL, `cmdb.sql` and `users.sql`, and the `config.json` file (see below) must be configured with the correct database server hostname and port, database user and password, application listener port, and other preferences. Once these tasks are complete, the daemon can be started with the following command:
```sh
systemctl start cmdbd
```
Service access, system events, and errors are written to the following log files:
* **`system.log`** records significant, non-error events.
* **`access.log`** records client activity in Apache Combined Log Format.
* **`error.log`** records service and database errors.

### Configuration
The JSON configuration file, [**`config.json`**](https://github.com/jscherff/cmdbd/blob/master/config.json), is mostly self-explanatory. The default settings are sane and should not have to be changed for most use cases.

**Server Settings**
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

**Database Settings**
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

**Logger Settings**
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
  * **`date`** -- date of the event in `YYYY/MM/DD` format.
  * **`time`** -- local time of the event in `HH:MM:SS` 24-hour clock format.
  * **`utc`** -- time in UTC rather than local time.
  * **`standard`** -- shorthand for `date` and `time`.
  * **`longfile`** -- long filename of the source code file that generated the event.
  * **`shortfile`** -- short filename of the source code file that generated the event.
* **`Stdout`** causes the daemon to write log entries to standard output (console) in addition to other destinations.
* **`Stderr`** causes the daemon to write log entries to standard error in addition to other destinations.
* **`Syslog`** causes the daemon to write log entries to a local or remote syslog daemon using the `Syslog` configuration settings, below.

**Syslog Settings**
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
* **`Facility`** specifies the type of program that is logging the message:
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
* **`Severity`** specifies the severity of the event:
  * **`LOG_EMERG`** -- system is unusable
  * **`LOG_ALERT`** -- action must be taken immediately
  * **`LOG_CRIT`** -- critical conditions
  * **`LOG_ERR`** -- error conditions
  * **`LOG_WARNING`** -- warning conditions
  * **`LOG_NOTICE`** -- normal but significant conditions
  * **`LOG_INFO`** -- informational messages
  * **`LOG_DEBUG`** -- debug-level messages

**Log Directory Settings**
```json
"LogDir": {
    "Windows": "log",
    "Linux": "/var/log/cmdbd"
}
```
* **`Windows`** is the log directory to use for Windows installations.
* **`Linux`** is the log directory to use for Linux installations.

**Options**
```json
"Options": {
    "Stdout": false,
    "Stderr": false,
    "Syslog": false,
    "RecoveryStack": false
}
```
* **`Stdout`** causes all logs to be written to standard output in addition to other destinations.
* **`Stderr`** causes all logs to be written to standard error in addition to other destinations.
* **`Syslog`** causes all logs to be written to the configured syslog daemon in addition to other destinations.
* **`RecoveryStack`** enables or suppresses writing of the stack track to the error log on panic conditions.

### API Endpoints
| Endpoint | Method | Purpose
| :------ | :------ | :------ |
| `/usbci/checkin/{host}/{vid}/{pid}` | POST | Submit configuration information for a new device or update information for an existing device. |
| `/usbci/checkout/{host}/{vid}/{pid}/{sn}` | GET | Obtain configuration information for a previously-registered, serialized device in order to perform a change audit. |
| `/usbci/audit/{host}/{vid}/{pid}/{sn}` | POST | Submit the results of a change audit on a serialized device. Results include the attribute name, previous value, and new value for each modified attribute.
| `/usbci/newsn/{host}/{vid}/{pid}` | POST | Obtain a new unique serial number from the server for assignment to the attached device. |

### API Parameters
* **`host`** is the _hostname_ of the workstation to which the device is attached.
* **`vid`** is the _vendor ID_ of the device.
* **`pid`** is the _product ID_ of the device.
* **`sn`** is the _serial number_ of the device.
