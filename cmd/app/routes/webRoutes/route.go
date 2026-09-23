package webRoute

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func MainRoute(g *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler
	cfg := handler.Config

	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.WebRateLimit))))

	g.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         cfg.WebXSSProtection,
		ContentTypeNosniff:    cfg.WebContentTypeNosniff,
		XFrameOptions:         cfg.WebXFrameOptions,
		HSTSMaxAge:            cfg.WebHSTSMaxAge,
		HSTSExcludeSubdomains: cfg.WebHSTSExcludeSubdomains,
		HSTSPreloadEnabled:    cfg.WebHSTSPreloadEnabled,
		ContentSecurityPolicy: cfg.WebContentSecurityPolicy,
		ReferrerPolicy:        cfg.WebReferrerPolicy,
		CSPReportOnly:         cfg.WebCSPReportOnly,
	}))

	web := g.Group("")

	web.GET("health-check", handler.HealthCheck)
	web.GET("file/*", handler.GetFile)
	web.GET("404-not-found", handler.NotFoundPage)
	web.GET("500-server-error", handler.ServerErrorPage)
	web.GET("privacy-policy", handler.PrivacyPolicy)

	web.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieName:     cfg.WebCookieName,
		CookiePath:     cfg.WebCookiePath,
		CookieDomain:   cfg.WebCookieDomain,
		CookieHTTPOnly: cfg.WebCookieHTTPOnly,
		CookieSecure:   cfg.WebCookieSecure,
		CookieSameSite: helpers.StrToSameSite(cfg.WebCookieSameSite),
		CookieMaxAge:   cfg.WebCookieMaxAge,
		TokenLookup:    cfg.WebCSRFTokenLookup,
		TokenLength:    uint8(cfg.WebCSRFTokenLength),
		ErrorHandler: func(err error, c echo.Context) error {
			return handler.CSRFError(c)
		},
	}))

	LoginRoutes(web, appData)
	UserRoutes(web, appData)
	ProfileRoutes(web, appData)
	DashboardRoutes(web, appData)
	SettingRoutes(web, appData)
	ReferenceRoutes(web, appData)
	CommonRoutes(web, appData)
}
