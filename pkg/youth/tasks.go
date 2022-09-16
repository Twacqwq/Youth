package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	uu "net/url"
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
	client := &http.Client{}
	url := fmt.Sprintf("https://%s/apih5/api/young/chapter/new", YouthStudy)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("failed %v", err)
		return "", err
	}
	req.Header.Add("X-Litemall-IdentiFication", "young")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed response %v", err)
		return "", err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	id := gjson.Get(string(data), "data.entity.id")
	return id.String(), nil
}

func Sign(memberId int) (string, error) {
	url := fmt.Sprintf("https://%s/questionnaire/getYouthLearningUrl?mid=%d", APIHost, memberId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("failed %v", err)
		return "", err
	}
	req.Header.Add("User-Agent", UserAgent)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed response %v", err)
		return "", err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if gjson.Get(string(data), "status").String() == "10002" {
		return "", errors.New("memberId not found")
	}
	val := gjson.Get(string(data), "youthLearningUrl")
	sign := val.String()[strings.Index(val.String(), "=")+1:]
	return sign, nil
}

func Token(sign string) (string, error) {
	url := fmt.Sprintf("https://%s/apih5/api/user/get", YouthStudy)
	client := &http.Client{}
	reqBody := fmt.Sprintf("sign=%s", uu.QueryEscape(sign))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		log.Fatalf("failed %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Litemall-IdentiFication", "young")
	req.Header.Add("User-Agent", UserAgent)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed response %v", err)
		return "", err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if gjson.Get(string(data), "errno").String() == strconv.Itoa(-4) {
		return "", errors.New("error token")
	}
	token := gjson.Get(string(data), "data.entity.token")
	return token.String(), nil
}

func Do(token, memberId string) (bool, error) {
	if token == "" {
		return false, errors.New("token error")
	}
	url := fmt.Sprintf("https://%s/apih5/api/young/course/chapter/saveHistory", YouthStudy)
	client := &http.Client{}
	reqBody := fmt.Sprintf("chapterId=%s", memberId)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		log.Fatalf("failed %v", err)
		return false, nil
	}
	req.Header.Add("X-Litemall-Token", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Litemall-IdentiFication", "young")
	req.Header.Add("User-Agent", UserAgent)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed response %v", err)
		return false, err
	}
	data, _ := io.ReadAll(res.Body)
	return gjson.Get(string(data), "errno").String() == strconv.Itoa(0), nil
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