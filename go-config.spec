%global debug_package %{nil}
%global goipath         github.com/farsightsec/go-config

%gometa

%global common_description %{expand:
Contains types useful for validating, parsing, and loading values of
some useful types in configuration files.}

Name:           go-config
Version:        0.1.1
Release:        1%{?dist}
Summary:        Minimalist Go config library

License:        MPL-2.0
URL:            %{gourl}
Source0:        https://github.com/farsightsec/go-config/archive/v%{version}/%{name}-%{version}.tar.gz

BuildRequires:  golang
BuildArch:      noarch

%description
%{common_description}

%package -n golang-github-farsightsec-go-config-devel
Summary:        %{summary}
Provides:       golang(%{goipath}) = %{version}

%description -n golang-github-farsightsec-go-config-devel
%{common_description}

This package provides the Go source code for importing.

%prep
%setup -q -n %{name}-%{version}

%build
# Go library — nothing to compile.

%install
install -d -p %{buildroot}/%{gopath}/src/%{goipath}
for file in $(find . -iname "*.go" \! -iname "*_test.go") ; do
    install -d -p %{buildroot}/%{gopath}/src/%{goipath}/$(dirname $file)
    install -p -m 0644 $file %{buildroot}/%{gopath}/src/%{goipath}/$file
done
for meta in go.mod go.sum ; do
    [ -f "$meta" ] && install -p -m 0644 $meta %{buildroot}/%{gopath}/src/%{goipath}/$meta
done

%files -n golang-github-farsightsec-go-config-devel
%license LICENSE
%doc README.md
%{gopath}/src/%{goipath}

%changelog
* Tue Jul 08 2025 DomainTools RelEng <releng@domaintools.com> - 0.1.1-1
- Initial RPM packaging (synthesized from test branches)
