%global debug_package %{nil}
# https://github.com/farsightsec/go-config
%global goipath         github.com/farsightsec/go-config

%gometa

Name:		go-config		
Version:	0.1.1
Release:	1%{?dist}
Summary:	Minimalist go config library

License:	MPLv2.0
URL:		https://github.com/farsightsec/go-config
#Source0:	https://github.com/farsightsec/go-config/archive/%{name}-%{version}.tar.gz
Source0:	https://github.com/farsightsec/go-config/archive/%{name}.tar.gz

BuildRequires:	%{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang} 
	
%description
Contains types useful for validating, parsing, and loading values of
some useful types in configuration files.

%package -n %{goname}-devel
Summary:	%{summary}
BuildArch:  noarch
%description -n %{goname}-devel
%{common_description}

%prep
%setup -q

%install
for file in $(find . -iname "*.go" \! -iname "*_test.go" \! -iname "main.go" ) ; do
    echo "%%dir %%{gopath}/src/%%{goipath}/$(dirname $file)" >> devel.file-list
    install -d -p %{buildroot}/%{gopath}/src/%{goipath}/$(dirname $file)
    cp -pav $file %{buildroot}/%{gopath}/src/%{goipath}/$file
    echo "%%{gopath}/src/%%{goipath}/$file" >> devel.file-list
done
sort -u -o devel.file-list devel.file-list

#define license tag if not already defined
%{!?_licensedir:%global license %doc}

# Not sure how this should be done right in rhel8
#%license LICENSE 
%doc README.md
%files -n %{goname}-devel -f devel.file-list

%changelog
