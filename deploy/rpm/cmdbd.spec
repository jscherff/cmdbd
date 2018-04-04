# =============================================================================
%define		name	cmdbd
%define		version	3.4.0
%define		release	1
%define		branch  master
%define		gecos	CMDBd Service
%define		summary	Configuration Management Database Daemon
%define		author	John Scherff <jscherff@24hourfit.com>
%define		package	github.com/jscherff/%{name}
%define		gopath	%{_builddir}/go
%define		docdir	%{_docdir}/%{name}-%{version}
%define		logdir	%{_var}/log/%{name}
%define		syslib	%{_prefix}/lib/systemd/system
%define		lmtdir  %{_sysconfdir}/security/limits.d
%define		confdir %{_sysconfdir}/%{name}
# =============================================================================

Name:		%{name}
Version:	%{version}
Release:	%{release}%{?dist}
Summary:	%{summary}

Group:		Applications/System
License:	ASL 2.0
URL:		https://www.24hourfitness.com
Vendor:		24 Hour Fitness, Inc.
Prefix:		%{_sbindir}
Packager: 	%{packager}
BuildRoot:	%{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
Distribution:	el

BuildRequires:    golang >= 1.8.0
Requires(pre):    %{_sbindir}/useradd, %{_bindir}/getent
Requires(postun): %{_sbindir}/userdel

%description
The Configuration Management Database Daemon, %{name}, is a lightweight HTTP
daemon that provides a REST API for clients installed on Windows endpoints.
The clients collect information about attached devices and send it to the
server for storage in the database. Clients can register attached devices
with the server, obtain unique serial numbers from the server for devices
that support serial number configuration, perform audits against previous
device configurations, and report any configuration changes found during
the audit to the server for later analysis.

%prep

%build

  export GOPATH=%{gopath}
  export GIT_DIR=%{gopath}/src/%{package}/.git

  go get %{package}
  git checkout %{branch}

  go build -ldflags='-X main.version=%{version}-%{release}' %{package}
  go build -ldflags='-X main.version=%{version}-%{release}' %{package}/bcrypt

%install

  test %{buildroot} != / && rm -rf %{buildroot}/*

  mkdir -p %{buildroot}{%{_sbindir},%{_bindir}}
  mkdir -p %{buildroot}{%{confdir},%{syslib},%{logdir},%{lmtdir},%{docdir}}

  install -s -m 755 %{_builddir}/%{name} %{buildroot}%{_sbindir}/
  install -s -m 755 %{_builddir}/bcrypt %{buildroot}%{_bindir}/
  install -m 640 %{gopath}/src/%{package}/deploy/ddl/%{name}.sql %{buildroot}%{docdir}/
  install -m 644 %{gopath}/src/%{package}/deploy/os%{syslib}/* %{buildroot}%{syslib}/
  install -m 644 %{gopath}/src/%{package}/deploy/os%{lmtdir}/* %{buildroot}%{lmtdir}/
  install -m 644 %{gopath}/src/%{package}/{LICENSE,*.md} %{buildroot}%{docdir}/

  install -m 640 %{gopath}/src/%{package}/config/*.json %{buildroot}%{confdir}/
  install -m 640 %{gopath}/src/%{package}/config/model/*.json %{buildroot}%{confdir}/model/
  install -m 640 %{gopath}/src/%{package}/config/server/*.json %{buildroot}%{confdir}/server/
  install -m 640 %{gopath}/src/%{package}/config/service/*.json %{buildroot}%{confdir}/service/
  install -m 640 %{gopath}/src/%{package}/config/service/*.pem %{buildroot}%{confdir}/service/
  install -m 640 %{gopath}/src/%{package}/config/store/*.json %{buildroot}%{confdir}/store/

%clean

  test %{buildroot} != / && rm -rf %{buildroot}/*
  test %{_builddir} != / && rm -rf %{_builddir}/*

%files

  %defattr(-,root,root)
  %license %{docdir}/LICENSE
  %{_sbindir}/*
  %{_bindir}/*
  %{syslib}/*
  %{lmtdir}/*
  %{docdir}/*

  %defattr(640,%{name},%{name},750)
  %config %{confdir}/*

  %defattr(644,%{name},%{name},755)
  %{logdir}

%pre

  # Tasks to perform FROM NEW RPM before install (1) or upgrade (2)

  case ${1} in

    1)
      %{_sbindir}/useradd -Mrd %{_sbindir} -c '%{gecos}' -s /sbin/nologin %{name}
      ;;

    2)
      systemctl --quiet is-active %{name} && systemctl --quiet stop %{name}
      systemctl --quiet is-enabled %{name} && systemctl --quiet disable %{name}
      ;;

  esac

  : Force zero return code

%post

  # Tasks to perform FROM NEW RPM after install (1) or upgrade (2)

  case ${1} in

    1)
      systemctl --quiet is-enabled %{name} || systemctl --quiet enable %{name} 
      ;;

    2)
      systemctl --quiet is-enabled %{name} || systemctl --quiet enable %{name} 
      systemctl --quiet is-active %{name} || systemctl --quiet start %{name}
      ;;

  esac

  : Force zero return code

%preun

  # Tasks to perform FROM OLD RPM before uninstall (0) or upgrade (1)

  case ${1} in

    0)
      systemctl --quiet is-active %{name} && systemctl --quiet stop %{name}
      systemctl --quiet is-enabled %{name} && systemctl --quiet disable %{name}
      ;;

    1)
      ;;

  esac

  : Force zero return code

%postun

  # Tasks to perform FROM OLD RPM after uninstall (0) or upgrade (1)

  case ${1} in

    0)
      %{_sbindir}/userdel %{name}
      test %{logdir} != / && rm -rf %{logdir}
      ;;

    1)
      ;;

  esac

  : Force zero return code

%changelog
* Wed Apr 4 2018 - jscherff@24hourfit.com
- Refactored server configuration source code
- Refactored request router and handler logic
- Refactored database connection pool configuration
- Moved middleware chaining to server configuration source code
- Added comprehensive logging during server configuration
* Tue Apr 3 2018 - jscherff@24hourfit.com
- Created MaxConnectionsHandler middleware to manage server concurrency
- Created a MaxConnections property for the Server object with a default of 50
- Added a concurrency testing endpoint and handler for testing concurrency
* Mon Apr 2 2018 - jscherff@24hourfit.com
- Set default MaxOpenConns to 50 to limit concurrent open database connections
- Added /etc/security/limits.d/cmdbd.conf with hard and soft limits
* Mon Mar 12 2018 - jscherff@24hourfit.com
- Added 'SET FOREIGN_KEY_CHECKS = 0' for table truncates in reset.sql
- Changed table name cmdb_errors to cmdb_events in reset.sql
* Thu Mar 1 2018 - jscherff@24hourfit.com
- Modified database insert procs to include all non-key columns in updates
- Added api.AppendRequest() to decorate error messages with request info
* Wed Feb 14 2018 - jscherff@24hourfit.com
- Set AUTO_INCREMENT=1 in DDL and changed serial number format to 24Fxxxx
* Mon Feb 12 2018 - jscherff@24hourfit.com
- Renamed 'Endpoint(s)' to 'Route(s)' and renamed associated files
* Fri Jan 26 2018 - jscherff@24hourfit.com
- Separated DML and DDL
- Modified RPM spec file to enhance GO build process
* Wed Jan 17 2018 - jscherff@24hourfit.com
- Comprehensive refactor to make code resusable and easier to maintain
- Converted model to lightweight ORM using sqlx
- Segregated server components into 'server' package
- Segregated common services into 'service' package
- Segregated store components into 'store' package
- Segregated API components into 'api' package
- Created separate v1 and v2 APIs for backward compatibility
* Mon Nov 13 2017 - jscherff@24hourfit.com
- Modified queries to use tables directly versus views
- Added DATETIME columns to inserts with time.Now() as value
- Modified Loc (location) database parameter to 'Local'
- Removed unnecessary views from DDL
* Wed Nov 8 2017 - jscherff@24hourfit.com
- Added cmdb_users table for authentication
- Added authentication API to support basic authentication
- Added authentication JWT support for protected API endpoints
- Added authentication JWT validation middleware to protect API endpoints
* Thu Oct 19 2017 - jscherff@24hourfit.com
- Added SQL script to truncate all tables
* Fri Oct 13 2017 - jscherff@24hourfit.com
- Refactored and streamlined
- Added API endpoints for device information lookups
* Mon Oct 9 2017 - jscherff@24hourfit.com
- Modified table, view, and stored procedure names
- Added column to each table for the JSON object
- Modified changes column in changes table to be datatype JSON
* Sat Oct 7 2017 - jscherff@24hourfit.com
- Added v1 prefix to URLs and handlers
* Sat Sep 30 2017 - jscherff@24hourfit.com
- Tightened file permissions mode on config.json
