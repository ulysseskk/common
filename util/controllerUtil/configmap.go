package controllerUtil

import corev1 "k8s.io/api/core/v1"

func ConfigMapDataEqual(cm1, cm2 *corev1.ConfigMap) bool {
	if cm1 == nil && cm2 == nil {
		return true
	}
	if cm1 == nil || cm2 == nil {
		return false
	}
	if len(cm1.Data) != len(cm2.Data) {
		return false
	}
	for k, v := range cm1.Data {
		if _, ok := cm2.Data[k]; !ok {
			return false
		}
		if cm2.Data[k] != v {
			return false
		}
	}
	return true
}
