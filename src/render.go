package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// BlessingResult 抽签结果
type BlessingResult struct {
	BackgroundImage string
	TextImage       string
	Dordas          string
	DordasColor     string
	ColorHex        string
	Blessing        string
	Entry           string
}

// colorMap 颜色映射表
var colorMap = map[string]string{
	"缘彩：丹色":   "#fb5731",
	"缘彩：杏黄":   "#ffa40b",
	"缘彩：琥珀黄":  "#ffcd02",
	"缘彩：麦苗绿":  "#62dba1",
	"缘彩：宫殿绿":  "#21a95e",
	"缘彩：星蓝":   "#a9cfec",
	"缘彩：雪青":   "#9f94c2",
	"缘彩：桃夭":   "#ffc7cd",
	"缘彩：霁青":   "#71d6ed",
	"缘彩：梧枝绿":  "#7fc8b2",
	"缘彩：酡颜":   "#ff9d72",
	"缘彩：明黄":   "#ffdc01",
	"缘彩：梅子青":  "#c5d7bc",
	"缘彩：田园绿":  "#7ad8a6",
	"缘彩：藤萝紫":  "#7c729b",
	"缘彩：藤黄":   "#ffe108",
	"缘彩：霞光红":  "#ff8aa1",
	"缘彩：鸢尾蓝":  "#0ea1d6",
	"缘彩：钴蓝":   "#18abdc",
	"缘彩：栀子黄":  "#ffd005",
	"缘彩：蛙绿":   "#4fd69e",
	"缘彩：石绿":   "#62e1de",
	"缘彩：竹月":   "#7d99a9",
	"缘彩：棟色":   "#b8abe5",
	"缘彩：桃红":   "#ffc3b4",
	"缘彩：蔚蓝":   "#28d1e9",
	"缘彩：蝶翅蓝":  "#5994c0",
	"缘彩：青矾绿":  "#2db691",
	"缘彩：明茶褐":  "#bf9c7d",
	"缘彩：丁香褐":  "#f0ccea",
	"缘彩：春梅红":  "#ff9b9e",
	"缘彩：金盏黄":  "#ffd600",
	"缘彩：天水碧":  "#69c2cc",
	"缘彩：景泰蓝":  "#2588d5",
	"缘彩：银朱":   "#ed4d48",
	"缘彩：苕荣":   "#ff823e",
	"缘彩：酡红":   "#ee3626",
	"缘彩：赤金":   "#ffd547",
	"缘彩：鹅掌黄":  "#ffcb25",
	"缘彩：松霜绿":  "#9cc6a9",
	"缘彩：螺钿紫":  "#8b8cba",
	"缘彩：品蓝":   "#2c88cf",
	"缘彩：豇豆红":  "#ffa8b9",
	"缘彩：雷雨垂":  "#9b9b98",
	"缘彩：缁色":   "#7c696c",
}

// textImageMap 签文图片映射
var textImageMap = map[string]string{
	"大吉": "text0.png",
	"中吉": "text1.png",
	"小吉": "text2.png",
	"吉":  "text3.png",
	"奇":  "text4.png",
}

// backgroundImageMap 背景图片映射
var backgroundImageMap = map[string]string{
	"backgroundimg0": "background0.png",
	"backgroundimg1": "background1.png",
	"backgroundimg2": "background2.png",
	"backgroundimg3": "background3.png",
}

// drawRandomItem 从带权重的项目列表中随机选择一个
func drawRandomItem(items []DrawItem) DrawItem {
	if len(items) == 0 {
		return DrawItem{}
	}

	// 创建权重索引数组
	indices := []int{}
	for i, item := range items {
		for j := 0; j < item.Weight; j++ {
			indices = append(indices, i)
		}
	}

	// 随机打乱
	shuffleArray(indices)

	// 随机选择
	selectedIndex := indices[rng.Intn(len(indices))]
	return items[selectedIndex]
}

// shuffleArray 随机打乱数组
func shuffleArray(arr []int) {
	for i := len(arr) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// getChildItems 获取指定父ID的子项
func getChildItems(parentID string) []DrawItem {
	var children []DrawItem
	for _, item := range drawItems {
		if item.ParentID == parentID {
			children = append(children, item)
		}
	}
	return children
}

// performDraw 执行抽签
func performDraw() *BlessingResult {
	result := &BlessingResult{}

	// 1. 抽取背景图
	bgItems := []DrawItem{}
	for _, item := range drawItems {
		if item.Remark == "backgroundimg" {
			bgItems = append(bgItems, item)
		}
	}
	bgItem := drawRandomItem(bgItems)
	if imgPath, ok := backgroundImageMap[bgItem.Name]; ok {
		result.BackgroundImage = imgPath
	}

	// 2. 抽取签文类型（大吉、中吉等）
	textItems := getChildItems(bgItem.ID)
	textItem := drawRandomItem(textItems)
	if imgPath, ok := textImageMap[textItem.Name]; ok {
		result.TextImage = imgPath
	}

	// 3. 递归抽取下级项
	drawSubItems(textItem.ID, result)

	return result
}

// drawSubItems 递归抽取子项
func drawSubItems(parentID string, result *BlessingResult) {
	children := getChildItems(parentID)
	if len(children) == 0 {
		return
	}

	selectedItem := drawRandomItem(children)

	switch selectedItem.Remark {
	case "dordas":
		result.Dordas = selectedItem.Name
	case "dordascolor":
		result.DordasColor = selectedItem.Name
		if colorHex, ok := colorMap[selectedItem.Name]; ok {
			result.ColorHex = colorHex
		}
	case "blessing":
		result.Blessing = selectedItem.Name
	case "entry":
		result.Entry = selectedItem.Name
	}

	// 继续递归
	drawSubItems(selectedItem.ID, result)
}

// generateBlessingImage 生成祈福签图片
func generateBlessingImage() ([]byte, error) {
	// 执行抽签
	result := performDraw()

	// 如果是 debug 模式，打印抽签结果
	if config.Server.LogLevel == "debug" {
		log.Println("--- 抽签结果 (Debug) ---")
		log.Printf("背景图: %s", result.BackgroundImage)
		log.Printf("签文图: %s", result.TextImage)
		log.Printf("结缘物: %s", result.Dordas)
		log.Printf("缘彩: %s (%s)", result.DordasColor, result.ColorHex)
		log.Printf("祝福语: %s", result.Blessing)
		log.Printf("词条: %s", result.Entry)
		log.Println("--------------------------")
	}

	// 创建画布
	dc := gg.NewContext(config.Image.Width, config.Image.Height)

	

	// 1. 再绘制带颜色的背景（使用遮罩，叠在装饰图上）
	if err := drawColoredBackground(dc, result); err != nil {
		return nil, fmt.Errorf("绘制背景失败: %w", err)
	}

	// 2. 先绘制背景装饰图（底层，无遮罩）
	if err := drawBackgroundImage(dc, result); err != nil {
		return nil, fmt.Errorf("绘制背景图片失败: %w", err)
	}

	// 3. 绘制签文图片
	if err := drawTextImage(dc, result); err != nil {
		return nil, fmt.Errorf("绘制签文图片失败: %w", err)
	}

	// 4. 绘制文字
	if err := drawTexts(dc, result); err != nil {
		return nil, fmt.Errorf("绘制文字失败: %w", err)
	}

	// 转换为 PNG 字节数组
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("编码 PNG 失败: %w", err)
	}

	return buf.Bytes(), nil
}

// parseColor 解析十六进制颜色
func parseColor(hex string) color.RGBA {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 204} // 0.8 * 255 = 204
}

// drawColoredBackground 绘制带颜色和遮罩的背景
func drawColoredBackground(dc *gg.Context, result *BlessingResult) error {
    // 读取遮罩图片
    maskPath := getAssetPath("image/background.png")
    maskFile, err := os.Open(maskPath)
    if err != nil {
        return fmt.Errorf("打开遮罩图片失败: %v", err)
    }
    defer maskFile.Close()

    maskImg, err := png.Decode(maskFile)
    if err != nil {
        return fmt.Errorf("解码遮罩图片失败: %v", err)
    }

    // 获取遮罩图片的alpha通道
    bounds := maskImg.Bounds()
    alphaMask := image.NewAlpha(bounds)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            _, _, _, a := maskImg.At(x, y).RGBA()
            // RGBA()返回16位色，范围0-65535，转换为8位
            alphaMask.SetAlpha(x, y, color.Alpha{uint8(a >> 8)})
        }
    }

    // 解析颜色，alpha固定204
    col := parseColor(result.ColorHex) // 返回color.RGBA，A=204

    // 创建纯色图层
    colorLayer := image.NewRGBA(bounds)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            colorLayer.SetRGBA(x, y, col)
        }
    }

    // 创建一个空白透明图层
    tempCanvas := image.NewRGBA(bounds)

    // 用alphaMask作为mask，将colorLayer绘制到tempCanvas
    draw.DrawMask(tempCanvas, bounds, colorLayer, image.Point{}, alphaMask, image.Point{}, draw.Over)

    // 将tempCanvas合成到底层画布(dc)上
    // gg.Context没有直接操作image.RGBA的方法，但可以用DrawImage
    dc.DrawImage(tempCanvas, 0, 0)

    return nil
}


// drawBackgroundImage 绘制背景图片（装饰纹理）
// 直接绘制，让图片自身的 alpha 通道控制透明度
func drawBackgroundImage(dc *gg.Context, result *BlessingResult) error {
	if result.BackgroundImage == "" {
		return nil
	}

	bgImgPath := getAssetPath("image/" + result.BackgroundImage)
	bgImg, err := gg.LoadImage(bgImgPath)
	if err != nil {
		log.Printf("警告：加载背景图失败: %v", err)
		return nil
	}

	// 直接绘制背景装饰图，其 alpha 通道会自动处理透明度
	dc.DrawImage(bgImg, 0, 0)

	return nil
}

// drawTextImage 绘制签文图片
func drawTextImage(dc *gg.Context, result *BlessingResult) error {
	if result.TextImage == "" {
		return nil
	}

	textImgPath := getAssetPath("image/" + result.TextImage)
	
	textImg, err := gg.LoadImage(textImgPath)
	if err != nil {
		return fmt.Errorf("加载签文图片失败: %w", err)
	}

	// 签文图片位置：根据原版布局，应该在左侧中间位置
	// 原 CSS: #textimg { width: 296px; height: 80px; margin-top: 10px; margin-right: 80px; }
	// 放在画布水平中心偏左的位置
	x := float64(config.Image.Width)*0.204
	y := float64(config.Image.Height)*0.49

	dc.DrawImage(textImg, int(x), int(y))
	return nil
}

// drawTexts 绘制文字内容
func drawTexts(dc *gg.Context, result *BlessingResult) error {
	// 加载字体
	fontPath := getAssetPath("font/LXGWWenKaiMono-Medium.ttf")
	fontFace, err := loadFont(fontPath, float64(config.Image.FontSize))
	if err != nil {
		return fmt.Errorf("加载字体失败: %w", err)
	}
	dc.SetFontFace(fontFace)
	dc.SetColor(color.White)

	// --- 布局参数 ---
	// 文字区域靠右对齐，宽度约 45%
	textWidth := float64(config.Image.Width) * 0.45
	textX := float64(config.Image.Width) - textWidth - 50 // 右边距 50
	lineSpacing := 1.8 // 1.8 倍行距，增加间距

	// --- 准备要绘制的文字 ---
	lines := []string{
		result.Dordas,
		result.DordasColor,
		"", // 空行作为间距
		result.Blessing,
		"", // 空行作为间距
		result.Entry,
	}
	
	// 过滤掉空字符串
	var contentLines []string
	for _, line := range lines {
		if line != "" {
			contentLines = append(contentLines, line)
		} else if len(contentLines) > 0 && contentLines[len(contentLines)-1] != "" {
			// 只在非空行后添加一个空行作为间距
			contentLines = append(contentLines, "")
		}
	}
	// 如果最后一行是空行，则移除
	if len(contentLines) > 0 && contentLines[len(contentLines)-1] == "" {
		contentLines = contentLines[:len(contentLines)-1]
	}

	// --- 先测量总高度，再计算起始 Y 坐标 ---
	totalHeight := 0.0
	for _, line := range contentLines {
		if line == "" {
			totalHeight += float64(config.Image.FontSize) * (lineSpacing / 2) // 空行高度小一点
		} else {
			w, h := dc.MeasureString(line)
			// 如果文字超宽，需要计算换行后的高度
			if w > textWidth {
				wrappedLines := dc.WordWrap(line, textWidth)
				totalHeight += float64(len(wrappedLines)) * float64(config.Image.FontSize) * lineSpacing
			} else {
				totalHeight += h * lineSpacing
			}
		}
	}
	totalHeight -= float64(config.Image.FontSize) * (lineSpacing - 1) // 减去最后一行多余的行距

	// 计算起始 Y 坐标，使其垂直居中，并增加一个向下的偏移量
	startY := (float64(config.Image.Height)-totalHeight)/2 + 20.0

	// --- 开始绘制 ---
	currentY := startY
	for _, line := range contentLines {
		if line == "" {
			currentY += float64(config.Image.FontSize) * (lineSpacing / 2)
			continue
		}
		dc.DrawStringWrapped(line, textX, currentY, 0, 0, textWidth, lineSpacing, gg.AlignLeft)
		
		// 更新 Y 坐标
		w, _ := dc.MeasureString(line)
		if w > textWidth {
			wrappedLines := dc.WordWrap(line, textWidth)
			currentY += float64(len(wrappedLines)) * float64(config.Image.FontSize) * lineSpacing
		} else {
			currentY += float64(config.Image.FontSize) * lineSpacing
		}
	}

	return nil
}

// loadFont 加载字体
func loadFont(fontPath string, size float64) (font.Face, error) {
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}
