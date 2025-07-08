package handlers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const hashedAdminPassword = "$2a$10$hJNTP1RBLtUAZzkrTKKL4uxg2crobGP8dOge4k940GEAhLvS4quEC"

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) ShowAdminLoginPage(c *gin.Context) {
	session := sessions.Default(c)

	errMsg := verifyLoginAvailability(session)

	c.HTML(http.StatusOK, "admin_login.html", errMsg)
}

func verifyLoginAvailability(session sessions.Session) gin.H {
	lastFailed, _ := session.Get("lastFailedAttempt").(int64)
	if time.Since(time.Unix(lastFailed, 0)) > 5*time.Minute {
		session.Delete("loginAttempts")
		session.Delete("lastFailedAttempt")
		session.Delete("lockedOut")
		session.Save()
		return gin.H{}
	}

	attempts, _ := session.Get("loginAttempts").(int)
	if attempts >= 3 {
		session.Set("lockedOut", true)
		session.Save()
		return gin.H{
			"error":     "Доступ временно заблокирован после 3 неверных попыток.",
			"lockedOut": true,
		}
	}

	flashes := session.Flashes()
	session.Save()

	var errMsg string
	if len(flashes) > 0 {
		errMsg = flashes[0].(string)
	} else {
		return gin.H{}
	}
	return gin.H{
		"error": errMsg,
	}
}

func (h *AdminHandler) HandleAdminLogin(c *gin.Context) {
	session := sessions.Default(c)
	if locked, ok := session.Get("lockedOut").(bool); ok && locked {
		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	password := c.PostForm("password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedAdminPassword), []byte(password))
	if err != nil {
		attempts, _ := session.Get("loginAttempts").(int)
		session.Set("lastFailedAttempt", time.Now().Unix())
		session.Set("loginAttempts", attempts+1)

		session.AddFlash("Неверный пароль. Попробуйте снова.")
		session.Save()

		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	session.Set("isAdmin", true)
	session.Delete("loginAttempts")
	session.Delete("lastFailedAttempt")
	session.Delete("lockedOut")
	session.Save()

	c.Redirect(http.StatusFound, "/admin/home")
}

func (h *AdminHandler) ShowAdminHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_home.html", nil)
}
