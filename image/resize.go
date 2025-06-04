package image

import(
	"image"
	"golang.org/x/image/draw"
)

func ResizeControl(img image.Image, resize []int) image.Image {
	if len(resize)==2{
		img = resizeImage(img,resize[0],resize[1])
	}else if len(resize)==1{
		img = resizeImageKeepAspect(img, resize[0])
	}
	return img
}

func resizeImage(img image.Image, width, height int) image.Image {
	// 欲しいサイズの画像を新しく作る
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// サイズを変更しながら画像をコピーする
	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	return newImage
}

func resizeImageKeepAspect(img image.Image, size int) image.Image {
	// 画像のサイズを取得する
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	// 結果となる画像のサイズを計算する
	if width > height {
		height = height * size / width
		width = size
	} else {
		width = width * size / height
		height = size
	}

	// 先ほどの関数を使って画像をリサイズする
	return resizeImage(img, width, height)
}

