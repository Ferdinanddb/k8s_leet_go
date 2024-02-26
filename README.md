# K8s_leet_go

This is a project (hopefully not dead in one week) I do to learn the basics of Golang, in the context of building a backend API service proposing some of the basic functionalities of Leetcode.

My goal is for a user to be able to :
- Log in,
- Perform a POST request containing some code to be executed on a Kubernetes container,
- See the result of their requests,
- Access the history of their requests
- ...


I do this to :
- Improve my skills in :
    - golang,
    - kubernetes,
    - creating APIs,
- Getting familiar with :
    - Authentication (using token),
    - Caching,
    - Persisting users' results (DBs, Disks, ...),
    - Database schema organization,
    - How things could be optimized ?


## TODO List (unorderered) :

- [X] Explose the codebase in different subfolders
- [X] Implement middleware
- [X] Implement the token generator and token verification
  - [ ] Try to improve this ? Is it worth to cache the token ?
- [X] Implement the DB
  - [ ] Implement a table to host all the UserCodeRequest (userID, instanciationTime, requestUUID, languague, codeContent, workerStatus, outputResult)
- [X] Implement 2 go routines in _backend/api/postJob.go_ in order to :
  - [X] Push an event to Redis,
  - [X] Push the same event in the PostgreSQL table _UserCodeRequest_.
- [X] Start developping the asynq-worker go module to take care of suscribing to the Redis queue, handle the jobs, update corresp. rows in PostgreSQL table_UserCodeRequest_
  - [X] Use this [ASYNQ](https://github.com/hibiken/asynq) lib
  - [X] Code the logic in asynq_worker/task/utils/createJob.go in order to write the result in the DB and may be a cache per user ?
  - [ ] Code the logic in asynq_worker/task/utils/createJob.go in order to a cache per user ?
    - Real question here is : how can the user retrieve "automatically" the result of their POST request in this setup ? websocket ? cache ?
  - [ ] Implement another logic to handle another programming language like golang
  - [ ] Understand how to implement the UI that is showed in the ASYNQ repo
  - 
- [ ] May be try to implement a front-end in Javascript or Typescript ?


- [X] Implement helm chart to deploy postgres, redis, backend, asynq-worker, ... and their needed resources.
  - [ ] Improve both backend's and asynq_worker's `initContainers` to wait for the PostgreSQL and Redis pods to be up and running.
- [X] Respect the _least access privileges principle_ by creating a k8s service account for the backend and another service account with extended rights for creating k8s jobs, pods, retrieving logs, ... for the asynq-worker service.
- [X] Implement the history retrieval
  - [ ] Implement a GET endpoint to retrieve the last request made be the user
- [ ] Implement a caching layer
- [ ] Improve the scalability by adding a HPA (Horinzontal Pod Autoscaler)
- [ ] Create a Makefile for this project to be initialized, built, tested, and containerized



## Usage

### Prerequisites

- Install golang
- Have a kubernetes cluster deployed locally (I use k3s via Rancher Desktop)
    - Rancher uses containerd in my case behind the scene, all the container images are retrieve thanks to Rancher.
- Install helm


### Build the backend container image
- I use `nerdctl` which comes with Rancher Desktop :
```sh
nerdctl build --namespace=k8s.io -t backend ./backend/
```

### Build the asynq_worker container image
- I use `nerdctl` which comes with Rancher Desktop :
```sh
nerdctl build --namespace=k8s.io -t asynq_worker ./asynq_worker/
```

## Deploy the resources using Helm :

- Create the resources using
```sh
helm upgrade -i  k8s-leet-go ./_INFRA/k8s-leet-go
```

### Perform a test

- Once the helm chart got installed and every pods are up and running, we can port-forward the backend service to interact with it:
```sh
kubectl port-forward --namespace default svc/k8s-leet-go-backend 8080:80
```

- I did use Postman to perform a test because it is simpler since cookies are used to store the auth token.
  - On macOS : `brew install postman`

#### Inside Postman :

1. Make a POST request :
- URL :  `http://localhost:8080/auth/register`
- Inside _Body_ :
    - Select _raw_ and _JSON_ and fill in :
    ```json
    {
        "username": "test1",
        "password": "test1"
    }
    ```
- Click on Send

2. Make another POST request :
- URL :  `http://localhost:8080/auth/login`
- Inside _Body_ :
    - Select _raw_ and _JSON_ and fill in :
    ```json
    {
        "username": "test1",
        "password": "test1"
    }
    ```
- Click on Send
- If you take a look at the Response section, you should see that a new table (containing info including your token) appeared inside the _Cookies_ section.


3. Make another POST request :
- URL :  `http://localhost:8080/api/run_code`
- Inside _Body_ :
    - Select _raw_ and _JSON_ and fill in :
    ```json
    {
        "language": "python",
        "content": "class Solution:\n\tdef add(a,b):\n\t\treturn a + b\n\nprint(Solution.add(1,1))"
    }
    ```
- Click on Send

4. After a few seconds (~3secs), you can perform a GET request to see your results:
- URL : `http://localhost:8080/api/get_history`
- Click on Send
- You should get a response that looks like the following:
```json
{
    "data": [
        {
            "UserID": 1,
            "InstanciationTS": "2024-02-26T21:27:04.186798Z",
            "RequestUUID": "3ec55904-f647-4c4c-aaf6-b0ce32e9874a",
            "CodeContent": "class Solution:\n\tdef add(a,b):\n\t\treturn a + b\n\nprint(Solution.add(1,8))",
            "WorkerStatus": {
                "String": "success",
                "Valid": true
            },
            "OutputResult": {
                "String": "2\n",
                "Valid": true
            }
        }
    ]
}
```

## Delete the Helm chart to delete the resources :

- Delete the resources using
```sh
helm delete k8s-leet-go
kubectl delete pvc --all
```

---
<center><font color="blue">Et</font color="blue"> voi<font color="red">là</font color="red"></center>