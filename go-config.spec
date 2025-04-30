# Define backup go macros
%if %{rhel} == 8
%global gopkg %package -n %{goname}-devel \
Summary:	%{summary} \
BuildArch:  noarch \
%description -n %{goname}-devel \
%{common_description}
%global goprep(-A) %setup -q
%global generate_buildrequires echo "Need more specific macro on rhel8"
%global gopkginstall for file in $(find . -iname "*.go" \! -iname "*_test.go" \! -iname "main.go" ) ; do \
    echo "%%dir %%{gopath}/src/%%{goipath}/$(dirname $file)" >> devel.file-list ;\
    install -d -p %{buildroot}/%{gopath}/src/%{goipath}/$(dirname $file) ;\
    cp -pav $file %{buildroot}/%{gopath}/src/%{goipath}/$file ;\
    echo "%%{gopath}/src/%%{goipath}/$file" >> devel.file-list ;\
done ;\
sort -u -o devel.file-list devel.file-list
%global gopkgfiles %files -n %{goname}-devel -f devel.file-list
%global gocheck echo "skipping gocheck on rhel8"
# Specific BuildRequires macro
%global go_generate_buildrequires BuildRequires:	%{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
%endif

%global debug_package %{nil}
# https://github.com/farsightsec/go-config
%global goipath         github.com/farsightsec/go-config
%global common_description %{expand:
Contains types useful for validating, parsing, and loading values of
some useful types in configuration files.}

Name:		go-config		
Version:	0.1.1
Release:	1%{?dist}
Summary:	Minimalist go config library
%gometa
License:	MPLv2.0
URL:		https://github.com/farsightsec/go-config
Source0:	https://github.com/farsightsec/go-config/archive/%{name}-%{version}.tar.gz

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
