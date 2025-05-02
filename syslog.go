package main

import "log/syslog"

var Facilities = map[string]syslog.Priority{
	"KERN":     syslog.LOG_KERN,
	"USER":     syslog.LOG_USER,
	"MAIL":     syslog.LOG_MAIL,
	"DAEMON":   syslog.LOG_DAEMON,
	"AUTH":     syslog.LOG_AUTH,
	"SYSLOG":   syslog.LOG_SYSLOG,
	"LPR":      syslog.LOG_LPR,
	"NEWS":     syslog.LOG_NEWS,
	"UUCP":     syslog.LOG_UUCP,
	"CRON":     syslog.LOG_CRON,
	"AUTHPRIV": syslog.LOG_AUTHPRIV,
	"FTP":      syslog.LOG_FTP,
	"LOCAL0":   syslog.LOG_LOCAL0,
	"LOCAL1":   syslog.LOG_LOCAL1,
	"LOCAL2":   syslog.LOG_LOCAL2,
	"LOCAL3":   syslog.LOG_LOCAL3,
	"LOCAL4":   syslog.LOG_LOCAL4,
	"LOCAL5":   syslog.LOG_LOCAL5,
	"LOCAL6":   syslog.LOG_LOCAL6,
	"LOCAL7":   syslog.LOG_LOCAL7}
