package utils

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

func GetSnapshot(videoPath string) (snapshot *bytes.Buffer, err error) {
	snapshot = bytes.NewBuffer(nil)
	//ffmpeg.
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(snapshot, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return nil, err
	}
	//img, err := imaging.Decode(buf)
	//if err != nil {
	//	log.Fatal("生成缩略图失败：", err)
	//	return nil, err
	//}
	//
	//err = imaging.Save(img, snapshotPath+".png")
	//if err != nil {
	//	log.Fatal("生成缩略图失败：", err)
	//	return nil, err
	//}
	//
	//names := strings.Split(snapshotPath, "\\")
	//snapshotName = names[len(names)-1] + ".png"
	return
}
