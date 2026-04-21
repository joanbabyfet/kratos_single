package server

import (
	"fmt"
	"io"
	"kratos_single/internal/pkg/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// 只允许 POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 最大上传限制（10MB）
	const maxSize int64 = 10
	err := r.ParseMultipartForm(maxSize << 20)
	if err != nil {
		http.Error(w, "上传失败", http.StatusBadRequest)
		return
	}

	// 读取文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "请选择文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 读取附加参数
	dir := r.FormValue("dir") // image/avatar/article
	if dir == "" {
		dir = "image"
	}

	thumbW, _ := strconv.Atoi(r.FormValue("thumb_w"))
	thumbH, _ := strconv.Atoi(r.FormValue("thumb_h"))

	// 文件转 bytes
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "读取文件失败", http.StatusBadRequest)
		return
	}

	// 大小限制
	if int64(len(data)) > maxSize*1024*1024 {
		http.Error(w, "文件过大", http.StatusBadRequest)
		return
	}

	// 后缀限制
	ext := strings.ToLower(path.Ext(header.Filename))
	allowExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if !allowExt[ext] {
		http.Error(w, "文件格式不支持", http.StatusBadRequest)
		return
	}

	// 日期目录
	dirDate := time.Now().Format("20060102")

	uploadRoot := "../../uploads"
	fileURL := "/uploads"

	fullDir := filepath.Join(uploadRoot, dir, dirDate)

	if err := os.MkdirAll(fullDir, 0775); err != nil {
		http.Error(w, "创建目录失败", http.StatusInternalServerError)
		return
	}

	// 新文件名
	fileName := utils.GenFileName(ext)
	savePath := filepath.Join(fullDir, fileName)

	// 保存原图
	if err := os.WriteFile(savePath, data, 0644); err != nil {
		http.Error(w, "保存文件失败", http.StatusInternalServerError)
		return
	}

	fileLink := fmt.Sprintf("%s/%s/%s/%s", fileURL, dir, dirDate, fileName)

	// 缩略图
	if thumbW > 0 || thumbH > 0 {

		img, err := imaging.Open(savePath)
		if err != nil {
			http.Error(w, "图片解析失败", http.StatusBadRequest)
			return
		}

		thumb := imaging.Resize(
			img,
			thumbW,
			thumbH,
			imaging.Lanczos,
		)

		fileName = utils.GenFileName(ext)
		savePath = filepath.Join(fullDir, fileName)
		fileLink = fmt.Sprintf("%s/%s/%s/%s", fileURL, dir, dirDate, fileName)

		if err := imaging.Save(thumb, savePath); err != nil {
			http.Error(w, "缩略图生成失败", http.StatusInternalServerError)
			return
		}
	}

	// 返回 JSON
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, `{
		"real_name":"%s",
		"file_name":"%s/%s",
		"file_link":"%s"
	}`, header.Filename, dirDate, fileName, fileLink)
}