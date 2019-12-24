package controller

 import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

var err error
var bn int
var DB *bolt.DB
var Bk *bolt.Bucket
var message chan string

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
func init(){
	message=make(chan string,1024)
	bn=time.Now().YearDay()
	DB,err=bolt.Open("log.db",0600,nil)
}

func Logger(){
	err=DB.Update(func(tx *bolt.Tx) error {
		Bk,err = tx.CreateBucketIfNotExists(IntToBytes(time.Now().YearDay()))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		for{
			select {
			case msg,ok:=<-message:
				if ok&&bn!=time.Now().YearDay(){
					Bk,err = tx.CreateBucketIfNotExists(IntToBytes(time.Now().YearDay()))
					Bk.Put([]byte(time.Now().String()),[]byte(msg))
					bn=time.Now().YearDay()
					log.Println("New",msg)
				}else if ok{
					Bk.Put([]byte(time.Now().String()),[]byte(msg))
					log.Println(msg)
				}else {
					return errors.New("internal error")
				}
			}
		}
		return nil
	})
	if err!=nil{
		panic(err)
	}
}


func View(bk []byte)(res map[string]string,err error){
	err=DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bk)
		if b==nil{
			return errors.New("bucket does not exist")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			res[string(k)]=string(v)
		}
		return nil
	})
	return
}

func Log(log string){
	message<-log
}