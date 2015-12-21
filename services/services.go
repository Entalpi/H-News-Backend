package services

import (
	"bytes"
	"encoding/binary"
	"hnews/Godeps/_workspace/src/github.com/boltdb/bolt"
	"hnews/Godeps/_workspace/src/github.com/headzoo/surf"
	"log"
	"strconv"
	"time"
)

const (
	Debug = true
)

/** Login service **/
func Login(username string, password string) bool {
	bow := surf.NewBrowser()
	err := bow.Open("https://news.ycombinator.com/login?goto=news")
	if err != nil {
		return false
	}

	fm, _ := bow.Form("form")
	fm.Input("acct", username)
	fm.Input("pw", password)
	erro := fm.Submit()
	if erro != nil {
		log.Println(erro)
		return false
	}
	log.Println(fm)

	return true
}

/** Database Service **/
// Ints need to be a defined size since binary.Read does not support just 'int'
type News struct {
	Rank     int32     `json:"rank"`
	Title    string    `json:"title"`
	Link     string    `json:"link"`
	Author   string    `json:"author"`
	Points   int32     `json:"points"`
	Time     time.Time `json:"time"`
	Comments int32     `json:"comments"`
}

var (
	db, err = bolt.Open("a.db", 0644, nil)
)

func ReadNews(from int, to int) []News {
	var news []News
	db.View(func(tx *bolt.Tx) error {
		for i := from; i <= to; i++ {
			b := tx.Bucket([]byte(strconv.Itoa(int(i))))
			if b == nil {
				log.Println("Bucket", i, "not found.")
				continue
			}

			title := string(b.Get([]byte("title")))
			author := string(b.Get([]byte("author")))
			link := string(b.Get([]byte("link")))

			var rank int32
			binary.Read(bytes.NewReader(b.Get([]byte("rank"))), binary.LittleEndian, &rank)

			var t int64
			binary.Read(bytes.NewReader(b.Get([]byte("time"))), binary.LittleEndian, &t)

			var points int32
			binary.Read(bytes.NewReader(b.Get([]byte("points"))), binary.LittleEndian, &points)

			var comments int32
			binary.Read(bytes.NewReader(b.Get([]byte("comments"))), binary.LittleEndian, &comments)

			if err != nil {
				log.Println("ReadNews:", err)
				return err
			}
			news = append(news, News{rank, title, link, author, points, time.Unix(t, 0), comments})
		}
		return nil
	})
	return news
}

func SaveNews(news []News) {
	db.Update(func(tx *bolt.Tx) error {
		for _, aNews := range news {
			b, err := tx.CreateBucketIfNotExists([]byte(strconv.Itoa(int(aNews.Rank))))
			if err != nil {
				log.Println(err)
				return err
			}

			b.Put([]byte("title"), []byte(aNews.Title))
			b.Put([]byte("link"), []byte(aNews.Link))
			b.Put([]byte("author"), []byte(aNews.Author))

			var t bytes.Buffer
			binary.Write(&t, binary.LittleEndian, aNews.Time.Unix())
			b.Put([]byte("time"), []byte(t.Bytes()))

			var r bytes.Buffer
			binary.Write(&r, binary.LittleEndian, aNews.Rank)
			b.Put([]byte("rank"), []byte(r.Bytes()))

			var p bytes.Buffer
			binary.Write(&p, binary.LittleEndian, aNews.Points)
			b.Put([]byte("points"), []byte(p.Bytes()))

			var c bytes.Buffer
			binary.Write(&c, binary.LittleEndian, aNews.Comments)
			b.Put([]byte("comments"), []byte(c.Bytes()))

			if err != nil {
				log.Println("Err SaveNews:", err)
				return err
			}
		}
		return nil
	})
}
