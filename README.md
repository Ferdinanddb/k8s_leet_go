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
  - [ ] Code the logic in asynq_worker/task/utils/createJob.go in order to write the result in the DB and may be a cache per user ?
    - Real question here is : how can the user retrieve "automatically" the result of their POST request in this setup ? websocket ? cache ?
  - [ ] Implement another logic to handle another programming language like golang
  - [ ] Understand how to implement the UI that is showed in the ASYNQ repo
  - 
- [ ] May be try to implement a front-end in golang lol ?


- [ ] Implement helm chart to deploy postgres, redis, backend, asynq-worker, ... and their needed resources.
- [ ] Respect the _least access privileges principle_ by creating a k8s service account for the backend and another service account with extended rights for creating k8s jobs, pods, retrieving logs, ... for the asynq-worker service.
- [ ] Implement the history retrieval
- [ ] Implement a caching layer
- [ ] Improve the scalability
- [ ] Create a Makefile for this project to be initialized, built, tested, and containerized



## Usage

### Prerequisites

- Install golang
- Have a kubernetes cluster deployed locally (I use k3s via Rancher Desktop)
    - Rancher uses containerd in my case behind the scene, all the container images are retrieve thanks to Rancher.
- Install helm

### Do not forget to `cd` into backend lol !
```sh
cd backend
```

### Get the dependencies 
```sh
go get .
```

### PostgreSQL initialization
```sh
helm install postgresql-test oci://registry-1.docker.io/bitnamicharts/postgresql
export PGPASSWORD=$(kubectl get secret --namespace default postgresql-test -o jsonpath="{.data.postgres-password}" | base64 -d)
export PG_HOST=$(kubectl get svc --namespace default postgresql-test -o jsonpath='{.spec.clusterIP}')


# The following is not needed when deploying through YAML manifest file or Helm chart
kubectl port-forward --namespace default svc/postgresql-test 5432:5432 &
createdb --host 127.0.0.1 -U postgres  -p 5432 test -w
```

### Redis Cluster initialization
```sh
helm install redis-cluster-test oci://registry-1.docker.io/bitnamicharts/redis-cluster
export REDIS_PASSWORD=$(kubectl get secret --namespace "default" redis-cluster-test -o jsonpath="{.data.redis-password}" | base64 -d)
export REDIS_HOST=$(kubectl get svc --namespace default redis-cluster-test -o jsonpath='{.spec.clusterIP}')
```

### Create env variables
```sh
mkdir .do_not_push
cat <<EOF > .do_not_push/.env
# Database credentials
DB_HOST="127.0.0.1"
DB_USER="postgres"
PGPASSWORD="$PGPASSWORD"
DB_NAME="test"
DB_PORT="5432"

# Authentication credentials
TOKEN_TTL="2000"
JWT_PRIVATE_KEY="<TO_BE_CHANGED>"

REDIS_HOST="$REDIS_HOST"
REDIS_PORT="6379"
REDIS_PASSWORD="$REDIS_PASSWORD"
REDIS_DB="0"
EOF
```



### Build the project
```sh
go build -o ./backend
```

### Run the project
```sh
# After building the project, run:
./backend

# Without building the project, run:
go run .
```

### Perform a test

I did use Postman to perform a test because it is simpler since cookies are used to store the auth token.
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
- After a few seconds (~5secs), the following should appear in the Response section :
```
User is test1
Result is 2

Code executed is:

class Solution:
	def add(a,b):
		return a + b

print(Solution.add(1,1))
```

### Build the backend container image
- I use `nerdctl` which comes with Rancher Desktop :
```sh
nerdctl build --namespace=k8s.io -t backend .
```

### Build the asynq_worker container image
- I use `nerdctl` which comes with Rancher Desktop :
```sh
nerdctl build --namespace=k8s.io -t asynq_worker .
```

### Personal notes :
To create the resources for testing this repo on K8S (will only work for me as of now :$ ) :

- For backend
```sh
cd backend
kubectl apply -f .do_not_push/zz_test.yaml
kubectl port-forward --namespace default svc/test-k8s-svc 8080:80 &
kubectl delete -f .do_not_push/zz_test.yaml

kubectl exec -it redis-cluster-test-2 -- bash
# Inside the container, run :L
redis-cli -c -h redis-cluster-test -a $REDIS_PASSWORD
# Then to see the content of the first 15 events in the queue, run :
LRANGE "queue:new-code-request" -15 -1


helm delete redis-cluster-test
helm delete postgresql-test
```

- For asynq_worker
```sh
cd asynq_worker
kubectl apply -f .do_not_push/zz_test.yaml
```

---
<center><font color="blue">Et</font color="blue"> voi<font color="red">là</font color="red"></center>