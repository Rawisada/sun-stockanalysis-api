package auth

type LoginInput struct {
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

type RegisterInput struct {
	Body struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
}

type RefreshInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token"`
	}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type RegisterResult struct {
	UserID string
}
