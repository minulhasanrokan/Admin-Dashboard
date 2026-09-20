package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) PrivacyPolicy(c echo.Context) error {

	termsPrivacy, err := h.Service.PrivacyPolicy()

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	data := map[string]interface{}{
		"TermsPrivacy": termsPrivacy,
	}

	return h.RenderWithErrorCheck(c, data, err, "privacy-policy")
}
