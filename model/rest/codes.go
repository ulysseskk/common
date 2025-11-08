package rest

const (
    CodeSuccess int = 20000
    // api error code
    RequestParameterInvalid int = 40001
    RequestDataExists       int = 40002
    RequestDataNotExisted   int = 40003
    AuthFailed              int = 40004
    PermissionDeny          int = 40005
    UserNotExists           int = 40011
    InvalidOperation        int = 40016

    InternalError        int = 50000
    DatabaseError        int = 50001
    WorkflowError        int = 50002
    InternalServiceError int = 50003
    ExternalServiceError int = 50004
    RequestFatalError    int = 50008
    ClientError          int = 60001
)

var codeMessageMap = map[int]string{
    RequestParameterInvalid: "request parameter is invalid",
    RequestDataExists:       "request data already exists",
    RequestDataNotExisted:   "request data does not exists",
    AuthFailed:              "authorization failed",
    PermissionDeny:          "no permission",
    UserNotExists:           "user not exists",

    InternalError:        "internal error",
    DatabaseError:        "database operation error",
    WorkflowError:        "request workflow api error",
    InternalServiceError: "internal service error",
    ExternalServiceError: "external service error",
    InvalidOperation:     "invalid operation",
    RequestFatalError:    "service has a fatal error when execute request",
}

func GetErrorMessage(code int) (string, bool) {
    msg, ok := codeMessageMap[code]
    return msg, ok
}
