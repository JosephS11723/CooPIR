package logtypes

// this package is mostly used for specifying log message text

// default values
// LogDefaultCaseUUID is the default UUID for the case log
var LogDefaultCaseUUID = "00000000-0000-0000-0000-000000000000"
// LogDefaultUser is the default user UUID for the case log
var LogDefaultUserUUID string = "00000000-0000-0000-0000-000000000000"

// ping pong
var PingPong string = "PingPong"

// file logs
var FileDownloadAttempt string = "FileDownloadAttempt"
var FileDownload string = "FileDownload"
var FileDownloadFailure string = "FileDownloadFailure"

var FileUploadAttempt string = "FileUploadAttempt"
var FileUpload string = "FileUpload"
var FileUploadFailure string = "FileUploadFailure"

var FileDeleteAttempt string = "FileDeleteAttempt"
var FileDelete string = "FileDelete"
var FileDeleteFailure string = "FileDeleteFailure"

// case logs
var GetCaseInfoAttempt string = "GetCaseInfoAttempt"
var GetCaseInfo string = "GetCaseInfo"
var GetCaseInfoFailure string = "GetCaseInfoFailure"

var CreateCaseAttempt string = "CreateCaseAttempt"
var CreateCase string = "CreateCase"
var CreateCaseFailure string = "CreateCaseFailure"

var UpdateCaseAttempt string = "UpdateCaseAttempt"
var UpdateCase string = "UpdateCase"
var UpdateCaseFailure string = "UpdateCaseFailure"

var GetCases string = "GetCases"

var GetCaseFilesAttempt string = "GetCaseFilesAttempt"
var GetCaseFiles string = "GetCaseFiles"
var GetCaseFilesFailure string = "GetCaseFilesFailure"

// user info logs
var GetUserInfoAttempt string = "GetUserInfoAttempt"
var GetUserInfo string = "GetUserInfo"
var GetUserInfoFailure string = "GetUserInfoFailure"

var CreateUserAttempt string = "CreateUserAttempt"
var CreateUser string = "CreateUser"
var CreateUserFailure string = "CreateUserFailure"

var GetUserEditUser string = "GetUserEditUser"

var GetUserAddFile string = "GetUserAddFile"

var GetUserMakeCase string = "GetUserMakeCase"

var UpdateUserAttempt string = "UpdateUserAttempt"
var UpdateUser string = "UpdateUser"
var UpdateUserFailure string = "UpdateUserFailure"

// file info logs
var GetFileInfoAttempt string = "GetFileInfoAttempt"
var GetFileInfo string = "GetFileInfo"
var GetFileInfoFailure string = "GetFileInfoFailure"

// authentication
var RenewTokenAttempt string = "RenewTokenAttempt"
var RenewToken string = "RenewToken"
var RenewTokenFailure string = "RenewTokenFailure"

var LogoutAttempt string = "LogoutAttempt"
var Logout string = "Logout"
var LogoutFailure string = "LogoutFailure"

var LoginAttempt string = "LoginAttempt"
var Login string = "Login"
var LoginFailure string = "LoginFailure"

// worker events
var WorkerStart string = "WorkerStart"
var WorkerStop string = "WorkerStop"

var WorkerFileUpload string = "WorkerFileUpload"
var WorkerFileUploadFailure string = "WorkerFileUploadFailure"

var WorkerFileModify string = "WorkerFileModify"
var WorkerFileModifyFailure string = "WorkerFileModifyFailure"

var WorkerResult string = "WorkerResult"
var WorkerResultAttempt string = "WorkerResultAttempt"
var WorkerResultFailure string = "WorkerResultFailure"

var WorkerResultError string = "WorkerResultError"
var WorkerResultDone string = "WorkerResultDone"