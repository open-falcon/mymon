# 更新日志

## [version] - 2018-7-25

本次更新，主要是在代码规范、代码结构及代码测试上的更新，增加了少量必要的监控项。

### Added

- 新增了项目的单元测试和集成测试，目前整体测试覆盖率在83.3%。
- 新增了监控项`longest_transaction`、`Is_slave`、`binlog_file_size`等。
- 新增了`CHANGES.md`文件。
- 新增了`vendor`包管理。
- 新增了源代码文件头Copyright信息。
- 新增了`Makefile`文件。

### Fixed

- 修复了`binlogFileSize`返回空值的bug 。
- 修复了当没有`mysql.heartbeat`表或数据为空时，数组越界报错问题。

### Changed

- 修改了`README.md`文件。
- 修改了`NOTICE`文件。
- 修改了日志级别为RFC5424标准的7个日志级别。

### Refactored

- 整理了文件结构，合并了`SHOW`相关的监控变量。
- 抽出了`common`包，包含数据库配置、日志配置、配置文件以及utils相关代码。
- 更改了大量变量名、函数名、编码方式、注释，增强可读性并通过gometalinter语法检查。
- 通过正则重构了`falconignore`文件解析相关代码。
- 对部分重复代码进行了抽象重构。
