package local

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ByBullet/stroxy/logger"
	"github.com/ByBullet/stroxy/util"
	"go.uber.org/zap"
)

// IpLimit
// 国内IP范围加载到程序中的结构
type IpLimit struct {
	limit [][]int
	size  int
}

// Check
// 使用二分查找算法检测IP是否在国内范围
func (l *IpLimit) Check(domain string) bool {
	//获取域名的ip
	cname, err := net.LookupIP(domain)
	if err != nil || len(cname) == 0 {
		log.Println(err, cname)
		return false
	}
	i := ipToInt(cname[0].String())
	left, right := 0, l.size-1

	for left <= right {
		mid := ((right - left) >> 1) + left
		if l.limit[mid][0] == i {
			return true
		}
		if l.limit[mid][0] < i {
			left = mid + 1
			continue
		}
		if l.limit[mid][0] > i {
			right = mid - 1
			continue
		}
	}
	if right < 0 {
		right = 0
	}
	if l.limit[right][0] <= i && l.limit[right][1] >= i {
		return true
	}
	return false
}

var localIpLimit *IpLimit

// 把Ip地址转为int型数字 例："0.0.0.2" -> 2
func ipToInt(ip string) int {
	s := strings.Split(ip, ".")
	if len(s) != 4 {
		return 0
	}
	a, _ := strconv.Atoi(s[0])
	b, _ := strconv.Atoi(s[1])
	c, _ := strconv.Atoi(s[2])
	d, _ := strconv.Atoi(s[3])

	return (a << 24) + (b << 16) + (c << 8) + d
}

// InitLimit
// IpLimit初始化
func InitLimit() {
	localIpLimit = new(IpLimit)
	localIpLimit.limit = make([][]int, 0)
	file, err := os.OpenFile(util.GetFilePath(util.PathIp), os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	rd := bufio.NewReader(file)

	for {
		line, err := rd.ReadString('\n')
		if line == "" {
			break
		}
		line = strings.ReplaceAll(line, "\n", "")
		split := strings.Split(line, " ")
		item := make([]int, 2)
		item[0] = ipToInt(split[0])
		item[1] = ipToInt(split[1])
		localIpLimit.limit = append(localIpLimit.limit, item)
		if io.EOF == err {
			break
		}
	}
	localIpLimit.size = len(localIpLimit.limit)
	logger.PROD().Info("本地ip范围初始化成功", zap.Int("地址范围数", localIpLimit.size))
}
