package main
/*
#include "judge_client.c"
*/
import "C"
import (
	"unsafe"
	"fmt"
	"encoding/json"
	"scheduler/rabbitmq"
	"scheduler/mysql"
	"scheduler/redisConn"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB
var REDISPOOL *redis.Pool

func init() {
	DB = mysql.Pool
	REDISPOOL = redisConn.Pool
}

type Task struct {
	UserId int
	TestrunCount string
	SolutionId int
	ProblemId  int
	Language   int
	TimeLimit  int
	MemoryLimit	int
	Source     string 
	InputText  string
}

func work(runid int){
	rabbitmq.Receive("oj", "problem", "problem", func(msg *[]byte) {
		conn := REDISPOOL.Get()
		defer conn.Close()
		var task Task
		json.Unmarshal(*msg, &task)
		ctask := C.struct_task{
			problem_id:C.int(task.ProblemId),
			language:C.int(task.Language),
			time_limit:C.int(task.TimeLimit),
			memory_limit:C.int(task.MemoryLimit),
			source:C.CString(task.Source),
			input_text:C.CString(task.InputText),
		}

		ojHome := C.CString("/home/judge/hustoj_core")
		// fmt.Println(ctask)
		// fmt.Println(ojHome)
		judgeResult,_ := C.judge(&ctask,ojHome,C.int(runid))
		ceinfo := C.GoString(judgeResult.ceinfo)
		reinfo := C.GoString(judgeResult.reinfo)
		custom_out := C.GoString(judgeResult.custom_out)
		error_point := C.GoString(judgeResult.error_point)
		result := int(judgeResult.result)
		time_used := int(judgeResult.time_used)
		memory_used := int(judgeResult.memory_used)
		pass_rate := float64(judgeResult.pass_rate)
		// fmt.Println(ceinfo)
		// fmt.Println(reinfo)
		// fmt.Println(custom_out)
		// fmt.Println(error_point)
		// fmt.Println(task)
		C.free(unsafe.Pointer(ojHome))
		C.free(unsafe.Pointer(ctask.source))
		C.free(unsafe.Pointer(ctask.input_text))
		C.free(unsafe.Pointer(judgeResult.ceinfo))
		C.free(unsafe.Pointer(judgeResult.reinfo))
		C.free(unsafe.Pointer(judgeResult.custom_out))
		C.free(unsafe.Pointer(judgeResult.error_point))
		//说明是测试运行
		if task.ProblemId == 0 {
			conn.Send("MULTI")
			conn.Send("hset","testrun:"+task.TestrunCount,"ceinfo",ceinfo)
			conn.Send("hset","testrun:"+task.TestrunCount,"reinfo",reinfo)
			conn.Send("hset","testrun:"+task.TestrunCount,"custom_out",custom_out)
			conn.Send("hset","testrun:"+task.TestrunCount,"error_point",error_point)
			// 测试运行数据只保存5分钟
			conn.Send("expire", "testrun:"+task.TestrunCount, 60*5)
			conn.Do("exec")
		}else{
			// 更新数据库
			fmt.Println("UPDATE DATABASE1",result)
			var err error
			_,err = DB.Exec("update solution set result = ?,time = ?, memory = ?,judgetime = NOW(),pass_rate = ? where solution_id = ?",result,time_used,memory_used,pass_rate,task.SolutionId)
			fmt.Println(err)
			_,err = DB.Exec("update user set solved = (select count(distinct problem_id) from solution where user_id = ? and result = 4 and contest_id = 0) where id = ?",task.UserId,task.UserId)
			fmt.Println(err)
			_,err = DB.Exec("update user set submit = (select count(1) from solution where user_id = ? and problem_id > 0 and contest_id = 0) where id = ?",task.UserId,task.UserId)
			fmt.Println(err)

			// 编译错误
			if result == 11 {
				DB.Exec("insert into compileinfo (solution_id,error) values(?,?)",task.SolutionId,ceinfo)
			}
			// 运行错误
			if result == 10 {
				DB.Exec("insert into runtimeinfo (solution_id,error) values(?,?)",task.SolutionId,reinfo)
			}
			// 错误点信息 记录到 runtime info中
			if result >= 6 && result <= 8 {
				DB.Exec("insert into runtimeinfo (solution_id,error) values(?,?)",task.SolutionId,error_point)
			}
		}
	})
}

func main() {
	for i := 0; i < 10; i++ {
		go work(i)
	}
	select{}
}
