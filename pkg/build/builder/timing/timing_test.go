package timing

import (
	"context"
	"testing"

	"k8s.io/kubernetes/pkg/api/unversioned"

	"github.com/openshift/origin/pkg/build/api"
)

func TestRecordStageAndStepInfo(t *testing.T) {
	var stages []api.StageInfo

	stages = api.RecordStageAndStepInfo(stages, api.StageFetchInputs, api.StepFetchGitSource, unversioned.Now(), unversioned.Now())

	if len(stages) != 1 || len(stages[0].Steps) != 1 {
		t.Errorf("There should be 1 stage and 1 step, but instead there were %v stage(s) and %v step(s).", len(stages), len(stages[0].Steps))
	}

	stages = api.RecordStageAndStepInfo(stages, api.StagePullImages, api.StepPullBaseImage, unversioned.Now(), unversioned.Now())
	stages = api.RecordStageAndStepInfo(stages, api.StagePullImages, api.StepPullInputImage, unversioned.Now(), unversioned.Now())

	if len(stages) != 2 || len(stages[1].Steps) != 2 {
		t.Errorf("There should be 2 stages and 2 steps under the second stage, but instead there were %v stage(s) and %v step(s).", len(stages), len(stages[1].Steps))
	}

}

func TestAppendStageAndStepInfo(t *testing.T) {
	var stages []api.StageInfo
	var stagesToMerge []api.StageInfo

	stages = api.RecordStageAndStepInfo(stages, api.StagePullImages, api.StepPullBaseImage, unversioned.Now(), unversioned.Now())
	stages = api.RecordStageAndStepInfo(stages, api.StagePullImages, api.StepPullInputImage, unversioned.Now(), unversioned.Now())

	stagesToMerge = api.RecordStageAndStepInfo(stagesToMerge, api.StagePushImage, api.StepPushImage, unversioned.Now(), unversioned.Now())
	stagesToMerge = api.RecordStageAndStepInfo(stagesToMerge, api.StagePostCommit, api.StepExecPostCommitHook, unversioned.Now(), unversioned.Now())

	stages = api.AppendStageAndStepInfo(stages, stagesToMerge)

	if len(stages) != 3 {
		t.Errorf("There should be 3 stages, but instead there were %v stage(s).", len(stages))
	}

}

func TestTimingContextGetStages(t *testing.T) {
	ctx := NewContext(context.Background())

	RecordNewStep(ctx, api.StagePullImages, api.StepPullBaseImage, unversioned.Now(), unversioned.Now())
	RecordNewStep(ctx, api.StageFetchInputs, api.StepFetchGitSource, unversioned.Now(), unversioned.Now())
	RecordNewStep(ctx, api.StagePostCommit, api.StepExecPostCommitHook, unversioned.Now(), unversioned.Now())

	stages := GetStages(ctx)
	if len(stages) != 3 {
		t.Errorf("There should be 3 stages but instead there are %v stage(s).", len(stages))
	}
}
