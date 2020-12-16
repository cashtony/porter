package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"porter/requester"
	"porter/wlog"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func ReadFileData(path string) string {
	f, err := os.Open(path)
	if err != nil {
		wlog.Fatal("读取文件失败, 请检查", err)
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		wlog.Fatal("读取文件内容失败, 请检查", err)
	}

	return string(fd)
}

func genThumbnails(vpath, gpath string) error {
	os.MkdirAll(path.Dir(gpath), os.ModePerm)
	msg, err := Cmd("ffmpeg", []string{"-y", "-i", vpath, "-vframes", "1", gpath})
	if err != nil {
		return fmt.Errorf("命令执行错误:%s message:%s", err, msg)
	}

	return nil
}

func cutVideoLength(vpath, cutPath string) error {
	msg, err := Cmd("ffmpeg", []string{"-y", "-ss", "00:00:00", "-i", vpath, "-to", "00:04:50", "-c", "copy", cutPath})
	if err != nil {
		return fmt.Errorf("命令执行错误:%s message:%s", err, msg)
	}

	return nil
}

func getVideoDuration(vpath string) (string, error) {
	msg, err := Cmd("ffprobe", []string{"-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", vpath})
	if err != nil {
		return "", fmt.Errorf("获取视频长度失败:%s message:%s", err, msg)
	}
	msg = strings.TrimSpace(msg)
	for {
		lastChracter := msg[len(msg)-1]
		if string(lastChracter) == "0" {
			msg = strings.TrimRight(msg, "0")
			continue
		}
		break
	}

	return msg, nil
}
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	var out, errbuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errbuf
	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("命令执行错误:%s \n %s", err, errbuf.String())
	}
	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("命令执行错误:%s \n %s", err, errbuf.String())
	}

	return out.String(), nil
}

func Ase256(plaintext string, key string, iv string) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), aes.BlockSize, len(plaintext))
	block, _ := aes.NewCipher(bKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func popCheack(bduss string) error {
	// var ids []cdp.NodeID
	cookies := map[string]string{
		"BDUSS": bduss,
	}
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.UserAgent(requester.UserAgent),
	}

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocatorCancel()

	optsctx, optCancel := chromedp.NewContext(
		allocatorCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer optCancel()

	ctx, cancel := context.WithTimeout(optsctx, 12*time.Hour)
	defer cancel()

	// var executed *runtime.RemoteObject
	var ids []cdp.NodeID
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for k, v := range cookies {
				success, err := network.SetCookie(k, v).
					WithExpires(&expr).
					WithDomain(cookieDomain).
					WithHTTPOnly(false).
					Do(ctx)
				if err != nil {
					return err
				}
				if !success {
					return fmt.Errorf("could not set cookie %s to %s", k, v)
				}
			}
			return nil
		}),
		// navigate to site
		chromedp.Navigate(uploadURL),
		chromedp.WaitReady(`div > span > div > span > input[type=file]`, chromedp.ByQuery),
		chromedp.NodeIDs(`div > span > div > span > input[type=file]`, &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			dom.SetFileInputFiles([]string{"./check.mp4"}).WithNodeID(ids[0]).Do(ctx)
			return nil
		}),

		// chromedp.WaitReady(`#root`, chromedp.ByQuery),
		// chromedp.Evaluate(js, &executed),
		chromedp.WaitVisible(".vcode-mask", chromedp.ByQuery),
		chromedp.WaitNotPresent(".vcode-mask", chromedp.ByQuery),
		// chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		return fmt.Errorf("自动上传失败:%s", err)
	}

	return nil
}

var js = `var url = "https://quanmin.baidu.com/wise/video/pcpub/getuploadid?video_num=1";
var xhr = new XMLHttpRequest();
xhr.open("POST", url, true);
xhr.send(null);
console.log("send")`
