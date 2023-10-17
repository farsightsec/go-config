%global debug_package %{nil}

# https://github.com/farsightsec/go-config
%global goipath         github.com/farsightsec/go-config
Version:                0.1.1

%gometa

%global common_description %{expand:
Provides types useful for validating, parsing, and loading values of some useful types in config files.}

%global golicenses      LICENSE
%global godocs          README.md

Name:           %{goname}
Release:        %autorelease
Summary:        Minimalist Go config library

License:        MPLv2.0
URL:            %{gourl}
Source0:        %{gosource}

%description
%{common_description}

%gopkg

%prep
%goprep

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
%autochangelog
