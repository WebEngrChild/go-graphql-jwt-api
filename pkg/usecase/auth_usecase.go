package usecase

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/domain/repository"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	Login(c echo.Context, fv *FormValue) (string, error)
	JwtParser(auth string) (*jwt.MapClaims, error)
	IdentifyJwtUser(id string) error
	DeleteCookie(c echo.Context, ck *http.Cookie)
}

type AuthUseCase struct {
	config   *config.Config
	userRepo repository.User
}

type FormValue struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// 独自クレーム型
type jwtCustomClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthUseCase(userRepo repository.User, config *config.Config) Auth {
	AuthUseCase := AuthUseCase{
		userRepo: userRepo,
		config:   config,
	}
	return &AuthUseCase
}

func (a *AuthUseCase) JwtParser(auth string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, xerrors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.config.JwtSecret), nil
	})

	if err != nil {
		return nil, xerrors.Errorf("JwtParser failed: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, xerrors.Errorf("inValid claims: %v", claims)
	}

	return &claims, nil
}

func (a *AuthUseCase) Login(c echo.Context, fv *FormValue) (string, error) {
	// Email暗号化
	encEmail, err := a.userRepo.Encrypt(fv.Email)
	if err != nil {
		return "", fmt.Errorf("login failed at Encrypt err %w", err)
	}

	// 暗号化済みEmailで取得
	user, err := a.userRepo.GetUserByEmail(encEmail)
	if err != nil {
		return "", fmt.Errorf("login failed at GetUserByEmail err %w", err)
	}

	// DBから取得したpasswordを復号化
	pass, err := a.userRepo.Decrypt(user.Password)
	if err != nil {
		return "", fmt.Errorf("login failed at Decrypt err %w", err)
	}

	// フォームデータと照会
	if pass != fv.Password {
		return "", fmt.Errorf("login failed at compare pass err %w", err)
	}

	// 独自クレーム作成
	claims := &jwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// トークン作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		return "", xerrors.Errorf("login failed at NewWithClaims err %w", err)
	}

	// 日本時間
	time.Local = time.FixedZone("Local", 9*60*60)
	jst, err := time.LoadLocation("Local")
	if err != nil {
		return "", xerrors.Errorf("login failed at LoadLocation err %w", err)
	}
	nowJST := time.Now().In(jst)

	// CookieにJWT格納
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = nowJST.Add(1 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	return user.ID, nil
}

func (a *AuthUseCase) IdentifyJwtUser(id string) error {
	_, err := a.userRepo.GetUserById(id)
	if err != nil {
		return fmt.Errorf("IdentifyJwtUser err %w", err)
	}

	return nil
}

func (a *AuthUseCase) DeleteCookie(c echo.Context, ck *http.Cookie) {
	ck.Value = ""
	ck.MaxAge = -1
	c.SetCookie(ck)
}
