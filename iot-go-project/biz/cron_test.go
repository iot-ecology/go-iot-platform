package biz

import (
	"igp/glob"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestPreviousCronExecutionTime(t *testing.T) {
	// 定义一个 6 位 cron 表达式，格式为：秒 分 时 日 月 星期几
	schedule := "1 0/1 * * * *"

	// 使用 cron.Parse 函数解析 cron 表达式
	c, _ := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(schedule)

	// 可以进一步使用 c.Next(time.Now()) 来获取下一次执行的时间
	nextTime := c.Next(time.Now())
	glob.GLog.Sugar().Infof("Next run time: %+v", nextTime)
}
