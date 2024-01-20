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
- [ ] Create a Makefile for this project to be initialized, built, tested, and Dockerized
- [ ] Explose the codebase in different subfolders
- [ ] Implement the token generator and token verification
- [ ] Implement the DB
- [ ] Implement middleware
- [ ] Implement helm chart to deploy the project on k8s (+ postgres, ...)
- [ ] Implement the history retrieval
- [ ] Implement a caching layer
- [ ] Improve the scalability



## Usage

### Prerequisites

- Install golang
- Have a kubernetes cluster deployed locally (I use k3s via Rancher Desktop)
    - Rancher uses containerd in my case behind the scene, all the container images are retrieve thanks to Rancher.
- Install helm


### Get the dependencies 
```sh
go get .
```

### PostgreSQL initialization
```sh
helm install postgresql-test oci://registry-1.docker.io/bitnamicharts/postgresql\n
export PGPASSWORD=$(kubectl get secret --namespace default postgresql-test -o jsonpath="{.data.postgres-password}" | base64 -d)
kubectl port-forward --namespace default svc/postgresql-test 5432:5432 &
createdb --host 127.0.0.1 -U postgres  -p 5432 test -w
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
```sh
curl http://localhost:8080/run_code \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"userid": "777","language": "python","content": "class Solution:\n\tdef add(a,b):\n\t\treturn a + b\n\nprint(Solution.add(1,1))"}'
```