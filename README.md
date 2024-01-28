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
- [ ] Create a Makefile for this project to be initialized, built, tested, and containerized
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
helm install postgresql-test oci://registry-1.docker.io/bitnamicharts/postgresql
export PGPASSWORD=$(kubectl get secret --namespace default postgresql-test -o jsonpath="{.data.postgres-password}" | base64 -d)
export PG_HOST=$(kubectl get svc --namespace default postgresql-test -o jsonpath='{.spec.clusterIP}')


# The following is not needed when deploying through YAML manifest file or Helm chart
kubectl port-forward --namespace default svc/postgresql-test 5432:5432 &
createdb --host 127.0.0.1 -U postgres  -p 5432 test -w
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
JWT_PRIVATE_KEY="TO_BE_CHANGED"
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

### Build the container image
- I use `nerdctl` which comes with Rancher Desktop :
```sh
nerdctl build --namespace=k8s.io -t backend .
```

### Personal notes :
- To create the resources for testing this repo on K8S :
```sh
kubectl apply -f .do_not_push/zz_test.yaml
kubectl port-forward --namespace default svc/test-k8s-svc 8080:80 &
kubectl delete -f .do_not_push/zz_test.yaml
```

---
<center><font color="blue">Et</font color="blue"> voi<font color="red">là</font color="red"></center>