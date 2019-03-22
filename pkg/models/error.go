package models

var (
	ErrorUnknownError             = "Unknown error"
	ErrorInvalidRequestParameters = "Invalid request parameters"
	ErrorRequiredField            = "This is required field"
	ErrorAddAuthLog               = "Unable to add auth log"
	ErrorCreateCookie             = "Unable to create cookie"
	ErrorCreateUser               = "Unable to create user"
	ErrorCreateUserIdentity       = "Unable to create user identity"
	ErrorLoginIncorrect           = "Login is incorrect"
	ErrorCryptPassword            = "Unable to crypt password"
	ErrorUnableChangePassword     = "Unable to change password"
	ErrorUnableCreateOttSettings  = "Unable create ott settings"
	ErrorPasswordIncorrect        = "Password is incorrect"
	ErrorPasswordRepeat           = "Password repeat is not equal to password"
	ErrorUnableValidatePassword   = "Unable to validate password"
	ErrorClientIdIncorrect        = "Client ID is incorrect"
	ErrorConnectionIncorrect      = "Connection is incorrect"
	ErrorCannotCreateToken        = "Cannot create token"
	ErrorCannotUseToken           = "Cannot use this token"
	ErrorRedirectUriIncorrect     = "Redirect URI is incorrect"
	ErrorCaptchaRequired          = "Captcha required"
	ErrorCaptchaIncorrect         = "Captcha is incorrect"
	ErrorAuthTemporaryLocked      = "Temporary locked"
	ErrorProviderIdIncorrect      = "Provider ID is incorrect"
	ErrorGetSocialData            = "Unable to load social data"
	ErrorGetSocialSettings        = "Unable to load social settings"
	ErrorMfaRequired              = "MFA required"
	ErrorMfaClientAdd             = "Unable to add MFA"
	ErrorMfaCodeInvalid           = "Invalid MFA code"
	ErrorCsrfSignature            = "Invalid request, start authorization or registration again"
	ErrorLoginChallenge           = "Invalid login challenge"
)

type ErrorInterface interface {
	GetHttpCode() int
	GetCode() string
	GetMessage() string
	Error() string
}

type CommonError struct {
	HttpCode int    `json:"error,omitempty"`
	Code     string `json:"error,omitempty"`
	Message  string `json:"error_message,omitempty"`
	Csrf     string `json:"csrf,omitempty"`
}

func (m *CommonError) Error() string {
	return m.Message
}

func (m *CommonError) GetHttpCode() int {
	return m.HttpCode
}

func (m *CommonError) GetCode() string {
	return m.Code
}

func (m *CommonError) GetMessage() string {
	return m.Message
}
