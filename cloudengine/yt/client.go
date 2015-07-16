package yt

import (
	"github.com/shima-park/asr/cloudengine/yt/log"
	"github.com/shima-park/asr/cloudengine/yt/session"
	"github.com/shima-park/asr/cloudengine/yt/util"

	//"bytes"
	//"encoding/json"
	//"io/ioutil"
	"net"
	//"net/http"
	"errors"
	//	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	GlobalSessions       *session.Manager //GlobalSessions
	SessionProvider      string
	SessionName          string
	SessionGCMaxLifetime int64
	SessionSavePath      string
	Year                 string
	Month                string
	Day                  string
	AudioDirPath         string
	LastResult           string
)

func init() {
	log.W("client", "global init")
	t := time.Now()
	SessionProvider = "memory"
	SessionName = "YT"
	SessionGCMaxLifetime = 60
	SessionSavePath = ""
	Year = strconv.Itoa(t.Year())
	Month = t.Month().String()
	Day = strconv.Itoa(t.Day())
	AudioDirPath = Year + `/` + Month + `/` + Day
	GlobalSessions, _ = session.NewManager(SessionProvider, SessionGCMaxLifetime, SessionSavePath)
	go GlobalSessions.GC()
	os.MkdirAll(AudioDirPath, 0777)
}

type Client struct {
	SessionId    string
	AudioDirPath string
	PrefixText   string
	conn         net.Conn
}

type UpPack struct {
	PackLength int
	HeadLength int
	Data       []byte
	CheckSum   [2]byte
	Appendfix  [2]byte
	Version    byte
	Seq        [2]byte
	Type       byte
	Status     byte
	Encrypt    byte
}

func NewClient(conn net.Conn) *Client {
	return &Client{SessionId: GlobalSessions.SessionStart().SessionID(), AudioDirPath: AudioDirPath, conn: conn}
}

func (this *Client) Run() {
	this.handleClient()
}

func (this *Client) getNextByte() (byte, error) {
	var err error = nil
	var len int = 0
	buf := make([]byte, 1)
	for {
		len, err = this.conn.Read(buf)
		if err != nil {
			return buf[0], err
		}
		if len >= 1 {
			break
		}
	}

	return buf[0], nil

}

func (this *Client) getNextData(max int) ([]byte, error) {
	var err error = nil
	var len int = 0
	totalLen := 0
	buf := make([]byte, max)

	for {
		len, err = this.conn.Read(buf[totalLen:])
		totalLen += len
		if err != nil {
			return buf, err
		}
		if totalLen >= max {
			break
		}
	}

	return buf, nil

}

// processData thead
func (this *Client) processData() (int, error) {
	var status int = 0
	var frame UpPack
	//open file
	fout, err := os.OpenFile(this.AudioDirPath+`/`+this.SessionId+".pcm", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.W(this.PrefixText, err.Error())
	}
	defer fout.Close()

	for {
		switch status {
		case 0:

			data, err := this.getNextByte()

			if err != nil {
				return 0, err
			}
			if data == 'Y' {
				status = 1
			} else {
				status = 0
			}
			break

		case 1:

			data, err := this.getNextByte()

			if err != nil {
				return 1, err
			}

			if data == 'Y' {
				status = 1
				break
			} else if data == 'T' {
				status = 2
				break
			} else {
				status = 0
				break
			}
			break

		case 2:

			data, err := this.getNextData(2)
			if err != nil {
				return 2, err
			}

			frame.PackLength = util.BytesToShort(data)

			if frame.PackLength < 0 {
				return 2, errors.New("frame : PackLengthError")
			} else {
				status = 3
			}
			break

		case 3:

			data, err := this.getNextData((frame.PackLength))
			if err != nil {
				return 3, err
			}

			var buf []byte = make([]byte, 4)
			buf[0] = data[0]
			buf[1] = data[1]
			frame.HeadLength = util.BytesToInt(buf)
			frame.Version = data[2]
			frame.Seq[0] = data[3]
			frame.Seq[1] = data[4]
			frame.Type = data[5]
			frame.Status = data[6]
			frame.Encrypt = data[7]

			frame.CheckSum[0] = data[frame.PackLength-4]
			frame.CheckSum[1] = data[frame.PackLength-3]
			frame.Appendfix[0] = data[frame.PackLength-2]
			frame.Appendfix[1] = data[frame.PackLength-1]

			if frame.Appendfix[0] == 'C' && frame.Appendfix[1] == 'N' {
				// 数据存放到Frame中

				frame.Data = data[8 : frame.PackLength-4]
				//frame.Data.Append("")
				status = 4

			} else {
				log.W(this.PrefixText, "Frame: error, skip.")
				status = 0
			}

			break

		case 4:

			// 上面的桢有效 ， 做操作
			log.W(this.PrefixText, "Frame: data len"+strconv.Itoa(frame.PackLength))

			// 写入文件
			fout.Write(util.DecoderFix(frame.Data))
			//fout.Write(frame.Data)
			// 解压缩， 上传google
			// 向客户端写入

			status = 0
			break
		}

	}

	return 0, nil
}

func (this *Client) handleClient() {

	this.PrefixText = "UP: " + this.conn.RemoteAddr().String() + " - "
	log.W(this.PrefixText, " Connect OK.") //remote ip
	// set 2 minutes timeout
	this.conn.SetReadDeadline(time.Now().Add(30 * time.Minute))

	defer this.conn.Close()                             // close connection before exit
	defer GlobalSessions.SessionDestroy(this.SessionId) // close session before exit

	status, err := this.processData()
	if err != nil {
		log.W(this.PrefixText, "status:"+strconv.Itoa(status)+", error:"+err.Error())
	}

}

/*
func (this *Client) writeResult() {
	if this.Header.Decollator == 1 {
		body := this.getResult()
		res := []byte{1, 0, 1, 0}
		res = append(res, util.IntToBytes(len(body))...)
		res = append(res, body...)
		log.W("client", res, "----end")
		this.conn.Write(res)
	} else if len(this.Content) > 0 && len(this.Result) > 0 && LastResult != string(this.Result) {
		this.conn.Write(this.Result)
		LastResult = string(this.Result)
	}
}


func (this *Client) buildResultResp() {
	timer := time.NewTicker(GET_RESULT_INTERVAL * time.Millisecond)

	for {
		select {
		case <-timer.C:
			body := this.getResult()
			if len(body) > 0 {
				if this.Header.Decollator != 1 {
					res := []byte{0, 0, 0, 0}
					res = append(res, util.IntToBytes(len(body))...)
					res = append(res, body...)
					this.Result = res
				}
			}

			if this.Header.Decollator == 1 {
				return
			}
		default:
			continue
		}
	}
}

func (this *Client) getResult() []byte {
	if len(this.Content) > 0 {
		var outpStream []byte
		if this.AudioConf.CodeFormat == "AMR" {
			outpStream = util.DecoderFix(this.Content)
			util.Trace(len(outpStream), "----len")
			saveFile := strconv.Itoa(len(outpStream))
			fout, err := os.Create(saveFile)
			if err != nil {
				util.Trace(saveFile, err)
			}
			defer fout.Close()
			fout.Write(outpStream)

		} else {
			outpStream = this.Content
		}

		buf := bytes.NewBuffer(this.Content)
		url := "http://173.194.72.104/speech-api/v1/recognize?xjerr=1&client=chromium&lang=zh-CN&maxresults=1"
		ContentType := "audio/L16; rate=" + this.AudioConf.getSampleRate()
		response, _ := http.Post(url, ContentType, buf)
		if response != nil && response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(response.Body)
			log.W("client", string(body), "resp body")
			return body
		}
	}
	return nil
}

func (this *Client) appendContent(audioByteArr []byte) {
	this.Content = append(this.Content, audioByteArr...)
}

func (this *AudioConfig) getSampleRate() string {
	if this.SampleRate > 0 {
		return strconv.Itoa(this.SampleRate)
	}
	return strconv.Itoa(16000)
}

func (this *AudioConfig) getCodeFormat() string {
	switch this.CodeFormat {
	case "PCM":
		return "L16"
	default:
		return "L16"
	}
}
*/
