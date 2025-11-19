package taskstore

import "context"

// TaskHandlerWithContext is an optional interface that task handlers can implement
// to receive context for cancellation support. This is backward compatible - handlers
// that don't implement this will continue to work using the standard Handle() method.
//
// Example usage:
//
//	type MyHandler struct {
//	    TaskHandlerBase
//	}
//
//	func (h *MyHandler) HandleWithContext(ctx context.Context) bool {
//	    select {
//	    case <-ctx.Done():
//	        h.LogInfo("Task cancelled")
//	        return false
//	    case <-time.After(5 * time.Second):
//	        h.LogSuccess("Task completed")
//	        return true
//	    }
//	}
type TaskHandlerWithContext interface {
	TaskHandlerInterface
	HandleWithContext(ctx context.Context) bool
}
