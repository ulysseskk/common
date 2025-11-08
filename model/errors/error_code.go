package errors

const (
	// api error code
	RequestParameterInvalid int = 40001
	RequestDataExists       int = 40002
	RequestDataNotExisted   int = 40003
	AuthFailed              int = 40004
	PermissionDeny          int = 40005
	UserNotExists           int = 40011
	InvalidOperation        int = 40016
	InvalidArgument         int = 40017

	InternalError        int = 50000
	InvalidDataError     int = 50001
	WorkflowError        int = 50002
	InternalServiceError int = 50003
	ExternalServiceError int = 50004
	CodeDatabaseError        = 50005

	ClientError       int = 60001
	RedisError        int = 60002
	K8SOperationError int = 60003
	OpensearchError   int = 60004

	CodeInitializeError = 70001
	CodeLackOfConfig    = 70002

	// 8开头是k8s相关
	CodeStorageClassNotFound = 80001
	CodeControllerInitError  = 80002
)
