package models

// LauncherTokenSettings contains settings for stored launcher token.
type LauncherTokenSettings struct {
	//TTL is the expiration time for the token.
	TTL int `bson:"ttl" json:"ttl"`
}

const (
	LauncherAuth_InProgress = "in_progress"
	LauncherAuth_Success    = "success"
	LauncherAuth_Canceled   = "canceled"
)

type LauncherToken struct {
	// Challenge is login_challenge
	Challenge string `json:"challenge"`
	// UserIdentity stores user identity
	UserIdentity *UserIdentity `json:"ui"`
	// UserIdentitySocial stores user social profile data
	UserIdentitySocial *UserIdentitySocial `json:"uis"`
	// Name is the name of social provider
	Name string `json:"name"`
	// Status stores state of the login process
	Status string `json:"status"`
	// URL to finish
	URL string `json:"url"`
}
