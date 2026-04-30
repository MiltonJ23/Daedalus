package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Daedalus/orchestrator-agent/internal/adapters/publishers"
	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
)

type memGoalRepo struct{ goals map[string]domain.Goal }

func newMemGoalRepo() *memGoalRepo { return &memGoalRepo{goals: map[string]domain.Goal{}} }
func (r *memGoalRepo) Create(_ context.Context, g domain.Goal) (domain.Goal, error) {
	r.goals[g.ID] = g
	return g, nil
}
func (r *memGoalRepo) GetByID(_ context.Context, id string) (*domain.Goal, error) {
	g, ok := r.goals[id]
	if !ok {
		return nil, nil
	}
	return &g, nil
}
func (r *memGoalRepo) List(_ context.Context, userID string) ([]domain.Goal, error) {
	out := []domain.Goal{}
	for _, g := range r.goals {
		if userID == "" || g.UserID == userID {
			out = append(out, g)
		}
	}
	return out, nil
}
func (r *memGoalRepo) Update(_ context.Context, g domain.Goal) (domain.Goal, error) {
	r.goals[g.ID] = g
	return g, nil
}

type memTaskRepo struct{ tasks map[string]domain.SubTask }

func newMemTaskRepo() *memTaskRepo { return &memTaskRepo{tasks: map[string]domain.SubTask{}} }
func (r *memTaskRepo) BulkCreate(_ context.Context, ts []domain.SubTask) error {
	for _, t := range ts {
		r.tasks[t.ID] = t
	}
	return nil
}
func (r *memTaskRepo) GetByID(_ context.Context, id string) (*domain.SubTask, error) {
	t, ok := r.tasks[id]
	if !ok {
		return nil, nil
	}
	return &t, nil
}
func (r *memTaskRepo) ListByGoal(_ context.Context, gid string) ([]domain.SubTask, error) {
	out := []domain.SubTask{}
	for _, t := range r.tasks {
		if t.GoalID == gid {
			out = append(out, t)
		}
	}
	return out, nil
}
func (r *memTaskRepo) Update(_ context.Context, t domain.SubTask) (domain.SubTask, error) {
	r.tasks[t.ID] = t
	return t, nil
}

func newSvc() (*OrchestratorService, *publishers.InMemory, *memGoalRepo, *memTaskRepo) {
	g := newMemGoalRepo()
	tr := newMemTaskRepo()
	pub := publishers.NewInMemory()
	return NewOrchestratorService(g, tr, pub), pub, g, tr
}

func TestSubmitGoal_Validation(t *testing.T) {
	svc, _, _, _ := newSvc()
	_, err := svc.SubmitGoal(context.Background(), map[string]interface{}{"description": ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve *domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestSubmitGoal_DecomposesAndPublishes(t *testing.T) {
	svc, pub, _, _ := newSvc()
	graph, err := svc.SubmitGoal(context.Background(), map[string]interface{}{
		"user_id":     "alice",
		"project_id":  "proj-1",
		"description": "Build a 200-unit/day biscuit factory in Douala",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if graph.Goal.Status != domain.GoalStatusRunning {
		t.Errorf("expected running, got %s", graph.Goal.Status)
	}
	// factory keyword triggers procurement + layout + render + cost = 4 tasks
	if len(graph.Tasks) != 4 {
		t.Fatalf("expected 4 tasks, got %d", len(graph.Tasks))
	}
	if len(pub.Published) != 4 {
		t.Errorf("expected 4 published tasks, got %d", len(pub.Published))
	}
	ids := map[string]struct{}{}
	for _, tk := range graph.Tasks {
		if tk.ID == "" {
			t.Error("task missing ID")
		}
		if _, dup := ids[tk.ID]; dup {
			t.Errorf("duplicate task ID %s", tk.ID)
		}
		ids[tk.ID] = struct{}{}
		if tk.Status != domain.TaskStatusDispatched {
			t.Errorf("task %s status=%s, expected dispatched", tk.ID, tk.Status)
		}
	}
}

func TestSubmitGoal_NoFactoryKeyword(t *testing.T) {
	svc, _, _, _ := newSvc()
	graph, _ := svc.SubmitGoal(context.Background(), map[string]interface{}{
		"user_id":     "bob",
		"description": "Source a CNC milling machine",
	})
	// no factory ⇒ procurement + cost only
	if len(graph.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(graph.Tasks))
	}
}

func TestUpdateTaskStatus_CompletesGoal(t *testing.T) {
	svc, _, _, _ := newSvc()
	graph, _ := svc.SubmitGoal(context.Background(), map[string]interface{}{
		"user_id":     "alice",
		"description": "Source a forklift",
	})
	for _, tk := range graph.Tasks {
		if _, err := svc.UpdateTaskStatus(context.Background(), tk.ID, domain.TaskStatusCompleted, ""); err != nil {
			t.Fatalf("update task: %v", err)
		}
	}
	final, _ := svc.GetGoal(context.Background(), graph.Goal.ID)
	if final.Goal.Status != domain.GoalStatusCompleted {
		t.Errorf("expected goal completed, got %s", final.Goal.Status)
	}
}

func TestUpdateTaskStatus_InvalidStatus(t *testing.T) {
	svc, _, _, _ := newSvc()
	graph, _ := svc.SubmitGoal(context.Background(), map[string]interface{}{
		"user_id":     "x",
		"description": "Source a generator",
	})
	_, err := svc.UpdateTaskStatus(context.Background(), graph.Tasks[0].ID, "weird", "")
	var is *domain.InvalidStatusError
	if !errors.As(err, &is) {
		t.Errorf("expected InvalidStatusError, got %T", err)
	}
}
