# =============================================================================
%define		name	gocmdbd
%define		version	1.0.0
%define		release	1
%define		gecos	GoCMDBd Service
%define		summary	[Go] Configuration Management Database Daemon
%define		author	John Scherff <jscherff@24hourfit.com>
%define		prefix	/opt/gocmdbd
%define		sbindir	%{prefix}/sbin
%define		confdir	%{prefix}/etc
%define		docdir	%{prefix}/doc
%define		logdir	%{_var}/log/gocmdbd
%define		sysdir	%{_sysconfdir}/systemd/system
%define		syslib	%{_prefix}/lib/systemd/system
%define		gopath	%{_builddir}/go
%define		package	github.com/jscherff/gocmdbd
# =============================================================================

Name:		%{name}
Version:	%{version}
Release:	%{release}%{?dist}
Summary:	%{summary}

Group:		Applications/System
License:	ASL 2.0
URL:		https://www.24hourfitness.com
Vendor:		24 Hour Fitness, Inc.
Prefix:		%{prefix}
Packager: 	%{packager}
BuildRoot:	%{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
Distribution:	el

BuildRequires:    golang >= 1.8.0
Requires(pre):    /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel

%description
Go Configuration Management Database Daemon, gocmdbd, is a lightweight HTTP
daemon that provides a REST API for gocmdbcli clients installed on Windows
endpoints. The clients collect information about attached devices and send
that information to the server for storage in the database. Clients can
register attached devices with the server, obtain unique serial numbers from
the server for devices that support serial number configuration, perform
audits against previous device configurations, and report any configuration
changes found during the audit to the server for later analysis.

%prep

%build

  export GOPATH=%{gopath}
  go get %{package}
  go build %{package}

%install

  test %{buildroot} != / && rm -rf %{buildroot}

  mkdir -p %{buildroot}{%{sbindir},%{confdir},%{docdir},%{syslib}}
  mkdir -p %{buildroot}%{logdir}

  install -s -m 755 gocmdbd %{buildroot}%{sbindir}/
  install -m 644 go/src/%{package}/config.json %{buildroot}%{confdir}/
  install -m 644 %{gopath}/src/%{package}/{README.md,LICENSE} %{buildroot}%{docdir}/
  install -m 644 %{gopath}/src/%{package}/svc/gocmdbd.service %{buildroot}%{syslib}/
  install -m 640 %{gopath}/src/%{package}/ddl/{database.sql,users.sql} %{buildroot}%{confdir}/

%pre

  /usr/bin/getent passwd %{name} >/dev/null 2>&1 ||
    /usr/sbin/useradd -Mrd %{prefix} -c '%{gecos}' -s /sbin/nologin %{name}

%post

  systemctl --quiet is-enabled gocmdbd || systemctl --quiet enable gocmdbd


%preun

  systemctl --quiet is-active gocmdbd && systemctl --quiet stop gocmdbd
  systemctl --quiet is-enabled gocmdbd && systemctl --quiet disable gocmdbd

%postun

  /usr/bin/getent passwd %{name} >/dev/null 2>&1 &&
    /usr/sbin/userdel %{name}

%clean

  test %{buildroot} != / && rm -rf %{buildroot}
  test %{_builddir} != / && rm -rf %{_builddir}/*

%files

  %defattr(-,root,root)
  %{sbindir}/gocmdbd
  %{docdir}/README.md
  %{syslib}/gocmdbd.service
  %{confdir}/database.sql
  %{confdir}/users.sql
  %license %{docdir}/LICENSE
  %config %{confdir}/config.json
  
  %defattr(-,%{name},%{name})
  %{logdir}

%changelog
