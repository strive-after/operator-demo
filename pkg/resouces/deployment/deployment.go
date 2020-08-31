package deployment

import (
	demo  "demo/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func New(app *demo.Frigate) *appsv1.Deployment {
	labels := map[string]string{"frigate.example.com/v1":app.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return  &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: app.Name,
			Namespace: app.Namespace,
			Labels: labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app,schema.GroupVersionKind{
					Group: demo.GroupVersion.Group,
					Version: demo.GroupVersion.Version,
					Kind: "Frigate",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: newContainers(app),
				},

			},
		},
	}
}


func newContainers(app *demo.Frigate) []corev1.Container {
	containers := make([]corev1.Container,0,0)
	containers = append(containers,corev1.Container{
		Name: app.Name,
		Image: app.Spec.Image,
		ImagePullPolicy: corev1.PullAlways,
	})
	return  containers
}