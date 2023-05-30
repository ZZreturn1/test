package random

import (
	"math/rand"
	"time"
)

var numSeq [10]rune      // 数字序列
var lowerSeq [26]rune    // 小写字母序列
var upperSeq [26]rune    // 大写字母序列
var numLowerSeq [36]rune // 数字和小写字母序列
var numUpperSeq [36]rune // 数字和大写字母序列
var allSeq [62]rune      // 包含数字和字母的序列

func init() {
                // 初始化随机数生成器的种子，使用当前时间的纳秒数作为种子值
	rand.Seed(time.Now().UnixNano()) 

	for i := 0; i < 10; i++ {
                                // 填充数字序列
		numSeq[i] = rune('0' + i) 
	}
	for i := 0; i < 26; i++ {
                                // 填充小写字母序列
		lowerSeq[i] = rune('a' + i) 

                                // 填充大写字母序列
		upperSeq[i] = rune('A' + i) 
	}

                // 将数字序列复制到数字和小写字母序列中
	copy(numLowerSeq[:], numSeq[:])

                // 将小写字母序列复制到数字和小写字母序列中
	copy(numLowerSeq[len(numSeq):], lowerSeq[:])

                // 将数字序列复制到数字和大写字母序列中
	copy(numUpperSeq[:], numSeq[:])

                // 将大写字母序列复制到数字和大写字母序列中
	copy(numUpperSeq[len(numSeq):], upperSeq[:])   

                // 将数字序列复制到包含数字和字母的序列中
	copy(allSeq[:], numSeq[:])

                // 将小写字母序列复制到包含数字和字母的序列中
	copy(allSeq[len(numSeq):], lowerSeq[:])

                // 将大写字母序列复制到包含数字和字母的序列中
	copy(allSeq[len(numSeq)+len(lowerSeq):], upperSeq[:]) 
}

// Seq 生成指定长度的随机序列
func Seq(n int) string {
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
                                // 随机选择序列中的字符
		runes[i] = allSeq[rand.Intn(len(allSeq))] 
	}
	return string(runes)
}