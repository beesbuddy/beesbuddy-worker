package internal

const (
	WorkerKey              = "worker"
	TopicPath              = "apiary/%s/hive/%s"
	SessionAuthUserIDKey   = "authUserID"
	SessionFlashMessageKey = "flashMessage"
	RememberCookieNameKey  = "remember"
	SessionIntendedURLKey  = "intendedURL"
)

type contextKey string

const AuthUserKey = contextKey("authUser")
