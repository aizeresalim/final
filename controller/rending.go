package controller

import (
	"github.com/gofiber/fiber/v2"
	"html/template"
	"log"
)

type Templates struct {
	register   *template.Template
	login      *template.Template
	allPost    *template.Template
	createBlog *template.Template
}

var TemplatesInstance *Templates

// LoadTemplates loads and parses the HTML html
func LoadTemplates() {
	TemplatesInstance = &Templates{}

	// Parse the register template
	registerTmpl, err := template.ParseFiles("ui/html/register.tmpl")
	if err != nil {
		log.Fatalf("Error parsing register template: %v", err)
	}
	TemplatesInstance.register = registerTmpl

	// Parse the login template
	loginTmpl, err := template.ParseFiles("ui/html/login.tmpl")
	if err != nil {
		log.Fatalf("Error parsing login template: %v", err)
	}
	TemplatesInstance.login = loginTmpl

	allPostTmpl, err := template.ParseFiles("ui/html/posts.tmpl")
	if err != nil {
		log.Fatalf("Error parsing login template: %v", err)
	}
	TemplatesInstance.allPost = allPostTmpl

	creatBlogTmpl, err := template.ParseFiles("ui/html/create_blog.tmpl")
	if err != nil {
		log.Fatalf("Error parsing login template: %v", err)
	}
	TemplatesInstance.createBlog = creatBlogTmpl
}

func RenderRegisterPage(c *fiber.Ctx) error {
	c.Type("html")
	// Render the register template
	return TemplatesInstance.register.Execute(c.Response().BodyWriter(), nil)
}

func RenderLoginPage(c *fiber.Ctx) error {
	// Set the Content-Type header
	c.Type("html")

	// Render the login template
	return TemplatesInstance.login.Execute(c.Response().BodyWriter(), nil)
}
func RenderAllPostPage(c *fiber.Ctx) error {
	// Set the Content-Type header
	c.Type("html")

	// Render the login template
	return TemplatesInstance.allPost.Execute(c.Response().BodyWriter(), nil)
}
func RenderCreateBlogPage(c *fiber.Ctx) error {
	// Set the Content-Type header
	c.Type("html")

	// Render the login template
	return TemplatesInstance.createBlog.Execute(c.Response().BodyWriter(), nil)
}
