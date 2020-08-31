package service

import (
	demo "demo/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

func New(app *demo.Frigate) *corev1.Service {
	svctype := corev1.ServiceTypeClusterIP
	if app.Spec.NodePort != nil {
		svctype = corev1.ServiceTypeNodePort
	}
	labels := map[string]string{"frigate.example.com/v1":app.Name}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app,schema.GroupVersionKind{
					Group: demo.GroupVersion.Group,
					Version: demo.GroupVersion.Version,
					Kind: "Frigate",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: svctype,
			Selector: labels,
			Ports: getPort(app),
		},
	}
}

func getPort(app *demo.Frigate) []corev1.ServicePort {
	ports :=  make([]corev1.ServicePort,0,0)
	port := corev1.ServicePort{
		Port: app.Spec.Port,
		NodePort: *app.Spec.NodePort,
	}
	ports = append(ports,port)
	return ports
}














