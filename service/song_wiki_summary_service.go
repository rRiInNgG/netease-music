// song_wiki_summary_service.go

package service

import (
	"net/http"
	"net/url"
)

// SongWikiSummaryService 音乐百科基础信息服务
type SongWikiSummaryService struct {
	client *http.Client
}

// SongWikiSummaryResponse 定义音乐百科基础信息的响应结构
type SongWikiSummaryResponse struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

// Data 响应数据主体
type Data struct {
	Cursor  string  `json:"cursor"`
	Blocks  []Block `json:"blocks"`
	HasMore bool    `json:"hasMore"`
}

// Block 数据块
type Block struct {
	ShowType   string     `json:"showType"`
	ID         string     `json:"id"`
	Channel    string     `json:"channel"`
	UIElement  UIElement  `json:"uiElement"`
	Code       string     `json:"code"`
	Creatives  []Creative `json:"creatives"`
	CanRefresh bool       `json:"canRefresh"`
	HasMore    bool       `json:"hasMore"`
	HideTitle  bool       `json:"hideTitle"`
	OpRcmd     int        `json:"opRcmd"`
}

// UIElement UI元素
type UIElement struct {
	MainTitle MainTitle  `json:"mainTitle"`
	SubTitles []SubTitle `json:"subTitles,omitempty"`
	Images    []Image    `json:"images,omitempty"`
	TextLinks []TextLink `json:"textLinks,omitempty"`
}

// MainTitle 主标题
type MainTitle struct {
	Title       string      `json:"title"`
	TitleImgId  interface{} `json:"titleImgId"`
	TitleImgUrl interface{} `json:"titleImgUrl"`
	Action      interface{} `json:"action"`
}

// SubTitle 副标题
type SubTitle struct {
	Title       string      `json:"title"`
	TitleImgId  interface{} `json:"titleImgId"`
	TitleImgUrl interface{} `json:"titleImgUrl"`
	Action      interface{} `json:"action"`
}

// Image 图片
type Image struct {
	Tag      interface{} `json:"tag"`
	Title    interface{} `json:"title"`
	ImageId  int64       `json:"imageId"`
	ImageUrl string      `json:"imageUrl"`
	Width    int         `json:"width"`
	Height   int         `json:"height"`
	Action   interface{} `json:"action"`
}

// TextLink 文本链接
type TextLink struct {
	Tag  interface{} `json:"tag"`
	Text string      `json:"text"`
	URL  string      `json:"url"`
}

// Creative 创作内容
type Creative struct {
	CreativeType string      `json:"creativeType"`
	UIElement    interface{} `json:"uiElement"`
	Resources    []Resource  `json:"resources"`
}

// Resource 资源
type Resource struct {
	ResourceType string      `json:"resourceType"`
	ResourceId   string      `json:"resourceId"`
	ResourceUrl  interface{} `json:"resourceUrl"`
	UIElement    interface{} `json:"uiElement"`
	Valid        bool        `json:"valid"`
	Alg          string      `json:"alg"`
}

// NewSongWikiSummaryService 创建音乐百科基础信息服务实例
func NewSongWikiSummaryService(client *http.Client) *SongWikiSummaryService {
	if client == nil {
		client = http.DefaultClient
	}
	return &SongWikiSummaryService{
		client: client,
	}
}

// GetSongWikiSummary 获取音乐百科基础信息
// 对应原JS项目中的 /api/song/play/about/block/page 接口
// 返回状态码和响应数据，与项目中其他服务保持一致
func (s *SongWikiSummaryService) GetSongWikiSummary(songId string) (float64, []byte) {
	// 构造请求数据
	data := url.Values{}
	data.Set("songId", songId)

	// 发送请求到 /api/song/play/about/block/page
	resp, err := s.client.PostForm("https://music.163.com/api/song/play/about/block/page", data)
	if err != nil {
		// 错误情况下返回默认值
		return 500, []byte("{}")
	}
	defer resp.Body.Close()

	// 读取响应体
	var result []byte
	if resp.StatusCode == 200 {
		result, err = readAll(resp.Body)
		if err != nil {
			return 500, []byte("{}")
		}
	} else {
		result, _ = readAll(resp.Body)
	}

	return float64(resp.StatusCode), result
}

// readAll 读取所有响应数据的辅助函数
func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	buf := make([]byte, 0, 512)
	for {
		n, err := r.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if len(buf) == cap(buf) {
			// 扩容
			newBuf := make([]byte, len(buf), 2*cap(buf))
			copy(newBuf, buf)
			buf = newBuf
		}
	}
	return buf, nil
}
