package main

import (
	"fmt"

	"github.com/m1a9s9a4/autofill"
	"github.com/m1a9s9a4/autofill/rules"
)

// Custom status rule that generates status values based on index
type StatusRule struct{}

func (r *StatusRule) Generate(ctx rules.Context) (interface{}, error) {
	statuses := []string{"active", "inactive", "pending", "suspended"}
	return statuses[ctx.Index()%len(statuses)], nil
}

func (r *StatusRule) Validate(v interface{}) error {
	return nil
}

// Custom priority rule
type PriorityRule struct{}

func (r *PriorityRule) Generate(ctx rules.Context) (interface{}, error) {
	priorities := []string{"low", "medium", "high", "critical"}
	return priorities[ctx.Index()%len(priorities)], nil
}

func (r *PriorityRule) Validate(v interface{}) error {
	return nil
}

type Task struct {
	ID       int64 `autofill:"seq"`
	Title    string
	Status   string `autofill:"rule=status"`
	Priority string `autofill:"rule=priority"`
	Assignee string
}

func main() {
	// Create a RuleSet with custom rules
	ruleSet := rules.DefaultRuleSet()
	ruleSet.Add("status", &StatusRule{})
	ruleSet.Add("priority", &PriorityRule{})

	// Create autofill with custom rules
	af := autofill.New().WithRules(ruleSet).WithSeed(12345)

	fmt.Println("=== Custom Rules Example ===")
	tasks := make([]Task, 8)
	if err := af.FillSlice(&tasks); err != nil {
		panic(err)
	}

	for i, task := range tasks {
		fmt.Printf("%d: ID=%d Status=%-10s Priority=%-10s\n",
			i+1, task.ID, task.Status, task.Priority)
	}
}
