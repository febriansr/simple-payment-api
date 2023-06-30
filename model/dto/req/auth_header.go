package req

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}
