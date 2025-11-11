package http_filters

import (
	"avileads-web/http/routers"
	"avileads-web/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func GuardFilter() beego.FilterFunc {
	return func(ctx *context.Context) {
		// support only GET methods
		if ctx.Request.Method == "GET" {
			tokenString := ctx.GetCookie("auth")

			if len(tokenString) <= 0 {
				return
			}

			user, err := utils.CurrentUserJWT(tokenString)
			if err != nil {
				return
			}

			path := ctx.Request.URL.Path
			if !routers.HasAccessToPath(path, user.Rules) {
				ctx.Redirect(302, "/")
				return
			}
		}
	}
}
