# CMDBd
The _**Configuration Management Database Daemon**_ is a lightweight HTTP server that provides a RESTful JSON API for workstations to register and manage information about attached devices. The _**Configuration Management Database Client**_ or [**CMDBc**](https://github.com/jscherff/cmdbc/blob/master/README.md) is the complementary component that collects configuration information for attached devices and reports that information to the server for storage in the database. **CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis. **CMDBd** stores the information in a back-end database.

### System Requirements
**CMDBd** is written in **Go** and can be compiled for any operating system and architecture. This document assumes **CMDBd** will be installed on **Red Hat Enterprise Linux** or **CentOS Release 7** -- or an equivalent operating system that supports **RPM** package management and uses **SystemD** initialization. It requires **MySQL 5.7** or **MariaDB 10.2** or higher for the back-end database.

### Installation
You can build the RPM package with only the RPM spec file, [`cmdbd.spec`](https://github.com/jscherff/cmdbd/blob/master/rpm/cmdbd.spec), using the following commands:
```sh
wget https://raw.githubusercontent.com/jscherff/cmdbd/master/rpm/cmdbd.spec
rpmbuild -bb --clean cmdbd.spec
```
You will need to install the `git`, `golang`, `libusbx`, `libusbx-devel`, and `rpm-build` packages (and their dependencies) in order to perform the build. Once you've built the RPM, you can install it with the RPM command. If you're installing the package for the first time, use the `-i` (install) flag to install the package:
```sh
rpm -i ${HOME}/rpmbuild/RPMS/{arch}/cmdbd-{version}-{release}.{arch}.rpm
```
If you're upgrading the package to a newer version, use the `-U` (upgrade) flag to upgrade the package:
```sh
rpm -i ${HOME}/rpmbuild/RPMS/{arch}/cmdbd-{version}-{release}.{arch}.rpm
```
In the above examples, `{arch}` is the system architecture (e.g. `x86_64`), `{version}` is the package version, (e.g. `1.0.0`), and `{release}` is the package release (e.g. `1.el7.centos`). 

The package will install the following files:
* **`/usr/sbin/cmdbd`** is the **CMDBd** daemon.
* **`/etc/cmdbd/config.json`** is the master configuration file.
* **`/etc/cmdbd/database.json`** contains settings for the database.
* **`/etc/cmdbd/queries.json`** contains SQL queries used by the model.
* **`/etc/cmdbd/syslog.json`** contains settings for the syslog daemon.
* **`/etc/cmdbd/logger.json`** contains settings for the HTTP server logs.
* **`/etc/cmdbd/router.json`** contains settings for the HTTP handlers.
* **`/etc/cmdbd/server.json`** contains settings for the HTTP server.
* **`/etc/cmdbd/metausb.json`** contains information about known USB devices.
* **`/usr/lib/systemd/system/cmdbd.service`** is the SystemD service configuration.
* **`/usr/share/doc/cmdbd-x.y.z/LICENSE`** is the Apache 2.0 license.
* **`/usr/share/doc/cmdbd-x.y.z/README.md`** is this documentation file.
* **`/usr/share/doc/cmdbd-x.y.z/cmdbd.sql`** is the database creation SQL.
* **`/usr/share/doc/cmdbd-x.y.z/users.sql`** is the application user creation SQL.
* **`/usr/share/doc/cmdbd-x.y.z/reset.sql`** resets (truncates) all database tables.
* **`/var/log/cmdbd`** is the directory where CMDBd writes its log files.

Once the package is installed, you must create the database schema, objects, and user account on the target database server using the provided SQL, `cmdbd.sql` and `users.sql`. You must also modify `database.json` configuration file to reflect the correct database hostname, port, user, and password; modify `server.json` to reflect the desired application listener port; and modify other configuration files as necessary and as desired (see below). By default, the config files are owned by the daemon user account and are not _'world-readable'_ as they contain potentially sensitive information. You should not relax the permissions mode of these files.

### Configuration
The JSON configuration files are mostly self-explanatory. The default settings are sane and you should not have to change them in most use cases.

#### Master Config (`config.json`)
The master configuration file contains global parameters and file names of other configuration files.
```json
{
    "SerialFmt": "24F%04X",
    "Configs": {
        "Database": "database.json",
        "Queries": "queries.json",
        "Syslog": "syslog.json",
        "Logger": "logger.json",
        "Router": "router.json",
        "Server": "server.json",
        "MetaUsb": "metausb.json"
    }
}
```
* **`SerialFmt`** is the C _printf-style_ format string for generating serial numbers from a seed integer.
* **`Configs`** is a collection of configurations for different components of the server along with their associated file names. The files **must** be located in the same directory as the master configuration file, above, or the server will fail to start. These components and their configuration settings are covered in more detail  later in this document.
    * **`Database`** names the file that contains database settings.
    * **`Queries`** names the file that contains SQL queries used by the server.
    * **`Syslog`** names the file that contains settings for the syslog daemon.
    * **`Logger`** names the file that contains settings for the HTTP server logs.
    * **`Router`** names the file that contains settings for the HTTP handlers.
    * **`Server`** names the file that contains settings for the HTTP server.
    * **`MetaUsb`** names the file that contains information about known USB devices.

#### Database Settings (`database.json`)
The database configuration file contains parameters required for the server to authenticate and communicate with the database:
```json
{
    "Driver": "mysql",
    "Config": {
        "User": "cmdbd",
        "Passwd": "K2Cvg3NeyR",
        "Net": "",
        "Addr": "localhost",
        "DBName": "gocmdb",
        "Params": null
    }
}
```
* **`Driver`** is the database driver. Only _'mysql'_ is supported.
* **`User`** is the database user the daemon uses to access the database.
* **`Passwd`** is the database user password. Change this to a new, unique value in production.
* **`Net`** is the port on which the database is listening. If blank, the daemon will use the MySQL default port, 3306.
* **`Addr`** is the database hostname or IP address.
* **`DBName`** is the database schema used by the application.
* **`Params`** are additional parameters to pass to the driver (advanced).

#### Query Settings (`queries.json`)
The query configuration file contains SQL queries used by the server to interact with the database. Do not change anything in this file unless directed to do so by a qualified database administrator.

#### Syslog Settings (`syslog.json`)
The syslog configuration file contains parameters for communicating with an optional local or remote syslog server:
```json
{
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

#### Logger Settings (`logger.json`)
The logger configuration file contains parameters that determine log file names and logging behavior:
```json
{
    "LogDir": "/var/log/cmdbd",
    "Stdout": false,
    "Stderr": false,
    "Syslog": false,
    "Logs": {
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
            "Syslog": false
        },
        "error": {
            "LogFile": "error.log",
            "LogFlags": ["date","time","shortfile"],
            "Stdout": false,
            "Stderr": false,
            "Syslog": false
        }
    }
}
```
* **`LogDir`** is the directory where all log files are written.
* **`Stdout`** causes the daemon to write log entries to standard output (console) in addition to other destinations. This overrides the same setting for individual logs, below.
* **`Stderr`** causes the daemon to write log entries to standard error in addition to other destinations. This overrides the same setting for individual logs, below.
* **`Syslog`** causes the daemon to write log entries to a local or remote syslog daemon using the syslog configuration settings, above. This overrides the same setting for individual logs, below.
* **`Logs`** describes each log used by the application. Each log has the following settings:
    * **`LogFile`** is the filename of the log file.
    * **`LogFlags`** is a comma-separated list of attributes to include in the prefix of each log entry. If it is empty, log entries will have no prefix. (This is appropriate for the HTTP access log which is written in the Apache combined log format and already includes relevent attributes in the prefix.) The following [case-sensitive] flags are supported:
        * **`date`** includes date of the event in `YYYY/MM/DD` format.
        * **`time`** includes local time of the event in `HH:MM:SS` 24-hour clock format.
        * **`utc`** includes time in UTC rather than local time.
        * **`standard`** is shorthand for `date` and `time`.
        * **`longfile`** includes the long filename of the source file of the code that generated the event.
        * **`shortfile`** includes the short filename of the source file of the code that generated the event.
    * **`Stdout`** causes the daemon to write log entries to standard output (console) in addition to other destinations.
    * **`Stderr`** causes the daemon to write log entries to standard error in addition to other destinations.
    * **`Syslog`** causes the daemon to write log entries to a local or remote syslog daemon using the syslog configuration settings, above.

#### Router/Handler Settings (`router.json`)
Contains parameters for the HTTP mux router recovery handler:
```json
{
        "RecoveryStack": true
}
```
* **`RecoveryStack`** enables or suppresses writing of the stack track to the error log on panic conditions.

#### Server Settings (`server.json`)
The server configuration file contains parameters for the HTTP server:
```json
{
        "Addr": ":8080",
        "ReadTimeout": 10,
        "WriteTimeout": 10,
        "MaxHeaderBytes": 1048576,
        "HttpBodySizeLimit": 1048576,
        "AllowedContentTypes": ["application/json"]
}
```
* **`Addr`** is the hostname (or IP address) and port of the listener, separated by a colon. If the hostname/address component is blank, the daemon will listen on all network interfaces.
* **`ReadTimeout`** is the maximum duration in seconds for reading the entire HTTP request, including the body.
* **`WriteTimeout`** is the maximum duration in seconds before timing out writes of the response.
* **`MaxHeaderBytes`** is the maximum size in bytes of the request header.
* **`HttpBodySizeLimit`** is the maximum size in bytes of the request body.
* **`AllowedContentTypes`** is a comma-separated list of allowed media types.

#### USB Metadata Settings (`metausb.json`)
The USB metadata configuration file contains vendor names, product names, class descriptions, subclass descriptions, and protocol descriptions for known USB devices. This file is generated from information provided by `http://www.linux-usb.org/usb.ids` and is updated automatically from that site every 30 days. You can force a refresh of this file in two ways:
1. Execute the daemon binary, `cmdbd`, with the `-refresh` flag (preferred). This will download a fresh copy of the metadata, store it in the configuration file, and update relevant metadata tables in the database.
2. Modify the `"Updated":` parameter in the configuration file to a date more than 30 days prior.
```json
    "Updated": "2017-10-17T16:54:09.4910059-07:00"
```
### Startup
The installation package configures the daemon to start automatically when on system startup. On initial package installation, you will have to start the daemon manually because there are post-installation steps required (e.g., configuration and database setup) for the daemon to start successfully. On subssequent package upgrades, the RPM package will shutdown and restart the daemon automatically.

To start the daemon manually, use the `systemctl` utilty with the `start` command:
```sh
systemctl start cmdbd
```
To shut down the daemon, use the `stop` command:
```sh
systemctl stop cmdbd
```
Refer to the `systemctl` man page for other options, such as `restart` and `reload`.

The daemon can also be started from the command line. The following command-line options are available:
* **`-config`** specifies the master configuration file, `config.json`. It is located in  `/etc/cmdbd` by default. All other configuration files will be loaded from the same location.
* **`-console`** causes _all logs_ to be written to standard output; it overrides `Stdout` setting for individual logs.
* **`-refresh`** causes metadata files to be refreshed from source URLs. It overwrites both local configuration files and corresponding database tables.
* **`-version`** displays the server version, `M.m.p-R`, where:
    * **`M`** = MAJOR version with incompatible API changes
    * **`m`** = MINOR version with backwards-compatible new functionality
    * **`p`** = PATCH version with backward-compatible bug fixes.
    * **`R`** = RELEASE number and optional metadata (e.g., `1.beta`)
* **`-help`** displays the above options with a short description.

Starting the daemon manually with console logging (using the `console` _option flag_) is good for troubleshooting. You must start the daemon in the context of the `cmdbd` user account or it will not be able to write to its log files:
```sh
sudo -u cmdbd /usr/sbin/cmdbd -console
```
You can also start the daemon directly as `root`:
```sh
/usr/sbin/cmdbd -console
```
However, doing so can hide permissions-base issues when troubleshooting. (_For security reasons, the daemon should never run as `root` in production; it should always run in the context of a nonprivileged account._) Manual startup example:
```sh
[root@sysadm-dev-01 ~]# sudo -u cmdbd /usr/sbin/cmdbd -help
Usage of /usr/sbin/cmdbd:
  -config <file>
        Master config <file> (default "/etc/cmdbd/conf.json")
  -console
        Enable logging to console
  -refresh
        Refresh application metadata
  -version
        Display application version

[root@sysadm-dev-01 ~]# sudo -u cmdbd /usr/sbin/cmdbd -console
system 2017/10/18 19:43:39 main.go:52: Database version 10.2.9-MariaDB (cmdbd@localhost/gocmdb)
system 2017/10/18 19:43:39 main.go:53: Server version 1.1.0-6.el7.centos started and listening on ":8080"
```
### Logging
Service access, system events, and errors are written to the following log files:
* **`system.log`** records significant, non-error events.
* **`access.log`** records client activity in Apache Combined Log Format.
* **`error.log`** records service and database errors.


### Database Structure
#### Tables
The following tables comprise the database:
* **Device Checkins** contains all device registrations. Multiple check-ins will create multiple records. This provides the ability to track device configuration changes over time. 
* **Serialized Devices** contains devices with serial numbers. It is populated automatically upon device check-in. It uses a unique index based on _Vendor ID_, _Product ID_, and _Serial Number_, and has only one record per serialized device. The first check-in creates the record; subsequent check-ins update modified configuration settings (if any), update the record's _Last Seen_ timestamp, and increment the record's _Checkins_ counter.
* **Unserialized Devices** contains devices without serial numbers. It is populated automatically upon device check-in. It uses a unique index based on _Hostname_, _Vendor ID_, _Product ID_, _Bus Number_, _Bus Address_, and _Port Number_, and strives to have as few records as possible per unserialized device, though this cannot be guaranteed if a device is move to a different workstation or to a different port on the same workstation. The first check-in creates the record; subsequent check-ins update modified configuration settings (if any), update the record's _Last Seen_ timestamp, and increment the record's _Checkins_ counter.
* **Serial Number Requests** contains all requests for a new serial number. **CMDBd** updates new request records with the issued serial number after it is generated. Multiple requests will create multiple records. This provides the ability to detect failures in device serial number configuration and also detect fraudulent usage and abuse.
* **Device Changes** contains configuration changes detected during device audits. Each device configuration attribute change detected during an audit creates one record.

#### Columns
The **Device Checkins**, **Serialized Devices**, **Unserialized Devices**, and **Serial Number Requests** tables have the following columns:

* Hostname
* Vendor ID
* Product ID
* Serial Number
* Vendor Name
* Product Name
* Product Version
* Firmware Version
* Software ID
* Port Number
* Bus Number
* Bus Address
* Buffer Size
* Max Packet Size
* USB Specification
* USB Class
* USB Subclass
* USB Protocol
* Device Speed
* Device Version
* Factory Serial Number
* Object Type
* Object JSON
* Remote Address

The **Device Checkins** table includes the following additional column:
* Checkin Date

The **Serial Number Requests** table includes the following additional column:
* Request Date

The **Serialized Devices** and **Unserialized Devices** tables both include the following additional columns:
* First Seen
* Last Seen
* Checkins

The **Device Changes** table has the following columns:
* Host Name
* Vendor ID
* Product ID
* Serial Number
* Changes
* Audit Date

For a given **Device Changes** record, the _Changes_ column contains a JSON object that represents a collection of one or more changes. Each change element in the collection has the following fields:
* Property Name
* Old Value
* New Value

### API Endpoints
| Endpoint | Method | Purpose
| :------ | :------ | :------ |
| /v1/usbci/checkin/`host`/`vid`/`pid` | `POST` | Submit configuration information for a new device or update information for an existing device. |
| /v1/usbci/checkout/`host`/`vid`/`pid`/`sn` | `GET` | Obtain configuration information for a previously-registered, serialized device in order to perform a change audit. |
| /v1/usbci/newsn/`host`/`vid`/`pid` | `POST` | Obtain a new unique serial number from the server for assignment to the attached device. |
| /v1/usbci/audit/`host`/`vid`/`pid`/`sn` | `POST` | Submit the results of a change audit on a serialized device. Results include the attribute name, previous value, and new value for each modified attribute. |
| /v1/usbmeta/vendor/`vid` | `GET` | Obtain the USB vendor name given the vendor ID |
| /v1/usbmeta/product/`vid`/`pid` | `GET` | Obtain the USB vendor and product names given the vendor and product IDs | 
| /v1/usbmeta/class/`cid` | `GET` | Obtain the USB class description given the class ID | 
| /v1/usbmeta/subclass/`cid`/`sid` | `GET` | Obtain the USB class and subclass descriptions given the class and subclass IDs |
| /v1/usbmeta/protocol/`cid`/`sid`/`pid` | `GET` | Obtain the USB class, subclass, and protocol descriptions given the class, subclass, and protocol IDs |

### API Parameters
* **`host`** is the _hostname_ of the workstation to which the device is attached.
* **`vid`** is the _vendor ID_ of the device.
* **`pid`** is the _product ID_ of the device.
* **`sn`** is the _serial number_ of the device.
* **`vid`** is the _vendor ID_ of the device.
* **`pid`** is the _product ID_ of the device.
* **`cid`** is the _class ID_ of the device.
* **`sid`** is the _subclass ID_ of the device.
* **`pid`** is the _protocol ID_ of the device.
