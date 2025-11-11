package http_filters

import (
	"avileads-web/utils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func AuthFilter() beego.FilterFunc {
	return func(ctx *context.Context) {
		tokenString := ctx.GetCookie("auth")

		if len(tokenString) <= 0 {
			redirect(ctx)
			return
		}

		user, err := utils.CurrentUserJWT(tokenString)
		if err != nil {
			redirect(ctx)
			return
		}

		ctx.Input.SetData("userName", user.UserName)
		ctx.Input.SetData("userRules", user.Rules)
		ctx.Input.SetData("login", "hidden")
	}
}

func redirect(ctx *context.Context) {
	if !strings.HasPrefix(ctx.Input.URL(), "/login") && ctx.Input.URL() != "/" {
		ctx.Redirect(302, "/login")
	}

	ctx.Input.SetData("logout", "hidden")
}
