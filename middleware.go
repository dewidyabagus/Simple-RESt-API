package main

import (
	"net"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

func MiddWhitelistRequest(whitelist Whitelist, allowMethods []string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if whitelist.Enabled {
				whitelistAddr := slices.ContainsFunc(whitelist.IpAddr, func(ip string) bool {
					ipAddr, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
					if strings.Contains(ip, "/") {
						_, subnet, _ := net.ParseCIDR(ip)
						return subnet.Contains(net.ParseIP(ipAddr))
					}
					return ip == strings.TrimSpace(ipAddr)
				})
				if !whitelistAddr {
					return c.JSON(http.StatusForbidden, echo.Map{"message": "Unregistered Ip Address"})
				}
			}

			if !slices.Contains(allowMethods, c.Request().Method) {
				return c.JSON(http.StatusMethodNotAllowed, echo.Map{"message": "Request Method Not Allowed"})
			}

			return next(c)
		}
	}
}
