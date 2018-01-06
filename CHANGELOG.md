# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### TODO
- more test case, and more covering.
- Range for sql in (${u.Cities}).
- Expression for select fields.
- Rewrite types
- Rewrite sql parser
- id/page/offset/size null wrap not required
- ${(page-1)*size} Support
### Added
- Intelligent guess, use ? not ${...}
- Parse go source file by go/types(get signatures) and go/parser(get comments).
- Parse sql (comment) to generate Statement AST.
- Support CREATE statement.
- Replacers #{...} like Variables ${...}.
- Generate implemented file by method signature and sql AST.
