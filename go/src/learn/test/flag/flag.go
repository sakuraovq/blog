package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

// 自定义类型
type interval []time.Duration

func (i *interval) String() string {
	return fmt.Sprintf("%v,emmmm", *i)
}

func (i *interval) Set(str string) error {

	if len(*i) > 0 {
		return errors.New("interval flag already set")
	}
	for _, duration := range strings.Split(str, ",") {
		parseDuration, err := time.ParseDuration(duration)
		if err != nil {
			return err
		}
		*i = append(*i, parseDuration)
	}
	return nil
}

var intervalFlags interval

func init() {
	//自定义类型解析命名行参数
	flag.Var(&intervalFlags, "intervals", "please input intervals")

}

func main() {
	flag.Parse()

	fmt.Println(intervalFlags.String())
}
