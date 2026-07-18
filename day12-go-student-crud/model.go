package main

type Student struct {
	ID        int    `json:"id"`
	StudentNo string `json:"student_no"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Status    string `json:"status"`
}

type StudentRequest struct {
	StudentNo string `json:"student_no"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Status    string `json:"status"`
}

type StudentStatusRequest struct {
	Status string `json:"status"`
}
