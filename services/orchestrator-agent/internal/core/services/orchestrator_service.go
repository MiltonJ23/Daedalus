package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
	"github.com/Daedalus/orchestrator-agent/internal/core/ports"
)

// OrchestratorService — decomposes goals and dispatches sub-tasks.
//
// PB-025 acceptance:
//   - task graph visible in agent log
//   - each sub-task has unique ID and status
type OrchestratorService struct {
	goals     ports.GoalRepository
	tasks     ports.TaskRepository
	publisher ports.TaskPublisher
}

func NewOrchestratorService(g ports.GoalRepository, t ports.TaskRepository, p ports.TaskPublisher) *OrchestratorService {
	return &OrchestratorService{goals: g, tasks: t, publisher: p}
}

// SubmitGoal validates → persists the goal → decomposes it into a sub-task
// graph → publishes each task → returns the full graph.
func (s *OrchestratorService) SubmitGoal(ctx context.Context, data map[string]interface{}) (domain.TaskGraph, error) {
	if err := s.validate(data); err != nil {
		return domain.TaskGraph{}, err
	}
	userID := strings.TrimSpace(data["user_id"].(string))
	projectID, _ := data["project_id"].(string)
	description := strings.TrimSpace(data["description"].(string))

	goal := domain.NewGoal(userID, projectID, description)
	created, err := s.goals.Create(ctx, goal)
	if err != nil {
		return domain.TaskGraph{}, fmt.Errorf("failed to persist goal: %w", err)
	}

	tasks := s.decompose(created, description)

	if err := s.tasks.BulkCreate(ctx, tasks); err != nil {
		return domain.TaskGraph{}, fmt.Errorf("failed to persist tasks: %w", err)
	}

	logTaskGraph(created, tasks)

	for i := range tasks {
		streamID, err := s.publisher.Publish(ctx, tasks[i])
		if err != nil {
			log.Printf("publish task %s (%s) failed: %v", tasks[i].ID, tasks[i].Type, err)
			_ = tasks[i].ApplyStatus(domain.TaskStatusFailed, err.Error())
		} else {
			tasks[i].StreamID = streamID
			_ = tasks[i].ApplyStatus(domain.TaskStatusDispatched, "")
		}
		if _, uerr := s.tasks.Update(ctx, tasks[i]); uerr != nil {
			log.Printf("update task %s after publish failed: %v", tasks[i].ID, uerr)
		}
	}

	created.Status = domain.GoalStatusRunning
	updated, err := s.goals.Update(ctx, created)
	if err != nil {
		return domain.TaskGraph{}, fmt.Errorf("failed to update goal status: %w", err)
	}

	return domain.TaskGraph{Goal: updated, Tasks: tasks}, nil
}

// GetGoal returns a goal and all its sub-tasks.
func (s *OrchestratorService) GetGoal(ctx context.Context, id string) (domain.TaskGraph, error) {
	g, err := s.goals.GetByID(ctx, id)
	if err != nil {
		return domain.TaskGraph{}, err
	}
	if g == nil {
		return domain.TaskGraph{}, &domain.GoalNotFoundError{GoalID: id}
	}
	tasks, err := s.tasks.ListByGoal(ctx, id)
	if err != nil {
		return domain.TaskGraph{}, err
	}
	return domain.TaskGraph{Goal: *g, Tasks: tasks}, nil
}

// ListGoals returns goals for a user (or all if userID is empty).
func (s *OrchestratorService) ListGoals(ctx context.Context, userID string) ([]domain.Goal, error) {
	return s.goals.List(ctx, userID)
}

// UpdateTaskStatus is called by specialised agents when they progress.
func (s *OrchestratorService) UpdateTaskStatus(ctx context.Context, taskID, status, errMsg string) (domain.SubTask, error) {
	t, err := s.tasks.GetByID(ctx, taskID)
	if err != nil {
		return domain.SubTask{}, err
	}
	if t == nil {
		return domain.SubTask{}, &domain.TaskNotFoundError{TaskID: taskID}
	}
	if err := t.ApplyStatus(status, errMsg); err != nil {
		return domain.SubTask{}, err
	}
	updated, err := s.tasks.Update(ctx, *t)
	if err != nil {
		return domain.SubTask{}, err
	}
	if err := s.maybeCompleteGoal(ctx, t.GoalID); err != nil {
		log.Printf("complete-goal check failed: %v", err)
	}
	return updated, nil
}

// ── decomposition rules ────────────────────────────────────────────

// decompose maps a goal description to an ordered set of sub-tasks.
// Rule-based today; the seam is ready for an LLM planner later.
func (s *OrchestratorService) decompose(goal domain.Goal, description string) []domain.SubTask {
	desc := strings.ToLower(description)
	tasks := []domain.SubTask{}

	procurement := domain.NewSubTask(goal.ID, domain.TaskTypeProcurementSearch, map[string]interface{}{
		"project_id":  goal.ProjectID,
		"description": description,
	}, nil)
	tasks = append(tasks, procurement)

	if containsAny(desc, "factory", "plant", "workshop", "usine", "atelier", "warehouse", "entrepot") {
		layout := domain.NewSubTask(goal.ID, domain.TaskTypeLayoutGeneration, map[string]interface{}{
			"project_id":  goal.ProjectID,
			"description": description,
		}, []string{procurement.ID})
		tasks = append(tasks, layout)

		render := domain.NewSubTask(goal.ID, domain.TaskTypeThreeDRender, map[string]interface{}{
			"project_id": goal.ProjectID,
		}, []string{layout.ID})
		tasks = append(tasks, render)
	}

	cost := domain.NewSubTask(goal.ID, domain.TaskTypeCostAnalysis, map[string]interface{}{
		"project_id": goal.ProjectID,
	}, []string{procurement.ID})
	tasks = append(tasks, cost)

	return tasks
}

func (s *OrchestratorService) maybeCompleteGoal(ctx context.Context, goalID string) error {
	tasks, err := s.tasks.ListByGoal(ctx, goalID)
	if err != nil {
		return err
	}
	allDone := true
	anyFailed := false
	for _, t := range tasks {
		if t.Status == domain.TaskStatusFailed {
			anyFailed = true
		}
		if t.Status != domain.TaskStatusCompleted && t.Status != domain.TaskStatusFailed {
			allDone = false
		}
	}
	if !allDone {
		return nil
	}
	g, err := s.goals.GetByID(ctx, goalID)
	if err != nil || g == nil {
		return err
	}
	if anyFailed {
		g.Status = domain.GoalStatusFailed
	} else {
		g.Status = domain.GoalStatusCompleted
	}
	_, err = s.goals.Update(ctx, *g)
	return err
}

func (s *OrchestratorService) validate(data map[string]interface{}) error {
	errs := map[string]string{}
	for _, f := range []string{"user_id", "description"} {
		v, ok := data[f].(string)
		if !ok || strings.TrimSpace(v) == "" {
			errs[f] = fmt.Sprintf("%s is required", f)
		}
	}
	if len(errs) > 0 {
		return &domain.ValidationError{Errors: errs}
	}
	return nil
}

func containsAny(s string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// logTaskGraph satisfies the PB-025 acceptance criterion ("Task graph visible
// in agent log"). Logged in a stable, parseable format.
func logTaskGraph(g domain.Goal, tasks []domain.SubTask) {
	log.Printf("ORCHESTRATOR goal=%s status=%s tasks=%d", g.ID, g.Status, len(tasks))
	for _, t := range tasks {
		log.Printf("ORCHESTRATOR  └ task id=%s type=%s deps=%v status=%s",
			t.ID, t.Type, t.DependsOn, t.Status)
	}
}
