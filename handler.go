package main

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func welcome(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "Request Success"})
}

func accounts(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.QueryParam("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}
		var account Account
		if err := db.Debug().Find(&account, "id = ?", id).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, account)
	}
}
