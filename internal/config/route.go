package config

type Route struct {
	Method   string
	Pattern  string
	Handler  string
}

func LoadRoutes() []Route {
	return []Route{
		{Method: "GET", Pattern: "^/user/(\\d+)$", Handler: "User:GetIndex"},
		{Method: "GET", Pattern: "^/user/(\\d+)/task$", Handler: "User:GetUserTask"},
		{Method: "POST|PUT", Pattern: "^/user/(\\d+)/task/(\\d+)$", Handler: "Task:AddTaskToUser"},
		{Method: "POST|PUT", Pattern: "^/task$", Handler: "Task:AddTask"},
		{Method: "POST|PUT", Pattern: "^/task/(\\d+)$", Handler: "Task:EditTask"},
		{Method: "DELETE", Pattern: "^/task/(\\d+)$", Handler: "Task:DeleteTask"},
		{Method: "DELETE", Pattern: "^/user/(\\d+)/task/(\\d+)$", Handler: "Task:DeleteUserTask"},
	}
}
