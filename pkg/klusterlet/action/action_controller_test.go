package controllers

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	actionv1beta1 "github.com/stolostron/cluster-lifecycle-api/action/v1beta1"
	"github.com/stolostron/multicloud-operators-foundation/pkg/utils/rest"
)

var c client.Client

const (
	actionName      = "name"
	actionNamespace = "default"
)

func TestControllerReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{
		Metrics: metricserver.Options{BindAddress: "0"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	c = mgr.GetClient()

	kubework := NewKubeWorkSpec()
	action := NewAction(actionName, actionNamespace, actionv1beta1.DeleteActionType, kubework)
	ar := NewActionReconciler(c, ctrl.Log.WithName("controllers").WithName("ManagedClusterAction"), mgr.GetScheme(), nil, rest.NewFakeKubeControl(), false)

	ar.SetupWithManager(mgr)

	SetupTestReconcile(ar)
	ar.SetupWithManager(mgr)

	cancel, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		cancel()
		mgrStopped.Wait()
	}()

	// Create the object and expect the Reconcile
	err = c.Create(context.TODO(), action)

	g.Expect(err).NotTo(gomega.HaveOccurred())

	defer c.Delete(context.TODO(), action)

	time.Sleep(time.Second * 1)
}
