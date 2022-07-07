package profiles

type Profile struct {
	Username  string
	Bio       string
	Image     string
	Following *bool
}

func NewProfile(username string, bio string, image string, following *bool) Profile {
	return Profile{
		Username:  username,
		Bio:       bio,
		Image:     image,
		Following: following,
	}
}
