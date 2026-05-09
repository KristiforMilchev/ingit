package models

type PushPlanExecution struct {
	Branch  string
	Base    string
	Message string
	Files   []string
}
