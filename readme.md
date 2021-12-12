# MEWS API

[![Go](https://github.com/fahminlb33/devoria1-wtc-backend/actions/workflows/go.yml/badge.svg)](https://github.com/fahminlb33/devoria1-wtc-backend/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/fahminlb33/devoria1-wtc-backend/branch/master/graph/badge.svg?token=hRNbJKqQgM)](https://codecov.io/gh/fahminlb33/devoria1-wtc-backend)

> Great news comes from great backend service :wink:

This is my submission for DEVORIA WTC.

In this repository you'll find my preferred codebase style which are highly influenced by my other project, [TASI](https://github.com/fahminlb33/tasi-backend). This project structure is also used in one of my squad service (LGC-Visibility). I hope I can demonstrate a clear and concise clean architecture in functional way (I'm more of an OOP guy).

More explanation will come later on after I finish testing this project.

## Pretty Coverage Report

I already included Codecov report if you want fancy report, but if you want a pretty coverage report that you can generate locally, I have the script in `./make coverage_pretty`, but you'll need to install these additional dependencies.

```bash
dotnet tool install -g dotnet-reportgenerator-globaltool
go install github.com/axw/gocov/gocov
go install github.com/AlekSi/gocov-xml
go install github.com/matm/gocov-html
```

Those tools will export the coverage as Cobertura XML and HTML based report. The XML file is used for ReportGenerator. You can then use this XML file as coverage report for SonarQube or Jenkins.
