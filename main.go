package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	fontKai *truetype.Font // 字体
	fontTtf *truetype.Font // 字体
)

func main() {
	// 根据路径打开模板文件
	templateFile, err := os.Open("/Users/robinfan/gowork/sample/notex2-img/back.png")
	if err != nil {
		panic(err)
	}
	defer templateFile.Close()
	// 解码
	templateFileImage, err := png.Decode(templateFile)
	if err != nil {
		panic(err)
	}
	// 新建一张和模板文件一样大小的画布
	newTemplateImage := image.NewRGBA(templateFileImage.Bounds())
	// 将模板图片画到新建的画布上
	draw.Draw(newTemplateImage, templateFileImage.Bounds(), templateFileImage, templateFileImage.Bounds().Min, draw.Over)

	// 加载字体文件  这里我们加载两种字体文件
	fontKai, err = loadFont("/Users/robinfan/gowork/sample/notex2-img/fonts/FZRunYYSJW.TTF")
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	fontTtf, err = loadFont("/Users/robinfan/gowork/sample/notex2-img/fonts/FZZHUNYSJW.TTF")
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	// 向图片中写入文字
	writeWord2Pic(newTemplateImage)
	saveFile(newTemplateImage)
}

func writeWord2Pic(newTemplateImage *image.RGBA) {
	// 在写入之前有一些准备工作
	content := freetype.NewContext()
	content.SetClip(newTemplateImage.Bounds())
	content.SetDst(newTemplateImage)
	content.SetSrc(image.Black) // 设置字体颜色
	content.SetDPI(72)          // 设置字体分辨率

	content.SetFontSize(120) // 设置字体大小
	content.SetFont(fontKai) // 设置字体样式，就是我们上面加载的字体
	DrawText(content, "劝学", 200, 200, 50)
	content.SetFontSize(64)  // 设置字体大小
	content.SetFont(fontKai) // 设置字体样式，就是我们上面加载的字体
	DrawText(content, "唐 颜真卿", 350, 520, 30)
	DrawText(content, "三更灯火五更鸡", 500, 300, 30)
	DrawText(content, "正是男儿读书时", 600, 300, 30)
	DrawText(content, "黑发不知勤学早", 700, 300, 30)
	DrawText(content, "白首方悔读书迟", 800, 300, 30)

}

// 根据文字转成竖排文字
func DrawText(content *freetype.Context, text string, left int, top int, lh int) {

	for i, Value := range text {
		content.DrawString(string(Value), freetype.Pt(left, top+i*lh))
	}

}

// 根据路径加载字体文件
// path 字体的路径
func loadFont(path string) (font *truetype.Font, err error) {
	var fontBytes []byte
	fontBytes, err = ioutil.ReadFile(path) // 读取字体文件
	if err != nil {
		err = fmt.Errorf("加载字体文件出错:%s", err.Error())
		return
	}
	font, err = freetype.ParseFont(fontBytes) // 解析字体文件
	if err != nil {
		err = fmt.Errorf("解析字体文件出错,%s", err.Error())
		return
	}
	return
}

func saveFile(pic *image.RGBA) {
	dstFile, err := os.Create("./dst.png")
	if err != nil {
		fmt.Println(err)
	}
	defer dstFile.Close()
	png.Encode(dstFile, pic)
}

// 根据地址获取图片内容
func getDataByUrl(url string) (img image.Image, err error) {
	res, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("[%s]通过url获取数据失败,err:%s", url, err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	// 读取获取的[]byte数据
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取数据失败,err:%s", err.Error())
		return
	}

	if !strings.HasSuffix(url, ".jpg") &&
		!strings.HasSuffix(url, ".jpeg") &&
		!strings.HasSuffix(url, ".png") {
		err = fmt.Errorf("[%s]不支持的图片类型,暂只支持.jpg、.png文件类型", url)
		return
	}

	// []byte 转 io.Reader
	reader := bytes.NewReader(data)
	if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") {
		// 此处jgeg.decode 有坑，明明是.jpg的图片但 会报 invalid JPEG format: missing SOI marker 错误
		// 所以当报错时我们再用 png.decode 试试
		img, err = jpeg.Decode(reader)
		if err != nil {
			fmt.Printf("jpeg.Decode err:%s", err.Error())
			reader2 := bytes.NewReader(data)
			img, err = png.Decode(reader2)
			if err != nil {
				err = fmt.Errorf("===>png.Decode err:%s", err.Error())
				return
			}
		}
	}

	if strings.HasSuffix(url, ".png") {
		img, err = png.Decode(reader)
		if err != nil {
			err = fmt.Errorf("png.Decode err:%s", err.Error())
			return
		}
	}

	return
}

func lineToPic(transparentImg *image.RGBA) {
	gc := draw2dimg.NewGraphicContext(transparentImg)
	gc.SetStrokeColor(color.RGBA{ // 线框颜色
		R: uint8(36),
		G: uint8(106),
		B: uint8(96),
		A: 0xff})
	gc.SetFillColor(color.RGBA{})
	gc.SetLineWidth(5) // 线框宽度
	gc.BeginPath()
	gc.MoveTo(0, 0)
	gc.LineTo(float64(transparentImg.Bounds().Dx()), 0)
	gc.LineTo(float64(transparentImg.Bounds().Dx()), float64(transparentImg.Bounds().Dy()))
	gc.LineTo(0, float64(transparentImg.Bounds().Dy()))
	gc.LineTo(0, 0)
	gc.Close()
	gc.FillStroke()
}
