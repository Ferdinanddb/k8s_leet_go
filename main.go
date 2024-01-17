package main

import (
	// "bufio"
	"context"
	"flag"
	"fmt"
	// "os"
	"path/filepath"
	"bytes"
	"io"
	"time"

	"net/http"

    "github.com/gin-gonic/gin"
	"log"

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
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.tmpl", gin.H{
			"title": "Welcome to my Project !",
		})
	})

	router.POST("/run_code", postK8sJob)

	log.Fatal(router.Run(":8080"))
}

type codeRequest struct {
    UserID string  `json:"userid"`
	Language     string  `json:"language"`
    Content  string  `json:"content"`
}

func postK8sJob(c *gin.Context) {
	var codeRequestToExec codeRequest

	if parsingErr := c.BindJSON(&codeRequestToExec); parsingErr != nil {
		return
	}

	jobResponse := createK8sJob(codeRequestToExec.Language, codeRequestToExec.Content)
	c.String(http.StatusOK, fmt.Sprintf("Result is %s\nCode executed is:\n\n%s\n", jobResponse, codeRequestToExec.Content))
}




func createK8sJob(language string, inputCode string) string {
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

	podsClient, jobsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault), clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	var containerName, containerImage string
	var containerCommand []string
	if language == "python" {
		containerName = "python"
		containerImage = "python:3.11-slim-bookworm"

		formattedExecCode := fmt.Sprintf("exec('''%v''')", inputCode)
		containerCommand = []string{"python", "-c", formattedExecCode}
	}
	

	job_req := &batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: containerName,
							Image: containerImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Command: containerCommand,
							
							// Command: []string{"python", "-c", "exec('''class Solution:\n\tdef add(a,b):\n\t\treturn a + b\n\nprint(Solution.add(1,1))''')"},
							// Command: []string{"python", "-c", "exec('''classdef add(a,b):\n  return a + b\n\nprint(add(1,1))''')"},
							
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

	
	log.Println("Creating job...")
	job_res, job_err := jobsClient.Create(context.TODO(), job_req, metav1.CreateOptions{})

	
	// result, err := podsClient.Create(context.TODO(), pod_req, metav1.CreateOptions{})
	if job_err != nil {
		log.Println(job_err.Error())
		panic(job_err)
	}
	log.Printf("Created job %q.\n", job_res.GetObjectMeta().GetName())

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

				if job.Status.Failed > 0 || job.Status.Succeeded > 0 {   
					log.Printf("The POD \"%s\" finished\n\n", job_res.GetObjectMeta().GetName())
					break looping;
				} else if job.Status.Active > 0 || (job.Status.Active == 0 && job.Status.Succeeded == 0 && job.Status.Failed == 0) {
					log.Printf("Status of the job is : %v\n", job.Status)
					time.Sleep(1 * time.Second)
				} else {
					log.Println("Something weird occured, job finished but not failed nor succeeded.")
					break looping;
				}

			case <-context.TODO().Done():
				log.Printf("Exit from waitPodRunning for POD \"%s\" because the context is done", job_res.GetObjectMeta().GetName())
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
		log.Println(list_pods_err.Error())
		panic(list_pods_err)
	}
	var strLogs string
	for _, v := range list_pods.Items {
		podLogResp := podsClient.GetLogs(v.Name, &apiv1.PodLogOptions{})
		podLogs, podlogs_err := podLogResp.Stream(context.TODO())
		if podlogs_err != nil {
			panic(podlogs_err)
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, bug_err := io.Copy(buf, podLogs)
		if bug_err != nil {
			panic(bug_err)
		}
		strLogs = buf.String()

		log.Println(strLogs)

		// prompt()
		log.Println("Deleting pod...")
		deletePolicy := metav1.DeletePropagationForeground
		if err := podsClient.Delete(context.TODO(), v.Name, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		log.Println("Deleted pod.")

		log.Println("Deleting job...")
		job_deletePolicy := metav1.DeletePropagationForeground
		if job_err := jobsClient.Delete(context.TODO(), job_res.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &job_deletePolicy,
		}); job_err != nil {
			panic(err)
		}
		log.Println("Deleted pod.")
	}

	return strLogs
}

// func prompt() {
// 	fmt.Printf("-> Press Return key to continue.")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		break
// 	}
// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println()
// }
