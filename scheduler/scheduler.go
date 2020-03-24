package main

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"os/exec"
	"scheduler/mysql"
	"scheduler/rabbitmq"
	"scheduler/redisConn"
	"strconv"
	"strings"
)

var langExts []string
var DB *sqlx.DB
var REDISPOOL *redis.Pool
var coreDir string
var schedulerDir string

func init() {
	DB = mysql.Pool
	REDISPOOL = redisConn.Pool

	configFilePath := "config.ini"
	cfg, _ := goconfig.LoadConfigFile(configFilePath)
	projectcfg, _ := cfg.GetSection("project")
	coreDir = projectcfg["coredir"]
	schedulerDir = projectcfg["schedulerdir"]
	langExts = []string{"c", "cc", "pas", "java", "rb", "sh", "py", "php", "pl", "cs", "m", "bas", "scm", "c", "cc", "lua", "js", "go"}
}

type Task struct {
	UserId       int
	TestrunCount string
	SolutionId   int
	ProblemId    int
	Language     int
	TimeLimit    int
	MemoryLimit  int
	Source       string
	InputText    string
}

type Result struct {
	ResultCode int
	TimeUsed   int
	MemoryUsed int
	CEInfo     string
	REInfo     string
	CustomOut  string
	ErrorPoint string
	PassRate float64
}

// 清理目录
func cleanDir(runid int) {
	workDir := coreDir + "/run" + strconv.Itoa(runid)
	var cmd *exec.Cmd
	cmd = exec.Command("rm", "-rf", workDir+"/log")
	cmd.Run()
	cmd = exec.Command("mkdir", workDir+"/log")
	cmd.Run()
	// 把当前目录下除了log文件夹的所有内容转移到log文件夹下面
	cmd = exec.Command("bash","-c"," ls | grep -v log | xargs -i mv {} log")
	cmd.Dir = workDir
	cmd.Run()
}

// 准备判题相关文件
func prepareEnv(runid int, task *Task) {
	// 代码
	workDir := coreDir + "/run" + strconv.Itoa(runid)
	codeFilePath := workDir + "/Main." + langExts[task.Language]
	fmt.Println(codeFilePath)
	ioutil.WriteFile(codeFilePath, []byte(task.Source), 0755)
	// 用户输入
	customInputPath := workDir + "/data.in"
	ioutil.WriteFile(customInputPath, []byte(task.InputText), 0755)
}

// 运行判题程序
func runJudge(runid int, task *Task) {
	fmt.Printf("run judge")
	var cmd *exec.Cmd
	var err error
	// 入参 第二个参数 ojhome 第三个参数 问题id 第四个参数 语言 第五个参数 时间限制 第六个参数内存限制 第七个参数 runid
	cmd = exec.Command("./judge_client", coreDir, strconv.Itoa(task.ProblemId), strconv.Itoa(task.Language), strconv.Itoa(task.TimeLimit), strconv.Itoa(task.MemoryLimit), strconv.Itoa(runid))
	cmd.Dir = schedulerDir
	err = cmd.Run()
	fmt.Println(err, "判题程序执行成功")
}

func getResult(runid int) *Result {
	workDir := coreDir + "/run" + strconv.Itoa(runid)
	result := &Result{}
	resultText, err := ioutil.ReadFile(workDir + "/result.txt")
	if err == nil {
		resultList := strings.Split(string(resultText),"\n")
		result.ResultCode,_ = strconv.Atoi(resultList[0])
		result.TimeUsed,_ = strconv.Atoi(resultList[1])
		result.MemoryUsed,_ = strconv.Atoi(resultList[2])
		result.PassRate,_ =  strconv.ParseFloat(resultList[3],64)
		result.ErrorPoint =  resultList[4]
	}
	ceInfo,err := ioutil.ReadFile(workDir+"/ce.txt")
	if	err == nil{
		result.CEInfo = string(ceInfo)
	}
	reInfo,err := ioutil.ReadFile(workDir+"/error.out")
	if	err == nil{
		result.REInfo = string(reInfo)
	}
	customOut,err := ioutil.ReadFile(workDir+"/user.out")
	if	err == nil{
		result.CustomOut = string(customOut)
	}
	fmt.Println("%+v",result)
	return result
}

func work(runid int) {
	rabbitmq.Receive("oj", "problem", "problem", func(msg *[]byte) {
		conn := REDISPOOL.Get()
		defer conn.Close()
		fmt.Println("runid", runid)
		var task Task
		json.Unmarshal(*msg, &task)
		// 清理目录
		cleanDir(runid)
		// 准备文件
		prepareEnv(runid, &task)
		//运行判题程序
		runJudge(runid, &task)
		// 获取结果
		result := getResult(runid)
		// 清理目录
		cleanDir(runid)
		if task.ProblemId == 0 {
			conn.Send("MULTI")
			conn.Send("hset","testrun:"+task.TestrunCount,"ceinfo",result.CEInfo)
			conn.Send("hset","testrun:"+task.TestrunCount,"reinfo",result.REInfo)
			conn.Send("hset","testrun:"+task.TestrunCount,"custom_out",result.CustomOut)
			conn.Send("hset","testrun:"+task.TestrunCount,"error_point",result.ErrorPoint)
			// 测试运行数据只保存5分钟
			conn.Send("expire", "testrun:"+task.TestrunCount, 60*5)
			conn.Do("exec")
		} else {
			// 更新数据库
			var err error
			_,err = DB.Exec("update solution set result = ?,time = ?, memory = ?,judgetime = NOW(),pass_rate = ? where solution_id = ?",result.ResultCode,result.TimeUsed,result.MemoryUsed,result.PassRate,task.SolutionId)
			fmt.Println(err)
			_,err = DB.Exec("update user set solved = (select count(distinct problem_id) from solution where user_id = ? and result = 4 and contest_id = 0) where id = ?",task.UserId,task.UserId)
			fmt.Println(err)
			_,err = DB.Exec("update user set submit = (select count(1) from solution where user_id = ? and problem_id > 0 and contest_id = 0) where id = ?",task.UserId,task.UserId)
			fmt.Println(err)

			// 编译错误
			if result.ResultCode == 11 {
				DB.Exec("insert into compileinfo (solution_id,error) values(?,?)",task.SolutionId,result.REInfo)
			}
			// 运行错误
			if result.ResultCode == 10 {
				DB.Exec("insert into runtimeinfo (solution_id,error) values(?,?)",task.SolutionId,result.REInfo)
			}
			// 错误点信息 记录到 runtime info中
			if result.ResultCode >= 6 && result.ResultCode <= 8 {
				DB.Exec("insert into runtimeinfo (solution_id,error) values(?,?)",task.SolutionId,result.ErrorPoint)
			}
		}
	})
}

func main() {
	f, _ := os.OpenFile("./schedule.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	os.Stdout = f
	for i := 0; i < 10; i++ {
		go work(i)
	}
	select {}
}
