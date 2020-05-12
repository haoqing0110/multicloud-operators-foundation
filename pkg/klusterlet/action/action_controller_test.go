package controllers

import (
	"testing"
	"time"

	tlog "github.com/go-logr/logr/testing"
	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	actionv1beta1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/action/v1beta1"
	"github.com/open-cluster-management/multicloud-operators-foundation/pkg/utils/rest"
)

var c client.Client

const (
	clusterActionName      = "name"
	clusterActionNamespace = "namespace"
)

func TestControllerReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	c = mgr.GetClient()

	kubework := NewKubeWorkSpec()
	clusteraction := NewClusterAction(clusterActionName, clusterActionNamespace, actionv1beta1.DeleteActionType, kubework)
	ar := NewActionReconciler(c, tlog.NullLogger{}, mgr.GetScheme(), nil, rest.NewFakeKubeControl(), false)

	ar.SetupWithManager(mgr)

	SetupTestReconcile(ar)
	ar.SetupWithManager(mgr)

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the object and expect the Reconcile
	err = c.Create(context.TODO(), clusteraction)

	g.Expect(err).NotTo(gomega.HaveOccurred())

	defer c.Delete(context.TODO(), clusteraction)

	time.Sleep(time.Second * 1)
}