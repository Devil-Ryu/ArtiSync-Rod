package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// ImgPattern Markdown语法正则表达式
const ImgPattern string = "(\\!\\[(.*?)\\]\\((.*?)\\))"

// ImageInfo 图片信息结构体
type ImageInfo struct {
	Line      int    // 图片在原文的索引
	Title     string // 图片的名称[title](title)
	URL       string // 图片的链接[title](url)
	CurURL    string // 图片的当前链接
	UploadURL string // 上传后的链接
}

// MarkdownTool 解析器结构体
type MarkdownTool struct {
	MarkdownLines []string    // 文章内容
	ImagesInfo    []ImageInfo // 图片信息
	MarkdownPath  string      // 文章路径
	ImagePath     string      // 图片路径
}

// AnalyzeMarkdown 分析markdown文章
func (m *MarkdownTool) AnalyzeMarkdown() (err error) {
	// 读取markdown文章
	m.MarkdownLines, err = m.readFile(m.MarkdownPath)
	if err != nil {
		return fmt.Errorf("AnalyzeMarkdown: %w", err)
	}
	// 提取文章中的图片
	m.ImagesInfo, err = m.extractImages(m.MarkdownLines)
	if err != nil {
		return fmt.Errorf("AnalyzeMarkdown: %w", err)
	}
	// 从本地读取图片字节
	err = m.loadImages()
	if err != nil {
		return fmt.Errorf("AnalyzeMarkdown: %w", err)
	}

	return nil
}

// extractImages 提取文章中图片信息
func (m *MarkdownTool) extractImages(markdownLines []string) (imagesInfo []ImageInfo, err error) {
	re, err := regexp.Compile(ImgPattern)
	if err != nil {
		return imagesInfo, fmt.Errorf("提取图片失败: %w", err)
	}
	for index, line := range markdownLines {
		results := re.FindAllStringSubmatch(line, -1)
		for _, result := range results {
			imagesInfo = append(imagesInfo, ImageInfo{Line: index, Title: result[2], URL: result[3], CurURL: result[3]})
		}
	}
	return imagesInfo, nil
}

func (m *MarkdownTool) isURL(inputURL string) bool {
	_, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return false
	}
	return true
}

// 导入本地图片
func (m *MarkdownTool) loadImages() error {
	for index, imgInfo := range m.ImagesInfo {
		imagePath := imgInfo.URL

		isURL := m.isURL(imagePath)

		if isURL { // 如果为URL则下载图片后返回地址
			savePath, err := m.DownloadImage(imagePath)
			if err != nil {
				return fmt.Errorf("图片链接下载错误: %w", err)
			}
			m.ImagesInfo[index].URL = savePath
			m.ImagesInfo[index].CurURL = savePath

		} else { // 如果为路径，则不变
			m.ImagesInfo[index].URL = imagePath
			m.ImagesInfo[index].CurURL = imagePath
		}
	}
	return nil
}

// ReplaceImages 替换图片
func (m *MarkdownTool) ReplaceImages() {
	for imgIndex, imgInfo := range m.ImagesInfo {
		index := imgInfo.Line
		m.MarkdownLines[index] = strings.Replace(m.MarkdownLines[index], imgInfo.CurURL, imgInfo.UploadURL, 1) // 替换为上传链接
		m.ImagesInfo[imgIndex].CurURL = imgInfo.UploadURL                                                      // 替换后更新当前链接
	}
}

// 读取文件
func (m *MarkdownTool) readFile(fileName string) ([]string, error) {
	var fileList []string
	file, err := os.Open(fileName)

	if err != nil {
		// m.PushLog("error", fmt.Sprintf("Open File Error: %s", err))
		return nil, fmt.Errorf("readFile open: %w", err)
	}

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fileList = append(fileList, line)
	}

	return fileList, nil
}

// SaveToMarkdown 保存文件
func (m *MarkdownTool) SaveToMarkdown() (savePath string, err error) {

	baseDir := filepath.Dir(m.MarkdownPath)
	fileName := filepath.Base(m.MarkdownPath)

	saveDir := filepath.Join(baseDir, "ArtiSync_convert")
	savePath = filepath.Join(baseDir, "ArtiSync_convert", fileName)

	// 创建文件夹
	err = os.Mkdir(saveDir, 0755)

	// 当不是文件存在错误的时候则抛出错误
	if err != nil && !os.IsExist(err) {
		fmt.Println("ArtiSync_convert: ", err)
	}

	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println("保存Markdown失败: %w", err)
		return savePath, fmt.Errorf("保存Markdown失败: %w", err)
	}
	bw := bufio.NewWriter(file)
	_, err = bw.WriteString(strings.Join(m.MarkdownLines, "\n"))
	if err != nil {

		fmt.Println("SaveJSONFile write content error: %w", err)
		return savePath, fmt.Errorf("SaveJSONFile write content error: %w", err)
	}
	bw.Flush()
	return savePath, nil
}

// ReadImage 读取图片
func (m *MarkdownTool) ReadImage(imageName string) ([]byte, error) {
	file, err := os.Open(imageName)
	if err != nil {
		// m.PushLog("error", fmt.Sprintf("Open File Error: %s", err))
		return nil, fmt.Errorf("打开图片错误: %w", err)

	}

	// 读取文件内容
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		// m.PushLog("error", fmt.Sprintf("Read Image Error: %s", err))
		return nil, fmt.Errorf("读取图片文件错误: %w", err)
	}

	return imageBytes, nil

}

// DownloadImage 下载图片
func (m *MarkdownTool) DownloadImage(imageURL string) (string, error) {
	// 请求图片
	response, err := http.Get(imageURL)

	if err != nil {
		err = fmt.Errorf("网络请求错误: %w", err)
		return "", err
	}

	// 获取图片名称
	fileExt := filepath.Ext(imageURL)
	fileName := fmt.Sprintf("%x.%s", md5.Sum([]byte(imageURL)), fileExt)

	// 保存文件
	savePath := path.Join(m.ImagePath, "download")
	filePathRelative := path.Join("download", fileName)

	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(savePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filePathRelative, nil
}
