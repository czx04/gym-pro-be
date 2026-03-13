package main

import (
	"gym-pro-2026-ptit/internal/bootstrap"
	_ "gym-pro-2026-ptit/docs"
)

// @title Gym Pro API
// @version 1.0
// @description Backend API for Gym Pro - A fitness tracking mobile application
// @contact.name API Support
// @contact.email support@gympro.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	app := bootstrap.NewApp()
	app.Run()
}
