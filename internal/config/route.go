package config

type Route struct {
	Method   string
	Pattern  string
	Handler  string
	ParamMap map[string]string
}

func LoadRoutes() []Route {
	return []Route{
		{
			Method:   "GET",
			Pattern:  "^/user/(\\d+)$",
			Handler:  "User:GetIndex",
			ParamMap: map[string]string{"1": "id"},
		},
		{
			Method:   "GET",
			Pattern:  "^/user/(\\d+)/task$",
			Handler:  "User:GetUserTask",
			ParamMap: map[string]string{"1": "id"},
		},
		{
			Method:   "POST|PUT",
			Pattern:  "^/user/(\\d+)/task/(\\d+)$",
			Handler:  "Task:AddTaskToUser",
			ParamMap: map[string]string{"1": "userId", "2": "taskId"},
		},
		{
			Method:   "POST|PUT",
			Pattern:  "^/task$",
			Handler:  "Task:AddTask",
		},
		{
			Method:   "POST|PUT",
			Pattern:  "^/task/(\\d+)$",
			Handler:  "Task:EditTask",
			ParamMap: map[string]string{"1": "id"},
		},
		{
			Method:   "DELETE",
			Pattern:  "^/task/(\\d+)$",
			Handler:  "Task:DeleteTask",
			ParamMap: map[string]string{"1": "id"},
		},
		{
			Method:   "DELETE",
			Pattern:  "^/user/(\\d+)/task/(\\d+)$",
			Handler:  "Task:DeleteUserTask",
			ParamMap: map[string]string{"1": "userId", "2": "taskId"},
		},
	}
}
