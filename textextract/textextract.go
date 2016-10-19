package textextract

import (
	"regexp"
	"strings"
)

const (
	BLOCKSWIDTH = 3
	/* 当待抽取的网页正文中遇到成块的新闻标题未剔除时，只要增大此阈值即可。*/
	/* 阈值增大，准确率提升，召回率下降；值变小，噪声会大，但可以保证抽到只有一句话的正文 */
	THRESHOLD = 86
)

type ExtractServer struct {
	source    string
	threshold int
}

type blockInfo struct {
	indexs    []int
	maxIndex  int
	threshold int
}

func NewExtract(source string) *ExtractServer {
	return &ExtractServer{source, THRESHOLD}
}

//	设置threshold值
func (this *ExtractServer) SetThreshold(threshold int) *ExtractServer {
	this.threshold = threshold
	return this
}

//	提取标题
func (this *ExtractServer) ExtractTitle() string {
	s := regexp.MustCompile("(?is)<title>(.*?)</title>").FindString(this.source)
	return regexp.MustCompile("</?title>").ReplaceAllString(s, "")
}

//	提取正文
func (this *ExtractServer) ExtractText() string {
	source := cleanData(this.source)
	lines := removeAllSpace(source)
	if len(lines) < BLOCKSWIDTH {
		return ""
	}

	blockLenIndexs, maxIndex := countBlockInfo(lines)
	b := &blockInfo{blockLenIndexs, maxIndex, this.threshold}
	startIndex := b.findStart()
	endIndex := b.findEnd()
	text := ""
	for i := startIndex; i <= endIndex; i++ {
		text += lines[i] + "\n"
	}

	return text

}

//	删除掉文本集合中一行中的空白
func removeAllSpace(source string) []string {
	lines := strings.Split(source, "\n")
	reg := regexp.MustCompile("\\s+")
	for i, line := range lines {
		lines[i] = reg.ReplaceAllString(line, "")
	}
	return lines
}

//	块 长度统计信息及最大索引
func countBlockInfo(lines []string) ([]int, int) {
	var wordsNum int
	blockLenIndexs := make([]int, len(lines)-BLOCKSWIDTH+1)
	for i := 0; i < len(lines)-BLOCKSWIDTH; i++ {
		wordsNum = 0
		for j := i; j < i+BLOCKSWIDTH; j++ {
			wordsNum += len(lines[j])
		}
		blockLenIndexs[i] = wordsNum
	}

	//	长度最大的块的索引
	count := len(blockLenIndexs)
	maxIndex := 0

	for k, v := range blockLenIndexs {
		if v > count {
			count = v
			maxIndex = k
		}
	}

	return blockLenIndexs, maxIndex
}

//	从maxIndex向前寻找起点
func (b *blockInfo) findStart() int {
	i := b.maxIndex - 1
	for ; i >= 0; i-- {
		if b.indexs[i] < b.threshold {
			break
		}
	}

	return i + 1
}

//	从maxIndex向后寻找终点
func (b *blockInfo) findEnd() int {
	i := b.maxIndex + 1

	for ; i < len(b.indexs); i++ {
		if b.indexs[i] < b.threshold {
			break
		}
	}
	return i - 1

}

//	清洗掉非正文的数据
func cleanData(source string) string {
	source = regexp.MustCompile("(?is)<!DOCTYPE.*?>").ReplaceAllString(source, "")
	source = regexp.MustCompile("(?is)<!--.*?-->").ReplaceAllString(source, "")
	source = regexp.MustCompile("(?is)<script.*?>.*?</script>").ReplaceAllString(source, "")
	source = regexp.MustCompile("(?is)<style.*?>.*?</style>").ReplaceAllString(source, "")
	source = regexp.MustCompile("(?is)<.*?>").ReplaceAllString(source, "")         //	remove HTML Tags
	source = regexp.MustCompile("&.{2,5};|&#.{2,5};").ReplaceAllString(source, "") //	remove some special charcaters
	source = regexp.MustCompile("\r\n").ReplaceAllString(source, "\n")
	return source
}
