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

// seaweedfs config
var UpTaskCount int = 5          // the number of tasks performed on a file when it is uploaded
var DoUploadLogging bool = false // whether or not the upload process should be logged with hashes
var ReadOnlyFiles bool = true    // whether files uploaded to seaweedfs should be read-only
