package image

import(
	"os"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

// 画像を読み取るための関数。
// ファイルパスを指定すると、画像データを返してくれる。
func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}