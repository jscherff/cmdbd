# CMDBd
The _Configuration Management Database Daemon_ is a lightweight HTTP server that provides a RESTful JSON API for workstations to register and manage information about attached devices. The [_Configuration Management Database Client_](https://github.com/jscherff/cmdbc/blob/master/README.md) or **CMDBc** is the complementary component that collects configuration information for attached devices and reports that information to the server for storage in the database. **CMDBc** can register attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis. **CMDBd** stores the information in a back-end database. It requires **MySQL 5.7** or **MariaDB 10.2** or higher.

### System Requirements
**CMDBd** is written in **Go** and can be compiled for any operating system and architecture. This document assumes **CMDBd** will be installed on **Red Hat Enterprise Linux** or **CentOS** release 7 or equivalent operating system that supports the RPM package management and the SystemD init system.

### Installation
The RPM package can be built using only the RPM spec file, [`cmdbd.spec`](https://raw.githubusercontent.com/jscherff/cmdbd/master/rpm/cmdbd.spec), using the following commands:
```sh
wget https://raw.githubusercontent.com/jscherff/cmdbd/master/rpm/cmdbd.spec
rpmbuild -bb --clean cmdbd.spec
```
The package will install the following files:
* `/usr/sbin/cmdbd` -- The CMDBd daemon
* `/etc/cmdbd/config.json` -- The CMDBd configuration file
* `/usr/lib/systemd/system/cmdbd.service` -- The SystemD service configuration
* `/usr/share/doc/cmdbd-1.0.1/LICENSE` -- The Apache 2.0 License
* `/usr/share/doc/cmdbd-1.0.1/README.md` -- This documentation file
* `/usr/share/doc/cmdbd-1.0.1/cmdbd.sql` -- The database creation SQL
* `/usr/share/doc/cmdbd-1.0.1/users.sql` -- The application user creation SQL
* `/var/log/cmdbd` -- The directory where CMDBd writes its log files.

Once the package is installed, the database schema, objects, and user account must be created on the target database server using the provided SQL, `cmdb.sql` and `users.sql`, and the `config.json` file (see below) must be configured with the correct database server hostname and port, database user and password, application listener port, and other preferences. Once these tasks are complete, the daemon can be started with the following command:
```sh
systemctl start cmdbd
```
Service access, system events, and errors are written to the following log files:
* `system.log` -- Significant, non-error events
* `access.log` -- Client activity in Apache Combined Log Format
* `error.log` -- Service and database errors

### Configuration
The JSON configuration file, `config.json`, is mostly self-explanatory. The default settings are sane and should not have to be changed for most use cases.

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
* `Addr` is the hostname or IP address and port of the listener, separated by a colon. If blank, the daemon will listen on all network interfaces.
* `ReadTimeout` is the maximum duration in seconds for reading the entire HTTP request, including the body.
* `WriteTimeout` is the maximum duration in seconds before timing out writes of the response.
* `MaxHeaderBytes` is the maximum size in bytes of the request header.
* `HttpBodySizeLimit` is the maximum size in bytes of the request body.
* `AllowedContentTypes` is a comma-separated list of allowed media types.

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
* `Driver` is the database driver. Only `mysql` is supported.
* `User` is the database user the daemon uses to access the database.
* `Passwd` is the database user password. The default, shown, should be changed in production.
* `Net` is the port on which the database is listening. If blank, the daemon will use the MySQL default port, 3306.
* `Addr` is the database hostname or IP address.
* `DBName` is the database schema used by the application.
* `Params` are additional parameters to pass to the driver (advanced).

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
* `LogFile` is the filename of the log file.
* `LogFlags` specify information to include in the prefix of each log entry. The following [case-sensitive] flags are supported:
  * `date` -- date of the event in `YYYY/MM/DD` format.
  * `time` -- local time of the event in `HH:MM:SS` 24-hour clock format.
  * `longfile` -- long filename of the source code file that generated the event.
  * `shortfile` -- short filename of the source code file that generated the event.
  * `standard` -- shorthand for `date` and `time`.
  * `UTC` -- use UTC rather than local time.
* `Stdout` causes the daemon to write log entries to standard output (console) in addition to other destinations.
* `Stderr` causes the daemon to write log entries to standard error in addition to other destinations.
* `Syslog` causes the daemon to write log entries to a local or remote syslog daemon using the `Syslog` configuration settings, below.

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
* `Protocol` --
* `Port` --
* `Host` --
* `Tag` --
* `Facility` --
* `Severity` --
```json
"LogDir": {
    "Windows": "log",
    "Linux": "/var/log/cmdbd"
}
```
* `Windows` --
* `Linux` --
```json
"Options": {
    "Stdout": false,
    "Stderr": false,
    "Syslog": false,
    "RecoveryStack": false
}
```

### API Endpoints
| Endpoint | Method | Purpose
| :------ | :------ | :------ |
| `/usbci/checkin/{host}/{vid}/{pid}` | POST | Submit configuration information for a new device or update information for an existing device. |
| `/usbci/checkout/{host}/{vid}/{pid}/{sn}` | GET | Obtain configuration information for a previously-registered, serialized device in order to perform a change audit. |
| `/usbci/audit/{host}/{vid}/{pid}/{sn}` | POST | Submit the results of a change audit on a serialized device. Results include the attribute name, previous value, and new value for each modified attribute.
| `/usbci/newsn/{host}/{vid}/{pid}` | POST | Obtain a new unique serial number from the server for assignment to the attached device. |

### API Parameters
* `host` -- The _hostname_ of the workstation to which the device is attached.
* `vid` -- The _vendor ID_ of the device.
* `pid` -- The _product ID_ of the device.
* `sn` -- The _serial number_ of the device.
