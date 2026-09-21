package app

type link struct {
	ID    string
	Label string
	URL   string
}

type profile struct {
	Name      string
	Username  string
	Bio       string
	Avatar    string
	AvatarAlt string
	Links     []link
}

var siteProfile = profile{
	Name:      "Sasha",
	Username:  "parrahex",
	Bio:       "Trained on Counter-Strike voice chat",
	Avatar:    "/assets/avatar.jpg",
	AvatarAlt: "Sasha’s avatar: blurred city lights at night",
	Links: []link{
		{ID: "github", Label: "GitHub", URL: "https://github.com/parrahex"},
		{ID: "telegram", Label: "Telegram", URL: "https://t.me/parrrahex"},
		{ID: "kofi", Label: "Ko-fi", URL: "https://ko-fi.com/parrahex"},
	},
}
