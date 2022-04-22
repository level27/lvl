#!/usr/bin/env pwsh

# Needs https://github.com/cloudbase/powershell-yaml.git

Import-Module powershell-yaml

$dat = Get-Content ~/.lvl.yaml | ConvertFrom-Yaml
$env:L27_TEST_KEY = $dat.apikey
$env:L27_TEST_API = $dat.apiurl
