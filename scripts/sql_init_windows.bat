@echo off
REM 用法示例：
REM sql_init_windows.bat root mypass

chcp 65001 >nul

if "%~2"=="" (
    echo 用法: %~nx0 用户名 密码
    echo 示例: %~nx0 root mypass
    exit /b 1
)

set MYSQL_USER=%~1
set MYSQL_PWD=%~2
set MYSQL_HOST=localhost
set SQL_FILE=maze_init.sql

REM 检查 mysql 是否存在
where mysql >nul 2>nul
if errorlevel 1 (
    echo 错误：未找到 mysql 命令，请检查 MySQL 是否安装并配置了 PATH。
    exit /b 1
)

mysql -u %MYSQL_USER% -p%MYSQL_PWD% -h %MYSQL_HOST% %MYSQL_DB% < "%SQL_FILE%"

if errorlevel 1 (
    echo 执行失败。
) else (
    echo 执行成功。
)
