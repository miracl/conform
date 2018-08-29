# CONFORM

_Library providing routines to help supporting multiple versions of configuration files_

[![Build Status](https://secure.travis-ci.org/miracl/conform.png?branch=master)](https://travis-ci.org/miracl/conform?branch=master)
[![Coverage Status](https://coveralls.io/repos/miracl/conform/badge.svg?branch=master&service=github)](https://coveralls.io/github/miracl/conform?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/miracl/conform)](https://goreportcard.com/report/github.com/miracl/conform)

## Description

This library provides some primitive operations to transform data between different format versions, as defined by one or more JSON schemas.
The library works closely with [conflate](github.com/miracl/conflate).

The library operates on data in the same form as that returned by `json.Unmarshal` i.e. hierarchical data containing any number of nested
`map[string]interface{}` and/or `[]interface[]` collections.

The overall aim of the library is to make it easier to support backwards-compatibility with older versions of configuration files, when you
release a new version of your application.

For each new configuration version, you create a JSON schema, as for [conflate](github.com/miracl/conflate), and you create a `Conformer` object
to convert old configuration files to the new schema version. The `Conformer` is a functional-style chain of operations to
move, rename, copy, delete or change the values of items in the configuration object. Multiple `Conformer` objects can be chained together so that
multiple historical versions of your configuration files, can be supported by your application.

This library is in the early stages of development and will be further developed in a fork.
