/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/26 4:51 下午
 * @Desc: TODO
 */

package internal

import (
	"math/rand"
	"time"
)

var seedStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandStr generate a string of specified length.
func RandStr(length int) (lastStr string) {
	rand.Seed(time.Now().UnixNano())
	
	pos, seedLen := 0, len(seedStr)
	for i := 0; i < length; i++ {
		pos = rand.Intn(seedLen)
		lastStr += seedStr[pos : pos+1]
	}
	
	return lastStr
}
