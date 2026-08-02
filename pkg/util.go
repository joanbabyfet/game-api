package pkg

import (
	"game-api/internal/config"
	"strings"
	"time"
)

const (
	DefaultPage     = 1   // 第几页
	DefaultPageSize = 20  // 每页几条
	MaxPageSize     = 100 // 最多几条

	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04:05"
)

var defaultLocation = time.Local

// SetLocation 设置默认时区
func SetLocation(name string) error {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return err
	}
	defaultLocation = loc
	return nil
}

// Now 当前时间
func Now() time.Time {
	return time.Now().In(defaultLocation)
}

func NowUTC() time.Time {
	return time.Now().UTC()
}

// Timestamp 当前 Unix 时间戳（秒）
func Timestamp() int64 {
	return time.Now().Unix()
}

// TimestampMilli 当前 Unix 时间戳（毫秒）
func TimestampMilli() int64 {
	return time.Now().UnixMilli()
}

func Format(t time.Time) string {
	return t.In(defaultLocation).Format(DateTimeLayout)
}

func FormatDate(t time.Time) string {
	return t.In(defaultLocation).Format(DateLayout)
}

// UnixToDateTime 时间戳转日期时间
func UnixToDateTime(ts int64) string {
	return Format(time.Unix(ts, 0))
}

// UnixToDate 时间戳转日期
func UnixToDate(ts int64) string {
	return FormatDate(time.Unix(ts, 0))
}

// DateTimeToUnix 日期时间转时间戳
func DateTimeToUnix(str string) (int64, error) {
	t, err := time.ParseInLocation(DateTimeLayout, str, defaultLocation)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// DateToUnix 日期转时间戳（00:00:00）
func DateToUnix(str string) (int64, error) {
	t, err := time.ParseInLocation(DateLayout, str, defaultLocation)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// Page 规范分页参数
func Page(page, pageSize int) (int, int) {

	if page <= 0 {
		page = DefaultPage
	}

	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return page, pageSize
}

// Offset 返回分页偏移量
func Offset(page, pageSize int) int {

	page, pageSize = Page(page, pageSize)

	return (page - 1) * pageSize
}

// FileURL 返回静态资源完整地址
func FileURL(path string) string {
	if path == "" {
		return ""
	}

	base := strings.TrimRight(config.Cfg.Provider.FileURL, "/")
	path = strings.TrimLeft(path, "/")

	return base + "/" + path
}

// ImageURL 返回图片完整地址
func ImageURL(filename string) string {
	if filename == "" {
		return ""
	}

	return FileURL("image/" + filename)
}

// VideoURL 返回视频完整地址
func VideoURL(filename string) string {
	if filename == "" {
		return ""
	}

	return FileURL("video/" + filename)
}

// 生成业务ID
func GenBizID(prefix string) string {
	if snowflake == nil {
		panic("snowflake not initialized")
	}

	return prefix + snowflake.NextIDString()
}

// 生成注单号
func GenOrderNo() string {
	return GenBizID("S")
}

// 生成局号
func GenRoundID() string {
	return GenBizID("R")
}

// 生成 FreeSpin 批次号
func GenFreeSpinID() string {
	return GenBizID("FS")
}

// 生成回滚单号
func GenRollbackNo() string {
	return GenBizID("RB")
}

// GenTransferNo 生成转账单号
func GenTransferNo() string {
	return GenBizID("T")
}
