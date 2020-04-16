package request

type Casbin struct {
	Rolename string `json:"rolename" binding:"required,max=40"`
	Path     string `json:"path" binding:"required,max=40"`
	Method   string `json:"method" binding:"required,max=40"`
}

type Contest struct {
	Name        string `json:"name" binding:"required,max=20"`
	StartTime   string `json:"start_time"  binding:"required"`
	EndTime     string `json:"end_time"  binding:"required"`
	Description string `json:"description"`
	Problems    string `json:"problems"`
	LangMask    int    `json:"langmask"`
	Private     int    `json:"private"  binding:"gte=0,lte=1"`
	TeamMode    int    `json:"team_mode" ltefield=Private`
}

type Issue struct {
	Title     string `json:"title" binding:"required,max=20"`
	ProblemId int    `json:"problem_id"  binding:"gte=0"`
}

type Reply struct {
	Content     string `json:"content" binding:"required"`
	ReplyId     int    `json:"reply_id"  binding:"gte=0"`
	ReplyUserId int    `json:"reply_user_id"  binding:"gte=0"`
}

type New struct {
	Title   string `json:"title" binding:"required,max=20"`
	Content string `json:"content"`
}

type Problem struct {
	Title        string        `json:"title" binding:"required,max=20"`
	TimeLimit    int           `json:"time_limit"  binding:"required"`
	MemoryLimit  int           `json:"memory_limit"  binding:"required"`
	Description  string        `json:"description"`
	Input        string        `json:"input"`
	Output       string        `json:"output"`
	SampleInput  string        `json:"sample_input"`
	SampleOutput string        `json:"sample_output"`
	Spj          int           `json:"spj"`
	Hint         string        `json:"hint"`
	Level        int           `json:"level",gte=0, lte=2`
	Tags         []interface{} `json:"tags"`
}

type Series struct {
	Name        string `json:"name" binding:"required,max=20"`
	Description string `json:"description"`
	TeamMode    int    `json:"team_mode" binding:"gte=0,lte=1"`
}

type Tag struct {
	Name string `json:"name" binding:"required,max=20"`
}

type Team struct {
	Name string `json:"name" binding:"required,max=20"`
}
type AssignRole struct {
	UserName string `json:"username" binding:"required"`
	RoleId   int    `json:"role_id" binding:"required"`
}

type UnassignRole struct {
	UserId int `json:"user_id" binding:"required"`
}
