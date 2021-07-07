package controllers

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("AzureClusterReconciler", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Reconcile an AzureCluster", func() {
		It("should reconcile and exit early due to the cluster not having an OwnerRef", func() {
			//ctx := context.TODO()
			////logListener := record.NewListener(testEnv.LogRecorder)
			////del := logListener.Listen()
			////defer del()
			//
			//instance := &capzv1alpha3.AzureCluster{ObjectMeta: metav1.ObjectMeta{Name: "randName", Namespace: "default"}}
			//Expect(testEnv.Create(ctx, instance)).To(Succeed())
			//defer func() {
			//	err := testEnv.Delete(ctx, instance)
			//	Expect(err).NotTo(HaveOccurred())
			//}()
			//
			//// Make sure the Cluster exists.
			//Eventually(logListener.GetEntries, 10*time.Second).
			//	Should(ContainElement(record.LogEntry{
			//		LogFunc: "Info",
			//		Values: []interface{}{
			//			"namespace",
			//			instance.Namespace,
			//			"azureCluster",
			//			randName,
			//			"msg",
			//			"Cluster Controller has not yet set OwnerRef",
			//		},
			//	}))
		})
	})
})
