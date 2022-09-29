package pkg

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	uu "net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

const (
	APIHost    = "tuanapi.12355.net"
	YouthStudy = "youthstudy.12355.net"
	UserAgent  = "Mozilla/5.0 (Linux; Android 10; Pixel 4 XL Build/QQ2A.200305.004.A1; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3140 MMWEBSDK/20211001 Mobile Safari/537.36 MMWEBID/8391 MicroMessenger/8.0.16.2040(0x2800103A) Process/toolsmp WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64"
)

func NewChapterId() (string, error) {
	url := fmt.Sprintf("https://%s/apih5/api/young/chapter/new", YouthStudy)
	headers := map[string]string{
		"X-Litemall-IdentiFication": "young",
	}
	code, details, err := Get(url, WithHeaders(headers))
	if err != nil || code != http.StatusOK {
		return "", err
	}
	id := gjson.Get(details, "data.entity.id").String()
	return id, nil
}

func Sign(memberId int) (string, error) {
	url := fmt.Sprintf("https://%s/questionnaire/getYouthLearningUrl?mid=%d", APIHost, memberId)
	headers := map[string]string{
		"User-Agent": UserAgent,
	}
	code, details, err := Get(url, WithHeaders(headers))
	if err != nil || code != http.StatusOK {
		return "", err
	}
	if gjson.Get(details, "status").String() == "10002" {
		return "", errors.New("memberId not found")
	}
	val := gjson.Get(details, "youthLearningUrl").String()
	sign := val[strings.Index(val, "=")+1:]
	return sign, nil
}

func Token(sign string) (string, error) {
	url := fmt.Sprintf("https://%s/apih5/api/user/get", YouthStudy)
	reqBody := fmt.Sprintf("sign=%s", uu.QueryEscape(sign))
	headers := map[string]string{
		"Content-Type":              "application/x-www-form-urlencoded",
		"X-Litemall-IdentiFication": "young",
		"User-Agent":                UserAgent,
	}
	code, details, err := Post(url, WithHeaders(headers), WithData(reqBody))
	if err != nil || code != http.StatusOK {
		return "", err
	}
	if gjson.Get(details, "errno").String() == strconv.Itoa(-4) {
		return "", errors.New("error token")
	}
	token := gjson.Get(details, "data.entity.token").String()
	return token, nil
}

func CompleteJPG(fileName, filePath string) {
	headers := map[string]string{
		"X-Litemall-IdentiFication": "young",
		"User-Agent":                UserAgent,
	}
	chapterUrl := fmt.Sprintf("https://%s/apih5/api/young/chapter/new", YouthStudy)
	code, details, err := Get(chapterUrl, WithHeaders(headers))
	if err != nil || code != http.StatusOK {
		log.Println("bad request", err, code)
		return
	}
	magicStr := strings.Split(gjson.Get(details, "data.entity.url").String(), "/")[5]
	url := fmt.Sprintf("https://h5.cyol.com/special/daxuexi/%s/images/end.jpg", magicStr)
	code, details, err = Get(url, WithHeaders(headers))
	if err != nil || code != http.StatusOK {
		log.Println("bad request", err, code)
		return
	}
	name := fmt.Sprintf("%s/%s.jpg", filePath, fileName)
	if err := os.WriteFile(name, []byte(details), 0644); err != nil {
		log.Println("write file error,", err)
		return
	}
}

func Do(token, memberId string) (bool, error) {
	url := fmt.Sprintf("https://%s/apih5/api/young/course/chapter/saveHistory", YouthStudy)
	reqBody := fmt.Sprintf("chapterId=%s", memberId)
	headers := map[string]string{
		"X-Litemall-Token":          token,
		"Content-Type":              "application/x-www-form-urlencoded",
		"X-Litemall-IdentiFication": "young",
		"User-Agent":                UserAgent,
	}
	code, details, err := Post(url, WithHeaders(headers), WithData(reqBody))
	if err != nil || code != http.StatusOK {
		return false, err
	}
	return gjson.Get(details, "errno").String() == strconv.Itoa(0), nil
}

func Push(chYouth chan Member, queue []Member) {
	for _, member := range queue {
		chYouth <- member
	}
	close(chYouth)
}

func CheckStatus(length int, results chan Member) {
	for i := 0; i < length; i++ {
		member := <-results
		if member.Status {
			color.Green(fmt.Sprintf("[成功] %d 完成最新一期学习!", member.MemberId))
		}
	}
	close(results)
}

func Worker(chYouth chan Member, results chan Member) {
	for member := range chYouth {
		sign, err := Sign(member.MemberId)
		if err != nil {
			color.Red("[失败] %d 好像无法获取捏~ 请检查memberId是否正确!!!", member.MemberId)
			results <- member
			continue
		}
		token, err := Token(sign)
		if err != nil || token == "" {
			color.Red("[失败] 获取Token失败 memberId: %d", member.MemberId)
			results <- member
			continue
		}
		chapterId, err := NewChapterId()
		if err != nil {
			color.Red("[失败] 系统错误")
			results <- member
			continue
		}
		status, err := Do(token, chapterId)
		if err != nil || !status {
			color.Red("[失败] %d 学习无效.", member.MemberId)
			results <- member
			continue
		}
		member.Status = true
		results <- member
	}
}
