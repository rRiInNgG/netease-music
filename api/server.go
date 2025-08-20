package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-musicfox/netease-music/service"
)

// Server 是 HTTP 服务器结构体
type Server struct {
	httpServer *http.Server
}

// NewServer 创建新的 HTTP 服务器
func NewServer() *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
			Handler: nil,
		},
	}
}

// Start 启动 HTTP 服务器
func (s *Server) Start() error {
	// 设置路由
	s.setupRoutes()

	// 设置默认端口
	if s.httpServer.Addr == ":" {
		s.httpServer.Addr = ":3000"
	}

	log.Printf("Starting server on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop 停止 HTTP 服务器
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 登录相关接口
	http.HandleFunc("/api/captcha/sent", captchaSentHandler)
	http.HandleFunc("/api/login/cellphone", loginCellphoneHandler)
	http.HandleFunc("/api/login/qr/key", qrKeyHandler)
	http.HandleFunc("/api/login/qr/check", qrCheckHandler)
	http.HandleFunc("/api/login/status", loginStatusHandler)
	http.HandleFunc("/api/logout", logoutHandler)
	http.HandleFunc("/login/qr", qrLoginPageHandler)

	// 用户相关接口
	http.HandleFunc("/api/user/detail", userDetailHandler)
	http.HandleFunc("/api/user/playlist", userPlaylistHandler)

	// 音乐相关接口
	http.HandleFunc("/api/song/url", songUrlHandler)
	http.HandleFunc("/api/song/detail", songDetailHandler)

	// 歌单相关接口
	http.HandleFunc("/api/playlist/detail", playlistDetailHandler)

	// 推荐相关接口
	http.HandleFunc("/api/recommend/resource", recommendResourceHandler)
	http.HandleFunc("/api/recommend/songs", recommendSongsHandler)

	// 搜索相关接口
	http.HandleFunc("/api/search", searchHandler)

	// banner
	http.HandleFunc("/api/banner", bannerHandler)
}

// setCORSHeaders 设置跨域请求头
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// captchaSentHandler 发送验证码
func captchaSentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setCORSHeaders(w)

	phone := r.URL.Query().Get("phone")
	if phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	service := service.CaptchaSentService{Cellphone: phone}
	_, result := service.CaptchaSent()

	log.Printf("[OK] /api/captcha/sent?phone=%s", phone)
	w.Write(result)
}

// loginCellphoneHandler 手机号登录
func loginCellphoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setCORSHeaders(w)

	phone := r.URL.Query().Get("phone")
	password := r.URL.Query().Get("password")
	captcha := r.URL.Query().Get("captcha")

	// 密码登录
	if password != "" {
		service := service.LoginCellphoneService{
			Phone:    phone,
			Password: password,
		}
		_, result := service.LoginCellphone()
		log.Printf("[OK] /api/login/cellphone?phone=%s&password=***", phone)
		w.Write(result)
		return
	}

	// 验证码登录
	if captcha != "" {
		service := service.CaptchaVerifyService{
			Cellphone: phone,
			Captcha:   captcha,
		}
		_, result := service.CaptchaVerify()
		log.Printf("[OK] /api/login/cellphone?phone=%s&captcha=%s", phone, captcha)
		w.Write(result)
		return
	}

	http.Error(w, "Password or captcha is required", http.StatusBadRequest)
}

// qrKeyHandler 获取二维码登录密钥
func qrKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setCORSHeaders(w)

	service := service.LoginQRService{}
	_, _, _ = service.GetKey()

	// 构造包含unikey的响应
	response := map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"unikey": service.UniKey,
		},
	}

	// 将响应转换为JSON
	respBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[OK] /api/login/qr/key")
	w.Write(respBytes)
}

// qrCheckHandler 检查二维码扫描状态
func qrCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setCORSHeaders(w)

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	service := service.LoginQRService{UniKey: key}
	_, result := service.CheckQR()

	log.Printf("[OK] /api/login/qr/check?key=%s", key)
	w.Write(result)
}

// qrLoginPageHandler 二维码登录页面
func qrLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	// 从文件中读取HTML模板
	htmlContent, err := ioutil.ReadFile("api/templates/qr_login.html")
	if err != nil {
		http.Error(w, "无法加载登录页面", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlContent)
}

// loginStatusHandler 登录状态
func loginStatusHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	service := service.UserAccountService{}
	_, result := service.AccountInfo()

	log.Printf("[OK] /api/login/status")
	w.Write(result)
}

// logoutHandler 退出登录
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	service := service.LogoutService{}
	_, result := service.Logout()

	log.Printf("[OK] /api/logout")
	w.Write(result)
}

// userDetailHandler 用户详情
func userDetailHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "uid is required", http.StatusBadRequest)
		return
	}

	service := service.UserDetailService{Uid: uid}
	_, result := service.UserDetail()

	log.Printf("[OK] /api/user/detail?uid=%s", uid)
	w.Write(result)
}

// userPlaylistHandler 用户歌单
func userPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "uid is required", http.StatusBadRequest)
		return
	}

	service := service.UserPlaylistService{Uid: uid}
	_, result := service.UserPlaylist()

	log.Printf("[OK] /api/user/playlist?uid=%s", uid)
	w.Write(result)
}

// songUrlHandler 歌曲播放地址
func songUrlHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	service := service.SongUrlService{ID: id}
	_, result := service.SongUrl()

	log.Printf("[OK] /api/song/url?id=%s", id)
	w.Write(result)
}

// songDetailHandler 歌曲详情
func songDetailHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	ids := r.URL.Query().Get("ids")
	if ids == "" {
		http.Error(w, "ids is required", http.StatusBadRequest)
		return
	}

	service := service.SongDetailService{Ids: ids}
	_, result := service.SongDetail()

	log.Printf("[OK] /api/song/detail?ids=%s", ids)
	w.Write(result)
}

// playlistDetailHandler 歌单详情
func playlistDetailHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	service := service.PlaylistDetailService{Id: id}
	_, result := service.PlaylistDetail()

	log.Printf("[OK] /api/playlist/detail?id=%s", id)
	w.Write(result)
}

// recommendResourceHandler 推荐资源
func recommendResourceHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	service := service.RecommendResourceService{}
	_, result := service.RecommendResource()

	log.Printf("[OK] /api/recommend/resource")
	w.Write(result)
}

// recommendSongsHandler 推荐歌曲
func recommendSongsHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	service := service.RecommendSongsService{}
	_, result := service.RecommendSongs()

	log.Printf("[OK] /api/recommend/songs")
	w.Write(result)
}

// searchHandler 搜索
func searchHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	keywords := r.URL.Query().Get("keywords")
	if keywords == "" {
		http.Error(w, "keywords is required", http.StatusBadRequest)
		return
	}

	service := service.SearchService{S: keywords}
	_, result := service.Search()

	log.Printf("[OK] /api/search?keywords=%s", keywords)
	w.Write(result)
}

// bannerHandler 首页横幅
func bannerHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	service := service.BannerService{}
	_, result := service.Banner()

	log.Printf("[OK] /api/banner")
	w.Write(result)
}

func main() {
	// 创建 API 服务器
	server := NewServer()
	server.httpServer.Addr = ":8080" // 在这里设置端口

	// 设置信号处理
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received interrupt signal, shutting down...")

		// 创建一个5秒的超时上下文
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Stop(shutdownCtx); err != nil {
			log.Printf("Server shutdown failed: %v\n", err)
		}
	}()

	// 启动服务器
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
