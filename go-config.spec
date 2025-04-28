%global debug_package %{nil}
# https://github.com/farsightsec/go-config
%global goipath         github.com/farsightsec/go-config

Name:		go-config		
Version:	0.1.1
Release:	1%{?dist}
Summary:	Minimalist go config library

%gometa

%global common_description %{expand:
Contains types useful for validating, parsing, and loading values of
some useful types in configuration files.}

License:	MPLv2.0
URL:		https://github.com/farsightsec/go-config
Source0:	https://github.com/farsightsec/go-config/archive/%{name}-%{version}.tar.gz

BuildRequires:	%{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang} 

%description %{common_description}

%gopkg

%prep
%goprep -A
%autopatch -p1

%generate_buildrequires
%go_generate_buildrequires

%install
%gopkginstall

%if %{with check}
%check
%gocheck
%endif

%gopkgfiles

%changelog
