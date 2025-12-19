package service

import (
	"net/http"
	"net/url"

	"github.com/go-musicfox/netease-music/util"
)

type LogoutService struct {
}

// Logout 注销登录
func (service *LogoutService) Logout() (float64, []byte, error) {
	api := "https://music.163.com/weapi/logout"
	data := make(map[string]interface{})
	cookiejar := util.GetGlobalCookieJar()
	csrfToken := util.GetCsrfToken(cookiejar)
	data["csrf_token"] = csrfToken
	code, bodyBytes, err := util.CallWeapi(api, data)
	if err != nil {
		return code, bodyBytes, err
	}

	// 登出成功后，清除本地 Cookie Jar 中的认证信息
	if cookiejar != nil {
		cookies := []string{"MUSIC_U", "__csrf", "NMTID", "_ntes_nuid", "__remember_me"}
		emptyCookies := make(map[string]string)
		for _, cookieName := range cookies {
			emptyCookies[cookieName] = ""
		}
		util.AddCookiesToJar(cookiejar, emptyCookies, "https://music.163.com")

		// 同时尝试从 cookie jar 中删除这些 cookies
		musicURL, _ := url.Parse("https://music.163.com")
		existingCookies := cookiejar.Cookies(musicURL)
		var validCookies []*http.Cookie

		// 只保留非认证相关的 cookies
		for _, cookie := range existingCookies {
			isAuthCookie := false
			for _, authCookieName := range cookies {
				if cookie.Name == authCookieName {
					isAuthCookie = true
					break
				}
			}
			if !isAuthCookie {
				validCookies = append(validCookies, cookie)
			}
		}

		// 更新 cookie jar
		cookiejar.SetCookies(musicURL, validCookies)
	}

	return code, bodyBytes, nil
}
