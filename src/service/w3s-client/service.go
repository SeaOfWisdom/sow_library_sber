package w3scli

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"io/fs"
// 	"os"

// 	"github.com/ipfs/go-cid"
// 	"github.com/web3-storage/go-w3s-client"
// 	w3fs "github.com/web3-storage/go-w3s-client/fs"
// )

// const tokenMVP = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDkyNDk3OTY2NjY3N2I1MjIwRGQ1M0VDNjE5OTU2NTEwOEE1ZDM4MTciLCJpc3MiOiJ3ZWIzLXN0b3JhZ2UiLCJpYXQiOjE2ODI2OTM2OTkyNjksIm5hbWUiOiJQUC1NVlAifQ.HxUuc0lqFPsiIqIxKzlc7jFIfZWZ-pd-P7KtEaQGZo8"

// type W3sSrv struct {
// 	cli w3s.Client
// }

// func CreateW3sSrv() *W3sSrv {
// 	c, err := w3s.NewClient(w3s.WithToken(tokenMVP))
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &W3sSrv{cli: c}
// }

// func (ws *W3sSrv) PutJsonFile() (string, error) {
// 	fmt.Println(os.Getwd())
// 	workFile, err := os.Open("./src/service/w3s-client/example/paper_1.json")
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println(workFile.Name())
// 	// TODO with deadline
// 	cid, err := ws.cli.Put(context.Background(), workFile)
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println("LINK")
// 	fmt.Printf("https://%v.ipfs.dweb.link\n", cid)

// 	return "", nil
// }

// func putMultipleFiles(c w3s.Client) cid.Cid {
// 	f0, err := os.Open("images/donotresist.jpg")
// 	if err != nil {
// 		panic(err)
// 	}
// 	f1, err := os.Open("images/pinpie.jpg")
// 	if err != nil {
// 		panic(err)
// 	}
// 	dir := w3fs.NewDir("comic", []fs.File{f0, f1})
// 	return putFile(c, dir)
// }

// func putFile(c w3s.Client, f fs.File, opts ...w3s.PutOption) cid.Cid {
// 	cid, err := c.Put(context.Background(), f, opts...)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("https://%v.ipfs.dweb.link\n", cid)
// 	return cid
// }

// func (ws *W3sSrv) GetStatusForCid(cidStr string) (string, error) {
// 	cid, err := cid.Parse(cidStr)
// 	if err != nil {
// 		return "", err
// 	}
// 	s, err := ws.cli.Status(context.Background(), cid)
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Printf("Status: %+v", s)
// 	return "", nil
// }

// func (ws *W3sSrv) GetFiles(cidStr string) error {
// 	// cid, err := cid.Parse(cidStr)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// res, err := ws.cli.Get(context.Background(), cid)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // f, fsys, err := res.Files()
// 	// // if err != nil {
// 	// // 	return err
// 	// // }

// 	// type Resp struct {
// 	// 	Title   string `json:"title"`
// 	// 	Author  string `json:"author"`
// 	// 	RawData string `json:"raw_data"`
// 	// }

// 	// resp := Resp{}
// 	// b1 := []byte{}
// 	// if _, err := f.Read(b1); err != nil {
// 	// 	fmt.Println("HERE")
// 	// 	return err
// 	// }

// 	// json.Unmarshal(b1, &resp)
// 	// fmt.Println("JSON :", resp)

// 	// info, err := f.Stat()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// ls, err := fsys.Open(info.Name())
// 	// if err != nil {
// 	// 	fmt.Println("xxx")
// 	// 	return err
// 	// }

// 	// b2 := []byte{}
// 	// if _, err := ls.Read(b2); err != nil {
// 	// 	fmt.Println("HER1231E")
// 	// 	return err
// 	// }

// 	// if info.IsDir() {
// 	// 	err = fs.WalkDir(fsys, "/", func(path string, d fs.DirEntry, err error) error {
// 	// 		info, _ := d.Info()
// 	// 		fmt.Printf("%s (%d bytes)\n", path, info.Size())
// 	// 		return err
// 	// 	})
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// } else {
// 	// 	fmt.Println("JERE")
// 	// 	fmt.Printf("%s (%d bytes)\n", cid.String(), info.Size())
// 	// }
// 	return nil
// }

// func (ws *W3sSrv) ListUploads() error {
// 	uploads, err := ws.cli.List(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		u, err := uploads.Next()
// 		if err != nil {
// 			// finished successfully
// 			if err == io.EOF {
// 				break
// 			}
// 			return err
// 		}

// 		fmt.Printf("%s	%s	Size: %d	Deals: %d	Pins: %d\n", u.Created.Format("2006-01-02 15:04:05"), u.Cid, u.DagSize, len(u.Deals), len(u.Pins))
// 	}
// 	return nil
// }
