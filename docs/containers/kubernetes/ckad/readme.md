---
id: application_design_and_build
title: Application Design and Build
sidebar_label: Application Design and Build
sidebar_position: 1
---


## Building Container Images

### What is a container image?

An image is a lightweight, standalone file that contains the software and executables needed to run a container. Once we packaged the application into an image we can use it to run any number of containers.

Docker is one tool that can be used to create our own images. 

+ A __Dockerfile__ defines what is contained in the image
+ The `docker build` command builds an image using the Dockerfile


#### Example 1. Create a custom nginx image

1. Create a custom `index.html` file

  ```bash
  $ mkdir my-website;cd my-website; echo "Hello World!" >> my-website/index.html
  ```

2. Create a custom `Dockerfile`

  ```bash
  FROM nginx:stable
  
  COPY index.html /usr/share/nginx/html/
  ```

3. Build the image

  ```bash
  docker build -t my-website:0.0.1 .
  ```

4. Run the image

  ```
  docker run --rm --name my-website -d -p 8080:80 my-website:0.0.1
  ```

5. Validate

  ```
  curl localhost:8080
  ```

## Running Jobs and CronJobs

+ __Jobs__ are designed to run a containerized task succesfully to completion
+ __CronJobs__ run Jobs periodically according to a schedule
+ The `restartPolicy` for a Job or CronJob must be `OnFailure` or `Never`.
+ Use `activeDeadlineSeconds` in the Job spec to terminate the Job if it runs too long. 

#### Example 1. Create a simple Job that runs the command `echo This is a test!`

```bash
$ vi my-job.yaml
```

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: my-job
spec:
  template:
    spec:
      containers:
      - name: print
        image: busybox:stable
        command: ["echo", "This is a test!"]
    restartPolicy: Never
backoffLimit: 4
activeDeadlineSeconds: 10
```

+ `backoffLimit`: Number of retries if the job fails
+ `activeDeadlineSeconds`: Max number of seconds that job is allowed to run


Create the job with

```bash
kubectl apply -f my-job.yml
```

Check the status of the Job

```bash
kubectl get jobs
```
View the Job output. First, you will need to find the name of the Job's Pod

```
kubectl get pods
kubectl get logs $JOB_POD_NAME
```

#### Example 2. Create a simple CronJob that will run the Job tasks every minute

```bash
vi my-cronjob.yml
```

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: my-cronjob
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: print
            image: busybox:stable
            command: ["echo", "This is a test!"]
          restartPolicy: Never
      backoffLimit: 4
      activeDeadlineSeconds: 10
```

```bash
kubectl apply -f my-cronjob.yml
```

Check the CronJob status

```bash
kubectl get cronjob
```

Get the jobs associated as part of the cronjob executions

```bash
kubectl get jobs
```

## Links

[Dockerfile reference](https://docs.docker.com/engine/reference/builder/)