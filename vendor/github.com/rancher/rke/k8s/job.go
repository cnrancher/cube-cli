package k8s

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobStatus struct {
	Completed bool
	Created   bool
}

func ApplyK8sSystemJob(jobYaml, kubeConfigPath string, k8sWrapTransport WrapTransport, timeout int, addonUpdated bool) error {
	job := v1.Job{}
	if err := decodeYamlResource(&job, jobYaml); err != nil {
		return err
	}
	if job.Namespace == metav1.NamespaceNone {
		job.Namespace = metav1.NamespaceSystem
	}
	k8sClient, err := NewClient(kubeConfigPath, k8sWrapTransport)
	if err != nil {
		return err
	}
	jobStatus, err := getK8sJobStatus(k8sClient, job.Name, job.Namespace)
	if err != nil {
		return err
	}
	// if the addon configMap is updated, or the previous job is not completed,
	// I will remove the existing job first, if any
	if addonUpdated || (jobStatus.Created && !jobStatus.Completed) {
		logrus.Debugf("[k8s] replacing job %s.. ", job.Name)
		if err := deleteK8sJob(k8sClient, job.Name, job.Namespace); err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
		} else { // ignoring NotFound errors
			time.Sleep(time.Second * 5)
		}
	}

	if _, err = k8sClient.BatchV1().Jobs(job.Namespace).Create(&job); err != nil {
		if apierrors.IsAlreadyExists(err) {
			logrus.Debugf("[k8s] Job %s already exists..", job.Name)
			return nil
		}
		return err
	}
	logrus.Debugf("[k8s] waiting for job %s to complete..", job.Name)
	return retryToWithTimeout(ensureJobCompleted, k8sClient, job, timeout)
}

func ensureJobCompleted(k8sClient *kubernetes.Clientset, j interface{}) error {
	job := j.(v1.Job)

	jobStatus, err := getK8sJobStatus(k8sClient, job.Name, job.Namespace)
	if err != nil {
		return fmt.Errorf("Failed to get job complete status: %v", err)
	}
	if jobStatus.Completed {
		logrus.Debugf("[k8s] Job %s completed successfully..", job.Name)
		return nil
	}
	return fmt.Errorf("Failed to get job complete status: %v", err)
}

func deleteK8sJob(k8sClient *kubernetes.Clientset, name, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return k8sClient.BatchV1().Jobs(namespace).Delete(
		name,
		&metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
}

func getK8sJob(k8sClient *kubernetes.Clientset, name, namespace string) (*v1.Job, error) {
	return k8sClient.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})
}

func getK8sJobStatus(k8sClient *kubernetes.Clientset, name, namespace string) (JobStatus, error) {
	existingJob, err := getK8sJob(k8sClient, name, namespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return JobStatus{}, nil
		}
		return JobStatus{}, err
	}
	for _, condition := range existingJob.Status.Conditions {
		if condition.Type == v1.JobComplete && condition.Status == corev1.ConditionTrue {
			return JobStatus{
				Created:   true,
				Completed: true,
			}, err
		}
	}
	return JobStatus{
		Created:   true,
		Completed: false,
	}, nil
}
