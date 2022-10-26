package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func FetchGFWlist(proxyUrl string) ([]string, error) {
	var gfwUrl = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"
	uri, _ := url.Parse(fmt.Sprintf("socks5h://%s", proxyUrl))
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	resp, err := client.Get(gfwUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	urlList := strings.Split(string(data), "\n")

	newUrlList := parseGFWlist(urlList)
	return newUrlList, nil
}

func parseGFWlist(urlList []string) []string {
	var newUrlList []string
	for _, line := range urlList {
		//过滤空行
		if len(line) == 0 {
			continue
		}
		//过滤注释和直连的
		re, err := regexp.Compile(`^[!,\[,@].*`)
		if err != nil {
			continue
		}
		match := re.MatchString(line)
		if match {
			continue
		}
		//过滤关键字类型的 不含点的
		re, err = regexp.Compile(`[\.]`)
		if err != nil {
			continue
		}
		match = re.MatchString(line)
		if !match {
			continue
		}
		//过滤链接
		re, err = regexp.Compile(`[\/]`)
		if err != nil {
			continue
		}
		match = re.MatchString(line)
		if match {
			continue
		}

		//转换
		reg, err := regexp.Compile(`^\|{1,2}|^\.`)
		if err != nil {
			continue
		}

		newUrlList = append(newUrlList, reg.ReplaceAllString(line, ""))
	}
	return newUrlList
}
