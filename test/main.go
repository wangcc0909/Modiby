package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

func main() {
	/*runtime.GOMAXPROCS(1)
	f,_ := os.Create("trace.dat")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	ctx, task := trace.NewTask(context.Background(), "sumTask")
	defer task.End()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			//标记region
			trace.WithRegion(ctx,region, func() {
				var sum,k int64
				for ; k < 1000000000; k++ {
					sum += k
				}
				fmt.Println(region,sum)
			})
		}(fmt.Sprintf("region_%02d",i))
	}
	wg.Wait()*/
	//fmt.Println(Inc())
	/*go demo()
	sg := make(chan os.Signal)
	signal.Notify(sg,syscall.SIGINT,syscall.SIGQUIT,syscall.SIGKILL)
	select {
	case s := <-sg:
		log.Println("退出： ",s.String())
	}*/
	//t := "&&&===&&"
	//result := parseString(t)
	//fmt.Println(result)
	//ReversePrint()
	/*num := 6
	for i := 0; i < num; i++ {
		resp,_ := http.Get("https://www.baidu.com")
		_,_ = ioutil.ReadAll(resp.Body)
	}
	fmt.Printf("此时goroutine个数=%d",runtime.NumGoroutine())*/
	/*c := 10
	b := make([]byte,c)
	_,err := rand.Read(b)
	if err != nil {
		fmt.Println("error: ",err)
		return
	}
	fmt.Println(bytes.Equal(b,make([]byte,c)))*/
	/*b := requestBTC()
	fmt.Println(len(b))
	fmt.Println(string(b))*/
	/*c := make(chan int,1)
	close(c)
	c <- 1*/
	/*r := Rect{Point{1,2},Point{3,4}}
	fmt.Println(r)*/
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()
	fmt.Println(total)
}

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i <= 100; i++ {
		atomic.AddUint64(&total, uint64(i))
	}
}

func requestBTC() []byte {
	resp, err := http.Get("https://sochain.com/api/v2/get_info/DOGE")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return bs
}

/*func test() {
	t := struct {
		time.Time
		N int
	}{
		time.Date(2020,12,20,0,0,0,0,time.UTC),
		5,
	}
	result,_ := json.Marshal(t) //因为time.Time已经实现了json.Marshaler interface
	fmt.Printf("%s",result) // "2020-12-20T00:00:00Z"
}

func hello(h http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, rst *http.Request) {
		s := rst.RemoteAddr
		log.Println(s)
		h(rw,rst)
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json; charset=utf-8")
	w.Write([]byte(`{"name": "你好"}`))
}

func demo() {
	http.HandleFunc("/hello",hello(greet))
	http.ListenAndServe(":9090",nil)
}

func Inc() (v int) {
	defer func() {
		v++
	}()
	return 42
}

//请写一个字符串协议解析的实现，如a=b&c=d ——>map[a]=b,map[c]=d,要求尽量保持对异常的兼容如&&&===&a=b&
func parseString(url string) map[string]interface{} {
	params := make(map[string]interface{})
	result := strings.Split(url,"&")
	for _, kv := range result {
		if !strings.Contains(kv,"=") {
			continue
		}
		kv = strings.Trim(kv,"=")
		kvArr := strings.Split(kv,"=")
		if len(kvArr) > 2 || kvArr[0] == "" || kvArr[1] == "" {
			continue
		}
		params[kvArr[0]] = kvArr[1]
	}
	return params
}

//请实现一个高效的单向链表逆序输出（1->2->3 转 3->2->1）
type ListNode struct {
	Value int
	Next *ListNode
}

func ReversePrint() {
	head := &ListNode{1,&ListNode{2,&ListNode{Value: 3}}}
	reversePrint(head)
}

func reversePrint(head *ListNode)  {
	if head != nil {
		reversePrint(head.Next)
		fmt.Print(head.Value,"\t")
	}
}*/

type Point struct {
	x, y int
}

type Rect struct {
	a, b Point
}
