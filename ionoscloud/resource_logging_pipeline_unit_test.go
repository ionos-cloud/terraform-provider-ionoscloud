package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

// TestPipelineWaitFailedState checks that a FAILED pipeline stops the wait right away and that the error keeps its
// type through the retry helper, so the configured timeout is not reported for it.
func TestPipelineWaitFailedState(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceLoggingPipeline().Schema, map[string]any{})
	d.SetId("pipeline-id")

	failed := func(context.Context, *schema.ResourceData) (bool, error) {
		return false, fmt.Errorf("%w, pipeline ID: %s, status message: %q", logging.ErrPipelineFailed, "pipeline-id", "quota exceeded")
	}

	start := time.Now()
	err := utils.WaitForResourceToBeReady(context.Background(), d, failed)
	require.Error(t, err)
	assert.Less(t, time.Since(start), 5*time.Second, "the wait must stop on the first FAILED state")
	require.ErrorIs(t, err, logging.ErrPipelineFailed)
	assert.Contains(t, err.Error(), "quota exceeded")
	assert.Nil(t, pipelineWaitErrorContext(d, err, schema.TimeoutCreate), "a FAILED pipeline must not report the timeout")

	otherErr := errors.New("401 Unauthorized")
	errCtx := pipelineWaitErrorContext(d, otherErr, schema.TimeoutCreate)
	require.NotNil(t, errCtx)
	assert.Equal(t, d.Timeout(schema.TimeoutCreate).String(), errCtx.Timeout)
}
