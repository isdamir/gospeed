package config

import (
	"fmt"
	"os"
	"testing"
	/*"time"*/
)

type test struct {
	I int
	B bool
	S string
}

func TestJson(t *testing.T) {
	cf := GetConfig()
	te := &test{}
	cf.Register("test.json", "test", te)
	te.B = true
	//te.I=10
	//te.S = "测试信息"
	cf.Save("test")
	cf.Close("test")
	tet := &test{}
	cf.Register("test.json", "test", tet)
	if tet.B != te.B || tet.I != te.I || tet.S != te.S {
		fmt.Println("读取的信息不正确")
		t.FailNow()
	}
	tet.B = false
	tet.S = "Test"
	bd, err := GetConfig().Get("test")
	if err != nil {
		t.FailNow()
	}
	tx := bd.(*test)
	if tet.B != tx.B || tet.I != tx.I || tet.S != tx.S {
		fmt.Println("读取对象信息不正确")
		t.FailNow()
	}
	/*本段代码用于测试文件被修改*/
	/*测试时手动修改一下文件,查看结果*/
	/*fmt.Println(tx)*/
	/*time.Sleep(20000 * time.Millisecond)*/
	/*fmt.Println(tx)*/
	defer os.Remove("test.json")
}
func TestGob(t *testing.T) {
	cf := GetConfig()
	te := &test{}
	cf.Register("test.dat", "testgob", te)
	te.B = true
	te.I = 10
	te.S = "测试信息"
	cf.Save("testgob")
	cf.Close("testgob")
	tet := &test{}
	cf.Register("test.dat", "testgob", tet)
	if tet.B != te.B || tet.I != te.I || tet.S != te.S {
		fmt.Println("读取的信息不正确")
		t.FailNow()
	}
	tet.B = false
	tet.S = "Test"
	bd, err := GetConfig().Get("testgob")
	if err != nil {
		t.FailNow()
	}
	tx := bd.(*test)
	if tet.B != tx.B || tet.I != tx.I || tet.S != tx.S {
		fmt.Println("读取对象信息不正确")
		t.FailNow()
	}
	defer os.Remove("test.dat")
}
func TestXml(t *testing.T) {
	cf := GetConfig()
	te := &test{}
	cf.Register("test.xml", "testxml", te)
	te.B = true
	te.I = 10
	te.S = "测试信息"
	cf.Save("testxml")
	cf.Close("testxml")
	tet := &test{}
	cf.Register("test.xml", "testxml", tet)
	if tet.B != te.B || tet.I != te.I || tet.S != te.S {
		fmt.Println("读取的信息不正确")
		t.FailNow()
	}
	tet.B = false
	tet.S = "Test"
	bd, err := GetConfig().Get("testxml")
	if err != nil {
		t.FailNow()
	}
	tx := bd.(*test)
	if tet.B != tx.B || tet.I != tx.I || tet.S != tx.S {
		fmt.Println("读取对象信息不正确")
		t.FailNow()
	}
	defer os.Remove("test.xml")
}

func TestJsonNone(t *testing.T) {
	cf := GetConfig()
	te := &test{}
	cf.Register("test.json", "test", te)
	te.B = true
	te.I = 10
	te.S = "测试信息"
	cf.Save("test")
	cf.Close("test")
	tet := &test{}
	os.Rename("test.json", "json")
	cf.Register("json", "test", tet)
	if tet.B != te.B || tet.I != te.I || tet.S != te.S {
		fmt.Println("读取的信息不正确")
		t.FailNow()
	}
	defer os.Remove("json")
}

func TestGobNone(t *testing.T) {
	cf := GetConfig()
	te := &test{}
	cf.Register("test.dat", "testgob", te)
	te.B = true
	te.I = 10
	te.S = "测试信息"
	cf.Save("testgob")
	cf.Close("testgob")
	tet := &test{}
	os.Rename("test.dat", "dat")
	cf.Register("dat", "testgob", tet)
	if tet.B != te.B || tet.I != te.I || tet.S != te.S {
		fmt.Println("读取的信息不正确")
		t.FailNow()
	}
	defer os.Remove("dat")
}

func TestXmlNone(t *testing.T) {
	cf := GetConfig()
	te := &test{}
	cf.Register("test.xml", "testxml", te)
	te.B = true
	te.I = 10
	te.S = "测试信息"
	cf.Save("testxml")
	cf.Close("testxml")
	tet := &test{}
	os.Rename("test.xml", "xml")
	cf.Register("xml", "testxml", tet)
	if tet.B != te.B || tet.I != te.I || tet.S != te.S {
		fmt.Println("读取的信息不正确")
		t.FailNow()
	}
	defer os.Remove("xml")
}
