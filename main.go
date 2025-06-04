package main

import (
	conimg "concurrency/image"
	"strings"
	"bufio"
	"log"
	"os"
	"image"
	"fmt"
	"strconv"
	"path/filepath"
	"sync"
	"runtime"
)

func main(){
	numWorkers := runtime.NumCPU()

	// createpath,imagepath,sizeを取得
	fmt.Print("input create of image directory: ")
	input := GetInput()
	createpath := strings.TrimSpace(input)

	fmt.Print("input path of image directory: ")
	input = GetInput()
	imagepath := strings.TrimSpace(input)

	fmt.Print("input size: ")
	input = GetInput()
	input = strings.TrimSpace(input)
	strsize := strings.Fields(input)
	var intsize []int
	for _,s:= range strsize{
		size, _ := strconv.Atoi(s)
		intsize = append(intsize,size)
	}

	pngFiles, err1 := filepath.Glob(imagepath+"/*.png")
	jpegFiles, err2 := filepath.Glob(imagepath+"/*.jpeg")

	var allFiles []string
	if err1==nil && err2==nil{
		allFiles = append(pngFiles, jpegFiles...)
	} else if err1==nil && err2!=nil{
		allFiles = pngFiles
	} else if err1!=nil && err2==nil{
		allFiles = jpegFiles
	}
	if (err1!=nil && err2!=nil) || allFiles==nil{
		panic(fmt.Errorf("読み込めませんでした。\nerr1: %v\nerr2: %v\n読み込んだファイル: %v",err1,err2,allFiles))
	}


	type ImageInfo struct{
		imagepath string
		image image.Image
	}
	var filePathChan = make(chan string, len(allFiles))
	var imgdata = make(chan ImageInfo,len(allFiles))
	var resizedimgdata = make(chan ImageInfo,len(allFiles))

	//filepathの送信
	go func() {
		for _, fp := range allFiles {
			filePathChan <- fp
		}
		close(filePathChan) // 全て送信したらクローズ
	}()
	
	// imageのloading
	var wg1 sync.WaitGroup
	for i:=0;i<numWorkers;i++{
		wg1.Add(1)
		go func(){
			defer wg1.Done()
			for fp:= range filePathChan{
				img,err := conimg.LoadImage(fp)
				if err!=nil{
					panic(err)
				}
				imgdata<-ImageInfo{imagepath: fp,image: img}
			}
		}()
	}
	// 読み込み完了後にimgdataをクローズするためのゴルーチン
	go func() {
		wg1.Wait()
		close(imgdata)
	}()

	// resize
	var wg2 sync.WaitGroup
	for i:=0;i<numWorkers;i++{
		wg2.Add(1)
		go func(){
			defer wg2.Done()
			for img:= range imgdata{
				img.image = conimg.ResizeControl(img.image,intsize)
				resizedimgdata<-img
			}
		}()
	}
	// リサイズ完了後にresizedimgdataをクローズするためのゴルーチン
	go func() {
		wg2.Wait()
		close(resizedimgdata)
	}()

	
	//save
	var wg3 sync.WaitGroup
	for i:=0;i<numWorkers;i++{
		wg3.Add(1)
		go func(){
			defer wg3.Done()
			for resizedimg:= range resizedimgdata{
				err := conimg.SaveImage(createpath,resizedimg.imagepath,resizedimg.image)
				if err!=nil{
					panic(err)
				}
			}
		}()
	}
	wg3.Wait()
}




func GetInput () string { 
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return input
}