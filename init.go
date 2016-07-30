package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/joyrexus/buckets"
	"github.com/naoina/toml"
)

var (
	cfg config
	bx  *buckets.DB
	lg  = log.New()
)

type config struct {
	DeveloperID string
	AppID       string
	Token       string
	SecruteID   string
	DBPath      string
}

func init() {
	//初始化logger
	f, err := os.Open("/root/weixin/weixin.log")
	if err != nil {
		panic(err)
	}
	lg.Out = f

	//初始化配置
	buf, err := ioutil.ReadFile("/root/weixin/config.toml")
	if err != nil {
		lg.WithField("filepath", "/root/weixin/config.toml").Fatalln("打开配置文件失败")
	}

	if err := toml.Unmarshal(buf, &cfg); err != nil {
		lg.Errorln(err)
	}

	//设置数据库
	bx, err = buckets.Open(cfg.DBPath)
	if err != nil {
		lg.Errorln(err)
	}

}

//验证函数
func validateURL(signature, timestamp, nonce, token string) bool {
	//排序参数并合并
	ss := []string{token, timestamp, nonce}
	sort.Strings(ss)
	s := strings.Join(ss, "")

	//计算sha1的值
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	//比较计算的signature与获取值比较
	if signatureHex := fmt.Sprintf("%x", bs); signatureHex != signature {
		return false
	}
	return true
}
