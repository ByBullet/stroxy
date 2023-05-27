package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ByBullet/stroxy/logger"

	"github.com/PuerkitoBio/goquery"
)

const ipRangeWebSite = "https://ip.bczs.net/country/CN"

// 从ipRangeWebSite上爬取中国ip地址范围，并保存到本地ip.txt中
func GotIpRange() error {
	res, err := http.Get(ipRangeWebSite)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("国内ip范围地址爬取失败, http错误%d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(GetFilePath(PathIp))
	if err != nil {
		return err
	}
	defer f.Close()

	doc.Find("thead tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			logger.PROD().Debug(s.Text())
		}
	})

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		build := strings.Builder{}
		s.Find("td").Each(func(j int, s *goquery.Selection) {
			switch j {
			case 0:
				build.WriteString(s.Text())
				build.WriteByte(' ')
			case 1:
				build.WriteString(s.Text())
				build.WriteByte('\n')
			case 2:
				return
			}
		})
		f.Write([]byte(build.String()))
		build.Reset()
	})
	return nil
}
