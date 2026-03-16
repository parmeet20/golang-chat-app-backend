package user

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, controller *UserController) chi.Router {

	r.Route("/users", func(r chi.Router) {

		r.Post("/register", controller.Register)
		r.Post("/login", controller.Login)

		r.Group(func(r chi.Router) {
			r.Use(controller.authService.JwtMiddleware)

			r.Get("/me", controller.GetMeByToken)
		})
	})

	return r
}