package image

import(
	"os"
	"image"
	"fmt"
	"strings"
	"path/filepath"
	"image/jpeg"
	"image/png"
)
func SaveImage(createpath,path string, img image.Image) error {
	name := filepath.Base(path)
	fp := filepath.Ext(path)//拡張子を取得
	ext := strings.ToLower(fp)//全て小文字に
	switch ext {
	case ".jpg", ".jpeg":
		f, err := os.Create(createpath+"/"+name)
		if err != nil {
			return err
		}
		defer f.Close()
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
	case ".png":
		f, err := os.Create(createpath+"/"+name)
		if err != nil {
			return err
		}
		defer f.Close()
		return png.Encode(f, img)
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}
}