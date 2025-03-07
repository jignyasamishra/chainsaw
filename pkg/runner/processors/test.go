package processors

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(context.Context, namespacer.Namespacer)
	CreateStepProcessor(namespacer.Namespacer, binding.Bindings, *cleaner, v1alpha1.TestSpecStep) StepProcessor
}

func NewTestProcessor(
	config v1alpha1.ConfigurationSpec,
	clusters clusters,
	clock clock.PassiveClock,
	summary *summary.Summary,
	testReport *report.TestReport,
	test discovery.Test,
	shouldFailFast *atomic.Bool,
	bindings binding.Bindings,
) TestProcessor {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return &testProcessor{
		config:         config,
		clusters:       clusters,
		clock:          clock,
		summary:        summary,
		testReport:     testReport,
		test:           test,
		shouldFailFast: shouldFailFast,
		bindings:       bindings,
		timeouts:       config.Timeouts.Combine(test.Spec.Timeouts),
	}
}

type testProcessor struct {
	config         v1alpha1.ConfigurationSpec
	clusters       clusters
	clock          clock.PassiveClock
	summary        *summary.Summary
	testReport     *report.TestReport
	test           discovery.Test
	shouldFailFast *atomic.Bool
	bindings       binding.Bindings
	timeouts       v1alpha1.Timeouts
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer) {
	t := testing.FromContext(ctx)
	if p.testReport != nil {
		t.Cleanup(func() {
			if t.Failed() {
				p.testReport.NewFailure("test failed")
			}
			p.testReport.MarkTestEnd()
		})
	}
	size := len("@cleanup")
	for i, step := range p.test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		if size < len(name) {
			size = len(name)
		}
	}
	if p.summary != nil {
		t.Cleanup(func() {
			if t.Skipped() {
				p.summary.IncSkipped()
			} else {
				if t.Failed() {
					p.summary.IncFailed()
				} else {
					p.summary.IncPassed()
				}
			}
		})
	}
	if p.test.Spec.Concurrent == nil || *p.test.Spec.Concurrent {
		t.Parallel()
	}
	if p.test.Spec.Skip != nil && *p.test.Spec.Skip {
		t.SkipNow()
	}
	if p.config.FailFast {
		if p.shouldFailFast.Load() {
			t.SkipNow()
		}
	}
	bindings, err := registerBindings(ctx, p.bindings, p.test.Spec.Bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	setupLogger := logging.NewLogger(t, p.clock, p.test.Name, fmt.Sprintf("%-*s", size, "@setup"))
	cleanupLogger := logging.NewLogger(t, p.clock, p.test.Name, fmt.Sprintf("%-*s", size, "@cleanup"))
	var namespace *corev1.Namespace
	_, cluster := p.clusters.client(p.test.Spec.Cluster)
	if cluster != nil {
		bindings = bindings.Register("$client", binding.NewBinding(cluster))
		if nspacer == nil || p.test.Spec.Namespace != "" {
			var ns corev1.Namespace
			if p.test.Spec.Namespace != "" {
				ns = client.Namespace(p.test.Spec.Namespace)
			} else {
				ns = client.PetNamespace()
			}
			namespace = &ns
		}
		if namespace != nil {
			object := client.ToUnstructured(namespace)
			bindings = bindings.Register("$namespace", binding.NewBinding(object.GetName()))
			if p.test.Spec.NamespaceTemplate != nil && p.test.Spec.NamespaceTemplate.Value != nil {
				template := v1alpha1.Any{
					Value: p.test.Spec.NamespaceTemplate.Value,
				}
				if merged, err := mutate.Merge(ctx, object, bindings, template); err != nil {
					t.FailNow()
				} else {
					object = merged
				}
				bindings = bindings.Register("$namespace", binding.NewBinding(object.GetName()))
			}
			nspacer = namespacer.New(cluster, object.GetName())
			setupCtx := logging.IntoContext(ctx, setupLogger)
			cleanupCtx := logging.IntoContext(ctx, cleanupLogger)
			if err := cluster.Get(setupCtx, client.ObjectKey(&object), object.DeepCopy()); err != nil {
				if !errors.IsNotFound(err) {
					// Get doesn't log
					setupLogger.Log(logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					t.FailNow()
				}
				if !cleanup.Skip(p.config.SkipDelete, p.test.Spec.SkipDelete, nil) {
					t.Cleanup(func() {
						operation := newOperation(
							false,
							timeout.Get(nil, p.timeouts.CleanupDuration()),
							opdelete.New(cluster, object, nspacer, false),
							nil,
							bindings,
						)
						operation.execute(cleanupCtx)
					})
				}
				if err := cluster.Create(logging.IntoContext(setupCtx, setupLogger), object.DeepCopy()); err != nil {
					t.FailNow()
				}
			}
		}
	}
	delay := p.config.DelayBeforeCleanup
	if p.test.Spec.DelayBeforeCleanup != nil {
		delay = p.test.Spec.DelayBeforeCleanup
	}
	cleaner := newCleaner(nspacer, delay)
	t.Cleanup(func() {
		cleaner.run(logging.IntoContext(ctx, cleanupLogger))
	})
	for i, step := range p.test.Spec.Steps {
		processor := p.CreateStepProcessor(nspacer, bindings, cleaner, step)
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		processor.Run(logging.IntoContext(ctx, logging.NewLogger(t, p.clock, p.test.Name, fmt.Sprintf("%-*s", size, name))))
	}
}

func (p *testProcessor) CreateStepProcessor(nspacer namespacer.Namespacer, bindings binding.Bindings, cleaner *cleaner, step v1alpha1.TestSpecStep) StepProcessor {
	var stepReport *report.TestSpecStepReport
	if p.testReport != nil {
		stepReport = report.NewTestSpecStep(step.Name)
		p.testReport.AddTestStep(stepReport)
	}
	return NewStepProcessor(p.config, p.clusters, nspacer, p.clock, p.test, step, stepReport, cleaner, bindings)
}
