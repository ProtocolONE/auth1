package manager

import (
	"fmt"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/models"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/service"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/validator"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type ChangePasswordManager struct {
	redis                   *redis.Client
	r                       service.InternalRegistry
	userIdentityService     *models.UserIdentityService
	identityProviderService *service.AppIdentityProviderService
}

func NewChangePasswordManager(db *mgo.Session, r *redis.Client, ir service.InternalRegistry) *ChangePasswordManager {
	m := &ChangePasswordManager{
		redis:                   r,
		r:                       ir,
		userIdentityService:     models.NewUserIdentityService(db),
		identityProviderService: service.NewAppIdentityProviderService(),
	}

	return m
}

func (m *ChangePasswordManager) ChangePasswordStart(form *models.ChangePasswordStartForm) *models.GeneralError {
	app, err := m.r.ApplicationService().Get(bson.ObjectIdHex(form.ClientID))
	if err != nil {
		return &models.GeneralError{Code: "client_id", Message: models.ErrorClientIdIncorrect, Err: errors.Wrap(err, "Unable to load application")}
	}

	ipc := m.identityProviderService.FindByTypeAndName(app, models.AppIdentityProviderTypePassword, models.AppIdentityProviderNameDefault)
	if ipc == nil {
		return &models.GeneralError{Code: "client_id", Message: models.ErrorUnknownError, Err: errors.New("Unable to get identity provider")}
	}

	ui, err := m.userIdentityService.Get(app, ipc, form.Email)
	if err != nil {
		return &models.GeneralError{Code: "email", Message: models.ErrorUnknownError, Err: errors.Wrap(err, "Unable to get user identity by email")}
	}

	if ui == nil || ui.ID == "" {
		// INFO: Do not need to disclose the login
		return nil
	}

	ottSettings := &models.OneTimeTokenSettings{
		Length: app.PasswordSettings.TokenLength,
		TTL:    app.PasswordSettings.TokenTTL,
	}
	token, err := m.r.OneTimeTokenService().Create(&models.ChangePasswordTokenSource{Email: form.Email}, ottSettings)
	if err != nil {
		return &models.GeneralError{Code: "common", Message: models.ErrorUnableCreateOttSettings, Err: errors.Wrap(err, "Unable to create OneTimeToken")}
	}

	if err := m.r.Mailer().Send(form.Email, "Change password token", fmt.Sprintf("Token: %s", token.Token)); err != nil {
		return &models.GeneralError{Code: "common", Message: models.ErrorUnknownError, Err: errors.Wrap(err, "Unable to send mail with change password token")}
	}

	return nil
}

func (m *ChangePasswordManager) ChangePasswordVerify(form *models.ChangePasswordVerifyForm) *models.GeneralError {
	if form.PasswordRepeat != form.Password {
		return &models.GeneralError{Code: "password_repeat", Message: models.ErrorPasswordRepeat, Err: errors.New(models.ErrorPasswordRepeat)}
	}

	app, err := m.r.ApplicationService().Get(bson.ObjectIdHex(form.ClientID))
	if err != nil {
		return &models.GeneralError{Code: "client_id", Message: models.ErrorClientIdIncorrect, Err: errors.Wrap(err, "Unable to load application")}
	}

	if false == validator.IsPasswordValid(app, form.Password) {
		return &models.GeneralError{Code: "password", Message: models.ErrorPasswordIncorrect, Err: errors.New(models.ErrorPasswordIncorrect)}
	}

	ts := &models.ChangePasswordTokenSource{}
	if err := m.r.OneTimeTokenService().Use(form.Token, ts); err != nil {
		return &models.GeneralError{Code: "common", Message: models.ErrorCannotUseToken, Err: errors.Wrap(err, "Unable to use OneTimeToken")}
	}

	ipc := m.identityProviderService.FindByTypeAndName(app, models.AppIdentityProviderTypePassword, models.AppIdentityProviderNameDefault)
	if ipc == nil {
		return &models.GeneralError{Code: "common", Message: models.ErrorUnknownError, Err: errors.New("Unable to get identity provider")}
	}

	ui, err := m.userIdentityService.Get(app, ipc, ts.Email)
	if err != nil {
		return &models.GeneralError{Code: "common", Message: models.ErrorUnknownError, Err: errors.Wrap(err, "Unable to get user identity")}
	}

	be := models.NewBcryptEncryptor(&models.CryptConfig{Cost: app.PasswordSettings.BcryptCost})
	ui.Credential, err = be.Digest(form.Password)
	if err != nil {
		return &models.GeneralError{Code: "password", Message: models.ErrorCryptPassword, Err: errors.Wrap(err, "Unable to crypt password")}
	}

	if err = m.userIdentityService.Update(ui); err != nil {
		return &models.GeneralError{Code: "password", Message: models.ErrorUnableChangePassword, Err: errors.Wrap(err, "Unable to update password")}
	}

	return nil
}
