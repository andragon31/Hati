package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	logger *log.Logger
	server *server.MCPServer
	root   string
}

func NewServer(logger *log.Logger, root string) *Server {
	srv := server.NewMCPServer("hati", "4.0.0")

	s := &Server{
		logger: logger,
		server: srv,
		root:   root,
	}

	s.registerAllTools()

	return s
}

func (s *Server) registerAllTools() {
	s.registerPlanTools()
	s.registerCheckpointTools()
	s.registerPhaseTools()
	s.registerFeedbackTools()
	s.registerRecordTools()
	s.registerSystemTools()
}

func (s *Server) registerPlanTools() {
	s.server.AddTool(mcp.NewTool("plan_create",
		mcp.WithDescription("Create a new plan with phases"),
		mcp.WithString("request", mcp.Required(), mcp.Description("User request")),
		mcp.WithString("module", mcp.Description("Primary module")),
	), s.handlePlanCreate)

	s.server.AddTool(mcp.NewTool("plan_get",
		mcp.WithDescription("Get plan details"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
	), s.handlePlanGet)

	s.server.AddTool(mcp.NewTool("plan_revise",
		mcp.WithDescription("Revise an existing plan"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
		mcp.WithString("changes", mcp.Description("Description of changes")),
	), s.handlePlanRevise)

	s.server.AddTool(mcp.NewTool("plan_abandon",
		mcp.WithDescription("Abandon a plan"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
		mcp.WithString("reason", mcp.Description("Reason for abandonment")),
	), s.handlePlanAbandon)

	s.server.AddTool(mcp.NewTool("plan_completeness",
		mcp.WithDescription("Check plan completeness score"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
	), s.handlePlanCompleteness)

	s.server.AddTool(mcp.NewTool("plan_quality",
		mcp.WithDescription("Get plan quality score"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
	), s.handlePlanQuality)
}

func (s *Server) registerCheckpointTools() {
	s.server.AddTool(mcp.NewTool("checkpoint_open",
		mcp.WithDescription("Open a checkpoint for approval"),
		mcp.WithString("type", mcp.Required(), mcp.Description("Checkpoint type"), mcp.Enum("pre", "post", "plan", "final")),
		mcp.WithString("phase_id", mcp.Description("Phase ID for phase checkpoints")),
	), s.handleCheckpointOpen)

	s.server.AddTool(mcp.NewTool("checkpoint_decide",
		mcp.WithDescription("Developer decision on checkpoint"),
		mcp.WithString("decision", mcp.Required(), mcp.Description("Decision"), mcp.Enum("approved", "rejected", "needs_revision")),
		mcp.WithString("feedback", mcp.Description("Feedback if rejected")),
	), s.handleCheckpointDecide)

	s.server.AddTool(mcp.NewTool("checkpoint_status",
		mcp.WithDescription("Get checkpoint status"),
		mcp.WithString("plan_id", mcp.Description("Plan ID")),
	), s.handleCheckpointStatus)
}

func (s *Server) registerPhaseTools() {
	s.server.AddTool(mcp.NewTool("phase_start",
		mcp.WithDescription("Start executing a phase"),
		mcp.WithString("phase_id", mcp.Required(), mcp.Description("Phase ID")),
	), s.handlePhaseStart)

	s.server.AddTool(mcp.NewTool("phase_report",
		mcp.WithDescription("Report phase completion"),
		mcp.WithString("phase_id", mcp.Required(), mcp.Description("Phase ID")),
		mcp.WithString("summary", mcp.Description("Phase summary")),
		mcp.WithString("why_this_approach", mcp.Description("Explainability report")),
		mcp.WithArray("files_changed", mcp.Description("Files changed"), mcp.WithStringItems()),
	), s.handlePhaseReport)
}

func (s *Server) registerFeedbackTools() {
	s.server.AddTool(mcp.NewTool("feedback_request",
		mcp.WithDescription("Request feedback from developer"),
		mcp.WithString("question", mcp.Required(), mcp.Description("Question for developer")),
	), s.handleFeedbackRequest)

	s.server.AddTool(mcp.NewTool("feedback_receive",
		mcp.WithDescription("Receive developer feedback"),
		mcp.WithString("feedback", mcp.Required(), mcp.Description("Feedback text")),
	), s.handleFeedbackReceive)

	s.server.AddTool(mcp.NewTool("feedback_escalate",
		mcp.WithDescription("Escalate to human review"),
		mcp.WithString("reason", mcp.Required(), mcp.Description("Reason for escalation")),
	), s.handleFeedbackEscalate)
}

func (s *Server) registerRecordTools() {
	s.server.AddTool(mcp.NewTool("record_list",
		mcp.WithDescription("List approval records"),
		mcp.WithString("status", mcp.Description("Filter by status")),
	), s.handleRecordList)

	s.server.AddTool(mcp.NewTool("record_get",
		mcp.WithDescription("Get approval record details"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
	), s.handleRecordGet)

	s.server.AddTool(mcp.NewTool("record_export",
		mcp.WithDescription("Export approval record"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
		mcp.WithString("format", mcp.Description("Format"), mcp.Enum("markdown", "json")),
	), s.handleRecordExport)
}

func (s *Server) registerSystemTools() {
	s.server.AddTool(mcp.NewTool("hati_status",
		mcp.WithDescription("Get Hati system status"),
	), s.handleHatiStatus)

	s.server.AddTool(mcp.NewTool("hati_stats",
		mcp.WithDescription("Get Hati statistics"),
		mcp.WithNumber("days", mcp.Description("Days to look back")),
	), s.handleHatiStats)

	s.server.AddTool(mcp.NewTool("hati_commit_info",
		mcp.WithDescription("Get commit info for plan traceability"),
		mcp.WithString("plan_id", mcp.Description("Plan ID")),
	), s.handleHatiCommitInfo)

	s.server.AddTool(mcp.NewTool("hati_register_commit",
		mcp.WithDescription("Register a commit with plan ID"),
		mcp.WithString("plan_id", mcp.Required(), mcp.Description("Plan ID")),
		mcp.WithString("commit_hash", mcp.Required(), mcp.Description("Commit hash")),
	), s.handleHatiRegisterCommit)

	s.server.AddTool(mcp.NewTool("module_hints",
		mcp.WithDescription("Get active module hints from Skoll"),
		mcp.WithArray("modules", mcp.Description("Module paths to check"), mcp.WithStringItems()),
	), s.handleModuleHints)

	s.server.AddTool(mcp.NewTool("spec_impact",
		mcp.WithDescription("Get specs affected by the plan"),
		mcp.WithString("plan_id", mcp.Description("Plan ID")),
	), s.handleSpecImpact)

	s.server.AddTool(mcp.NewTool("quality_snapshot",
		mcp.WithDescription("Get Quality Snapshot from Tyr"),
	), s.handleQualitySnapshot)

	s.server.AddTool(mcp.NewTool("learning_answer",
		mcp.WithDescription("Submit learning mode answer"),
		mcp.WithString("question_id", mcp.Required(), mcp.Description("Question ID")),
		mcp.WithString("answer", mcp.Required(), mcp.Description("Answer")),
	), s.handleLearningAnswer)
}

func (s *Server) handlePlanCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"plan_id":      "pln-" + fmt.Sprint(time.Now().Unix()),
		"phases":       []map[string]interface{}{},
		"completeness": 0.0,
		"spec_impact":  map[string]interface{}{},
		"module_hints": []map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePlanGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"plan_id":    "",
		"status":     "draft",
		"phases":     []map[string]interface{}{},
		"created_at": "",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePlanRevise(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"revised": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePlanAbandon(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"abandoned": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePlanCompleteness(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"score":   1.0,
		"checks":  []string{},
		"missing": []string{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePlanQuality(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"score":       1.0,
		"explanation": "",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleCheckpointOpen(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	checkpointType := getString(request.GetArguments(), "type")

	result := map[string]interface{}{
		"checkpoint_type":  checkpointType,
		"status":           "open",
		"quality_snapshot": map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleCheckpointDecide(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	decision := getString(request.GetArguments(), "decision")

	result := map[string]interface{}{
		"decision": decision,
		"recorded": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleCheckpointStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"status": "none",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePhaseStart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"status":     "started",
		"started_at": "",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handlePhaseReport(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"reported": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleFeedbackRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"question_id": "",
		"question":    "",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleFeedbackReceive(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"received": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleFeedbackEscalate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"escalated": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRecordList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"records": []map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRecordGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"plan_id":     "",
		"checkpoints": []map[string]interface{}{},
		"decisions":   []map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRecordExport(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"format":  "markdown",
		"content": "",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleHatiStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"version":        "4.0.0",
		"active_plans":   0,
		"quality_source": "tyr",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleHatiStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"total_plans":    0,
		"completed":      0,
		"abandoned":      0,
		"fast_approvals": 0,
		"avg_quality":    0.0,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleHatiCommitInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"plan_id": "",
		"message": "",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleHatiRegisterCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"registered": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleModuleHints(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"hints": []map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleSpecImpact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"affected_specs": []map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleQualitySnapshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"overall_quality_score": 1.0,
		"unit_tests":            []map[string]interface{}{},
		"e2e_tests":             []map[string]interface{}{},
		"sast":                  []map[string]interface{}{},
		"source":                "tyr",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleLearningAnswer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"accepted": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) RunStdio() error {
	return server.ServeStdio(s.server)
}

func getString(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getStringOrDefault(args map[string]interface{}, key, defaultVal string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getBoolOrDefault(args map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}
