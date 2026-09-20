package renderer

import (
	"admin-dashboard/cmd/internal/session"
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func loadTemplates(root string) *template.Template {

	tmpl := template.New("")

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".gohtml" {
			return nil
		}

		name, err := filepath.Rel(root, path)

		if err != nil {
			return err
		}

		b, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		_, err = tmpl.New(filepath.ToSlash(name)).Parse(string(b))

		return err
	})

	if err != nil {
		panic(err)
	}

	if len(tmpl.Templates()) == 0 {
		panic("no templates found in " + root)
	}

	return tmpl
}

func NewTemplate() (*Template, error) {

	t := loadTemplates("web/html")

	if t == nil {
		return nil, errors.New("failed to load templates")
	}

	return &Template{
		templates: t,
	}, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if t.templates == nil {

		return errors.New("templates not initialized")
	}

	userSess, _, _ := session.GetUserSession(c)

	routeActionDetails := userSess.RouteActionDetails
	categoryRouteAction := userSess.CategoryRouteAction
	routeAction := userSess.RouteAction

	viewData := map[string]any{}

	switch v := data.(type) {
	case map[string]interface{}:
		viewData = v
	default:
		viewData["Data"] = data
	}

	viewData["RouteActionDetails"] = routeActionDetails
	viewData["CategoryRouteAction"] = categoryRouteAction
	viewData["RouteAction"] = routeAction

	if _, exists := viewData["ConfigData"]; !exists || viewData["ConfigData"] == nil {
		if c != nil {
			if cfgVal := c.Get("config"); cfgVal != nil {
				viewData["ConfigData"] = cfgVal
			}
		}
	}

	var csrfToken string

	if token := c.Get("csrf"); token != nil {
		if s, ok := token.(string); ok {
			csrfToken = s
		}
	} else if token := c.Get(middleware.DefaultCSRFConfig.ContextKey); token != nil {
		if s, ok := token.(string); ok {
			csrfToken = s
		}
	}

	if csrfToken != "" {
		c.Response().Header().Set("X-CSRF-Token", csrfToken)
		viewData["_csrf"] = csrfToken
	}
	
	return t.templates.ExecuteTemplate(w, filepath.ToSlash(name), viewData)
}
