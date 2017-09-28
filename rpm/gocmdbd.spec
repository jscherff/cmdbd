%define		name	gocmdbd
%define		version	1.0.0
%define		release	1
%define		author	John Scherff <jscherff@24hourfit.com>
%define		prefix	/opt/gocmdbd
%define		sbindir	%{prefix}/sbin
%define		confdir	%{prefix}/etc
%define		docdir	%{prefix}/doc
%define		gopath	%{_builddir}/go
%define		sysdir	%{_sysconfdir}/systemd/system
%define		syslib	%{_prefix}/lib/systemd/system
%define		gecos	GoCMDBd Service
%define		package	github.com/jscherff/gocmdbd

Name:		%{name}
Version:	%{version}
Release:	%{release}%{?dist}
Summary:	Go Configuration Management Database Daemon

Group:		Applications/System
License:	ASL 2.0
URL:		https://www.24hourfitness.com
Vendor:		24 Hour Fitness, Inc.
Prefix:		%{prefix}
Packager: 	%{packager}
BuildRoot:	%{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:    golang >= 1.8.0
Requires(pre):    /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel

%description
Go Configuration Management Database Daemon, gocmdbd, is a lightweight HTTP
daemon that provides a REST API for the gocmdbcli client running on Windows
endpoints. The client collects information about attached devices and sends
that information to the server for storage in the database. The client can
perform audits and report on addition or removal of devices or configuration
changes to existing devices.

%prep

%build

  export GOPATH=%{gopath}
  go get %{package}
  go build %{package}

%install

  test %{buildroot} != / && rm -rf %{buildroot}

  mkdir -p %{buildroot}{%{sbindir},%{confdir},%{docdir},%{syslib}}

  install -s -m 755 gocmdbd %{buildroot}%{sbindir}/
  install -m 640 go/src/%{package}/config.json %{buildroot}%{confdir}/
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
  test %{gopath} != / && rm -rf %{gopath}

%files

  %defattr(-,root,root)
  %{sbindir}/gocmdbd
  %{docdir}/README.md
  %{syslib}/gocmdbd.service
  %{confdir}/database.sql
  %{confdir}/users.sql
  %config %{confdir}/config.json
  %license %{docdir}/LICENSE

%changelog
