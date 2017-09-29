# =============================================================================
%define		name	cmdbd
%define		version	1.0.1
%define		release	1
%define		gecos	CMDBd Service
%define		summary	Configuration Management Database Daemon
%define		author	John Scherff <jscherff@24hourfit.com>
%define		package	github.com/jscherff/%{name}
%define		gopath	%{_builddir}/go
%define		docdir	%{_docdir}/%{name}-%{version}
%define		logdir	%{_var}/log/%{name}
%define		syslib	%{_prefix}/lib/systemd/system
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
  go get %{package}
  go build %{package}

%install

  test %{buildroot} != / && rm -rf %{buildroot}

  mkdir -p %{buildroot}{%{_sbindir},%{confdir},%{syslib},%{logdir},%{docdir}}

  install -s -m 755 %{name} %{buildroot}%{_sbindir}/
  install -m 644 go/src/%{package}/config.json %{buildroot}%{confdir}/
  install -m 644 go/src/%{package}/svc/%{name}.service %{buildroot}%{syslib}/
  install -m 644 go/src/%{package}/{README.md,LICENSE} %{buildroot}%{docdir}/
  install -m 640 go/src/%{package}/ddl/{%{name}.sql,users.sql} %{buildroot}%{docdir}/

%pre

  %{_bindir}/getent passwd %{name} >/dev/null 2>&1 ||
    %{_sbindir}/useradd -Mrd %{_sbindir} -c '%{gecos}' -s /sbin/nologin %{name}

%post

  systemctl --quiet is-enabled %{name} || systemctl --quiet enable %{name} 


%preun

  systemctl --quiet is-active %{name} && systemctl --quiet stop %{name}
  systemctl --quiet is-enabled %{name} && systemctl --quiet disable %{name}

%postun

  %{_bindir}/getent passwd %{name} >/dev/null 2>&1 &&
    %{_sbindir}/userdel %{name}

%clean

  test %{buildroot} != / && rm -rf %{buildroot}
  test %{_builddir} != / && rm -rf %{_builddir}/*

%files

  %defattr(-,root,root)
  %{_sbindir}/%{name}
  %{syslib}/%{name}.service
  %{docdir}/%{name}.sql
  %{docdir}/README.md
  %{docdir}/users.sql
  %license %{docdir}/LICENSE
  %config %{confdir}/config.json
  
  %defattr(-,%{name},%{name})
  %{logdir}

%changelog
