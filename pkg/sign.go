package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

//生成签名
func GenerateSign(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//验证签名
func VerifySign(data, sign, secret string) bool {
	return hmac.Equal([]byte(GenerateSign(data, secret)), []byte(sign))
}

// BuildSignData 构建签名字串
//
// 会自动：
//   - 忽略 nil
//   - 忽略空字符串
//   - Key 按 ASCII 排序
//   - 输出 key=value&key=value
func BuildSignData(values map[string]any) string {

	keys := make([]string, 0, len(values))

	for k, v := range values {

		if v == nil {
			continue
		}

		switch value := v.(type) {

		case string:
			if value == "" {
				continue
			}

		case *string:
			if value == nil || *value == "" {
				continue
			}

		case *uint32:
			if value == nil {
				continue
			}

		case *uint64:
			if value == nil {
				continue
			}

		case *int:
			if value == nil {
				continue
			}

		case *int64:
			if value == nil {
				continue
			}
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)

	items := make([]string, 0, len(keys))

	for _, k := range keys {

		v := values[k]

		switch value := v.(type) {

		// string
		case string:
			items = append(items, fmt.Sprintf("%s=%s", k, value))

		case *string:
			items = append(items, fmt.Sprintf("%s=%s", k, *value))

		// signed int
		case int:
			items = append(items, fmt.Sprintf("%s=%d", k, value))

		case int32:
			items = append(items, fmt.Sprintf("%s=%d", k, value))

		case int64:
			items = append(items, fmt.Sprintf("%s=%d", k, value))

		case *int:
			items = append(items, fmt.Sprintf("%s=%d", k, *value))

		case *int64:
			items = append(items, fmt.Sprintf("%s=%d", k, *value))

		// unsigned int
		case uint:
			items = append(items, fmt.Sprintf("%s=%d", k, value))

		case uint32:
			items = append(items, fmt.Sprintf("%s=%d", k, value))

		case uint64:
			items = append(items, fmt.Sprintf("%s=%d", k, value))

		case *uint32:
			items = append(items, fmt.Sprintf("%s=%d", k, *value))

		case *uint64:
			items = append(items, fmt.Sprintf("%s=%d", k, *value))

		// JSON Number -> float64
		case float64:
			if value == float64(int64(value)) {
				items = append(items, fmt.Sprintf("%s=%d", k, int64(value)))
			} else {
				items = append(items, fmt.Sprintf("%s=%v", k, value))
			}

		case float32:
			if value == float32(int64(value)) {
				items = append(items, fmt.Sprintf("%s=%d", k, int64(value)))
			} else {
				items = append(items, fmt.Sprintf("%s=%v", k, value))
			}

		// bool
		case bool:
			items = append(items, fmt.Sprintf("%s=%t", k, value))

		default:
			items = append(items, fmt.Sprintf("%s=%v", k, value))
		}
	}

	return strings.Join(items, "&")
}
