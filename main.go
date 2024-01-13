package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"bytes"
	"io"
	"time"

	// appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	// "k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)

	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job_req := &batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: "python",
							Image: "python",
							Command: []string{"python", "-c", "exec('''def add(a,b):\n  return a + b\n\nprint(add(1,1))''')"},
							
						},
					},
					RestartPolicy: "Never",
					
					
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-job",
		},
		
	}

	
	fmt.Println("Creating job...")
	job_res, job_err := jobsClient.Create(context.TODO(), job_req, metav1.CreateOptions{})

	
	// result, err := podsClient.Create(context.TODO(), pod_req, metav1.CreateOptions{})
	if job_err != nil {
		fmt.Println(job_err.Error())
		panic(job_err)
	}
	fmt.Printf("Created job %q.\n", job_res.GetObjectMeta().GetName())

	// wait for the pod to be ready

	watcher, watcher_err := jobsClient.Watch(context.TODO(), metav1.ListOptions{
		TypeMeta: job_req.TypeMeta,
		
	})
    if watcher_err != nil {
        panic(watcher_err)
    }

    defer watcher.Stop()

    looping: for {
		select {
			case event := <-watcher.ResultChan():
				job := event.Object.(*batchv1.Job)

				if job.Status.Failed >= 1 {   
					fmt.Printf("The POD \"%s\" failed", job_res.GetObjectMeta().GetName())
					break looping;
				} else if job.Status.Succeeded >= 1 {
					fmt.Printf("The POD \"%s\" succeeded", job_res.GetObjectMeta().GetName())
					break looping;
				} else {
					fmt.Println("Sleeping for 5 seconds waiting for pods to be running...")
					time.Sleep(5 * time.Second)
				}

			case <-context.TODO().Done():
				fmt.Printf("Exit from waitPodRunning for POD \"%s\" because the context is done", job_res.GetObjectMeta().GetName())
				break looping;
		}
	}


	
	list_pods, list_pods_err := podsClient.List(
		context.TODO(), 
		metav1.ListOptions{
			LabelSelector: "batch.kubernetes.io/job-name=test-job",
		},
	)
	if list_pods_err != nil {
		fmt.Println(list_pods_err.Error())
		panic(list_pods_err)
	}
    
	for _, v := range list_pods.Items {
		log := podsClient.GetLogs(v.Name, &apiv1.PodLogOptions{})
		podLogs, podlogs_err := log.Stream(context.TODO())
		if podlogs_err != nil {
			panic(podlogs_err)
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, bug_err := io.Copy(buf, podLogs)
		if bug_err != nil {
			panic(bug_err)
		}
		str := buf.String()

		fmt.Println(str)

		prompt()
		fmt.Println("Deleting pod...")
		deletePolicy := metav1.DeletePropagationForeground
		if err := podsClient.Delete(context.TODO(), v.Name, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		fmt.Println("Deleted pod.")

		fmt.Println("Deleting job...")
		job_deletePolicy := metav1.DeletePropagationForeground
		if job_err := jobsClient.Delete(context.TODO(), job_res.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &job_deletePolicy,
		}); job_err != nil {
			panic(err)
		}
		fmt.Println("Deleted pod.")
	}
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
