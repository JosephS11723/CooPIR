package config

// network configuration
var FtpIP string = "10.10.10.2"
var ApiIP string = "10.10.10.4"
var DBIP string = "mongodb://10.10.10.3:27017"

// authentication credentials
var FtpUsername string = "api"
var FtpPassword string = "password"

// ftp config
var FtpConnectionsPerHost int = 1
var FtpTimeout int = 100 // seconds